package controllers

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
)

type RecipeController interface {
	FindAllWithEmail(email string) []map[string]interface{}
	FindAll() []map[string]interface{}
	FindByID(id string) map[string]interface{}
	Save(recipe *models.Recipe)
	Delete(id string)
	Update(recipe *models.Recipe)
}

type controller struct {
	service services.RecipeService
}

func NewRecipeController(serv services.RecipeService) RecipeController {
	return &controller{
		service: serv,
	}
}

func (c *controller) FindAll() []map[string]interface{} {
	return c.service.FindAll()
}

func (c *controller) Save(recipe *models.Recipe) {
	c.service.Save(recipe)
}

func (c *controller) FindAllWithEmail(email string) []map[string]interface{} {
	return c.service.FindAllWithEmail(email)
}

func (c *controller) Update(recipe *models.Recipe) {
	c.service.Update(recipe)
}

func (c *controller) FindByID(id string) map[string]interface{} {
	return c.service.FindByID(id)
}

func (c *controller) Delete(id string) {
	c.service.Delete(id)
}
