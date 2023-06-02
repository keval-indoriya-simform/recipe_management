package services

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"strconv"
)

type RecipeService interface {
	Save(recipe *models.Recipe)
	FindAll() []models.Recipe
	FindAllWithEmail(email string) []models.Recipe
	FindByID(id string) models.Recipe
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

func (service *recipeService) FindAll() []models.Recipe {
	return models.GetAllRecipe()
}

func (service *recipeService) FindAllWithEmail(email string) []models.Recipe {
	return models.GetRecipeByEmail(email)
}

func (service *recipeService) Update(recipe *models.Recipe) {
	models.EditRecipeData(recipe)
}

func (service *recipeService) FindByID(id string) models.Recipe {
	ID, _ := strconv.Atoi(id)
	return models.GetRecipeByID(ID)
}

func (service *recipeService) Delete(id string) {
	ID, _ := strconv.Atoi(id)
	models.DeleteRecipeByID(ID)
}
