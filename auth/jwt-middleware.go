package auth

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return echo.ErrUnauthorized
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return echo.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.ErrUnauthorized
		}

		userID := claims["user_id"].(string)
		c.Set("user_id", userID)
		return next(c)
	}
}
