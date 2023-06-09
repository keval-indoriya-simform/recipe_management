package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Rating   int    `form:"rating-count" json:"rating,omitempty"`
	Comment  string `form:"rate" json:"comment,omitempty"`
	EmailID  string `form:"email" json:"emailID,omitempty"`
	RecipeID int    `form:"recipe_id" json:"recipeID,omitempty"`
	Recipe   Recipe `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func InsertReviewData(review *Review) {
	DB.Create(review)
}

func GetReviewByEmailID(email string, id int) map[string]interface{} {
	var review map[string]interface{}
	DB.Model(Review{}).Select("rating, comment").Where("email_id=? AND recipe_id=?", email, id).Find(&review)
	return review
}

func GetReviewByRecipeID(email string, id int) []map[string]interface{} {
	var review []map[string]interface{}
	DB.Model(Review{}).Where("recipe_id=? and not email_id=?", id, email).Limit(3).Order("random()").Find(&review)
	return review
}
