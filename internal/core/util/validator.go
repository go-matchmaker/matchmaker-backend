package util

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	EmailRegexp     = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	IsValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	IsValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !IsValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

// PasswordPolicy is a function that checks if a password meets a certain criteria.
type PasswordPolicy func(string) error

func MinLength(n int) PasswordPolicy {
	return func(password string) error {
		if len(password) < n {
			return errors.New("password is too short")
		}
		return nil
	}
}

func MaxLength(n int) PasswordPolicy {
	return func(password string) error {
		if len(password) > n {
			return errors.New("password is too long")
		}
		return nil
	}
}

func MustContainLowercase(password string) error {
	for _, c := range password {
		if c >= 'a' && c <= 'z' {
			return nil
		}
	}

	return errors.New("password does not contain lowercase")
}

func MustContainUppercase(password string) error {
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			return nil
		}
	}

	return errors.New("password does not contain uppercase")
}

func MustContainNumber(password string) error {
	for _, c := range password {
		if c >= '0' && c <= '9' {
			return nil
		}
	}

	return errors.New("password does not contain number")
}

func MustContainSpecialChar(password string) error {
	for _, c := range password {
		if (c >= '!' && c <= '/') || (c >= ':' && c <= '@') || (c >= '[' && c <= '`') || (c >= '{' && c <= '~') {
			return nil
		}
	}

	return errors.New("password does not contain special character")
}
