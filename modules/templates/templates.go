package templates

import (
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type TemplatesRepository interface {
	Create(c context.Context, req CreateTemplateRequest) error
	List(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error)
	FindById(c context.Context, req int) (*Template, error)
}

func (e *EnvVarsSchema) Scan(value interface{}) error {
	if value == nil {
		*e = make(EnvVarsSchema)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, e)
}

// Implement driver.Valuer for EnvVarsSchema (write to DB)
func (e EnvVarsSchema) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

func (c *Cursor) Encode() string {
	jsonBytes, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

// Decode cursor from base64
func DecodeCursor(encodedCursor string) (*Cursor, error) {
	if encodedCursor == "" {
		return nil, nil
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	var cursor Cursor
	if err := json.Unmarshal(jsonBytes, &cursor); err != nil {
		return nil, fmt.Errorf("invalid cursor data: %w", err)
	}

	return &cursor, nil
}
