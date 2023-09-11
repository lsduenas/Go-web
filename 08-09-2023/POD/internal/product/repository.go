package product

import (
	"app/POD/internal/domain"
	"errors"
	"fmt"
)

type ProductRepository struct {
	db     []*domain.Product
	lastId int
}

// NewRepositoryProductInMemory creates a new storage for products
func NewRepositoryProductInMemory(db []*domain.Product, lastId int) *ProductRepository {
	return &ProductRepository{db: db, lastId: lastId}
}

func (pr *ProductRepository) Save(product domain.Product)(pro domain.Product, err error) {

	// validate unique id and code_value before save it
	for _, prod := range pr.db {
		if prod.Id == product.Id {
			err = errors.New("Id must be unique")
		}
		if prod.Code_value == product.Code_value {
			err = errors.New("code value must be unique")
		}
	}
	product.Id = pr.lastId + 1
	// -> save in storage
	pr.db = append(pr.db, &product)
	pr.lastId++
	pro = product
	return
}

// Get product by id
func (pr *ProductRepository) GetById(id int) (product domain.Product, err error) {
	if len(pr.db) == 0 {
		panic("There are not products yet")
	}
	for _, prod := range pr.db {
		if prod.Id == id {
			product = *prod
		} else {
			err = errors.New("Product not found")
		}
	}
	return
}

// Delete product by id
func (pr *ProductRepository) DeleteById(id int) (deleted bool, err error){
	var index int
	for i := range pr.db {
		if pr.db[i].Id == id {
			index = i
			deleted = true
		}
	}
	if !deleted {
		err = fmt.Errorf("Product %d not found", id)
	}
	pr.db = append(pr.db[:index], pr.db[index+1:]...)
	return 
}

// Update product
func (pr *ProductRepository) Update(product domain.Product) (prod domain.Product, err error) {
	for key, value := range pr.db {
		if value.Id == product.Id {
			// update all products fields
			pr.db[key].Name = product.Name
			pr.db[key].Quantity = product.Quantity
			pr.db[key].Code_value = product.Code_value
			pr.db[key].Quantity = product.Quantity
			pr.db[key].Is_published = product.Is_published
			pr.db[key].Expiration = product.Expiration
			pr.db[key].Price = product.Price
			prod = *value
			break
		} else {
			err = fmt.Errorf("Product ID not found")
		}
	}
	return
}

func (pr *ProductRepository) UpdateName(name string, id int) (prod domain.Product, err error){
	for key, value := range pr.db {
		if value.Id == id {
			// update product name
			pr.db[key].Name = name
			prod = *value
			break
		} else {
			err = fmt.Errorf("Product ID not found")
		}
	}
	return
}

func (pr *ProductRepository) GetAll() ([] *domain.Product){
	return pr.db
}