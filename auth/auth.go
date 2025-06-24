package auth

import (
	"net/http"
	"practice/db/db"

	"github.com/labstack/echo"
)

type Req struct {
	Email     string `json:"email"`
	Passsword string `json:"password"`
}

// register user by creating user in db , generating jwt
func Register(q *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		hash, err := HashPassword(req.Passsword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		user, err := q.CreateUser(c.Request().Context(), db.CreateUserParams{
			Email:    req.Email,
			Password: hash,
		})
		if err != nil {
			return c.JSON(http.StatusConflict, "user already exists")
		}
		token, err := GenerateJWT(user.ID.String())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, echo.Map{"token": token})
	}
}

func Login(q *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		user, err := q.GetUserByEmail(c.Request().Context(), req.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}
		if err := CheckPasswordHash(user.Password, req.Passsword); err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}
		token, err := GenerateJWT(user.ID.String())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, echo.Map{"token": token})
	}

}
