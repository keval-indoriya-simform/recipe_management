package models

import (
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name    string   `json:"name,omitempty" form:"categories" gorm:"notnull,unique"`
	Recipes []Recipe `json:"recipes,omitempty" gorm:"->;many2many:recipe_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func FindAllCategoriesName() []string {
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	var categories []string
	db.Model(Category{}).Select("name").Find(&categories)
	return categories
}
