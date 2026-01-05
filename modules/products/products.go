package products

import (
	"context"
	"time"
)

type Product struct {
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Container struct {
	ContainerID       int               `json:"container_id"`
	ProductName       *string           `json:"product_name"`
	ProductID         int               `json:"product_id" validate:"required"`
	ContainerName     string            `json:"container_name" validate:"required"`
	DockerContainerID string            `json:"docker_container_id,omitempty"`
	Image             string            `json:"image" validate:"required"`
	Tag               string            `json:"tag"`
	HostPort          int               `json:"host_port"`
	ContainerPort     int               `json:"container_port"`
	Environment       map[string]string `json:"environment"`
	Volumes           []string          `json:"volumes"`
	Command           []string          `json:"command"`
	Status            string            `json:"status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
type ProductWithContainers struct {
	Product
	Containers []Container `json:"containers"`
}

type ProductRepository interface {
	CreateProduct(c context.Context, req Product) (*Product, error)
	GetProductByID(c context.Context, id int) (*Product, error)
	GetAllProducts(c context.Context) ([]Product, error)
	UpdateProduct(c context.Context, id int, req Product) error
	DeleteProduct(c context.Context, id int) error

	// Containers
	CreateContainer(c context.Context, req Container) (*Container, error)
	FindAll(c context.Context) ([]Container, error)
	GetContainerByID(c context.Context, id int) (*Container, error)
	GetContainersByProductID(c context.Context, productID int) ([]Container, error)
	UpdateContainer(c context.Context, id int, req Container) error
	UpdateContainerStatus(c context.Context, id int, status, dockerID string) error
	DeleteContainer(c context.Context, id int) error
	GetProductWithContainers(c context.Context, id int) (*ProductWithContainers, error)
}
type Services interface {
	CreateProduct(c context.Context, req Product) (*Product, error)
	GetProductByID(c context.Context, id int) (*Product, error)
	GetAllProducts(c context.Context) ([]Product, error)
	UpdateProduct(c context.Context, id int, req Product) error
	DeleteProduct(c context.Context, id int) error

	// Containers
	CreateContainer(c context.Context, req Container) (*Container, error)
	GetContainerByID(c context.Context, id int) (*Container, error)
	GetContainersByProductID(c context.Context, productID int) ([]Container, error)
	UpdateContainer(c context.Context, id int, req Container) error
	DeleteContainer(c context.Context, id int) error

	// Docker Operations
	StartContainer(c context.Context, containerID int) error
	StopContainer(c context.Context, containerID int) error
	RestartContainer(c context.Context, containerID int) error
	GetContainerLogs(c context.Context, containerID int, tail int) (string, error)
	GetContainerStatus(c context.Context, containerID int) (map[string]interface{}, error)

	// Combined
	GetProductWithContainers(c context.Context, id int) (*ProductWithContainers, error)
}

type DockerConfig struct {
	Image         string            `json:"image"`
	Tag           string            `json:"tag"`
	Port          int               `json:"port"`           // deprecated, untuk backward compatibility
	HostPort      int               `json:"host_port"`      // port di host machine
	ContainerPort int               `json:"container_port"` // port di dalam container
	Environment   map[string]string `json:"environment"`
	Volumes       []string          `json:"volumes"`
	Command       []string          `json:"command"`
}
