package main

import (
	"app/POD/cmd/server/handler"
	"app/POD/internal/domain"
	"app/POD/internal/product"

	"github.com/gin-gonic/gin"
)

func main (){

	// repository
	rp := product.NewRepositoryProductInMemory(make([]*domain.Product, 0), 0)

	// service
	sp := product.NewServiceProduct(*rp)

	// handler
	hd := handler.NewHandlerProduct(sp)

	// server 
	server := gin.Default()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	group := server.Group("/products")

	// handlers
	group.GET("/", hd.GetAll())
	group.POST("/save", hd.Save())
	group.GET("/:id", hd.GetById())
	group.DELETE("/:id", hd.DeleteById())
	group.PUT("/:id", hd.Update())
	group.PATCH("/:id", hd.UpdateName())

	// localhost 
	server.Run(":8081")
}