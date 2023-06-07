package services

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"strconv"
)

type RecipeService interface {
	Save(recipe *models.Recipe)
	FindAll() []map[string]interface{}
	FindAllWithEmail(email string) []map[string]interface{}
	FindByID(id string) map[string]interface{}
	Update(recipe *models.Recipe)
	Delete(id string)
}

type recipeService struct {
	Recipes []models.Recipe
}

func NewRecipeService() *recipeService {
	return &recipeService{}
}

func (service *recipeService) Save(recipe *models.Recipe) {
	models.InsertRecipeData(recipe)
}

func (service *recipeService) FindAll() []map[string]interface{} {
	return models.GetAllRecipe()
}

func (service *recipeService) FindAllWithEmail(email string) []map[string]interface{} {
	return models.GetRecipeByEmail(email)
}

func (service *recipeService) Update(recipe *models.Recipe) {
	models.EditRecipeData(recipe)
}

func (service *recipeService) FindByID(id string) map[string]interface{} {
	ID, _ := strconv.Atoi(id)
	return models.GetRecipeByID(ID)
}

func (service *recipeService) Delete(id string) {
	ID, _ := strconv.Atoi(id)
	models.DeleteRecipeByID(ID)
}
