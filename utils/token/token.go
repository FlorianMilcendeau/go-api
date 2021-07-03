package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(userId uint) (string, error) {
	tokenHourLifespan, err := strconv.Atoi(os.Getenv(("TOKEN_HOUR_LIFESPAN")))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenHourLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ExtractTokenId(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userId"]), 10, 12)

		if err != nil {
			return 0, err
		}

		return uint(userId), nil
	}

	return 0, nil
}
