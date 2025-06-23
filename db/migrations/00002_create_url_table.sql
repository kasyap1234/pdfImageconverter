-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE urls(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    original_url TEXT NOT NULL,
    short_code TEXT UNIQUE NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE urls;

-- +goose StatementEnd
