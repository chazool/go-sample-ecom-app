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
	log.Printf("Starting GetRecentInteractions function for user %d with limit %d", user, limit)
	defer log.Println("Ending GetRecentInteractions function")

	log.Printf("Retrieving interactions for user %d\n", user)
	// Retrieve interactions from the database with given user id and limit
	interactions, err := service.interactionRepo.RecentInteractions(user, limit)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("No interactions found for user %d: %v", user, err)
			return interactions, ErrInteractionsNotFound
		}
		log.Printf("Failed to retrieve interactions for user %d: %v", user, err)
		return interactions, err
	}

	log.Printf("Retrieved %d interactions for user %d: %+v\n", len(interactions), user, interactions)
	return interactions, nil
}

func (service *interactionService) Create(interaction dto.Interaction) (dto.Interaction, error) {
	// Log function start
	log.Printf("Starting Create function for interaction: %+v", interaction)
	defer log.Printf("Ending Create function for interaction: %+v", interaction)

	// Create interaction in the database
	createdInteraction, err := service.interactionRepo.Create(interaction)
	if err != nil {
		log.Printf("Failed to create interaction: %+v, err: %v", interaction, err)
		return createdInteraction, ErrFailToCreateProduct
	}

	// Log successful creation of interaction
	log.Printf("Created interaction with ID %d", createdInteraction.ID)

	return createdInteraction, nil
}
