package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value,"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

var productsList = []Product{}

func createProductList() {
	file, err := os.Open("./products.json")
	if err != nil {
		panic("The file can´t be opened")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&productsList); err!= nil {
		panic("Something went wrong with decoder process")
	}
	for _, product := range productsList{
		fmt.Printf("ID: %d, Name: %s, Quantity: %d, Code value: %s, Is Published: %t, Expiration %s, Price: %.2f\n", product.Id, product.Name, product.Quantity, product.Code_value, product.Is_published, product.Expiration, product.Price)
	}
}

// functions for each endpoint
func Pong(ctx *gin.Context){
	ctx.String(http.StatusOK, "Pong")
}

func GetAllProducts(ctx *gin.Context){
	ctx.JSON(http.StatusAccepted, productsList)
}

func GetProductById(ctx *gin.Context){
	id, _ := strconv.Atoi(ctx.Param("id"))
	product := Product{}
	for _, p := range productsList {
		if p.Id == id {
			product = p
		}
	}
	ctx.JSON(http.StatusAccepted, product)
}

func GetProductsByPrice(ctx *gin.Context){
	priceParam, _ := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	proList := [] Product{}
	for _, p := range productsList{
		if p.Price > priceParam {
			proList = append(proList, p)
		}
	}
	ctx.JSON(http.StatusAccepted, proList)
}

func main() {
	createProductList()
	server := gin.Default()
	// handlers
	server.GET("/ping", Pong)
	server.GET("/products", GetAllProducts)
	server.GET("/products/:id", GetProductById)
	server.GET("/products/search", GetProductsByPrice)

	// localhost 
	server.Run(":8081")
}