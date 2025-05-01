package repository

import (
	"context"
	"database/sql"
	"errors"
	"session-app/internal/models"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresRepository interface {
	CreateUser(ctx context.Context, user *models.RegisterRequest) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	Close() error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(postgresURL string) (PostgresRepository, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Ensure tables exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &postgresRepository{
		db: db,
	}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	_, err := db.Exec(query)
	return err
}

func (r *postgresRepository) CreateUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2",
		req.Username, req.Email).Scan(&count)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("username or email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = r.db.ExecContext(ctx,
		"INSERT INTO users (id, username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1",
		id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}
