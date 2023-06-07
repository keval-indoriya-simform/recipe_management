package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name    string   `gorm:"notnull,unique" json:"name,omitempty" form:"categories"`
	Recipes []Recipe `gorm:"->;many2many:recipe_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"recipes,omitempty"`
}

func FindAllCategoriesName() []string {
	var categories []string
	DB.Model(Category{}).Select("name").Find(&categories)
	return categories
}

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
