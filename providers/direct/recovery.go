package direct

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pmoieni/auth/config"
	"github.com/pmoieni/auth/models"
	"github.com/pmoieni/auth/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordResetTokenExpiry = time.Minute * 10
)

type (
	PasswordResetReq struct {
		Email       string `json:"email"`
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	PasswordResetTokenReq struct {
		Email string `json:"email"`
	}

	PasswordResetMessage struct {
		Token string
	}
)

func ResetPassword(r *PasswordResetReq) error {
	getTokenRedisDTO := store.RedisDTO{
		Key: r.Token,
	}
	v, err := getTokenRedisDTO.Get(store.PasswordResetTokenDB)
	if err != nil {
		switch err {
		case redis.Nil:
			return &models.ErrorResponse{Status: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		default:
			return err
		}
	}

	// remove token since it's already used
	delTokenRedisDTO := store.RedisDTO{
		Key: r.Token,
	}
	err = delTokenRedisDTO.Delete(store.PasswordResetTokenDB)
	if err != nil {
		switch err {
		case redis.Nil:
			return &models.ErrorResponse{Status: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
		default:
			return err
		}
	}

	tInfo := PasswordResetTokenReq{}
	err = json.Unmarshal([]byte(v), &tInfo)
	if err != nil {
		return err
	}

	// check if the token belongs to the email which is trying to reset the token
	if r.Email != tInfo.Email {
		return &models.ErrorResponse{Status: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}
	}

	u := store.User{
		Email: r.Email,
	}
	userInfo, err := u.GetUser()
	if err != nil {
		return err
	}

	// check if new password is valid
	err = validatePassword(r.NewPassword)
	if err != nil {
		return &models.ErrorResponse{Status: http.StatusForbidden, Message: errBadPassword}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updateUserInfo := store.User{
		Password: string(hashedPassword),
	}

	err = userInfo.Update(&updateUserInfo)
	if err != nil {
		return err
	}

	return nil
}

func (r *PasswordResetTokenReq) SendPasswordResetToken() (err error) {
	// TODO: check if user exists. this approach is not working
	u := store.User{
		Email: r.Email,
	}
	exists, err := u.CheckIfExists()
	if err != nil {
		return
	}
	if !exists {
		// status 404 if user with such info doesn't exist
		err = &models.ErrorResponse{Status: http.StatusNotFound, Message: errUserNotFound}
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	token, err := createRandomNumericalToken(6)
	if err != nil {
		return
	}

	err = savePasswordResetToken(token, r.Email)
	if err != nil {
		return
	}

	smtpEmail := config.SMTPEmail
	smtpPassword := config.SMTPPassword

	receiver := []string{
		r.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	t, err := template.ParseFiles("./templates/password_recovery_template.html")
	if err != nil {
		return
	}

	body := bytes.Buffer{}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Bellamy Labs Account Password Reset \n%s\n\n", mimeHeaders)))

	t.Execute(&body, PasswordResetMessage{
		Token: token,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, receiver, body.Bytes())
	return
}

func savePasswordResetToken(token, email string) (err error) {
	t := PasswordResetTokenReq{
		Email: email,
	}
	b, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	redisDTO := store.RedisDTO{
		Key:   token,
		Value: string(b),
		Exp:   PasswordResetTokenExpiry,
	}
	err = redisDTO.Set(store.PasswordResetTokenDB)
	if err != nil {
		return err
	}

	return nil
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func createRandomNumericalToken(length int) (string, error) {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
