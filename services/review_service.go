package services

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"strconv"
)

type ReviewService interface {
	Save(review *models.Review)
	GetReviewByEmailID(email string, id string) map[string]interface{}
	GetReviewByRecipeID(id string) []map[string]interface{}
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

func (service *reviewService) GetReviewByEmailID(email string, id string) map[string]interface{} {
	ID, _ := strconv.Atoi(id)
	return models.GetReviewByEmailID(email, ID)
}

func (service *reviewService) GetReviewByRecipeID(id string) []map[string]interface{} {
	ID, _ := strconv.Atoi(id)
	return models.GetReviewByRecipeID(ID)
}
