package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var recipes []Recipe

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func NewRecipeHandler(context *gin.Context) {
	var recipe Recipe
	if err := context.ShouldBindJSON(&recipe); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)

	context.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(context *gin.Context) {
	context.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(context *gin.Context) {
	id := context.Param("id")

	var recipe Recipe
	if err := context.ShouldBindJSON(&recipe); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, recp := range recipes {
		if recp.ID == id {
			recipes[i] = recipe
			context.JSON(http.StatusOK, recipe)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found!",
	})
}

func DeleteRecipeHandler(context *gin.Context) {
	id := context.Param("id")

	for i, recp := range recipes {
		if recp.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			context.JSON(http.StatusOK, gin.H{
				"message": "Recipe Deleted Successfully!",
			})
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found!",
	})
}

func SearchRecipesHandler(context *gin.Context) {
	tag := context.Query("tag")
	matchingRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		for _, t := range recipe.Tags {
			if strings.EqualFold(tag, t) {
				matchingRecipes = append(matchingRecipes, recipe)
				break
			}
		}
	}

	context.JSON(http.StatusOK, matchingRecipes)
}

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	if err := json.Unmarshal(file, &recipes); err != nil {
		return
	}
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	if err := router.Run(); err != nil {
		return
	}
}
