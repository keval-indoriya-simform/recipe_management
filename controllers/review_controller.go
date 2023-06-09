package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"net/http"
	"strconv"
)

func AddReviewApi(context *gin.Context) {
	var rating models.Review
	context.Bind(&rating)
	models.InsertReviewData(&rating)
	url := "/fullrecipe/" + strconv.Itoa(rating.RecipeID)
	context.Redirect(http.StatusFound, url)
}

func GetReviewApi(context *gin.Context) {
	id := context.Param("id")
	ID, _ := strconv.Atoi(id)
	email, _ := context.Get("email")
	context.JSON(http.StatusOK, models.GetReviewByEmailID(email.(string), ID))
}

func GetAllReviewsByRecipeIDApi(context *gin.Context) {
	id := context.Param("id")
	ID, _ := strconv.Atoi(id)
	email, _ := context.Get("email")
	context.JSON(http.StatusOK, models.GetReviewByRecipeID(email.(string), ID))
}
