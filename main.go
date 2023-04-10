package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"nvsu-api/http"
	"nvsu-api/repository"
)

func main() {
	// create a database connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	// create a repository
	repository := repository.NewRepository(client.Database("nvsu_db"))

	// create an http server
	server := http.NewServer(repository)

	// create a gin router
	router := gin.Default()
	{
		// user
		router.GET("/users/:email", server.GetUser)
		router.POST("/users", server.CreateUser)
		router.PUT("/users/:email", server.UpdateUser)
		router.DELETE("/users/:email", server.DeleteUser)

		// product
		router.POST("/products", server.CreateProduct)
		router.GET("/product/:id", server.GetProduct)
		router.POST("/getProducts", server.GetProducts)

		// order
		router.POST("/createOrder", server.CreateOrder)
		router.GET("/order/:id", server.GetOrder)
		router.PUT("/order/:id", server.UpdateOrder)
	}

	// start the router
	router.Run(":3000")
}
