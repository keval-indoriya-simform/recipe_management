package main

import (
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"github.com/keval-indoriya-simform/recipe_management/models"
	"github.com/keval-indoriya-simform/recipe_management/router"
)

func main() {
	router.Server.Run("localhost:8080")
	defer initializers.CloseConnection(models.DB)
}
