package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/controllers"
	"net/http"
)

var (
	Server          = gin.Default()
	loginController = controllers.NewLoginController()
)

func init() {

	Server.Static("/css", "./templates/css")
	Server.LoadHTMLGlob("templates/*.html")

	Server.GET("/", controllers.LoginPage)
	Server.GET("/signup", controllers.SignupPage)
	Server.GET("/google/login", controllers.GoogleLogin)
	Server.GET("/home", controllers.GoogleCallback)
	Server.GET("/login", func(context *gin.Context) {
		var user controllers.User
		context.ShouldBind(&user)
		context.JSON(http.StatusOK, user)
		fmt.Println(user)
		//token := loginController.Login(context)
		//if token != "" {
		//	context.SetCookie("token", token, 3600, "/", "", false, true)
		//	context.Redirect(http.StatusFound, "view/videos")
		//} else {
		//	context.JSON(http.StatusForbidden, nil)
		//}
	})

}
