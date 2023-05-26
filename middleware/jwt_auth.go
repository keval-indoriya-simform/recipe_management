package middelware

import (
	"github.com/keval-indoriya-simform/recipe_management/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenstr, err := context.Cookie("token")
		if err != nil {
			log.Fatal(err)
		}
		if tokenstr == "" {
			context.Redirect(http.StatusFound, "/login")
		}

		token, err := services.NewJWTService().ValidateToken(tokenstr)

		if !token.Valid {
			log.Println(err)
			context.AbortWithStatus(http.StatusForbidden)
		}
	}
}
