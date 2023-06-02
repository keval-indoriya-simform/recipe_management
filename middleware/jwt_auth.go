package middelware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v5"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		tokenstr := session.Get("token")
		fmt.Println(tokenstr)
		if tokenstr == nil {
			context.Redirect(http.StatusTemporaryRedirect, "/")
		} else {
			token, err := services.NewJWTService().ValidateToken(fmt.Sprintf("%v", tokenstr))

			if !token.Valid {
				log.Println(err)
				context.AbortWithStatus(http.StatusForbidden)
			} else {
				claims := token.Claims.(jwt.MapClaims)
				context.Set("email", claims["email"])
			}
		}
	}
}
