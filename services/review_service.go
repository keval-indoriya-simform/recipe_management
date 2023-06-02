package services

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"strconv"
)

type ReviewService interface {
	Save(review *models.Review)
	GetReviewByEmailID(email string, id string) models.Review
}

type reviewService struct {
	Reviews []models.Review
}

func NewReviewService() *reviewService {
	return &reviewService{}
}

func (service *reviewService) Save(review *models.Review) {
	models.InsertReviewData(review)
}

func (service *reviewService) GetReviewByEmailID(email string, id string) models.Review {
	ID, _ := strconv.Atoi(id)
	return models.GetReviewByEmailID(email, ID)
}
