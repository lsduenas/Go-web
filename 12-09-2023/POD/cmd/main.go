package main

import (
	"app/cmd/server/handler"
	"app/cmd/server/middlewares"
	"app/internal/product"
	"app/pkg/store"
	"log"
	"os"
	docs "app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
// @title Products API - Bootcamp Go
// @version 1.0
// @description This API Handle MELI Products.
// @termsOfService https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /products
func main() {

	// ENV variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error trying to load .env file")
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
	server := gin.New()
	server.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	group := server.Group("/products")
	group.Use(middlewares.Authenticator())
	// docs swagger
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
