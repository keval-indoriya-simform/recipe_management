package services

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type RecipeForm struct {
	ID              string   `json:"ID,omitempty" form:"id"`
	Title           string   `json:"title,omitempty" form:"title"`
	Description     string   `json:"description,omitempty" form:"description"`
	CookingTime     int32    `json:"cookingTime,omitempty" form:"cookingtime"`
	Serving         int32    `json:"serving,omitempty" form:"serving"`
	Ingredients     string   `json:"ingredients,omitempty" form:"ingredients"`
	CookingSteps    string   `json:"cookingSteps,omitempty" form:"cookingsteps"`
	Type            string   `json:"type,omitempty" form:"type"`
	Meals           []string `json:"meals,omitempty" form:"meals[]"`
	DifficultyLevel string   `json:"difficultyLevel,omitempty" form:"difficulty"`
	VideoURL        string   `json:"videoURL,omitempty" form:"url"`
	EmailID         string   `json:"emailID,omitempty" form:"email"`
	Categories      []string `json:"category,omitempty" form:"categories[]"`
}

func StructFromRecipeForm(form RecipeForm, filenames []string) models.Recipe {
	var Recipe models.Recipe
	var Categories []models.Category
	for _, val := range form.Categories {
		Categories = append(Categories, models.Category{Name: val})
	}
	id, _ := strconv.Atoi(form.ID)
	Recipe = models.Recipe{
		Model:           gorm.Model{ID: uint(id)},
		Title:           form.Title,
		Description:     form.Description,
		Ingredients:     form.Ingredients,
		CookingTime:     form.CookingTime,
		Serving:         form.Serving,
		CookingSteps:    form.CookingSteps,
		Type:            form.Type,
		Meals:           strings.Join(form.Meals, ", "),
		DifficultyLevel: form.DifficultyLevel,
		ImageURL:        strings.Join(filenames, ", "),
		VideoURL:        form.VideoURL,
		EmailID:         form.EmailID,
		Categories:      Categories,
	}
	return Recipe
}
