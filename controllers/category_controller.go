package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"net/http"
)

func GetAllCategories(context *gin.Context) {
	context.JSON(http.StatusOK, models.FindAllCategoriesName())
}
