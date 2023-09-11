package main

import (
	"app/PracticaPost/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	// dependencies 
	db := make([]*handlers.Product, 0)
	ct := handlers.NewControllerProducts(db, 0)
	
	server := gin.Default()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// handlers
	server.POST("/products", ct.Save())
	server.GET("/products/:id", ct.GetById())

	// localhost 
	server.Run(":8081")
}