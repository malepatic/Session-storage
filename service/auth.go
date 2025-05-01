package service

import (
	"context"
	"errors"
	"session-app/internal/models"
	"session-app/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionExpired     = errors.New("session expired")
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	pgRepo          repository.PostgresRepository
	redisRepo       repository.RedisRepository
	tokenExpiration time.Duration
	jwtSecret       string
}

func NewAuthService(pgRepo repository.PostgresRepository, redisRepo repository.RedisRepository, tokenExpiration time.Duration) AuthService {
	return &authService{
		pgRepo:          pgRepo,
		redisRepo:       redisRepo,
		tokenExpiration: tokenExpiration,
		jwtSecret:       "your-jwt-secret", // In production, get this from config
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	return s.pgRepo.CreateUser(ctx, req)
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user from database
	user, err := s.pgRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	expiresAt := time.Now().Add(s.tokenExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expiresAt.Unix(),
		"jti": uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Store session in Redis
	session := &models.Session{
		Token:     tokenString,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	if err := s.redisRepo.StoreSession(ctx, session); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:   tokenString,
		Expires: expiresAt.Format(time.RFC3339),
		User:    *user,
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Get session from Redis first (faster and has built-in revocation)
	session, err := s.redisRepo.GetSession(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionExpired
	}

	// Validate JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrSessionExpired
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrSessionExpired
	}

	// Extract user ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrSessionExpired
	}

	// Get user from database
	user, err := s.pgRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.redisRepo.DeleteSession(ctx, token)
}
