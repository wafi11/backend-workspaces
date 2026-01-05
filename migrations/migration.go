package migrations

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS products (
			product_id SERIAL PRIMARY KEY,
			product_name VARCHAR(100) NOT NULL,
			metadata JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_product_name ON products(product_name)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	fmt.Println("✅ Migrations completed successfully")
	return nil
}
