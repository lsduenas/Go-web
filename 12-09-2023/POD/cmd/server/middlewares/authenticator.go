package middlewares

import (
	"app/cmd/server/handler"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Authenticator() gin.HandlerFunc {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error trying to load .env file")
	}
	token := os.Getenv("TOKEN")
	return func(ctx *gin.Context) {

		//  before handler
		tokenHeader := ctx.GetHeader("Token")
		if tokenHeader != token {
			code := http.StatusUnauthorized
			responseBody := handler.ResponseBody{
				Message: "Unauthorized",
				Data:    nil,
			}
			ctx.JSON(code, responseBody)
			ctx.Abort() // importante, de lo contrario continua la ejecución
			return
		}
		ctx.Next()

		// after handler 
	}
}
