package handler

import (
	"errors"
	"sample_app/app/dto"
	"sample_app/app/services"

	"github.com/gofiber/fiber/v2"
	///"github.com/golang-jwt/jwt"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler() AuthHandler {
	return AuthHandler{
		authService: services.NewAuthService(),
	}
}

// handle user registration
func (h *AuthHandler) register(c *fiber.Ctx) error {

	// Parse request body
	var user dto.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Message: "invalid request",
		})

	}

	token, err := h.authService.Register(user)
	if err != nil {
		status := fiber.StatusInternalServerError
		if errors.Is(err, services.ErrUserAlreadyExists) {
			status = fiber.StatusConflict

		}
		return c.Status(status).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(dto.Response{
		Data: token,
	})
}

// handle user authentication
func (h *AuthHandler) login(c *fiber.Ctx) error {

	var user dto.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Message: "invalid request body",
		})

	}

	token, err := h.authService.Login(user.Email, user.Password)
	if err != nil {
		status := fiber.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidCredentials) {
			status = fiber.StatusUnauthorized
		}
		return c.Status(status).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(dto.Response{
		Data: token,
	})

}

// // middleware for user authentication
// func (h *AuthHandler) authenticationMiddleware(c *fiber.Ctx) error {

// 	user, err := h.authService.ParseToken(c.Get("Authorization"), services.ReqRegAdmin)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(dto.Response{
// 			Message: err.Error(),
// 		})
// 	}

// 	c.Locals("user", &user)
// 	return c.Next()
// }

// middleware for user authentication
func (h *AuthHandler) searchAuthenticationMiddleware(c *fiber.Ctx) error {

	user, err := h.authService.ParseToken(c.Get("Authorization"), services.ReqSearch)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	c.Locals("user", &user)
	return c.Next()
}

// middleware for user authentication
func (h *AuthHandler) createAuthenticationMiddleware(c *fiber.Ctx) error {

	user, err := h.authService.ParseToken(c.Get("Authorization"), services.ReqCreate)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	c.Locals("user", &user)
	return c.Next()
}
