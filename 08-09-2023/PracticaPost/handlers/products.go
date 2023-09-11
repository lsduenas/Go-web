package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// store Product in the slice
type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

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

// response body
type ResponseBody struct {
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

// ControllerProducts is an struct that represents a controller for products
// exposing methods to handle products
type ControllerProducts struct {
	db     []*Product
	lastId int
}

func NewControllerProducts(db []*Product, lastId int) *ControllerProducts {
	return &ControllerProducts{
		db:     db,
		lastId: lastId,
	}
}

func Validator(pr *Product) (err error) { // Se valida al recibir en el request body

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

func (c *ControllerProducts) Save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// -> header
		token := ctx.GetHeader("Authorization")
		if token != "123" {
			code := http.StatusUnauthorized
			body := ResponseBody{Message: "Invalid token"}

			ctx.JSON(code, body)
			return
		}

		// -> body
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

		// process
		// -> deserialization
		pr := &Product{
			Name:         reqBody.Name,
			Quantity:     reqBody.Quantity,
			Code_value:   reqBody.Code_value,
			Is_published: reqBody.Is_published,
			Expiration:   reqBody.Expiration,
			Price:        reqBody.Price,
		}
		// autoincrement
		pr.Id = c.lastId + 1

		// -> validation
		if err := Validator(pr); err != nil {
			code := http.StatusConflict
			body := ResponseBody{Message: "invalid product"}
			ctx.JSON(code, body)
			return
		}

		// validate unique id and code_value before save it
		for _, prod := range c.db {
			if prod.Id == pr.Id {
				code := http.StatusConflict
				body := ResponseBody{Message: "id must be unique"}
				ctx.JSON(code, body)
				return
			}
			if prod.Code_value == pr.Code_value {
				code := http.StatusConflict
				body := ResponseBody{Message: "code value must be unique"}
				ctx.JSON(code, body)
				return
			}
		}

		// -> save in storage
		c.db = append(c.db, pr)

		c.lastId++

		// response
		code := http.StatusCreated
		body := ResponseBody{
			Message: "Product created",
			Data: &Data{ // serialization
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

// Get product by id
func (cp *ControllerProducts) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(cp.db) == 0 {
			code := http.StatusFailedDependency
			Message := "There are not products yet"
			ctx.JSON(code, Message)
			return
		}
		data := Data{}
		for _, prod := range cp.db {
			if id, _ := strconv.Atoi(ctx.Param("id")); prod.Id == id {
				data = Data{
					Id:           prod.Id,
					Name:         prod.Name,
					Quantity:     prod.Quantity,
					Code_value:   prod.Code_value,
					Is_published: prod.Is_published,
					Expiration:   prod.Expiration,
					Price:        prod.Price,
				}
			} else {
				code := http.StatusBadRequest
				Message := "Id was not founded"
				ctx.JSON(code, Message)
				return
			}
		}
		code := http.StatusOK
		body := ResponseBody{
			Message: "Product was founded",
			Data:    &data,
		}
		ctx.JSON(code, body)
	}
}
