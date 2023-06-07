package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	googleConf       = &services.ConfigureOAuth2{}
	microsoftConf    = &services.ConfigureOAuth2{}
	loginControllers = NewLoginController()
	recipeService    = services.NewRecipeService()
	recipeController = NewRecipeController(recipeService)
	reviewService    = services.NewReviewService()
	reviewController = NewReviewController(reviewService)
)

func init() {
	envError := godotenv.Load(".env")
	if envError != nil {
		log.Fatal("Error loading .env file", envError)
	}
	googleConf.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleConf.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleConf.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	googleConf.Scopes = strings.Split(os.Getenv("GOOGLE_SCOPES"), ",")
	googleConf.Endpoint = google.Endpoint
	googleConf.State = "randomstate"
	googleConf.GetInfoURL = os.Getenv("GOOGLE_GET_DETAIL_URL") + "?access_token="
	googleConf.RequestMethod = http.MethodGet
	googleConf.Body = http.NoBody
	microsoftConf.ClientID = os.Getenv("MICROSOFT_CLIENT_ID")
	microsoftConf.ClientSecret = os.Getenv("MICROSOFT_CLIENT_SECRET")
	microsoftConf.RedirectURL = os.Getenv("MICROSOFT_REDIRECT_URL")
	microsoftConf.Scopes = strings.Split(os.Getenv("MICROSOFT_SCOPES"), ",")
	microsoftConf.Endpoint = microsoft.AzureADEndpoint(os.Getenv("MICROSOFT_TENANT_ID"))
	microsoftConf.State = "randomstate"
	microsoftConf.GetInfoURL = os.Getenv("MICROSOFT_GET_DETAIL_URL")
	microsoftConf.RequestMethod = http.MethodGet
	microsoftConf.Body = http.NoBody
}

func LoginPage(Context *gin.Context) {
	Context.HTML(http.StatusOK, "login.html", nil)
}

func SignupPage(Context *gin.Context) {
	Context.HTML(http.StatusOK, "signup.html", nil)
}

func GoogleLogin(context *gin.Context) {
	url := services.Login(googleConf)
	context.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(context *gin.Context) {
	userInfo, getInfoError := services.Callback(
		context.Request,
		googleConf,
		services.QueryString,
	)
	if getInfoError != nil {
		log.Fatalln(getInfoError)
	}
	var TokenRes map[string]interface{}
	json.Unmarshal(userInfo, &TokenRes)
	user := models.Login{
		Name:  TokenRes["name"].(string),
		Email: TokenRes["email"].(string),
	}
	token := loginControllers.Login(user)
	if token != "" {
		session := sessions.Default(context)
		session.Set("token", token)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		context.Redirect(http.StatusFound, "/home")
	} else {
		context.JSON(http.StatusForbidden, nil)
	}
}

func HomePage(context *gin.Context) {
	//recipe := recipeController.FindAll()
	context.HTML(http.StatusOK, "home.html", gin.H{
		"categories": models.FindAllCategoriesName(),
		"recipes":    models.GetAllRecipe(),
	})
}

func AddRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	context.HTML(http.StatusOK, "add_recipe.html", gin.H{
		"email":      email,
		"categories": models.FindAllCategoriesName(),
	})
}

func MyRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	context.HTML(http.StatusOK, "my_recipe.html", gin.H{
		"categories": models.FindAllCategoriesName(),
		"recipes":    recipeController.FindAllWithEmail(email.(string)),
	})
}

func Logout(context *gin.Context) {
	session := sessions.Default(context)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Fatal(err)
	}
	context.Redirect(http.StatusFound, "/login")
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
	recipeController.Save(&Recipes)
	context.Redirect(http.StatusFound, "/home")
}

func EditRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	id := context.Param("id")
	categories := models.FindAllCategoriesName()
	context.HTML(http.StatusOK, "edit_recipe.html", gin.H{
		"ID":         id,
		"email":      email,
		"categories": categories,
	})
}

func EditRecipeApi(context *gin.Context) {
	var recipe services.RecipeForm
	err := context.Bind(&recipe)
	initializers.ErrorCheck(err)
	initializers.ErrorCheck(err)
	form, multipartError := context.MultipartForm()
	if multipartError != nil {
		context.String(http.StatusBadRequest, "get form err: %s", multipartError.Error())
		log.Fatal(multipartError)
	}
	files, getFileErr := form.File["files"]
	log.Println(getFileErr)
	if getFileErr == true {
		dst := "./upload/" + filepath.Base(files[0].Filename)
		if err := context.SaveUploadedFile(files[0], dst); err != nil {
			context.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}

	Recipes := services.StructFromRecipeForm(recipe, filepath.Base(files[0].Filename))
	recipeController.Update(&Recipes)
	context.Redirect(http.StatusFound, "/fullrecipe/"+recipe.ID)
}

func DeleteRecipeApi(context *gin.Context) {
	id := context.Param("id")
	recipeController.Delete(id)
	context.Redirect(http.StatusFound, "/myrecipe")
}

func FindRecipeByID(context *gin.Context) {
	id := context.Param("id")
	recipe := recipeController.FindByID(id)
	context.JSON(http.StatusOK, recipe)
}

func FullRecipePage(context *gin.Context) {
	id := context.Param("id")
	email, _ := context.Get("email")
	recipe := recipeController.FindByID(id)
	categories := models.FindAllCategoriesName()
	context.HTML(http.StatusOK, "full_recipe_page.html", gin.H{
		"email":      email,
		"recipe":     recipe,
		"categories": categories,
	})
}

func AddReviewApi(context *gin.Context) {
	var rating models.Review
	context.Bind(&rating)
	reviewController.Save(&rating)
	url := "/fullrecipe/" + strconv.Itoa(rating.RecipeID)
	context.Redirect(http.StatusFound, url)
}

func GetReviewApi(context *gin.Context) {
	id := context.Param("id")
	email, _ := context.Get("email")
	rating := reviewController.GetReviewByEmailID(email.(string), id)
	context.JSON(http.StatusOK, rating)
}

func SearchApi(context *gin.Context) {
	var searchStruct models.Search
	err := context.Bind(&searchStruct)
	if err != nil {
		log.Println(err)
	}
	recipe := models.SearchRecipe(searchStruct)
	categories := models.FindAllCategoriesName()
	context.HTML(http.StatusOK, "search.html", gin.H{
		"recipes":    recipe,
		"categories": categories,
	})
}

func GetAllReviewsByRecipeIDApi(context *gin.Context) {
	id := context.Param("id")
	reviews := reviewController.GetReviewByRecipeID(id)
	context.JSON(http.StatusOK, reviews)
}

func MicrosoftLogin(context *gin.Context) {
	url := services.Login(microsoftConf)
	context.Redirect(http.StatusFound, url)
}
func MicrosoftCallback(context *gin.Context) {
	userInfo, getInfoError := services.Callback(
		context.Request,
		microsoftConf,
		services.AuthorizationBearer)
	if getInfoError != nil {
		log.Fatalln(getInfoError)
	}
	var TokenRes map[string]interface{}
	json.Unmarshal(userInfo, &TokenRes)
	user := models.Login{
		Name:  TokenRes["displayName"].(string),
		Email: TokenRes["mail"].(string),
	}
	token := loginControllers.Login(user)
	if token != "" {
		session := sessions.Default(context)
		session.Set("token", token)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		context.Redirect(http.StatusFound, "/home")
	} else {
		context.JSON(http.StatusForbidden, nil)
	}
}

func GetAllCategories(context *gin.Context) {
	context.JSON(http.StatusOK, models.FindAllCategoriesName())
}
