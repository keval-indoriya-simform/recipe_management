package controllers

import (
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
)

var (
	JwtService = services.NewJWTService()
)

type LoginController interface {
	Login(user models.Login) string
}

type loginController struct {
	JwtService services.JWTService
}

func NewLoginController() LoginController {
	return &loginController{
		JwtService: JwtService,
	}
}

func (controller *loginController) Login(user models.Login) string {
	return controller.JwtService.GenerateToken(user)
}
