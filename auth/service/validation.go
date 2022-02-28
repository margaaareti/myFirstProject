package service

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/pkg/errors"
)

func isEmailValid(input string) error {
	email := ev.NewSyntaxValidator()
	check := email.Validate(ev.NewInput(evmail.FromString(input)))

	if !check.IsValid() {
		return errors.Errorf("Email '%s' is not valid: ", input)
	}

	return nil
}
