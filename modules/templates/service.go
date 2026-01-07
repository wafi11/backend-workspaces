package templates

import (
	"context"
)

type Service struct {
	repo TemplatesRepository
}

func NewService(repo TemplatesRepository) Service {
	return Service{repo: repo}
}

func (s Service) Create(c context.Context, req CreateTemplateRequest) error {
	return s.repo.Create(c, req)
}

func (s Service) List(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error) {
	return s.repo.List(ctx, req)
}
func (s Service) FindById(ctx context.Context, req int) (*Template, error) {
	return s.repo.FindById(ctx, req)
}
