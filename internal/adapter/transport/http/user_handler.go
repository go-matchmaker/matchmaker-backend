package http

import (
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/converter"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/dto"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"time"
)

func (s *server) RegisterUser(c fiber.Ctx) error {
	reqBody := new(dto.UserRegisterRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return s.errorResponse(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	hashedPassword, err := util.HashPassword(reqBody.Password)
	if err != nil {
		return s.errorResponse(c, "error while trying to hash password", err, nil, fiber.StatusBadRequest)
	}

	userModel, err := converter.UserRegisterRequestToModel(reqBody, "user", hashedPassword)
	if err != nil {
		return s.errorResponse(c, "error while trying to convert user register to model", err, nil, fiber.StatusBadRequest)
	}
	userID, err := s.userService.Register(s.ctx, userModel)
	if err != nil {
		return s.errorResponse(c, "error while trying to register user", err, nil, fiber.StatusBadRequest)
	}

	zap.S().Info("User Registered Successfully! User:", userID)
	return s.successResponse(c, nil, "user registered successfully", fiber.StatusOK)
}

func (s *server) Login(c fiber.Ctx) error {
	reqBody := new(dto.UserLoginRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return s.errorResponse(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	ip := c.IP()
	session, err := s.userService.Login(c.Context(), reqBody.Email, reqBody.Password, ip)
	if err != nil {
		return s.errorResponse(c, "error while trying to login", err, nil, fiber.StatusBadRequest)
	}

	sessionID, err := s.sessionService.AddSession(c.Context(), session)
	if err != nil {
		return s.errorResponse(c, "error while trying to add session", err, nil, fiber.StatusBadRequest)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session",
		Value:   sessionID.String(),
		Expires: time.Now().Add(24 * time.Hour),
	})

	userResponse := dto.NewUserLoginRequestResponse(session)
	return s.successResponse(c, userResponse, "user logged in successfully", fiber.StatusOK)
}

func (s *server) UpdatePassword(c fiber.Ctx) error {
	reqBody := new(dto.PasswordChangeRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return s.errorResponse(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	userSession, ok := c.Locals(AuthSession).(*aggregate.Session)
	if !ok {
		return s.errorResponse(c, "session not found in context", nil, nil, fiber.StatusUnauthorized)
	}

	err := util.ComparePassword(reqBody.OldPassword, userSession.PasswordHash)
	if err != nil {
		return s.errorResponse(c, "old password not match", err, nil, fiber.StatusBadRequest)
	}

	hashedPassword, err := util.HashPassword(reqBody.NewPassword)
	if err != nil {
		return s.errorResponse(c, "error while trying to hash password", err, nil, fiber.StatusBadRequest)
	}

	changedUser := new(entity.User)
	changedUser.PasswordHash = hashedPassword
	changedUser.ID = userSession.ID
	updatedUser, err := s.userService.UpdatePassword(c.Context(), changedUser.ID, changedUser.PasswordHash)
	if err != nil {
		return s.errorResponse(c, "error while trying to update password", err, nil, fiber.StatusBadRequest)
	}

	userSession.PasswordHash = updatedUser.PasswordHash
	_, err = s.sessionService.AddSession(c.Context(), userSession)
	if err != nil {
		return s.errorResponse(c, "error while trying to add session", err, nil, fiber.StatusBadRequest)
	}

	return s.successResponse(c, nil, "password updated successfully", fiber.StatusOK)
}
