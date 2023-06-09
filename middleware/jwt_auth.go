package middelware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/keval-indoriya-simform/recipe_management/services"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		tokenstr := session.Get("token")
		if tokenstr == nil {
			context.Redirect(http.StatusTemporaryRedirect, "/login")
		} else {
			token, _ := services.ValidateToken(tokenstr.(string))
			if !token.Valid {
				context.AbortWithStatus(http.StatusForbidden)
			} else {
				claims := token.Claims.(jwt.MapClaims)
				context.Set("email", claims["email"])
			}
		}
	}
}
