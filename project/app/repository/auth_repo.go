package repository

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/pkg/config/db"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user dto.User) (dto.User, error)
	FindById(id int) (dto.User, error)
	FindByEmail(email string) (dto.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: db.GetDBConnection(),
	}
}

func (r *userRepository) FindById(id int) (dto.User, error) {
	log.Printf("Finding user by ID %d...", id)
	var user dto.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with ID %d not found", id)
			return user, ErrRecordNotFound
		}
		log.Printf("Error finding user by ID %d: %s", id, result.Error)
		return user, result.Error
	}
	log.Printf("Found user by ID %d: %+v", id, user)
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (dto.User, error) {
	log.Printf("FindByEmail started for email: %s", email)

	var user dto.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found for email: %s", email)
			return user, ErrRecordNotFound
		}
		log.Printf("Error while finding user for email: %s, err: %v", email, err)
		return user, err
	}

	log.Printf("Found user for email: %s, user: %+v", email, user)
	return user, nil
}

func (r *userRepository) Create(user dto.User) (dto.User, error) {
	log.Println("Create method started")

	// Create the user
	if err := r.db.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return user, err
	}
	log.Printf("Created User: %+v", user)

	log.Println("Create method ended")
	return user, nil
}
