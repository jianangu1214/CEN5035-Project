// seedtestuser inserts a predictable account for local / E2E testing.
// The app has no admin role; this is a normal user used as a "test admin" fixture.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/travellog/backend/config"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("database: %v", err)
	}

	email := strings.ToLower(strings.TrimSpace(os.Getenv("SEED_EMAIL")))
	if email == "" {
		email = "admin@test.local"
	}
	password := os.Getenv("SEED_PASSWORD")
	if password == "" {
		password = "AdminTest123!"
	}
	if len(password) < 8 {
		log.Fatal("SEED_PASSWORD must be at least 8 characters")
	}

	db := database.GetDB()
	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		log.Printf("seed: user already exists (%s)", email)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash: %v", err)
	}
	u := models.User{Email: email, PasswordHash: string(hash)}
	if err := db.Create(&u).Error; err != nil {
		log.Fatalf("create user: %v", err)
	}
	log.Printf("seed: created test user %s (password from SEED_PASSWORD or default)", email)
}
