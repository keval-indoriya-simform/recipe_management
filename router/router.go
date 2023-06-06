package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/controllers"
	middleware "github.com/keval-indoriya-simform/recipe_management/middleware"
)

var (
	Server          = gin.Default()
	loginController = controllers.NewLoginController()
)

func init() {
	store := cookie.NewStore([]byte("secret"))
	Server.Static("/Upload", "./upload")

	Server.Use(sessions.Sessions("mysession", store))
	Server.Static("/css", "./templates/css")
	Server.LoadHTMLGlob("templates/*.html")
	Server.GET("/login", controllers.LoginPage)
	Server.GET("/logout", controllers.Logout)
	Server.GET("/google/login", controllers.GoogleLogin)
	Server.GET("/google/callback", controllers.GoogleCallback)
	Server.GET("/microsoft/login", controllers.MicrosoftLogin)
	Server.GET("/microsoft/callback", controllers.MicrosoftCallback)
	Server.GET("/home", middleware.AuthorizeJWT(), controllers.HomePage)
	Server.GET("/addrecipe", middleware.AuthorizeJWT(), controllers.AddRecipePage)
	Server.GET("/editrecipe/:id", middleware.AuthorizeJWT(), controllers.EditRecipePage)
	Server.POST("/findrecipe/:id", middleware.AuthorizeJWT(), controllers.FindRecipeByID)
	Server.GET("/fullrecipe/:id", middleware.AuthorizeJWT(), controllers.FullRecipePage)
	Server.GET("/myrecipe", middleware.AuthorizeJWT(), controllers.MyRecipePage)
	Server.POST("/search", middleware.AuthorizeJWT(), controllers.SearchApi)

	apiGroup := Server.Group("/api", middleware.AuthorizeJWT())
	recipeGroup := apiGroup.Group("recipe")
	recipeGroup.POST("Add_recipe", controllers.AddRecipeApi)
	recipeGroup.POST("Edit_recipe", controllers.EditRecipeApi)
	recipeGroup.GET("Delete_recipe/:id", controllers.DeleteRecipeApi)

	reviewGroup := apiGroup.Group("review")
	reviewGroup.POST("Add_review", controllers.AddReviewApi)
	reviewGroup.POST("Get_review/:id", controllers.GetReviewApi)
	reviewGroup.POST("GetAll_review/:id", controllers.GetAllReviewsByRecipeIDApi)
}
