package products

import (
	"context"
)

type Service struct {
	repo ProductRepository
}

func NewProductServices(repo ProductRepository) Service {
	return Service{repo: repo}
}
func (s *Service) CreateProduct(c context.Context, req Product) (*Product, error) {
	return s.repo.CreateProduct(c, req)
}

func (s *Service) GetProductByID(c context.Context, id int) (*Product, error) {
	return s.repo.GetProductByID(c, id)
}

func (s *Service) GetAllProducts(c context.Context) ([]Product, error) {
	return s.repo.GetAllProducts(c)
}

func (s *Service) UpdateProduct(c context.Context, id int, req Product) error {
	return s.repo.UpdateProduct(c, id, req)
}

func (s *Service) DeleteProduct(c context.Context, id int) error {
	return s.repo.DeleteProduct(c, id)
}
