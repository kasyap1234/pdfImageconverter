-- name: CreateUser :one
INSERT INTO users(id,email,password,created_at) VALUES ($1,$2,$3,$4) RETURNING *; 
-- name: GetURL :one 
SELECT * FROM urls WHERE  id=$1 LIMIT 1; 

-- name: GetURLByUserID :one
SELECT * FROM urls WHERE user_id=$1 LIMIT 10; 

-- name: CreateURL :one 
INSERT INTO urls(id,user_id,original_url,short_code) VALUES($1,$2,$3,$4) RETURNING *; 