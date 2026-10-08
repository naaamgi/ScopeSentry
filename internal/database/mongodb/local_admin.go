package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// SeedLocalAdmin adds a convenient account only for explicitly opted-in local
// installations. An existing admin account, including its password, is untouched.
func SeedLocalAdmin(ctx context.Context) error {
	if os.Getenv("SCOPE_SENTRY_LOCAL_DEFAULT_ADMIN") != "true" {
		return nil
	}

	users := DB.Collection("user")
	err := users.FindOne(ctx, bson.M{"username": "admin"}).Err()
	if err == nil {
		return nil
	}
	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("check local administrator: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash local administrator password: %w", err)
	}
	now := time.Now()
	_, err = users.InsertOne(ctx, models.User{
		ID:        primitive.NewObjectID(),
		Username:  "admin",
		Password:  string(hash),
		Role:      "admin",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("create local administrator: %w", err)
	}
	return nil
}
