package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SayHello(ctx *gin.Context){
	ctx.JSON(http.StatusOK, "Hello world!")
}
func SayHelloWithName(ctx *gin.Context){
	ctx.String(http.StatusOK, "Hello " + ctx.Param("name"))
}

func SayHelloWithQueryParams(ctx *gin.Context){
	ctx.String(http.StatusOK, "Hello, this query param is "+ ctx.Query("age"))
}

func main(){
	server := gin.Default()
	helloPaths := server.Group("/hello")

	helloPaths.GET("", SayHello)
	helloPaths.GET("/:name", SayHelloWithName)
	helloPaths.GET("/queryParams", SayHelloWithQueryParams)

	// localhost 
	server.Run(":8081")
}