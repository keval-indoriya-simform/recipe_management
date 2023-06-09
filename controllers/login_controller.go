package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	googleConf    = &services.ConfigureOAuth2{}
	microsoftConf = &services.ConfigureOAuth2{}
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

func GoogleLogin(context *gin.Context) {
	url := services.Login(googleConf)
	context.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(context *gin.Context) {
	userInfo, getInfoError := services.Callback(context.Request, googleConf, services.QueryString)
	if getInfoError != nil {
		log.Fatalln(getInfoError)
	}
	var TokenRes map[string]interface{}
	json.Unmarshal(userInfo, &TokenRes)
	user := models.Login{
		Name:  TokenRes["name"].(string),
		Email: TokenRes["email"].(string),
	}
	if token := services.GenerateToken(user); token != "" {
		session := sessions.Default(context)
		session.Set("token", token)
		if err := session.Save(); err != nil {
			log.Fatal(err)
		}
		context.Redirect(http.StatusFound, "/home")
	} else {
		context.JSON(http.StatusForbidden, nil)
	}
}

func MicrosoftLogin(context *gin.Context) {
	url := services.Login(microsoftConf)
	context.Redirect(http.StatusFound, url)
}

func MicrosoftCallback(context *gin.Context) {
	var TokenRes map[string]interface{}
	userInfo, getInfoError := services.Callback(context.Request, microsoftConf, services.AuthorizationBearer)
	if getInfoError != nil {
		log.Fatalln(getInfoError)
	}
	json.Unmarshal(userInfo, &TokenRes)
	user := models.Login{
		Name:  TokenRes["displayName"].(string),
		Email: TokenRes["mail"].(string),
	}
	if token := services.GenerateToken(user); token != "" {
		session := sessions.Default(context)
		session.Set("token", token)
		if err := session.Save(); err != nil {
			log.Fatal(err)
		}
		context.Redirect(http.StatusFound, "/home")
	} else {
		context.JSON(http.StatusForbidden, nil)
	}
}

func Logout(context *gin.Context) {
	session := sessions.Default(context)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Fatal(err)
	}
	context.Redirect(http.StatusFound, "/login")
}
