package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/travellog/backend/config"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email required and password must be at least 8 chars"})
			return
		}

		db := database.GetDB()

		// check existing
		var existing models.User
		if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		u := models.User{
			Email:        req.Email,
			PasswordHash: string(hash),
		}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
	}
}

func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
			return
		}

		db := database.GetDB()

		var u models.User
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"sub": u.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := t.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  gin.H{"id": u.ID, "email": u.Email},
		})
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	}
}
