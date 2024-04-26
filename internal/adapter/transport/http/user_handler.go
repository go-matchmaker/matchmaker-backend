package http

import (
	"github.com/bulutcan99/company-matcher/internal/adapter/converter"
	"github.com/bulutcan99/company-matcher/internal/core/domain/entity"
	"github.com/bulutcan99/company-matcher/internal/core/dto"
	"github.com/bulutcan99/company-matcher/internal/core/util"
	"go.uber.org/zap"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

func (s *server) RegisterUser(c fiber.Ctx) error {
	var reqBody *dto.UserRegister
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

	hashedPassword, err := util.HashPassword(reqBody.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while hashing user password: " + err.Error(),
		})
	}

	userModel, err := converter.UserRegisterToModel(reqBody, entity.UserRoleCustomer, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while getting user model: " + err.Error(),
		})
	}
	userID, err := s.userService.Register(s.ctx, userModel)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to register user: " + err.Error(),
		})
	}

	zap.S().Info("User Registered Successfully! User:", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "user registered successfully",
		"data":  userID,
	})
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
