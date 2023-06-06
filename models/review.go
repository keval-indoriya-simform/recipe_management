package models

import (
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Rating   int    `json:"rating,omitempty" form:"rating-count"`
	Comment  string `json:"comment,omitempty" form:"rate"`
	EmailID  string `json:"emailID,omitempty" form:"email"`
	RecipeID int    `json:"recipeID,omitempty" form:"recipe_id"`
	Recipe   Recipe `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func InsertReviewData(review *Review) {
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Create(review)
}

func GetReviewByEmailID(email string, id int) Review {
	var review Review
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Where("email_id=? AND recipe_id=?", email, id).Find(&review)
	return review
}

func GetReviewByRecipeID(id int) []Review {
	var review []Review
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Where("recipe_id=?", id).Limit(3).Order("random()").Find(&review)
	return review
}
