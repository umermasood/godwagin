// Recipes API
//
// This is the Recipe API
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: M Umer Masood <umermasood.dev@gmail.com> https://github.com/umermasood
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
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

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
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

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func ListRecipesHandler(context *gin.Context) {
	context.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid recipe ID
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
			recipe.ID = id
			recipe.PublishedAt = recp.PublishedAt
			recipes[i] = recipe
			context.JSON(http.StatusOK, recipe)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found!",
	})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
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

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
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
