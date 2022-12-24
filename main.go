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
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/solnsumei/recipe-api/config"
	"github.com/solnsumei/recipe-api/handlers"
	"github.com/solnsumei/recipe-api/services"
)

var recipeHandler *handlers.RecipesHandler

func init() {
	config.LoadEnvVariables()
	ctx := context.Background()

	// Initialize mongo collection
	collection, err := services.InitMongoDB(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize redis client
	redisClient := services.InitRedis(ctx)

	// SeedDB(ctx, collection) // comment out after seeding
	recipeHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	// enable CORS
	router.Use(cors.Default())

	// Routes
	router.POST("/recipes", recipeHandler.NewRecipeHandler)
	router.GET("/recipes", recipeHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipeHandler.GetRecipeHandler)
	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
