package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wafi11/backend-workspaces/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB  *sql.DB
	cfg config.Config
}

func NewRepository(db *sql.DB, cfg config.Config) UserRepository {
	return &Repository{DB: db, cfg: cfg}
}

func (repo *Repository) Create(c context.Context, req RegisterUser) error {
	// tx, err := repo.DB.BeginTx(c, nil)
	// if err != nil {
	// 	return fmt.Errorf("%s: %w", ErrTransactionFailed, err)
	// }
	// defer tx.Rollback()

	query := `
        INSERT INTO users (username, email, phone_number, password, password_hash) 
        VALUES ($1, $2, $3, $4, $5)
    `

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrPasswordHashFailed, err)
	}

	_, err = repo.DB.ExecContext(c, query, req.Username, req.Email, req.PhoneNumber, req.Password, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "idx_users_username") {
			return fmt.Errorf("%s", ErrUsernameExists)
		}
		if strings.Contains(err.Error(), "idx_users_email") {
			return fmt.Errorf("%s", ErrEmailExists)
		}
		if strings.Contains(err.Error(), "idx_users_phone_number") {
			return fmt.Errorf("%s", ErrPhoneNumberExists)
		}
		return fmt.Errorf("%s: %w", ErrQueryFailed, err)
	}

	// if err := tx.Commit(); err != nil {
	// 	return fmt.Errorf("%s: %w", ErrTransactionFailed, err)
	// }

	return nil
}

func (repo *Repository) Login(c context.Context, req LoginUser) (*LoginResponse, error) {
	var password_hash, password string
	var id int
	query := `
		select id,password_hash,password from users where username = $1

	`
	err := repo.DB.QueryRowContext(c, query, req.Username).Scan(&id, &password_hash, &password)
	if err != nil {
		return nil, fmt.Errorf("%s", ErrInvalidCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("%s", ErrInvalidCredentials)
	}

	accessToken, _, err := GenerateTokenPair(id, repo.cfg)
	if err != nil {
		log.Printf("%s", err.Error())
		return nil, fmt.Errorf("%s", ErrInvalidCredentials)
	}

	return &LoginResponse{
		Token: accessToken,
	}, nil
}
