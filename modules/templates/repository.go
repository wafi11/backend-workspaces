package templates

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/wafi11/backend-workspaces/pkg/types"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) TemplatesRepository {
	return &Repository{DB: DB}
}
func (r *Repository) Create(c context.Context, req CreateTemplateRequest) error {
	envVarsJSON, err := json.Marshal(req.EnvVarsSchema)
	if err != nil {
		return fmt.Errorf("failed to marshal env vars schema: %w", err)
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	features := req.Features
	if features == nil {
		features = []string{}
	}

	// Execute query
	var id int64
	var createdAt, updatedAt sql.NullTime

	err = r.DB.QueryRowContext(
		c,
		QueryCreate,
		req.Name,
		req.DisplayName,
		req.Description,
		req.Category,
		req.GitRepoURL,
		req.GitBranch,
		req.HelmChartPath,
		req.DockerfilePath,
		req.DefaultCPURequest,
		req.DefaultCPULimit,
		req.DefaultMemoryRequest,
		req.DefaultMemoryLimit,
		req.DefaultReplicas,
		req.RequiresDatabase,
		req.DefaultDatabaseType,
		req.RequiresRedis,
		req.RequiresRabbitMQ,
		req.DefaultPort,
		envVarsJSON,
		pq.Array(tags),
		pq.Array(features),
		req.IconURL,
		true,
		false,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error) {

	// Decode cursor
	var cursor *Cursor
	var err error
	if req.Cursor != nil && *req.Cursor != "" {
		cursor, err = DecodeCursor(*req.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
	}

	// Prepare query parameters
	var createdAt interface{}
	var cursorID int64

	if cursor != nil {
		createdAt = cursor.CreatedAt
		cursorID = cursor.ID
	}

	rows, err := r.DB.QueryContext(ctx, queryListWithCursor, createdAt, cursorID, req.Limit+1)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	templates := []ListType{}
	for rows.Next() {
		var tmpl ListType

		err := rows.Scan(
			&tmpl.Id,
			&tmpl.Name,
			&tmpl.DisplayName,
			&tmpl.Description,
			&tmpl.Category,
			&tmpl.Version,
			&tmpl.GitRepoURL,
			&tmpl.GitBranch,
			&tmpl.DefaultPort,
			&tmpl.IconURL,
			&tmpl.IsActive,
			&tmpl.IsFeatured,
			&tmpl.CreatedAt,
			&tmpl.UpdatedAt,
			&tmpl.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		templates = append(templates, tmpl)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	// Check if there's more data
	hasMore := len(templates) > req.Limit
	if hasMore {
		templates = templates[:req.Limit]
	}

	// Generate next cursor
	var nextCursor *string
	if hasMore && len(templates) > 0 {
		lastTemplate := templates[len(templates)-1]
		cursor := Cursor{
			CreatedAt: lastTemplate.CreatedAt,
			ID:        int64(lastTemplate.Id),
		}
		encoded := cursor.Encode()
		nextCursor = &encoded
	}

	return &ListTemplatesResponse{
		Data: templates,
		Pagination: types.PaginationCursor{
			NextCursor: nextCursor,
			HasMore:    hasMore,
			Count:      len(templates),
		},
	}, nil
}
func (r *Repository) FindById(c context.Context, ID int) (*Template, error) {
	var data Template

	// Variabel temporary untuk scan JSON fields
	var envVarsJSON []byte
	var tagsArray, featuresArray []string
	var screenshotURLsArray []string

	err := r.DB.QueryRowContext(c, queryGetByID, ID).Scan(
		&data.Id,
		&data.Name,
		&data.DisplayName,
		&data.Description,
		&data.Category,
		&data.Version,
		&data.GitRepoURL,
		&data.GitBranch,
		&data.HelmChartPath,
		&data.DockerfilePath,
		&data.DefaultCPURequest,
		&data.DefaultCPULimit,
		&data.DefaultMemoryRequest,
		&data.DefaultMemoryLimit,
		&data.DefaultReplicas,
		&data.RequiresDatabase,
		&data.DefaultDatabaseType,
		&data.RequiresRedis,
		&data.RequiresRabbitMQ,
		&data.DefaultPort,
		&envVarsJSON, // <- Scan ke []byte dulu
		pq.Array(&tagsArray),
		pq.Array(&featuresArray),
		&data.IconURL,
		pq.Array(&screenshotURLsArray), // <- Scan ke []byte dulu
		&data.IsActive,
		&data.IsFeatured,
		&data.CreatedAt,
		&data.UpdatedAt,
	)

	if err != nil {
		log.Printf("failed to get template : %s", err.Error())
		return nil, fmt.Errorf("failed to get templates : %w", err)
	}

	// Unmarshal semua JSON fields
	if len(envVarsJSON) > 0 {
		if err := json.Unmarshal(envVarsJSON, &data.EnvVarsSchema); err != nil {
			return nil, fmt.Errorf("failed to unmarshal env_vars_schema: %w", err)
		}
	}

	return &data, nil
}
