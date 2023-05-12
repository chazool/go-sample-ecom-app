package services

import (
	"errors"
	"fmt"
	"log"
	"sample_app/app/dto"
	"sample_app/app/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrFailToGenerateToken    = errors.New("failed to generate token")
	ErrUserAlreadyExists      = errors.New("user with the same email already exists")
	ErrFailToHashPassword     = errors.New("failed to hash password")
	ErrFailCreateUser         = errors.New("failed to create user record")
	ErrFailGenerateToken      = errors.New("failed to generate token")
	ErrFailToPassToken        = errors.New("failed to parse token")
	ErrInvalidToken           = errors.New("invalid token")
	ErrMissingAuthHeader      = errors.New("missing Authorization header")
	ErrInvalidAuthHeader      = errors.New("invalid Authorization header")
	ErrUnauthorizedAuthHeader = errors.New("unauthorized request")
)

type RequestType string

const (
	ReqSearch   RequestType = "search"
	ReqCreate   RequestType = "create"
	ReqRegAdmin RequestType = "reg_admin"
)

type AuthService interface {
	Register(user dto.User) (string, error)
	Login(email, password string) (string, error)
	ParseToken(authHeader string, reqType RequestType) (dto.User, error)
	hashPassword(password string) (string, error)
	verifyPassword(userPassword string, reqPassword string) error
	generateToken(user dto.User) (string, error)
	verifyToken(token string) (dto.UserRole, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService() AuthService {
	return &authService{
		userRepo: repository.NewUserRepository(),
	}
}

func (service *authService) Login(email, password string) (string, error) {

	// Find user with given email
	user, err := service.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Verify password
	err = service.verifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		return "", ErrFailToGenerateToken
	}

	return token, nil
}

func (service *authService) FindByEmail(email string) (dto.User, error) {
	log.Printf("Retrieving user by email %v\n", email)

	// Retrieve user from the database with given email
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return user, ErrProductNotFound
		}
		return user, err
	}

	log.Printf("Retrieved user with ID %v: %+v\n", email, user)
	return user, nil
}

func (service *authService) Register(user dto.User) (string, error) {

	// Check if user with the same email already exists
	_, err := service.FindByEmail(user.Email)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := service.hashPassword(user.Password)
	if err != nil {
		return "", ErrFailToHashPassword
	}

	// Create the user record
	user.Password = hashedPassword
	user, err = service.userRepo.Create(user)
	if err != nil {
		return "", ErrFailCreateUser
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		return "", ErrFailGenerateToken
	}

	return token, nil
}

func (service *authService) ParseToken(authHeader string, requestType RequestType) (dto.User, error) {

	var user dto.User
	if authHeader == "" {
		return user, ErrMissingAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return user, ErrInvalidAuthHeader
	}

	tokenString := parts[1]

	userRole, err := service.verifyToken(tokenString)
	if err != nil {
		return user, err
	}

	switch requestType {
	case ReqSearch:
		if userRole.Role == dto.RoleSystem {
			return user, ErrUnauthorizedAuthHeader
		}
	case ReqCreate:
		if userRole.Role != dto.RoleAdmin {
			return user, ErrUnauthorizedAuthHeader
		}
	}

	user, err = service.userRepo.FindById(int(userRole.ID))
	if err != nil {
		return user, ErrInvalidToken
	}

	return user, nil
}

// Verify password
func (service *authService) verifyPassword(userPassword string, reqPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqPassword)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// Hash the password
func (service *authService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (service *authService) generateToken(user dto.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["sub"] = dto.UserRole{
		ID:   user.ID,
		Role: user.Role,
	}

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service *authService) verifyToken(reqToken string) (dto.UserRole, error) {
	var user dto.UserRole
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return user, ErrFailToPassToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userData, ok := claims["sub"]
		if !ok {
			return user, ErrInvalidToken
		}

		mapstructure.Decode(userData, &user)
		return user, nil
	}

	return user, ErrInvalidToken
}
