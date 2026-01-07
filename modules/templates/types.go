package templates

import (
	"encoding/json"
	"time"

	"github.com/wafi11/backend-workspaces/pkg/types"
)

type Template struct {
	Id          int    `json:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"displayName" db:"display_name"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
	Version     string `json:"version" db:"version"`

	// Git Repository
	GitRepoURL string `json:"gitRepoUrl" db:"git_repo_url"`
	GitBranch  string `json:"gitBranch" db:"git_branch"`

	// Kubernetes Deployment
	HelmChartPath  *string `json:"helmChartPath" db:"helm_chart_path"`
	DockerfilePath string  `json:"dockerfilePath" db:"dockerfile_path"`

	// Resource Requirements
	DefaultCPURequest    string `json:"defaultCpuRequest" db:"default_cpu_request"`
	DefaultCPULimit      string `json:"defaultCpuLimit" db:"default_cpu_limit"`
	DefaultMemoryRequest string `json:"defaultMemoryRequest" db:"default_memory_request"`
	DefaultMemoryLimit   string `json:"defaultMemoryLimit" db:"default_memory_limit"`
	DefaultReplicas      int    `json:"defaultReplicas" db:"default_replicas"`

	// Database Requirements
	RequiresDatabase    bool    `json:"requiresDatabase" db:"requires_database"`
	DefaultDatabaseType *string `json:"defaultDatabaseType" db:"default_database_type"`

	// Additional Services
	RequiresRedis    bool `json:"requiresRedis" db:"requires_redis"`
	RequiresRabbitMQ bool `json:"requiresRabbitmq" db:"requires_rabbitmq"`

	// Port Configuration
	DefaultPort int `json:"defaultPort" db:"default_port"`

	// Environment Variables Schema
	EnvVarsSchema EnvVarsSchema `json:"envVarsSchema" db:"env_vars_schema"`

	// Tags & Features
	Tags     []string `json:"tags" db:"tags"`
	Features []string `json:"features" db:"features"`

	// UI Display
	IconURL        *string  `json:"iconUrl" db:"icon_url"`
	ScreenshotURLs []string `json:"screenshotUrls" db:"screenshot_urls"`

	// Status
	IsActive   bool `json:"isActive" db:"is_active"`
	IsFeatured bool `json:"isFeatured" db:"is_featured"`

	// Audit
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

type ListType struct {
	Id          int    `json:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"displayName" db:"display_name"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
	Version     string `json:"version" db:"version"`
	GitRepoURL  string `json:"gitRepoUrl" db:"git_repo_url"`
	GitBranch   string `json:"gitBranch" db:"git_branch"`
	DefaultPort int    `json:"defaultPort" db:"default_port"`
	// Status
	IsActive   bool `json:"isActive" db:"is_active"`
	IsFeatured bool `json:"isFeatured" db:"is_featured"`
	// UI Display
	IconURL *string `json:"iconUrl" db:"icon_url"`

	// Audit
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

type EnvVarsSchema map[string]EnvVarProperty

type EnvVarProperty struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Required    bool            `json:"required"`
	Secret      bool            `json:"secret"`
	Default     json.RawMessage `json:"default,omitempty"`
}

type CreateTemplateRequest struct {
	Name                 string        `json:"name" validate:"required,max=100"`
	DisplayName          string        `json:"displayName" validate:"required,max=150"`
	Description          string        `json:"description"`
	Category             string        `json:"category" validate:"required,max=50"`
	GitRepoURL           string        `json:"gitRepoUrl" validate:"required,url"`
	GitBranch            string        `json:"gitBranch"`
	HelmChartPath        *string       `json:"helmChartPath"`
	DockerfilePath       string        `json:"dockerfilePath"`
	DefaultCPURequest    string        `json:"defaultCpuRequest"`
	DefaultCPULimit      string        `json:"defaultCpuLimit"`
	DefaultMemoryRequest string        `json:"defaultMemoryRequest"`
	DefaultMemoryLimit   string        `json:"defaultMemoryLimit"`
	DefaultReplicas      int           `json:"defaultReplicas"`
	RequiresDatabase     bool          `json:"requiresDatabase"`
	DefaultDatabaseType  *string       `json:"defaultDatabaseType"`
	RequiresRedis        bool          `json:"requiresRedis"`
	RequiresRabbitMQ     bool          `json:"requiresRabbitmq"`
	DefaultPort          int           `json:"defaultPort"`
	EnvVarsSchema        EnvVarsSchema `json:"envVarsSchema"`
	Tags                 []string      `json:"tags"`
	Features             []string      `json:"features"`
	IconURL              *string       `json:"iconUrl"`
}
type ListTemplatesRequest struct {
	Limit  int     `query:"limit"`
	Cursor *string `query:"cursor"` // base64 encoded cursor
}
type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

type ListTemplatesResponse struct {
	Data       []ListType             `json:"data"`
	Pagination types.PaginationCursor `json:"pagination"`
}
