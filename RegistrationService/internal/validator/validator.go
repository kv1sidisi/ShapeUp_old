package validator

import (
	regv1 "RegistrationService/api/pb"
	"errors"
	"github.com/asaskevich/govalidator"
	"strings"
)

func ValidateRegisterRequest(req *regv1.RegisterRequest) error {
	// Validate email
	if !govalidator.IsEmail(req.GetEmail()) {
		return errors.New("incorrect email address")
	}

	// Validate password length
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Validate spaces in password
	if strings.Contains(req.GetPassword(), " ") {
		return errors.New("password must not contain spaces")
	}

	return nil
}
