package direct

import (
	"net/http"
	"net/mail"
	"unicode"

	"github.com/pmoieni/auth/models"
)

const (
	errInvalidEmail = "invalid email"
	errBadPassword  = "password must be at least 8 characters with 1 numeric digit"
)

func validateRegisterInfo(u *UserRegisterInfo) (err error) {
	err = validateEmail(u.Email)
	if err != nil {
		return
	}

	err = validatePassword(u.Password)
	if err != nil {
		return
	}
	return
}

func validateEmail(e string) error {
	if _, err := mail.ParseAddress(e); err != nil {
		return &models.ErrorResponse{Status: http.StatusBadRequest, Message: errInvalidEmail}
	}
	return nil
}

func validatePassword(p string) error {
	var (
		hasMinLen = false
		hasNumber = false
	)
	if len(p) >= 8 {
		hasMinLen = true
	}
	for _, char := range p {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !(hasMinLen && hasNumber) {
		return &models.ErrorResponse{Status: http.StatusBadRequest, Message: errBadPassword}
	}

	return nil
}
