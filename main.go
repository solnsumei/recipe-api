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
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/solnsumei/recipe-api/config"
	"github.com/solnsumei/recipe-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipeHandler *handlers.RecipesHandler

func init() {
	config.LoadEnvVariables()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")

	// SeedDB(ctx, collection) // comment out after seeding
	recipeHandler = handlers.NewRecipesHandler(ctx, collection)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
}

func main() {
	router := gin.Default()

	// enable CORS
	router.Use(cors.Default())

	// Routes
	router.POST("/recipes", recipeHandler.NewRecipeHandler)
	router.GET("/recipes", recipeHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)

	router.Run(":8080")
}
