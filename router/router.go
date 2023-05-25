package router

import (
	"github.com/keval-indoriya-simform/recipe_management/controllers"
	"github.com/keval-indoriya-simform/recipe_management/initializers"
)

func init() {
	initializers.Server.Static("/css", "./templates/css")
	initializers.Server.LoadHTMLGlob("/home/keval/Desktop/practice/recipe management/templates/*.html")

	initializers.Server.GET("/login", controllers.LoginPage)
}
