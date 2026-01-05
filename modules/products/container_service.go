package products

import "context"

// Container Methods
func (s *Service) CreateContainer(c context.Context, req Container) (*Container, error) {
	return s.repo.CreateContainer(c, req)
}

func (s *Service) GetContainerByID(c context.Context, id int) (*Container, error) {
	return s.repo.GetContainerByID(c, id)
}
func (s *Service) FindAll(c context.Context) ([]Container, error) {
	return s.repo.FindAll(c)
}

func (s *Service) GetContainersByProductID(c context.Context, productID int) ([]Container, error) {
	return s.repo.GetContainersByProductID(c, productID)
}

func (s *Service) UpdateContainer(c context.Context, id int, req Container) error {
	return s.repo.UpdateContainer(c, id, req)
}

func (s *Service) DeleteContainer(c context.Context, id int) error {
	// Stop container if running
	container, err := s.repo.GetContainerByID(c, id)
	if err == nil && container.Status == "running" {
		s.StopContainer(c, id)
	}
	return s.repo.DeleteContainer(c, id)
}

func (s *Service) GetProductWithContainers(c context.Context, id int) (*ProductWithContainers, error) {
	return s.repo.GetProductWithContainers(c, id)
}
