package main

import (
	"app/cmd/server/handler"
	"app/internal/product"
	"app/pkg/store"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// ENV variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error al intentar cargar archivo .env")
	}
	token := os.Getenv("TOKEN")

	// json store
	file_path := os.Getenv("file_path")
	js := store.NewControllerStorage(file_path)

	// repository
	rp := product.NewRepositoryProductInMemory(*js)

	// service
	sp := product.NewServiceProduct(*rp)

	// handler (service and token as parameters)
	hd := handler.NewHandlerProduct(sp, token)

	// server
	server := gin.Default()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	group := server.Group("/products")

	// handlers
	group.GET("/", hd.GetAll())
	group.POST("/", hd.Save())
	group.GET("/:id", hd.GetById())
	group.DELETE("/:id", hd.DeleteById())
	group.PUT("/:id", hd.Update())
	group.PATCH("/:id", hd.UpdateName())

	// localhost
	server.Run(":8081")
}
