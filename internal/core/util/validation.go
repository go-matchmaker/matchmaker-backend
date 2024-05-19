package util

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	SQLInjectionRegexp = regexp.MustCompile(`(?i)[\s\p{Zs}]+|(select|union|insert|delete|update|where|drop|create|from|set|or|and|like|case|when|between|exists|in|not|order|by|group|having|limit|offset|truncate|alter|add|constraint|default|distinct|index|primary|key|references|foreign|check|cascade|inner|join|left|outer|right|cross|natural|on|using|as|asc|desc|into|values|having|asc|desc|sleep|')`)
	EmailRegexp        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	IsValidFullName    = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func SafeSQL(inputs ...string) {
	for _, input := range inputs {
		if SQLInjectionRegexp.MatchString(input) {
			input = ""
		}
	}
}

func ValidateName(name string) error {
	if err := ValidateString(name, 3, 30); err != nil {
		return err
	}
	if !IsValidFullName(name) {
		return fmt.Errorf("Must contain only letters or spaces")
	}
	return nil
}

func ValidateSurname(name string) error {
	if err := ValidateString(name, 3, 10); err != nil {
		return err
	}
	if !IsValidFullName(name) {
		return fmt.Errorf("Must contain only letters or spaces")
	}
	return nil
}

func ValidateEmail(email string) error {
	if !EmailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber[0] == '0' {
		phoneNumber = phoneNumber[1:]
	}

	if phoneNumber != "" {
		if err := ValidateString(phoneNumber, 10, 10); err != nil {
			return err
		}
	}
	return nil

}

// PasswordPolicy is a function that checks if a password meets a certain criteria.
type PasswordPolicy func(string) error

func ValidatePassword(password string) error {
	passwordPolicies := []PasswordPolicy{
		minLength(8),
		maxLength(32),
		MustContainLowercase,
		MustContainUppercase,
		MustContainNumber,
		MustContainSpecialChar,
	}
	for _, policy := range passwordPolicies {
		if err := policy(password); err != nil {
			return errors.New("password does not meet the requirements" + err.Error())
		}
	}

	return nil
}

func minLength(n int) PasswordPolicy {
	return func(password string) error {
		if len(password) < n {
			return errors.New("password is too short")
		}
		return nil
	}
}

func maxLength(n int) PasswordPolicy {
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

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("Must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}
