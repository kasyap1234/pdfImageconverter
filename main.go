package main

import (
	"context"
	"log"
	"os"
	"practice/auth"
	"practice/core"
	"practice/db/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
)

func InitDB() (*db.Queries, *pgxpool.Pool) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	ctx := context.Background()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("unable to parse DATABASE_URL: %v", err)
	}
	poolConfig.MaxConns = 10
	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	q := db.New(dbpool)
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("error pinging db %v", err)
	}

	return q, dbpool
}

func main() {

	e := echo.New()
	q, _ := InitDB()
	e.POST("/register", auth.Register(q))
	e.POST("/login", auth.Login(q))
	protected := e.Group("/user")
	protected.Use(auth.JWTMiddleware)
	protected.POST("/shorten", core.ShortenURL(q))
	e.Start(":8080")

}
