package http

import (
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/dto"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"sync"
)

var (
	validations *validator.Validate
	once        sync.Once
)

func (s *server) registerValidate(c fiber.Ctx) error {
	userDto := new(dto.UserRegisterRequest)
	body := c.Body()
	err := json.Unmarshal(body, &userDto)
	if err != nil {
		return s.errorResponse(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}
	util.SafeSQL(userDto.Name, userDto.Surname, userDto.Email, userDto.PhoneNumber, userDto.Password)
	validationMessage := make([]*ValidationMessage, 0)
	if err := util.ValidateName(userDto.Name); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "name",
			Tag:         "name",
			Message:     err.Error(),
		})
	}
	if err := util.ValidateSurname(userDto.Surname); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "surname",
			Tag:         "surname",
			Message:     err.Error(),
		})
	}
	if err := util.ValidateEmail(userDto.Email); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "email",
			Tag:         "email",
			Message:     err.Error(),
		})
	}
	if err := util.ValidatePhoneNumber(userDto.PhoneNumber); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "phone_number",
			Tag:         "phone_number",
			Message:     err.Error(),
		})
	}
	if err := util.ValidatePassword(userDto.Password); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "password",
			Tag:         "password",
			Message:     err.Error(),
		})
	}

	if len(validationMessage) > 0 {
		return s.errorResponse(c, "validation failed", nil, validationMessage, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}

func (s *server) loginValidate(c fiber.Ctx) error {
	data := new(dto.UserLoginRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return s.errorResponse(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	util.SafeSQL(data.Email, data.Password)
	validationMessage := make([]*ValidationMessage, 0)
	if err := util.ValidateEmail(data.Email); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "email",
			Tag:         "email",
			Message:     err.Error(),
		})
	}
	if err := util.ValidatePassword(data.Password); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "password",
			Tag:         "password",
			Message:     err.Error(),
		})
	}

	if len(validationMessage) > 0 {
		return s.errorResponse(c, "validation failed", nil, validationMessage, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}

func (s *server) adminLoginValidate(c fiber.Ctx) error {
	data := new(dto.UserLoginRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return s.errorResponse(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	util.SafeSQL(data.Email, data.Password)
	validationMessage := make([]*ValidationMessage, 0)
	if err := util.ValidateEmail(data.Email); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "email",
			Tag:         "email",
			Message:     err.Error(),
		})
	}
	if err := util.ValidatePassword(data.Password); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "password",
			Tag:         "password",
			Message:     err.Error(),
		})
	}

	if len(validationMessage) > 0 {
		return s.errorResponse(c, "validation failed", nil, validationMessage, fiber.StatusUnprocessableEntity)
	}

	c.Locals("admin-login", data)
	return c.Next()
}

func (s *server) updatePasswordValidate(c fiber.Ctx) error {
	data := new(dto.PasswordChangeRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return s.errorResponse(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	util.SafeSQL(data.OldPassword, data.NewPassword)
	validationMessage := make([]*ValidationMessage, 0)
	if err := util.ValidatePassword(data.OldPassword); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "old_password",
			Tag:         "old_password",
			Message:     err.Error(),
		})
	}
	if err := util.ValidatePassword(data.NewPassword); err != nil {
		validationMessage = append(validationMessage, &ValidationMessage{
			FailedField: "new_password",
			Tag:         "new_password",
			Message:     err.Error(),
		})
	}

	if len(validationMessage) > 0 {
		return s.errorResponse(c, "validation failed", nil, validationMessage, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}
