package main

import (
	"github.com/keval-indoriya-simform/recipe_management/router"
)

func main() {
	//categories := []models.Category{
	//	{Name: "Low Cal"},
	//	{Name: "Gluten Free"},
	//	{Name: "Diabetic"},
	//	{Name: "Low Carb"},
	//	{Name: "Low Cholesterol"},
	//	{Name: "Weight Loss"},
	//	{Name: "Balanced Diet"},
	//	{Name: "Zero Oil"},
	//	{Name: "Low Sodium"},
	//	{Name: "High Fiber"},
	//	{Name: "Pregnancy"},
	//	{Name: "High Protein"},
	//	{Name: "Dairy Free"},
	//	{Name: "Tree Nut Free"},
	//}
	//db := initializers.GetConnection()
	//db.Migrator().DropTable("recipe_categories", &models.Category{}, &models.Recipe{})
	//db.AutoMigrate(&models.Recipe{}, &models.Category{})
	//db.AutoMigrate(&models.Review{})
	//defer initializers.CloseConnection(db)
	//db.Create(&categories)
	router.Server.Run("localhost:8080")
}
