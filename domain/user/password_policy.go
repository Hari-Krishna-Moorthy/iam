package user

import (
	"context"
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort        = errors.New("password is too short")
	ErrPasswordMissingNumber   = errors.New("password must contain at least one number")
	ErrPasswordMissingUpper    = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingSpecial  = errors.New("password must contain at least one special character")
)

type PasswordPolicy struct {
	TenantID        string
	MinLength       int
	RequireNumber   bool
	RequireUppercase bool
	RequireSpecial   bool
}

func (p *PasswordPolicy) Validate(password string) error {
	if len(password) < p.MinLength {
		return ErrPasswordTooShort
	}

	var (
		hasNumber  bool
		hasUpper   bool
		hasSpecial bool
	)

	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if p.RequireNumber && !hasNumber {
		return ErrPasswordMissingNumber
	}
	if p.RequireUppercase && !hasUpper {
		return ErrPasswordMissingUpper
	}
	if p.RequireSpecial && !hasSpecial {
		return ErrPasswordMissingSpecial
	}

	return nil
}

type PasswordPolicyRepository interface {
	GetByTenantID(ctx context.Context, tenantID string) (*PasswordPolicy, error)
	Save(ctx context.Context, policy *PasswordPolicy) error
}
