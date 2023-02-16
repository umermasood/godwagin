package main

import (
	"embed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Picture     string       `json:"imageURL"`
}

type Ingredient struct {
	Quantity string `json:"quantity"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

var recipes []Recipe

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"recipes": recipes,
	})
}

func RecipeHandler(c *gin.Context) {
	for _, recipe := range recipes {
		if recipe.ID == c.Param("id") {
			c.HTML(http.StatusOK, "recipe.tmpl", gin.H{
				"recipe": recipe,
			})
			return
		}
	}
	c.File("404.html")
}

//go:embed assets/* templates/* 404.html recipes.json
var f embed.FS

func init() {
	recipes = make([]Recipe, 0)
	data, _ := f.ReadFile("recipes.json")
	if err := json.Unmarshal(data, &recipes); err != nil {
		panic(err)
	}
}

func main() {
	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))

	fsys, err := fs.Sub(f, "assets")
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.SetHTMLTemplate(templ)
	router.StaticFS("/assets", http.FS(fsys))
	router.GET("/", IndexHandler)
	router.GET("/recipes/:id", RecipeHandler)
	panic(router.Run(":4000"))
}
