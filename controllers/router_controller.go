package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	googleConf = &services.ConfigureOAuth2{}
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
	username := fmt.Sprintf("%v", TokenRes["name"])
	useremail := fmt.Sprintf("%v", TokenRes["email"])
	context.Set("jsonstr", gin.H{
		"name":  username,
		"email": useremail,
	})
	context.Redirect(http.StatusTemporaryRedirect, "/login")
}
