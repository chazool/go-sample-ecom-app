package services

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/app/repository"
)

var (
	ErrInteractionsNotFound = errors.New("interactions not found")
)

type InteractionService interface {
	Create(interaction dto.Interaction) (dto.Interaction, error)
	GetRecentInteractions(user, limit uint) ([]dto.Interaction, error)
}

type interactionService struct {
	interactionRepo repository.InteractionRepository
}

func NewInteractionService() InteractionService {
	return &interactionService{
		interactionRepo: repository.NewInteractionRepository(),
	}
}

func (service *interactionService) GetRecentInteractions(user, limit uint) ([]dto.Interaction, error) {

	log.Printf("Retrieving interaction with user %d\n", user)

	// Retrieve interaction from the database with given user id
	interactions, err := service.interactionRepo.RecentInteractions(user, limit)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return interactions, ErrInteractionsNotFound
		}
		return interactions, err
	}

	log.Printf("Retrieved interaction with user %d: %+v\n", user, interactions)

	return interactions, nil
}

func (service *interactionService) Create(interaction dto.Interaction) (dto.Interaction, error) {

	log.Printf("Creating interaction: %+v\n", interaction)

	product, err := service.interactionRepo.Create(interaction)
	if err != nil {
		return product, ErrFailToCreateProduct
	}

	log.Printf("Created interaction with ID %d\n", product.ID)

	return product, nil
}
