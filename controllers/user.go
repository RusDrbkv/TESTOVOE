package controllers

import (
	"TESTOVOE/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func FindUser(username string) (*models.User, error) {
	var user *models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username}
	user.Password = HashValue(input.Password)

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user.Username})
}

func LoginUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := FindUser(input.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	fmt.Println(input.Password)
	var passwordCheckError = CompareToHash(input.Password, user.Password)
	if passwordCheckError {
		token, err := GenerateToken(user.Username, input.Password)
		fmt.Println(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Token error"})
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})

	}

}

func GenerateToken(username, password string) (string, error) {
	user, err := FindUser(username)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
func HashValue(v string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func CompareToHash(v, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(v))
	fmt.Println(v, hash, err)

	return err == nil
}
