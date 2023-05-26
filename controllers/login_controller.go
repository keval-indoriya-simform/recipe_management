package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keval-indoriya-simform/recipe_management/services"
)

var (
	jwtService = services.NewJWTService()
	userInfo   User
)

type LoginController interface {
	Login(context *gin.Context) string
}

type loginController struct {
	jwtService services.JWTService
}

type User struct {
	name  string `json:"name,omitempty"`
	email string `json:"email,omitempty"`
}

func NewLoginController() LoginController {
	return &loginController{
		jwtService: jwtService,
	}
}

func (controller *loginController) Login(context *gin.Context) string {
	err := context.ShouldBind(&userInfo)
	if err != nil {
		return ""
	}
	return controller.jwtService.GenerateToken(userInfo.name, userInfo.email)
}
