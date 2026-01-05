package products

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) ProductRepository {
	return &Repository{DB: db}
}
func (r *Repository) CreateProduct(c context.Context, req Product) (*Product, error) {
	query := `INSERT INTO products (product_name, description) 
	          VALUES ($1, $2) 
	          RETURNING product_id, created_at, updated_at`

	var p Product
	err := r.DB.QueryRowContext(c, query, req.ProductName, req.Description).
		Scan(&p.ProductID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	p.ProductName = req.ProductName
	p.Description = req.Description
	return &p, nil
}

func (r *Repository) GetProductByID(c context.Context, id int) (*Product, error) {
	query := `SELECT product_id, product_name, description, created_at, updated_at 
	          FROM products WHERE product_id = $1`

	var p Product
	err := r.DB.QueryRowContext(c, query, id).Scan(
		&p.ProductID, &p.ProductName, &p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) GetAllProducts(c context.Context) ([]Product, error) {
	query := `SELECT product_id, product_name, description, created_at, updated_at 
	          FROM products ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ProductID, &p.ProductName, &p.Description, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) UpdateProduct(c context.Context, id int, req Product) error {
	query := `UPDATE products 
	          SET product_name = $1, description = $2, updated_at = $3 
	          WHERE product_id = $4`

	_, err := r.DB.ExecContext(c, query, req.ProductName, req.Description, time.Now(), id)
	return err
}

func (r *Repository) DeleteProduct(c context.Context, id int) error {
	query := `DELETE FROM products WHERE product_id = $1`
	_, err := r.DB.ExecContext(c, query, id)
	return err
}

// Container Methods
func (r *Repository) CreateContainer(c context.Context, req Container) (*Container, error) {
	envJSON, _ := json.Marshal(req.Environment)
	volJSON, _ := json.Marshal(req.Volumes)
	cmdJSON, _ := json.Marshal(req.Command)

	query := `INSERT INTO containers 
	          (product_id, container_name, image, tag, host_port, container_port, 
	           environment, volumes, command, status) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	          RETURNING container_id, created_at, updated_at`

	var cont Container
	err := r.DB.QueryRowContext(c, query,
		req.ProductID, req.ContainerName, req.Image, req.Tag,
		req.HostPort, req.ContainerPort, envJSON, volJSON, cmdJSON, "stopped",
	).Scan(&cont.ContainerID, &cont.CreatedAt, &cont.UpdatedAt)

	if err != nil {
		return nil, err
	}

	cont.ProductID = req.ProductID
	cont.ContainerName = req.ContainerName
	cont.Image = req.Image
	cont.Tag = req.Tag
	cont.HostPort = req.HostPort
	cont.ContainerPort = req.ContainerPort
	cont.Environment = req.Environment
	cont.Volumes = req.Volumes
	cont.Command = req.Command
	cont.Status = "stopped"

	return &cont, nil
}

func (r *Repository) FindAll(c context.Context) ([]Container, error) {
	query := `
		SELECT 
			c.container_id,
			p.product_name,
			c.product_id,
			c.container_name,
			c.docker_container_id,
			c.image,
			c.tag,
			c.host_port,
			c.container_port,
			c.environment,
			c.volumes,
			c.command,
			c.status,
			c.created_at,
			c.updated_at
		FROM containers c
		LEFT JOIN products p ON p.product_id = c.product_id
		ORDER BY c.created_at DESC
	`

	rst, err := r.DB.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	defer rst.Close()

	var results []Container
	for rst.Next() {
		var result Container
		var envJSON, volumesJSON, commandJSON []byte

		err := rst.Scan(
			&result.ContainerID,
			&result.ProductName,
			&result.ProductID,
			&result.ContainerName,
			&result.DockerContainerID,
			&result.Image,
			&result.Tag,
			&result.HostPort,
			&result.ContainerPort,
			&envJSON,
			&volumesJSON,
			&commandJSON,
			&result.Status,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Unmarshal JSON fields
		if len(envJSON) > 0 {
			if err := json.Unmarshal(envJSON, &result.Environment); err != nil {
				return nil, fmt.Errorf("failed to unmarshal environment: %w", err)
			}
		}

		if len(volumesJSON) > 0 {
			if err := json.Unmarshal(volumesJSON, &result.Volumes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal volumes: %w", err)
			}
		}

		if len(commandJSON) > 0 {
			if err := json.Unmarshal(commandJSON, &result.Command); err != nil {
				return nil, fmt.Errorf("failed to unmarshal command: %w", err)
			}
		}

		results = append(results, result)
	}

	if err := rst.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

func (r *Repository) GetContainerByID(c context.Context, id int) (*Container, error) {
	query := `SELECT container_id, product_id, container_name, docker_container_id, 
	          image, tag, host_port, container_port, environment, volumes, command, 
	          status, created_at, updated_at 
	          FROM containers WHERE container_id = $1`

	var cont Container
	var envJSON, volJSON, cmdJSON []byte
	var dockerID sql.NullString

	err := r.DB.QueryRowContext(c, query, id).Scan(
		&cont.ContainerID, &cont.ProductID, &cont.ContainerName, &dockerID,
		&cont.Image, &cont.Tag, &cont.HostPort, &cont.ContainerPort,
		&envJSON, &volJSON, &cmdJSON, &cont.Status, &cont.CreatedAt, &cont.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if dockerID.Valid {
		cont.DockerContainerID = dockerID.String
	}

	json.Unmarshal(envJSON, &cont.Environment)
	json.Unmarshal(volJSON, &cont.Volumes)
	json.Unmarshal(cmdJSON, &cont.Command)

	return &cont, nil
}

func (r *Repository) GetContainersByProductID(c context.Context, productID int) ([]Container, error) {
	query := `SELECT container_id, product_id, container_name, docker_container_id, 
	          image, tag, host_port, container_port, environment, volumes, command, 
	          status, created_at, updated_at 
	          FROM containers WHERE product_id = $1 ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(c, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []Container
	for rows.Next() {
		var cont Container
		var envJSON, volJSON, cmdJSON []byte
		var dockerID sql.NullString

		err := rows.Scan(
			&cont.ContainerID, &cont.ProductID, &cont.ContainerName, &dockerID,
			&cont.Image, &cont.Tag, &cont.HostPort, &cont.ContainerPort,
			&envJSON, &volJSON, &cmdJSON, &cont.Status, &cont.CreatedAt, &cont.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if dockerID.Valid {
			cont.DockerContainerID = dockerID.String
		}

		json.Unmarshal(envJSON, &cont.Environment)
		json.Unmarshal(volJSON, &cont.Volumes)
		json.Unmarshal(cmdJSON, &cont.Command)

		containers = append(containers, cont)
	}

	return containers, nil
}

func (r *Repository) UpdateContainerStatus(c context.Context, id int, status, dockerID string) error {
	query := `UPDATE containers 
	          SET status = $1, docker_container_id = $2, updated_at = $3 
	          WHERE container_id = $4`

	_, err := r.DB.ExecContext(c, query, status, dockerID, time.Now(), id)
	return err
}

func (r *Repository) UpdateContainer(c context.Context, id int, req Container) error {
	envJSON, _ := json.Marshal(req.Environment)
	volJSON, _ := json.Marshal(req.Volumes)
	cmdJSON, _ := json.Marshal(req.Command)

	query := `UPDATE containers 
	          SET container_name = $1, image = $2, tag = $3, host_port = $4, 
	              container_port = $5, environment = $6, volumes = $7, command = $8,
	              updated_at = $9 
	          WHERE container_id = $10`

	_, err := r.DB.ExecContext(c, query,
		req.ContainerName, req.Image, req.Tag, req.HostPort, req.ContainerPort,
		envJSON, volJSON, cmdJSON, time.Now(), id,
	)
	return err
}

func (r *Repository) DeleteContainer(c context.Context, id int) error {
	query := `DELETE FROM containers WHERE container_id = $1`
	_, err := r.DB.ExecContext(c, query, id)
	return err
}

func (r *Repository) GetProductWithContainers(c context.Context, id int) (*ProductWithContainers, error) {
	product, err := r.GetProductByID(c, id)
	if err != nil {
		return nil, err
	}

	containers, err := r.GetContainersByProductID(c, id)
	if err != nil {
		return nil, err
	}

	return &ProductWithContainers{
		Product:    *product,
		Containers: containers,
	}, nil
}
