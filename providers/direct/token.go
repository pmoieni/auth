package direct

import (
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pmoieni/auth/config"
)

type (
	AuthTokensRes struct {
		IDToken      string
		AccessToken  string
		RefreshToken string
	}
	RefreshTokenRes struct {
		AccessToken  string
		RefreshToken string
	}
	idTokenClaims struct {
		username string
		email    string
		//	emailVerified bool
		picture string
	}
	accessTokenClaims struct {
		email    string
		clientID string
		//	role     string
		// scope []string
	}
	refreshTokenClaims struct {
		clientID string
	}
)

const (
	IDTokenExpiry      = time.Hour * 10
	AccessTokenExpiry  = time.Minute * 5
	RefreshTokenExpiry = time.Hour * 24 * 7
)

func genIDToken(c *idTokenClaims) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer("http://localhost:8080").
		Subject(c.email).
		Audience([]string{"http://localhost:3000"}).
		IssuedAt(time.Now().UTC()).
		Expiration(time.Now().UTC().Add(IDTokenExpiry)).
		Claim("username", c.username).
		Claim("email", c.email).
		//		Claim("email_verified", c.emailVerified).
		Claim("picture", c.picture).
		Build()
	if err != nil {
		return "", err
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(config.JWTIDTokenSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func genAccessToken(c *accessTokenClaims) (string, error) {
	// scope := strings.Join(c.scope, " ")
	token, err := jwt.NewBuilder().
		Issuer("http://localhost:8080").
		Subject(c.clientID).
		Audience([]string{"http://localhost:3000"}).
		IssuedAt(time.Now().UTC()).
		Expiration(time.Now().UTC().Add(AccessTokenExpiry)).
		Claim("email", c.email).
		// Claim("scope", scope).
		// Claim("role", c.role).
		Build()
	if err != nil {
		return "", err
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(config.JWTAccessTokenSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func genRefreshToken(c *refreshTokenClaims) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer("http://localhost:8080").
		Subject(c.clientID).
		Audience([]string{"http://localhost:3000"}).
		IssuedAt(time.Now().UTC()).
		Expiration(time.Now().UTC().Add(RefreshTokenExpiry)).
		Build()
	if err != nil {
		return "", err
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(config.JWTRefreshTokenSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func parseAccessToken(token string) (payload jwt.Token, err error) {
	c, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	payload, err = jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, []byte(c.JWTAccessTokenSecret)))
	return
}

func parseRefreshToken(token string) (payload jwt.Token, err error) {
	c, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	payload, err = jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, []byte(c.JWTRefreshTokenSecret)))
	return
}

// checks access token validity and returns its payload
func parseAccessTokenWithValidate(token string) (payload jwt.Token, err error) {
	c, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	payload, err = jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, []byte(c.JWTAccessTokenSecret)), jwt.WithValidate(true))
	return
}

// checks access token validity and returns its payload
func parseRefreshTokenWithValidate(token string) (payload jwt.Token, err error) {
	c, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	payload, err = jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, []byte(c.JWTRefreshTokenSecret)), jwt.WithValidate(true))
	return
}
