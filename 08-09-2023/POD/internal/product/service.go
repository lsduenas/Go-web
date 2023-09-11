package product

import (
	"app/POD/internal/domain"
)

type ServiceProduct struct {
	rp ProductRepository
}

// constructor 
func NewServiceProduct(rp ProductRepository) *ServiceProduct{
	return &ServiceProduct{rp: rp}
}

// save product 
func (s *ServiceProduct) Save(product domain.Product) (domain.Product, error){ 
	return s.rp.Save(product)
}

// get product by id 
func (s *ServiceProduct) GetById(id int) (domain.Product, error)  {
	return s.rp.GetById(id)
}

// delete product by id
func (s *ServiceProduct) DeleteById(id int) (bool, error){
	return s.rp.DeleteById(id)
}

// update all product fields
func (s *ServiceProduct) Update(product domain.Product) (domain.Product, error){
	return s.rp.Update(product)
}

// update product name
func (s *ServiceProduct) UpdateName(name string, id int) (domain.Product, error){
	return s.rp.UpdateName(name, id)
}

// get all products
func (s *ServiceProduct) GetAll() ([] *domain.Product){
	return s.rp.GetAll()
}