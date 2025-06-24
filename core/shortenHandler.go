package core

import (
	"log"
	"net/http"
	"practice/db/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo"
)

func ToPGUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}
func ShortenURL(q *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDstr := c.Get("user_id").(string)
		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid user id")
		}
		var req struct {
			OriginalURL string `json:"original_url"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		code, err := GenerateShortCode(6)
		if err != nil {
			log.Fatal("failed to generate short code")
		}
		url, err := q.CreateURL(c.Request().Context(), db.CreateURLParams{
			UserID:      ToPGUUID(userID),
			OriginalUrl: req.OriginalURL,
			ShortCode:   code,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, url)
	}
}
