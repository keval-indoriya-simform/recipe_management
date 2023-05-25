package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginPage(Context *gin.Context) {
	//Context.HTML(http.StatusOK, "login.html", nil)
	Context.JSON(http.StatusOK, gin.H{"keval": "indoriya"})
}
