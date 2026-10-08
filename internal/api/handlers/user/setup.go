package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// The fixed ID makes competing first-user requests atomic without requiring
// MongoDB transactions. Existing installations already contain a user.
var firstAdministratorID = primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

type setupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SetupStatus(c *gin.Context) {
	count, err := mongodb.DB.Collection("user").CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "계정 상태를 확인할 수 없습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"required": count == 0}})
}

func Setup(c *gin.Context) {
	var req setupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "계정 이름과 비밀번호를 입력하세요"})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 || len(req.Username) > 64 || len(req.Password) < 12 || len(req.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "계정 이름은 3~64자, 비밀번호는 12~72자로 입력하세요"})
		return
	}
	collection := mongodb.DB.Collection("user")
	count, err := collection.CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "계정 상태를 확인할 수 없습니다"})
		return
	}
	if count != 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "관리자 계정이 이미 생성되었습니다"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "비밀번호를 저장할 수 없습니다"})
		return
	}
	now := time.Now()
	_, err = collection.InsertOne(c, models.User{
		ID: firstAdministratorID, Username: req.Username, Password: string(hash),
		Role: "admin", Status: "active", CreatedAt: now, UpdatedAt: now,
	})
	if mongo.IsDuplicateKeyError(err) {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "관리자 계정이 이미 생성되었습니다"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "관리자 계정을 만들 수 없습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "관리자 계정이 생성되었습니다"})
}
