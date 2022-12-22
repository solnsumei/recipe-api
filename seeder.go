package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type recipeData struct {
	//swagger:ignore
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}

func readJsonData() []interface{} {
	var recipes []recipeData

	file, err := os.ReadFile("recipe.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(file), &recipes); err != nil {
		panic(err)
	}

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	return listOfRecipes
}

func SeedDB(ctx context.Context, collection *mongo.Collection) {
	data := readJsonData()
	insertManyResult, err := collection.InsertMany(
		ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ",
		len(insertManyResult.InsertedIDs))
}
