package handler

import (
	"app/internal/domain"
	"app/internal/product"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// store product json from request
type RequestBody struct {
	Name         string  `json:"name" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	Code_value   string  `json:"code_value," binding:"required"`
	Is_published bool    `json:"is_published"` // it can be false
	Expiration   string  `json:"expiration" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
}

// store product json from response
type Data struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value,"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type HandlerProduct struct {
	sp *product.ServiceProduct
	token string
}

type NameProduct struct {
	Name string
}

// constructor
func NewHandlerProduct(sp *product.ServiceProduct, token string) *HandlerProduct {
	return &HandlerProduct{sp: sp, token: token}
}

// response body
type ResponseBody struct {
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

func Validator(pr domain.Product) (err error) { // Se valida al recibir en el request body

	// required
	if pr.Name == "" {
		err = errors.New("name is required")
		return
	}

	if pr.Code_value == "" {
		err = errors.New("Code value is required")
		return
	}

	if pr.Expiration == "" {
		err = errors.New("Expiration is required")
		return
	}
	if pr.Price == 0.0 {
		err = errors.New("Price is required")
		return
	}
	return
}

func ValidateToken(token string) (err error){

	if token != os.Getenv("TOKEN") {
		err = fmt.Errorf("Invalid token")
	}
	return
}

// CRUD
func (hd *HandlerProduct) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Token")
		// validate token
		
		if token != hd.token {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}
		
		// request
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalid ID", Data: nil}

			ctx.JSON(code, body)
			return
		}
		fmt.Println(id)
		// process
		pr, err := hd.sp.GetById(id)
		
		if err != nil {
			code := http.StatusNotFound
			body := ResponseBody{
				Message: "Product not found",
				Data:    nil,
			}
			ctx.JSON(code, body)
			return
		}

		// response
		code := http.StatusOK
		body := ResponseBody{
			Message: "Success",
			Data: &Data{
				Id:           pr.Id,
				Name:         pr.Name,
				Quantity:     pr.Quantity,
				Code_value:   pr.Code_value,
				Is_published: pr.Is_published,
				Expiration:   pr.Expiration,
				Price:        pr.Price,
			},
		}
		ctx.JSON(code, body)
	}
}

func (hd *HandlerProduct) Save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		// validate token
		er := ValidateToken(token)
		if er != nil {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}

		var reqBody RequestBody
		err := ctx.ShouldBindJSON(&reqBody)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalid request body"}
			ctx.JSON(code, body)
			return
		}

		// parse date
		expiration_date, err := time.Parse("02/01/2006", reqBody.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalidad date format "}
			ctx.JSON(code, body)
			return
		}
		year, _, _ := expiration_date.Date()

		if year < 2023 {
			code := http.StatusBadRequest
			message := "Expiration year can't be less than 2023"
			ctx.JSON(code, message)
			return
		}

		pr := &domain.Product{
			Name:         reqBody.Name,
			Quantity:     reqBody.Quantity,
			Code_value:   reqBody.Code_value,
			Is_published: reqBody.Is_published,
			Expiration:   reqBody.Expiration,
			Price:        reqBody.Price,
		}

		product, err := hd.sp.Save(*pr)
		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBody{
				Message: "Internal server error",
				Data:    nil,
			}
			ctx.JSON(code, body)
			return
		}

		code := http.StatusCreated
		body := ResponseBody{
			Message: "Success",
			Data: &Data{
				Id:           product.Id,
				Name:         product.Name,
				Quantity:     product.Quantity,
				Code_value:   product.Code_value,
				Is_published: product.Is_published,
				Expiration:   product.Expiration,
				Price:        product.Price,
			},
		}
		ctx.JSON(code, body)
	}
}

func (hd *HandlerProduct) DeleteById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		// validate token
		err := ValidateToken(token)
		if err != nil {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}

		// request
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Missing id", Data: nil}

			ctx.JSON(code, body)
			return
		}
		deleted, err := hd.sp.DeleteById(id)

		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBody{Message: "Product id not found to delete it", Data: nil}
			ctx.JSON(code, body)
			return
		}

		if deleted {
			code := http.StatusAccepted
			body := ResponseBody{Message: "Success, Product was deleted", Data: nil}
			ctx.JSON(code, body)
			return
		}
	}
}

func (hd *HandlerProduct) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		// validate token
		er := ValidateToken(token)
		if er != nil {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}

		var reqBody RequestBody
		err := ctx.ShouldBindJSON(&reqBody)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalid request body"}
			ctx.JSON(code, body)
			return
		}

		// parse date
		expiration_date, err := time.Parse("02/01/2006", reqBody.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalidad date format "}
			ctx.JSON(code, body)
			return
		}
		year, _, _ := expiration_date.Date()

		if year < 2023 {
			code := http.StatusBadRequest
			message := "Expiration year can't be less than 2023"
			ctx.JSON(code, message)
			return
		}
		id, err:= strconv.Atoi(ctx.Param("id")) // TO DO validate id - idempotencia
		if err != nil {
			panic("Id not found")
		}

		pr := &domain.Product{
			Id: 		  id,
			Name:         reqBody.Name,
			Quantity:     reqBody.Quantity,
			Code_value:   reqBody.Code_value,
			Is_published: reqBody.Is_published,
			Expiration:   reqBody.Expiration,
			Price:        reqBody.Price,
		}
		product, err := hd.sp.Update(*pr)
		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBody{
				Message: "Internal server error",
				Data:    nil,
			}
			ctx.JSON(code, body)
			return
		}

		code := http.StatusCreated
		body := ResponseBody{
			Message: "Success",
			Data: &Data{
				Id:           product.Id,
				Name:         product.Name,
				Quantity:     product.Quantity,
				Code_value:   product.Code_value,
				Is_published: product.Is_published,
				Expiration:   product.Expiration,
				Price:        product.Price,
			},
		}
		ctx.JSON(code, body)
	}
}


func (hd *HandlerProduct) UpdateName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		// validate token
		er := ValidateToken(token)
		if er != nil {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}

		var reqBody NameProduct
		err := ctx.ShouldBindJSON(&reqBody)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBody{Message: "Invalid request body"}
			ctx.JSON(code, body)
			return
		}
		
		id, err:= strconv.Atoi(ctx.Param("id")) // TO DO validate id - idempotencia
		if err != nil {
			panic("Id not found")
		}

		product, err := hd.sp.UpdateName(reqBody.Name, id)
		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBody{
				Message: "Internal server error",
				Data:    nil,
			}
			ctx.JSON(code, body)
			return
		}

		code := http.StatusCreated
		body := ResponseBody{
			Message: "Success",
			Data: &Data{
				Id:           product.Id,
				Name:         product.Name,
				Quantity:     product.Quantity,
				Code_value:   product.Code_value,
				Is_published: product.Is_published,
				Expiration:   product.Expiration,
				Price:        product.Price,
			},
		}
		ctx.JSON(code, body)
	}
}

func (hd *HandlerProduct) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		// validate token
		err := ValidateToken(token)
		if err != nil {
			code := http.StatusForbidden
			body := ResponseBody{Message: "Invalid Token", Data: nil}

			ctx.JSON(code, body)
			return
		}
		
		productList := hd.sp.GetAll()
		code := http.StatusOK
		ctx.JSON(code, productList)
	}
}

