package store

import (
	"app/internal/domain"
	"encoding/json"
	"errors"
	"os"
)

type ProductStore interface {
	GetAll() (productList []*domain.Product, err error)
	GetById(id int) (product domain.Product, err error)
	Save(product domain.Product) (domain.Product, error)
	DeleteById(id int) bool
	Update(product domain.Product) domain.Product
	UpdateName(name string, id int) domain.Product
	WriteJSONFile(productList []domain.Product) error
}

type ControllerStorage struct {
	filepath string
}

func NewControllerStorage(filepath string) *ControllerStorage {
	return &ControllerStorage{filepath: filepath}
}


func (c *ControllerStorage) GetAll() (productList []*domain.Product, err error) {
	file, err := os.ReadFile(c.filepath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		panic(err)
	}
	return productList, nil
}

func (c *ControllerStorage) GetById(id int) (product domain.Product, err error) {
	productList, err := c.GetAll()
	if err != nil {
		panic("Error retrieving product")
	}
	for _, value := range productList {
		if value.Id == id {
			product = *value
		} else {
			err = errors.New("Product not founded")
		}
	}
	return
}

func (c *ControllerStorage) Save(product domain.Product) (pro domain.Product, err error) {
	productList, err := c.GetAll()
	if err != nil {
		panic("Error getting product list")
	}
	productList = append(productList, &product)
	err = c.WriteJSONFile(productList)
	pro = product
	return 
}

func (c *ControllerStorage) DeleteById(id int) (deleted bool, err error) {
	productList, err := c.GetAll()
	if err != nil {
		err = errors.New("Error getting product list")
	}
	var index int
	for i := range productList {
		if productList[i].Id == id {
			index = i
			deleted = true
		}
	}
	if !deleted {
		err = errors.New("Product id not found")
	}
	productList= append(productList[:index], productList[index+1:]...)
	
	er := c.WriteJSONFile(productList)
	if er != nil{
		err = errors.New("Error updating json file")
	}
	return
}

func (c *ControllerStorage) Update(product domain.Product) (productUpdated domain.Product, err error){
	productList, err := c.GetAll()
	if err != nil {
		panic("Error getting product list")
	}
	for _, value := range productList {
		if value.Id == product.Id {
			value.Name = product.Name
			value.Quantity = product.Quantity
			value.Code_value = product.Code_value
			value.Is_published = product.Is_published
			value.Expiration = product.Expiration
			value.Price = product.Price
			productUpdated = *value
		}
	}
	// write in json file
	er := c.WriteJSONFile(productList)
	if er != nil {
		panic("Error writing in product list")
	}
	return
}

func (c *ControllerStorage) UpdateName(name string, id int) (productUpdated domain.Product, err error) {
	productList, err := c.GetAll()
	if err != nil {
		panic("Error getting product list")
	}
	for _, value := range productList {
		if value.Id == id {
			value.Name = name
			productUpdated = *value
		}
	}
	// write in json file
	er := c.WriteJSONFile(productList)
	if er != nil {
		panic("Error writing in product list")
	}
	return
}

func (c *ControllerStorage) WriteJSONFile(productList []*domain.Product) error {
	file, err := os.Create(c.filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)                    // crea el encoder
	if err := encoder.Encode(productList); err != nil { // a través del encoder convierte la lista de structs en el json file
		return err
	}
	return nil
}
