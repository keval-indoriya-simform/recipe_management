package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(user models.Login) string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: GetSecretKey(),
	}
}

func GetSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "keval"
	}
	return []byte(secret)
}

func (jwtSrv jwtService) GenerateToken(user models.Login) string {
	sleepTime := 72
	claims := &JwtCustomClaims{
		user.Name,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * time.Duration(sleepTime))},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSrv.secretKey)
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtSrv jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSrv.secretKey, nil
	})
}
