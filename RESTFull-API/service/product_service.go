package service

import (
	"be-golang-chapter-24/RESTFull-API/collections"
	"be-golang-chapter-24/RESTFull-API/repository"
)

type ProductService interface {
	GetAllProducts() ([]collections.Product, error)
	GetProductByID(id int) (*collections.Product, error)
	CreateProduct(product *collections.Product) error
	UpdateProduct(product *collections.Product) error
	DeleteProduct(id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAllProducts() ([]collections.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetProductByID(id int) (*collections.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) CreateProduct(product *collections.Product) error {
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product *collections.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
