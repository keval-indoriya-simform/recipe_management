package controllers

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
)

type ReviewController interface {
	Save(review *models.Review)
	GetReviewByEmailID(email string, id string) models.Review
}

type reviewcontroller struct {
	service services.ReviewService
}

func NewReviewController(serv services.ReviewService) ReviewController {
	return &reviewcontroller{
		service: serv,
	}
}

func (c *reviewcontroller) Save(review *models.Review) {
	c.service.Save(review)
}

func (c *reviewcontroller) GetReviewByEmailID(email string, id string) models.Review {
	return c.service.GetReviewByEmailID(email, id)
}
