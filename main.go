// Recipe API.
//
// This is a sample reciples API. You can find out more about the API at
// https://github.com/solnsumei/recipe-api
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Solomon Nsumei <solnsumei@gmail.com> https://solnsumei.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)

	file, err := os.ReadFile("recipe.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(file), &recipes); err != nil {
		panic(err)
	}
}

// swagger:operation POST /recipes recipes addRecipe
// Add a new recipe
// ---
// parameters:
//   - name: body
//     in: body
//     description: request body
//     required: true
//     type: object
//
// produces:
// - application/json
// responses:
//
//	  '200':
//		    description: Successful operation
//	  '400':
//	 	    description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	   '200':
//			description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//   - name: body
//     in: body
//     description: request body
//     required: true
//     type: object
//
// produces:
// - application/json
// responses:
//
//	  '200':
//		    description: Successful operation
//	  '400':
//	 	    description: Invalid input
//	  '404':
//	 	    description: Invalid recipe ID
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Resource not found",
		})
		return
	}

	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	  '200':
//		    description: Successful operation
//	  '404':
//	 	    description: Invalid recipe ID
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Resource not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Search recipe with tags
// ---
// parameters:
//   - name: tag
//     in: query
//     description: Fetch recipes that match tag
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	  '200':
//		    description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()

	// enable CORS
	router.Use(cors.Default())

	// Routes
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.Run(":8080")
}
