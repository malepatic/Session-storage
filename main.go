package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"context"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctx = context.Background()
	db  *gorm.DB
	rdb *redis.Client
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	token := fmt.Sprintf("token-%s-%d", username, time.Now().Unix())

	err := rdb.Set(ctx, token, username, time.Hour).Err()
	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Session Token: %s", token)
}

func main() {
	// Connect to PostgreSQL
	dsn := os.Getenv("DB_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	db.AutoMigrate(&User{})

	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	http.HandleFunc("/login", loginHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
