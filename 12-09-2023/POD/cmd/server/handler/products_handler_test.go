package handler

import (
	"app/internal/product"
	"app/pkg/store"
	"encoding/json"


	"net/http"
	"net/http/httptest"


	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

// Tests for HandlerProduct.GetByID
func CreateTestServerForProductsHandler(hd *HandlerProduct) *gin.Engine {
	r := gin.New()
	r.GET("/products/:id", hd.GetById())
	return r
}

func TestFunctional_ProductsHandler_GetByID(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		// arrange
		// -> expectations

		// -> setup
		// dependencies

		// json store
		
	
		js := store.NewControllerStorage("../../../productstesting.json")

		// repository
		rp := product.NewRepositoryProductInMemory(*js)
		sv := product.NewServiceProduct(*rp)

		token := "TOP SUPER SECRET"
		hd := NewHandlerProduct(sv, token)

		server := CreateTestServerForProductsHandler(hd)

		// -> input
		request := httptest.NewRequest("GET", "/products/7", nil)
		request.Header.Add("Token", "TOP SUPER SECRET")
		t.Log("HEADER", request.Header)
		response := httptest.NewRecorder()

		// act
		server.ServeHTTP(response, request)

		var responseBody ResponseBody
		err := json.NewDecoder(response.Body).Decode(&responseBody)
		assert.NoError(t, err)

		// assert
		expectedStatusCode := http.StatusOK
		expectedResponseBody := ResponseBody{
			Message: "Success",
			Data: &Data{
				Id:           7,
				Name:         "Melon - Honey Dew",
				Quantity:     165,
				Code_value:   "S52381G",
				Is_published: true,
				Expiration:   "01/06/2021",
				Price:        622.33,
			},
		}

		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedResponseBody, responseBody)
		//assert.Equal(t, string(expectedResponseBody), response.Body.String())
		assert.Equal(t, expectedHeaders, response.Header())
	})
}
