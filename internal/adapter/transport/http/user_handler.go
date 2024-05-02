package http

import (
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/converter"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/dto"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

func (s *server) DeleteUser(c fiber.Ctx) error {
	userID := c.Params("id")
	fmt.Println("userID:", userID)
	id, err := uuid.Parse(userID)
	if err != nil {
		return s.responseFactory.Response(err, "error while trying to parse user id", fiber.StatusBadRequest)
	}
	err = s.userService.DeleteUser(s.ctx, id)
	if err != nil {
		return s.responseFactory.Response(err, "error while trying to delete user", fiber.StatusBadRequest)
	}

	zap.S().Info("User Deleted Successfully! User:", userID)
	return s.responseFactory.Response(nil, "user deleted successfully", fiber.StatusOK)

}

func (s *server) RegisterUser(c fiber.Ctx) error {
	reqBody := new(dto.UserRegister)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return s.responseFactory.Response(err, "error while trying to parse body", fiber.StatusBadRequest)
	}

	hashedPassword, err := util.HashPassword(reqBody.Password)
	if err != nil {
		return s.responseFactory.Response(err, "error while trying to hash password", fiber.StatusBadRequest)
	}

	userModel, err := converter.UserRegisterToModel(reqBody, entity.UserRoleCustomer, hashedPassword)
	if err != nil {
		return s.responseFactory.Response(err, "error while trying to convert user register to model", fiber.StatusBadRequest)
	}
	userID, err := s.userService.Register(s.ctx, userModel)
	if err != nil {
		return s.responseFactory.Response(err, "error while trying to register user", fiber.StatusBadRequest)
	}

	zap.S().Info("User Registered Successfully! User:", userID)
	return s.responseFactory.Response(nil, "user registered successfully", fiber.StatusOK)
}

//
// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }
//
// func (a *UserController) Login(c fiber.Ctx) error {
// 	var reqBody LoginRequest
// 	body := c.Body()
// 	if err := json.Unmarshal(body, &reqBody); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   "error while trying to parse body",
// 		})
// 	}
//
// 	token, err := a.userService.Login(c.Context(), reqBody.Email, reqBody.Password)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   "error while trying to login: " + err.Message,
// 		})
// 	}
//
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "token",
// 		Value:   token,
// 		Expires: time.Now().Add(24 * time.Hour),
// 	})
//
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"error": false,
// 		"msg":   "login successful",
// 		"data":  token,
// 	})
// }
