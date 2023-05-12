package repository

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/pkg/config/db"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(product dto.User) (dto.User, error)
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
	var user dto.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, ErrRecordNotFound
		}
		return user, result.Error
	}
	log.Printf("Found user by ID %d: %+v", id, user)
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (dto.User, error) {
	var user dto.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrRecordNotFound
		}
		return user, err
	}

	return user, nil

}

func (r *userRepository) Create(user dto.User) (dto.User, error) {
	// Create the user
	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}
	log.Printf("Created User: %+v", user)
	return user, nil
}
