package product

import (
	"app/internal/domain"
	"app/pkg/store"
	"errors"
	"fmt"
)

type ProductRepository struct {
	jsonStore store.ControllerStorage
}

// NewRepositoryProductInMemory creates a new storage for products
func NewRepositoryProductInMemory(jsonStore store.ControllerStorage) *ProductRepository {
	return &ProductRepository{jsonStore: jsonStore}
}

// Save product
func (pr *ProductRepository) Save(product domain.Product)(pro domain.Product, err error) {
	productList, err := pr.jsonStore.GetAll()
	if err != nil {
		panic("error getting product list")
	}
	// validate unique id and code_value before save it
	for _, prod := range productList{
		if prod.Id == product.Id {
			panic("Id must be unique")
		}
		if prod.Code_value == product.Code_value {
			panic("code value must be unique")
		}
	}
	
	product.Id = len(productList) + 1 // setting id to the product
	// -> save in storage
	pro, err = pr.jsonStore.Save(product)
	if err != nil {
		panic("error saving the product")
	}
	return 
}

// Get product by id
func (pr *ProductRepository) GetById(id int) (product domain.Product, err error) {
	product, err = pr.jsonStore.GetById(id)
	if err != nil {
		err = errors.New("Product not founded")
	}
	fmt.Println("desde repo")
	fmt.Println(product)
	return
}

// Delete product by id
func (pr *ProductRepository) DeleteById(id int) (deleted bool, err error){
	deleted, err = pr.jsonStore.DeleteById(id)
	return 
}

// Update product
func (pr *ProductRepository) Update(product domain.Product) (prod domain.Product, err error) {
	prod, err = pr.jsonStore.Update(product)
	if err != nil {
		panic("Error updating product")
	}
	return
}

// Update product name
func (pr *ProductRepository) UpdateName(name string, id int) (prod domain.Product, err error){
	prod, err = pr.jsonStore.UpdateName(name, id) // utiliza el package de storage 
	if err !=nil {
		panic(err)
	}
	return
}

// Get all products
func (pr *ProductRepository) GetAll() ([] *domain.Product){
	productList, err := pr.jsonStore.GetAll()
	if err != nil {
		panic(err)
	}
	return productList
}