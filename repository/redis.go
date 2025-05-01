package repository

import (
	"context"
	"encoding/json"
	"session-app/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	StoreSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, token string) (*models.Session, error)
	DeleteSession(ctx context.Context, token string) error
	Close() error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(redisURL string) (RedisRepository, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisRepository{
		client: client,
	}, nil
}

func (r *redisRepository) StoreSession(ctx context.Context, session *models.Session) error {
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return nil // Don't store expired sessions
	}

	// Serialize the session to JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Store in Redis with TTL
	return r.client.Set(ctx, "session:"+session.Token, sessionData, ttl).Err()
}

func (r *redisRepository) GetSession(ctx context.Context, token string) (*models.Session, error) {
	// Get session from Redis
	data, err := r.client.Get(ctx, "session:"+token).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Session not found
		}
		return nil, err
	}

	// Deserialize the session from JSON
	var session models.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *redisRepository) DeleteSession(ctx context.Context, token string) error {
	return r.client.Del(ctx, "session:"+token).Err()
}

func (r *redisRepository) Close() error {
	return r.client.Close()
}
