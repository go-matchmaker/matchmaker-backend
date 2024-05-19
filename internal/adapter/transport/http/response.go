package http

import (
	"github.com/gofiber/fiber/v3"
)

type ValidationMessage struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Message          string               `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validation_errors,omitempty"`
	Error            string               `json:"error,omitempty"`
	Status           int                  `json:"status"`
}

type SuccessResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SuccessResponseWithoutData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type RedirectionResponse struct {
	Url string `json:"url"`
}

func (s *server) successResponse(c fiber.Ctx, data interface{}, message string, status int) error {
	if data != nil {
		return c.Status(status).JSON(successResponseBuilder(data, message, status))
	}

	return c.Status(status).JSON(successResponseWithoutDataBuilder(message, status))
}

func (s *server) errorResponse(c fiber.Ctx, message string, err error, validationErrors []*ValidationMessage, status int) error {
	if err != nil {
		return c.Status(status).JSON(errorResponseBuilder(message, err, status))
	}

	return c.Status(status).JSON(validationErrorsResponseBuilder(message, validationErrors, status))
}

func errorResponseBuilder(message string, err error, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Error:   err.Error(),
		Status:  status,
	}
}

func validationErrorsResponseBuilder(message string, validationErrors []*ValidationMessage, status int) *ErrorResponse {
	return &ErrorResponse{
		Message:          message,
		ValidationErrors: validationErrors,
		Status:           status,
	}
}

func successResponseBuilder[T any](data T, message string, status int) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		Data:    data,
		Message: message,
		Status:  status,
	}
}

func successResponseWithoutDataBuilder(message string, status int) *SuccessResponseWithoutData {
	return &SuccessResponseWithoutData{
		Message: message,
		Status:  status,
	}
}

func limitReachedFunc(c fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
		Message: fiber.ErrTooManyRequests.Message,
		Status:  fiber.StatusTooManyRequests,
	})
}
