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
	ReqSearch  RequestType = "search"
	ReqCreate  RequestType = "create"
	ReqRegUser RequestType = "reg_user"
)

type AuthService interface {
	Register(user, ctxUser dto.User) (string, error)
	Login(email, password string) (string, error)
	ParseToken(authHeader string, reqType RequestType) (dto.User, error)
	hashPassword(password string) (string, error)
	verifyPassword(userPassword string, reqPassword string) error
	generateToken(user dto.User) (string, error)
	verifyToken(token string) (dto.UserRole, error)
	checkUserRole(creatingRole, userRole dto.Role) error
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
	// Log function start
	log.Printf("Starting Login function for email: %v", email)
	defer log.Printf("Ending Login function for email: %v", email)

	// Find user with given email
	user, err := service.FindByEmail(email)
	if err != nil {
		log.Printf("Failed to find user with email %v: %v", email, err)
		return "", ErrInvalidCredentials
	}

	// Verify password
	err = service.verifyPassword(user.Password, password)
	if err != nil {
		log.Printf("Failed to verify password for user with email %v: %v", email, err)
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		log.Printf("Failed to generate token for user with email %v: %v", email, err)
		return "", ErrFailToGenerateToken
	}

	return token, nil
}

func (service *authService) FindByEmail(email string) (dto.User, error) {
	log.Printf("Starting FindByEmail for email %v\n", email)
	defer log.Printf("End FindByEmail for email %v\n", email)

	// Retrieve user from the database with given email
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("User not found with email %v\n", email)
			return user, ErrProductNotFound
		}
		log.Printf("Error retrieving user with email %v: %v\n", email, err)
		return user, err
	}

	log.Printf("Retrieved user with email %v: %+v\n", email, user)
	return user, nil
}

func (service *authService) Register(user, ctxUser dto.User) (string, error) {
	log.Println("Register function started")
	defer log.Println("Register function ended")

	// Check if the user creating the account has the appropriate role
	err := service.checkUserRole(user.Role, ctxUser.Role)
	if err != nil {
		log.Println("Failed to check user role")
		return "", err
	}

	// Check if user with the same email already exists
	_, err = service.FindByEmail(user.Email)
	if err == nil {
		log.Println("User with email already exists")
		return "", ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := service.hashPassword(user.Password)
	if err != nil {
		log.Println("Failed to hash password")
		return "", ErrFailToHashPassword
	}

	// Create the user record
	user.Password = hashedPassword
	user, err = service.userRepo.Create(user)
	if err != nil {
		log.Println("Failed to create user")
		return "", ErrFailCreateUser
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		log.Println("Failed to generate token")
		return "", ErrFailGenerateToken
	}

	log.Println("User successfully registered")
	return token, nil
}

func (service *authService) ParseToken(authHeader string, requestType RequestType) (dto.User, error) {
	log.Println("Starting ParseToken function")

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
		log.Println("Error verifying token:", err)
		return user, err
	}

	switch requestType {
	case ReqSearch:
		if userRole.Role == dto.RoleSystem {
			log.Println("Unauthorized auth header for request type:", requestType)
			return user, ErrUnauthorizedAuthHeader
		}
	case ReqCreate:
		if userRole.Role != dto.RoleAdmin {
			log.Println("Unauthorized auth header for request type:", requestType)
			return user, ErrUnauthorizedAuthHeader
		}
	case ReqRegUser:
		if !(userRole.Role == dto.RoleAdmin || userRole.Role == dto.RoleSystem) {
			log.Println("Unauthorized auth header for request type:", requestType)
			return user, ErrUnauthorizedAuthHeader
		}
	}

	user, err = service.userRepo.FindById(int(userRole.ID))
	if err != nil {
		log.Println("Error finding user by ID:", err)
		return user, ErrInvalidToken
	}

	log.Println("Ending ParseToken function")
	return user, nil
}

func (service *authService) verifyPassword(userPassword string, reqPassword string) error {
	log.Println("Verifying password...")
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("Incorrect password")
			return ErrInvalidCredentials
		}
		log.Printf("Error verifying password: %v", err)
		return ErrInvalidCredentials
	}
	log.Printf("Password verified")
	return nil
}

// Hash the password
func (service *authService) hashPassword(password string) (string, error) {
	log.Printf("Hashing password")
	defer log.Println("Finished hashing password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}

	return string(hashedPassword), nil
}

func (service *authService) generateToken(user dto.User) (string, error) {
	log.Printf("Generating token for user ID %d with role %s", user.ID, user.Role)
	defer log.Println("Token generation completed")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["sub"] = dto.UserRole{
		ID:   user.ID,
		Role: user.Role,
	}

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func (service *authService) verifyToken(reqToken string) (dto.UserRole, error) {
	log.Printf("Verifying token: %s", reqToken)

	var user dto.UserRole
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return user, ErrFailToPassToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userData, ok := claims["sub"]
		if !ok {
			log.Printf("Token claims do not contain 'sub' field")
			return user, ErrInvalidToken
		}

		if err := mapstructure.Decode(userData, &user); err != nil {
			log.Printf("Error decoding user data from token claims: %v", err)
			return user, ErrInvalidToken
		}

		log.Printf("Token verified, user: %+v", user)
		return user, nil
	}

	log.Printf("Token is not valid")
	return user, ErrInvalidToken
}

func (service *authService) checkUserRole(creatingRole, userRole dto.Role) error {
	log.Printf("Checking user role: %s, creating role: %s", userRole, creatingRole)
	switch userRole {
	case dto.RoleAdmin:
		if !(creatingRole == dto.RoleSystem || creatingRole == dto.RoleAdmin) {
			log.Printf("Unauthorized access: creating role %s not allowed for user role %s", creatingRole, userRole)
			return ErrUnauthorizedAuthHeader
		}
	case dto.RoleSystem:
		if creatingRole != dto.RoleUser {
			log.Printf("Unauthorized access: creating role %s not allowed for user role %s", creatingRole, userRole)
			return ErrUnauthorizedAuthHeader
		}
	default:
		log.Printf("Unauthorized access: user role %s not allowed", userRole)
		return ErrUnauthorizedAuthHeader
	}
	log.Printf("User role %s has authorized access", userRole)
	return nil
}
