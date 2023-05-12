package repository

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/pkg/config/db"

	"gorm.io/gorm"
)

type InteractionRepository interface {
	Create(product dto.Interaction) (dto.Interaction, error)
	RecentInteractions(user, limit uint) ([]dto.Interaction, error)
}

type interactionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository() InteractionRepository {
	return &interactionRepository{
		db: db.GetDBConnection(),
	}
}

func (r *interactionRepository) RecentInteractions(user, limit uint) ([]dto.Interaction, error) {

	// Get user's recent interactions
	var recentInteractions []dto.Interaction
	if err := r.db.Where("user_id = ?", user).Order("created_at desc").Limit(int(limit)).Find(&recentInteractions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return recentInteractions, ErrRecordNotFound
		}
		return recentInteractions, err
	}

	return recentInteractions, nil

}

func (r *interactionRepository) Create(interaction dto.Interaction) (dto.Interaction, error) {
	// Create the interaction record
	if err := r.db.Create(&interaction).Error; err != nil {
		return interaction, err
	}
	log.Printf("Created interaction: %+v", interaction)
	return interaction, nil
}
