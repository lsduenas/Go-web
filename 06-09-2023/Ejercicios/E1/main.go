package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Crea un router con gin
	router := gin.Default()

	// Captura la solicitud GET “/hello-world”
  	router.GET("/ping", func(c *gin.Context) {
		text := "pong"

		// encabezado content-type
		c.Header("Content-Type", "text/plain; charset=utf-8")
		
		// Escribir el texto en el body de la respuesta
		c.String(http.StatusOK, text)
   })
   router.Run(":8081") // Corremos nuestro servidor sobre el puerto 8081


}