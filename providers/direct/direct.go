package direct

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pmoieni/auth/models"
	"github.com/pmoieni/auth/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	UserLoginCreds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UserRegisterInfo struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RefreshTokenDTO struct {
		Email    string `json:"email"`
		ClientID string `json:"client_id"`
	}

	ClientIDDTO struct {
		Email string `json:"email"`
	}
)

func (c *UserRegisterInfo) Register() (err error) {
	err = validateRegisterInfo(c)
	if err != nil {
		return
	}

	u := store.User{}
	u.Username = c.Username
	u.Email = c.Email
	// u.EmailVerified = false // deafult value when user registers for the first time
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hashedPassword)
	u.Picture = ""

	err = u.Create()
	var mysqlErr *mysql.MySQLError

	if err != nil {
		if err != nil && errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err = &models.ErrorResponse{Status: http.StatusConflict, Message: errUserAlreadyExists}
		}
		return
	}
	return
}

func (c *UserLoginCreds) Login() (tokens AuthTokensRes, err error) {
	// check if provided email is valid
	err = validateEmail(c.Email)
	if err != nil {
		return
	}

	// create user model to find from database
	u := store.User{
		Email: c.Email,
	}

	// get user info from database to create new authentication tokens
	userInfo, err := u.GetUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: errUserNotFound}
		}
		return
	}

	// check user password
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(c.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: errWrongPassword}
		}
		if errors.Is(err, bcrypt.ErrHashTooShort) {
			err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: errWrongPassword}
		}
		return
	}

	clientID := uuid.New().String()
	idClaims := idTokenClaims{
		username: u.Username,
		email:    u.Email,
		// emailVerified: u.EmailVerified,
		picture: u.Picture,
	}
	idToken, err := genIDToken(&idClaims)
	if err != nil {
		return
	}

	accessClaims := accessTokenClaims{
		email:    u.Email,
		clientID: clientID,
		// role:     "",
		// scope:    []string{"email", "profile"},
	}
	accessToken, err := genAccessToken(&accessClaims)
	if err != nil {
		return
	}

	refreshClaims := refreshTokenClaims{
		email:    u.Email,
		clientID: clientID,
	}
	rt, err := genRefreshToken(&refreshClaims)
	if err != nil {
		return
	}

	tokens = AuthTokensRes{
		IDToken:      idToken,
		AccessToken:  accessToken,
		RefreshToken: rt,
	}
	return
}

func RefreshToken(rt string) (tokens RefreshTokenRes, err error) {
	fmt.Println(rt + "\n")

	// validate and parse the refresh token
	rtPayload, err := parseRefreshTokenWithValidate(rt)
	if err != nil {
		return
	}

	privateClaims := rtPayload.PrivateClaims()

	email, ok := privateClaims["email"].(string)
	if !ok {
		err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		return
	}

	// role, ok := privateClaims["role"].(string)
	// if !ok {
	// 	err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
	// 	return
	// }

	// scope, ok := privateClaims["scope"].(string)
	// if !ok {
	// 	err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
	// 	return
	// }

	// check for token reuse
	reuse, err := isTokenUsed(rt)
	if err != nil {
		return
	}
	if reuse {
		err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		return
	}

	// save used token in database to detect token reuse when this toke is used again
	// since we have passed the token reuse check, this is the first time the token is being used
	// so we save it in used refresh tokens list to detect future token reuses
	err = saveRefreshToken(rt, email, rtPayload.Subject())
	if err != nil {
		return
	}

	// if the client id is revoked then the token is invalid and is reused by malicious user
	revoked, err := isClientIDRevoked(rtPayload.Subject())
	if err != nil {
		return
	}
	if revoked {
		err = &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		return
	}

	// generate new access token from previous access token claims
	// scopeArr := strings.Fields(scope)
	newATClaims := accessTokenClaims{
		email:    email,
		clientID: rtPayload.Subject(),
		// role:     role,
		// scope:    scopeArr,
	}
	newAT, err := genAccessToken(&newATClaims)
	if err != nil {
		return
	}

	// generate new refresh token form previous access token claims
	newRTClaims := refreshTokenClaims{
		email:    email,
		clientID: rtPayload.Subject(),
	}
	newRT, err := genRefreshToken(&newRTClaims)
	if err != nil {
		return
	}

	tokens = RefreshTokenRes{
		AccessToken:  newAT,
		RefreshToken: newRT,
	}

	return
}

func saveRefreshToken(token, email, clientID string) error {
	t := RefreshTokenDTO{
		Email:    email,
		ClientID: clientID,
	}
	b, err := json.Marshal(&t)
	if err != nil {
		return err
	}

	payload, err := parseRefreshToken(token)
	if err != nil {
		return err
	}

	d := time.Now().UTC().Sub(payload.IssuedAt())

	redisDTO := store.RedisDTO{
		Key:   token,
		Value: string(b),
		Exp:   d,
	}
	err = redisDTO.Set(store.RefreshTokenDB)
	if err != nil {
		return err
	}

	return nil
}

func isTokenUsed(token string) (bool, error) {
	// check if token is available in redis database
	// if it's not then token is not reused
	redisDTO := store.RedisDTO{
		Key: token,
	}
	v, err := redisDTO.Get(store.RefreshTokenDB)
	if err != nil {
		switch err {
		case redis.Nil:
			return false, nil
		default:
			return false, err
		}
	}

	// token is available in redis database which means it's reused
	// get token information containing client id and email of user
	t := RefreshTokenDTO{}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return false, err
	}

	// save client id in redis database to deny any refresh token with the sub value of revoked client id
	err = revokeClientID(t.ClientID, t.Email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func revokeClientID(clientID, email string) error {
	redisDTO := store.RedisDTO{
		Key:   clientID,
		Value: email,
		Exp:   RefreshTokenExpiry,
	}
	err := redisDTO.Set(store.ClientIDDB)
	if err != nil {
		return err
	}
	return nil
}

func isClientIDRevoked(clientID string) (bool, error) {
	// check if a key with client id exists
	// if the key exists it means that the client id is revoked and token should be denied
	// we don't need the email value here
	redisDTO := store.RedisDTO{
		Key: clientID,
	}
	_, err := redisDTO.Get(store.ClientIDDB)
	if err != nil {
		switch err {
		case redis.Nil:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
