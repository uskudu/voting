package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"voting/internal/db"
	"voting/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context) {
	// get cookie off request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// validate data
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header)
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || token == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// check exp
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// find user with token sub
	var usr user.User
	if err = db.DB.First(&usr, "id = ?", claims["sub"]).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if _, err = uuid.Parse(usr.ID); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// attach to req
	c.Set("user", usr)
	c.Next()
}
