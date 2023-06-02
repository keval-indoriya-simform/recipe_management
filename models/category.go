package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name    string   `json:"name,omitempty" form:"categories" gorm:"notnull,unique"`
	Recipes []Recipe `json:"recipes,omitempty" gorm:"->;many2many:recipe_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
