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
	log.Printf("RecentInteractions started for user %d with limit %d", user, limit)

	// Get user's recent interactions
	var recentInteractions []dto.Interaction
	if err := r.db.Where("user_id = ?", user).Order("created_at desc").Limit(int(limit)).Find(&recentInteractions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No recent interactions found for user %d", user)
			return recentInteractions, ErrRecordNotFound
		}
		log.Printf("Error while finding recent interactions for user %d, err: %v", user, err)
		return recentInteractions, err
	}

	log.Printf("Found %d recent interactions for user %d: %+v", len(recentInteractions), user, recentInteractions)
	return recentInteractions, nil
}

func (r *interactionRepository) Create(interaction dto.Interaction) (dto.Interaction, error) {
	log.Printf("Create started for interaction: %+v", interaction)

	if err := r.db.Create(&interaction).Error; err != nil {
		log.Printf("Error while creating interaction: %+v, err: %v", interaction, err)
		return interaction, err
	}

	log.Printf("Created interaction: %+v", interaction)
	return interaction, nil
}
