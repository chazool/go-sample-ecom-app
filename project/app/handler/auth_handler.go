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

	// parse the request body
	// var req struct {
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }

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

// middleware for user authentication
func (h *AuthHandler) authenticationMiddleware(c *fiber.Ctx) error {

	user, err := h.authService.ParseToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	c.Locals("user", user)
	return c.Next()
}

/*
func authenticationMiddleware(c *fiber.Ctx) error {

	var dbcon *gorm.DB = db.GetDBConnection()

	// get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "missing Authorization header",
		})
	}

	// extract the JWT token from the header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid Authorization header",
		})
	}
	tokenString := parts[1]

	// parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil // replace with your JWT secret
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to parse token",
		})
	}

	// verify the JWT token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		// check if the user exists
		var user dto.User
		if err := dbcon.First(&user, uint(userID)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "invalid token",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to retrieve user",
			})
		}

		// set the user context
		c.Locals("user", user)

		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "invalid token",
	})
}
*/
