package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func HomePage(context *gin.Context) {
	context.HTML(http.StatusOK, "home.html", gin.H{
		"recipes": models.GetAllRecipe(),
	})
}

func AddRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	context.HTML(http.StatusOK, "add_recipe.html", gin.H{
		"email": email,
	})
}

func MyRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	context.HTML(http.StatusOK, "my_recipe.html", gin.H{
		"recipes": models.GetRecipeByEmail(email.(string)),
	})
}

func AddRecipeApi(context *gin.Context) {
	var recipe services.RecipeForm
	err := context.Bind(&recipe)
	form, multipartError := context.MultipartForm()
	if multipartError != nil {
		context.String(http.StatusBadRequest, "get form err: %s", err.Error())
		log.Fatal(err)
	}
	files := form.File["files"]
	dst := "./upload/" + filepath.Base(files[0].Filename)
	if err := context.SaveUploadedFile(files[0], dst); err != nil {
		context.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}
	Recipes := services.StructFromRecipeForm(recipe, filepath.Base(files[0].Filename))
	models.InsertRecipeData(&Recipes)
	context.Redirect(http.StatusFound, "/home")
}

func EditRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	id := context.Param("id")
	context.HTML(http.StatusOK, "edit_recipe.html", gin.H{
		"ID":    id,
		"email": email,
	})
}

func EditRecipeApi(context *gin.Context) {
	var recipe services.RecipeForm
	var filename string
	err := context.Bind(&recipe)
	initializers.ErrorCheck(err)
	form, multipartError := context.MultipartForm()
	if multipartError != nil {
		context.String(http.StatusBadRequest, "get form err: %s", multipartError.Error())
	}
	files, getFileErr := form.File["files"]
	if getFileErr == true {
		dst := "./upload/" + filepath.Base(files[0].Filename)
		if err := context.SaveUploadedFile(files[0], dst); err != nil {
			context.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		filename = filepath.Base(files[0].Filename)
	}
	Recipes := services.StructFromRecipeForm(recipe, filename)
	models.EditRecipeData(&Recipes)
	context.Redirect(http.StatusFound, "/fullrecipe/"+recipe.ID)
}

func DeleteRecipeApi(context *gin.Context) {
	id := context.Param("id")
	ID, _ := strconv.Atoi(id)
	models.DeleteRecipeByID(ID)
	context.Redirect(http.StatusFound, "/myrecipe")
}

func FindRecipeByID(context *gin.Context) {
	id := context.Param("id")
	ID, _ := strconv.Atoi(id)
	context.JSON(http.StatusOK, models.GetRecipeByID(ID))
}

func FullRecipePage(context *gin.Context) {
	id := context.Param("id")
	ID, _ := strconv.Atoi(id)
	email, _ := context.Get("email")
	context.HTML(http.StatusOK, "full_recipe_page.html", gin.H{
		"email":  email,
		"recipe": models.GetRecipeByID(ID),
	})
}

func SearchApi(context *gin.Context) {
	var searchStruct models.Search
	if err := context.Bind(&searchStruct); err != nil {
		log.Println(err)
	}
	context.HTML(http.StatusOK, "search.html", gin.H{
		"recipes":      models.SearchRecipe(searchStruct),
		"searchstruct": searchStruct,
	})
}
