package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
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
	if err := router.Run(); err != nil {
		return
	}
}
