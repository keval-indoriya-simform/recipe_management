package controllers

import (
	"database/sql"
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

	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	//recipe := recipeController.FindAll()
	recipe := models.GetAllRecipe()
	categories := models.FindAllCategoriesName()
	context.HTML(http.StatusOK, "home.html", gin.H{
		"categories": categories,
		"recipes":    recipe,
	})
}

func AddRecipePage(context *gin.Context) {
	//key, _ := services.NewJWTService().ValidateToken(token.(string))
	//claims := key.Claims.(jwt.MapClaims)
	email, _ := context.Get("email")
	categories := models.FindAllCategoriesName()
	context.HTML(http.StatusOK, "add_recipe.html", gin.H{
		"email":      email,
		"categories": categories,
	})
}

func MyRecipePage(context *gin.Context) {
	email, _ := context.Get("email")
	categories := models.FindAllCategoriesName()
	recipe := recipeController.FindAllWithEmail(email.(string))
	context.HTML(http.StatusOK, "my_recipe.html", gin.H{
		"categories": categories,
		"recipes":    recipe,
	})
}

func Logout(context *gin.Context) {
	session := sessions.Default(context)
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Fatal(err)
	} else {
		//context.String(http.StatusOK, "logout")
		context.Redirect(http.StatusFound, "/login")
	}
}

func AddRecipeApi(context *gin.Context) {
	var recipe services.RecipeForm
	err := context.Bind(&recipe)
	form, err := context.MultipartForm()
	if err != nil {
		context.String(http.StatusBadRequest, "get form err: %s", err.Error())
		log.Fatal(err)
	}
	files := form.File["files"]
	var filenames []string
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		filenames = append(filenames, filename)
		dst := "./upload/" + filename
		if err := context.SaveUploadedFile(file, dst); err != nil {
			context.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	Recipes := services.StructFromRecipeForm(recipe, filenames)
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
	form, err := context.MultipartForm()
	if err != nil {
		context.String(http.StatusBadRequest, "get form err: %s", err.Error())
		log.Fatal(err)
	}
	files := form.File["files"]
	var filenames []string
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		filenames = append(filenames, filename)
		dst := "./upload/" + filename
		if err := context.SaveUploadedFile(file, dst); err != nil {
			context.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	Recipes := services.StructFromRecipeForm(recipe, filenames)
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
	type Search struct {
		SearchKeyword string   `json:"search_keyword,omitempty" form:"search"`
		Categories    []string `json:"categories,omitempty" form:"search-categories[]"`
		MealTypes     []string `json:"meal_types,omitempty" form:"search-type[]"`
		Courses       []string `json:"courses,omitempty" form:"search-meals[]"`
		Sort          string   `json:"sort,omitempty" form:"search-sort"`
	}
	var searchStruct Search
	err := context.Bind(&searchStruct)
	if err != nil {
		log.Println(err)
	}
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	var recipe []models.Recipe
	var courseQuery, categoriesQuery, typeQuery = "", "", ""
	if len(searchStruct.Categories) != 0 {
		categoriesQuery = " AND c.name IN ('" + strings.Join(searchStruct.Categories, "', '") + "')"
	}
	if len(searchStruct.MealTypes) != 0 {
		typeQuery = " AND type IN ('" + strings.Join(searchStruct.MealTypes, "', '") + "')"
	}
	if len(searchStruct.Courses) != 0 {
		courseQuery = " AND ("
		for i, _ := range searchStruct.Courses {
			courseQuery += "meals ILike '%" + searchStruct.Courses[i] + "%'"
			if i != (len(searchStruct.Courses) - 1) {
				courseQuery += " OR "
			}
		}
		courseQuery += ")"
	}
	query := db.Debug().Model(models.Recipe{}).Distinct("recipes.*, avg(rating) as rating").
		Joins("left join recipe_categories as rc on recipes.id = rc.recipe_id"+
			" left join categories as c on rc.category_id = c.id"+
			" left join reviews as r on r.recipe_id = recipes.id").Group("r.recipe_id, recipes.id").Where("(title ILike @keyword OR ingredients ILike @keyword OR type ILike @keyword OR meals ILike @keyword OR c.name ILike @keyword)"+categoriesQuery+typeQuery+courseQuery,
		sql.Named("keyword", "%"+searchStruct.SearchKeyword+"%"))
	query.Order(searchStruct.Sort).Find(&recipe)
	for i, _ := range recipe {
		var avgRating int = 0
		db.Model(models.Review{}).Select("floor(avg(rating))").Where("recipe_id=?", recipe[i].ID).Group("recipe_id").Find(&avgRating)
		recipe[i].AvgRating = avgRating
	}
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
