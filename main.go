package main

import (
	"github.com/keval-indoriya-simform/recipe_management/router"
)

func main() {
	//categories := []models.Category{
	//	{Name: "Low Cal"},
	//	{Name: "Diabetic"},
	//	{Name: "Low Carb"},
	//	{Name: "Low Cholesterol"},
	//	{Name: "Balanced Diet"},
	//	{Name: "Low Sodium"},
	//	{Name: "High Fiber"},
	//	{Name: "High Protein"},
	//	{Name: "High Carb"},
	//	{Name: "North Indian"},
	//	{Name: "South Indian"},
	//	{Name: "Chinese"},
	//	{Name: "Asian"},
	//	{Name: "Italian"},
	//	{Name: "Mughlai"},
	//	{Name: "American"},
	//	{Name: "Punjabi"},
	//	{Name: "Mexican"},
	//	{Name: "Gujarati"},
	//	{Name: "Fusion"},
	//	{Name: "Thai"},
	//	{Name: "Sattvik"},
	//	{Name: "Maharashtrian"},
	//	{Name: "Rajasthani"},
	//	{Name: "Malabar Cuisine"},
	//	{Name: "Bengali"},
	//	{Name: "Bihari"},
	//	{Name: "Konkani"},
	//	{Name: "Goan"},
	//	{Name: "Kashmiri"},
	//	{Name: "Vietnamese"},
	//	{Name: "Middle Eastern"},
	//	{Name: "Parsi"},
	//	{Name: "French"},
	//	{Name: "Australian"},
	//}
	//db := initializers.GetConnection()
	//db.Migrator().DropTable(&models.Category{})
	//defer initializers.CloseConnection(db)
	//db.AutoMigrate(&models.Recipe{}, &models.Category{}, &models.Review{})
	//db.Create(&categories)
	router.Server.Run("localhost:8080")
}
