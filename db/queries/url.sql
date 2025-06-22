-- name: CreateShortURL :one
INSERT INTO url(original_url,short_url) VALUES($1,$2) RETURNING * ; 

-- name: GetURLByShortURL :one 
SELECT * FROM url WHERE short_url = $1; 
