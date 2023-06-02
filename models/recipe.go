package models

import (
	"fmt"
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Title           string     `json:"title,omitempty" gorm:"notnull"`
	Description     string     `json:"description,omitempty" gorm:"notnull"`
	CookingTime     int32      `json:"cookingTime,omitempty" gorm:"notnull"`
	Serving         int32      `json:"serving,omitempty" gorm:"notnull"`
	Ingredients     string     `json:"ingredients,omitempty" gorm:"null"`
	CookingSteps    string     `json:"cookingSteps,omitempty" gorm:"notnull"`
	Type            string     `json:"type,omitempty" gorm:"notnull"`
	Meals           string     `json:"meals,omitempty" gorm:"notnull"`
	DifficultyLevel string     `json:"difficultyLevel,omitempty" gorm:"notnull"`
	ImageURL        string     `json:"imageURL,omitempty" gorm:"notnull"`
	VideoURL        string     `json:"videoURL,omitempty"`
	EmailID         string     `json:"emailID,omitempty" gorm:"notnull,unique"`
	Categories      []Category `json:"category,omitempty" gorm:"->;many2many:recipe_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func InsertRecipeData(recipe *Recipe) {
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Create(recipe)
	for index, _ := range recipe.Categories {
		db.Where("name=?", recipe.Categories[index].Name).Find(&recipe.Categories[index])
		fmt.Println(recipe.Categories[index].ID)
		if recipe.Categories[index].ID == 0 {
			db.Create(&recipe.Categories[index])
		}
		db.Table("recipe_categories").Create(map[string]interface{}{
			"category_id": recipe.Categories[index].ID,
			"recipe_id":   recipe.ID,
		})
	}
}

func GetAllRecipe() []Recipe {
	var recipes []Recipe
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Preload("Categories").Find(&recipes)
	return recipes
}

func GetRecipeByEmail(email string) []Recipe {
	var recipes []Recipe
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Preload("Categories").Where("email_id=?", email).Find(&recipes)
	return recipes
}

func GetRecipeByID(id int) Recipe {
	var recipe Recipe
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Preload("Categories").Where("id=?", id).Find(&recipe)
	return recipe
}

func DeleteRecipeByID(id int) Recipe {
	var recipe Recipe
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Where("id=?", id).Delete(&recipe)
	return recipe
}

func EditRecipeData(recipe *Recipe) {
	db := initializers.GetConnection()
	defer initializers.CloseConnection(db)
	db.Model(Recipe{}).Where("id=?", recipe.ID).Updates(recipe)
	db.Table("recipe_categories").Where("recipe_id=?", recipe.ID).Delete("recipe_id, category_id")
	for index, _ := range recipe.Categories {
		db.Where("name=?", recipe.Categories[index].Name).Find(&recipe.Categories[index])
		if recipe.Categories[index].ID == 0 {
			db.Create(&recipe.Categories[index])
		}
		db.Table("recipe_categories").Create(map[string]interface{}{
			"category_id": recipe.Categories[index].ID,
			"recipe_id":   recipe.ID,
		})
	}
}
