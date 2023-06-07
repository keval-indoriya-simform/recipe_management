package models

import (
	"database/sql"
	"fmt"
	"github.com/keval-indoriya-simform/recipe_management/initializers"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Recipe struct {
	gorm.Model
	Title           string `gorm:"notnull" json:"title,omitempty"`
	Description     string `gorm:"notnull" json:"description,omitempty"`
	CookingTime     int32  `gorm:"notnull" json:"cookingTime,omitempty"`
	Serving         int32  `gorm:"notnull" json:"serving,omitempty"`
	Ingredients     string `gorm:"null" json:"ingredients,omitempty"`
	CookingSteps    string `gorm:"notnull" json:"cookingSteps,omitempty"`
	Type            string `gorm:"notnull" json:"type,omitempty"`
	Meals           string `gorm:"notnull" json:"meals,omitempty"`
	DifficultyLevel string `gorm:"notnull" json:"difficultyLevel,omitempty"`
	ImageURL        string `gorm:"notnull" json:"imageURL,omitempty"`
	VideoURL        string `json:"videoURL,omitempty"`
	EmailID         string `gorm:"notnull,unique" json:"emailID,omitempty"`
	//AvgRating       int        `gorm:"-" json:"avg_rating,omitempty"`
	Categories []Category `gorm:"->;many2many:recipe_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category,omitempty"`
}

type Search struct {
	SearchKeyword string   `json:"search_keyword,omitempty" form:"search"`
	Categories    []string `json:"categories,omitempty" form:"search-categories[]"`
	MealTypes     []string `json:"meal_types,omitempty" form:"search-type[]"`
	Courses       []string `json:"courses,omitempty" form:"search-meals[]"`
	Sort          string   `json:"sort,omitempty" form:"search-sort"`
}

var DB *gorm.DB

func init() {
	DB = initializers.GetConnection()
}

func InsertRecipeData(recipe *Recipe) {
	DB.Create(recipe)
	for index := range recipe.Categories {
		DB.Where("name=?", recipe.Categories[index].Name).Find(&recipe.Categories[index])
		if recipe.Categories[index].ID == 0 {
			DB.Create(&recipe.Categories[index])
		}
		DB.Table("recipe_categories").Create(map[string]interface{}{
			"category_id": recipe.Categories[index].ID,
			"recipe_id":   recipe.ID,
		})
	}
}

func GetAllRecipe() []map[string]interface{} {
	var recipes []map[string]interface{}
	DB.Model(Recipe{}).Select("recipes.id, recipes.image_url, recipes.title, recipes.cooking_time, recipes.serving," +
		" recipes.type, recipes.meals, recipes.difficulty_level, floor(avg(rating)) as avg_rating").
		Joins(" left join reviews as r on r.recipe_id = recipes.id").Group("r.recipe_id, recipes.id").
		Find(&recipes)

	return recipes
}

func GetRecipeByEmail(email string) []map[string]interface{} {
	var recipes []map[string]interface{}
	DB.Model(Recipe{}).Select("recipes.id, recipes.image_url, recipes.title, recipes.cooking_time, recipes.serving,"+
		" recipes.type, recipes.meals, recipes.difficulty_level, floor(avg(rating)) as avg_rating").
		Joins("left join reviews as r on r.recipe_id = recipes.id").Group("r.recipe_id, recipes.id").
		Where("recipes.email_id=?", email).
		Find(&recipes)
	return recipes
}

func GetRecipeByID(id int) map[string]interface{} {
	var recipe map[string]interface{}
	log.Println(id)
	fmt.Println(DB.Debug().Model(Recipe{}).Select("recipes.*,FLOOR(AVG(reviews.rating)) AS avg_rating ,STRING_AGG(distinct(categories.name),', ') AS categories").
		Joins("left JOIN recipe_categories ON recipe_categories.recipe_id = recipes.id"+
			" left JOIN categories ON categories.id = recipe_categories.category_id"+
			" left JOIN reviews ON reviews.recipe_id = recipes.id").Where("recipes.id=?", id).
		Group("recipes.id, reviews.recipe_id").Find(&recipe).Error)
	//DB.Model(Recipe{}).Where("id=?", id).Find(&recipe)
	log.Println(recipe)
	//DB.Model(Review{}).Select("floor(avg(rating))").Where("recipe_id=?", recipe.ID).Group("recipe_id").Find(&recipe.AvgRating)
	return recipe
}

func DeleteRecipeByID(id int) {
	DB.Table("recipe_categories").Where("recipe_id=?", id).Delete("recipe_id, category_id")
	DB.Where("id=?", id).Delete(&Recipe{})
}

func EditRecipeData(recipe *Recipe) {
	DB.Model(Recipe{}).Where("id=?", recipe.ID).Updates(recipe)
	DB.Table("recipe_categories").Where("recipe_id=?", recipe.ID).Delete("recipe_id, category_id")
	for index := range recipe.Categories {
		DB.Where("name=?", recipe.Categories[index].Name).Find(&recipe.Categories[index])
		if recipe.Categories[index].ID == 0 {
			DB.Create(&recipe.Categories[index])
		}
		DB.Table("recipe_categories").Create(map[string]interface{}{
			"category_id": recipe.Categories[index].ID,
			"recipe_id":   recipe.ID,
		})
	}
}

func SearchRecipe(searchStruct Search) []map[string]interface{} {
	var recipe []map[string]interface{}
	var courseQuery, categoriesQuery, typeQuery = "", "", ""
	if len(searchStruct.Categories) != 0 {
		categoriesQuery = " AND c.name IN ('" + strings.Join(searchStruct.Categories, "', '") + "')"
	}
	if len(searchStruct.MealTypes) != 0 {
		typeQuery = " AND type IN ('" + strings.Join(searchStruct.MealTypes, "', '") + "')"
	}
	if len(searchStruct.Courses) != 0 {
		courseQuery = " AND ("
		for i := range searchStruct.Courses {
			courseQuery += "meals ILike '%" + searchStruct.Courses[i] + "%'"
			if i != (len(searchStruct.Courses) - 1) {
				courseQuery += " OR "
			}
		}
		courseQuery += ")"
	}
	query := DB.Model(Recipe{}).Distinct("recipes.id, recipes.image_url, recipes.title, recipes.cooking_time, recipes.serving,"+
		" recipes.type, recipes.meals, recipes.difficulty_level, floor(avg(rating)) as avg_rating").
		Joins("left join recipe_categories as rc on recipes.id = rc.recipe_id"+
			" left join categories as c on rc.category_id = c.id"+
			" left join reviews as r on r.recipe_id = recipes.id").Group("r.recipe_id, recipes.id").Where("(title ILike @keyword OR ingredients ILike @keyword OR type ILike @keyword OR meals ILike @keyword OR c.name ILike @keyword)"+categoriesQuery+typeQuery+courseQuery,
		sql.Named("keyword", "%"+searchStruct.SearchKeyword+"%"))
	query.Order(searchStruct.Sort).Find(&recipe)
	return recipe
}
