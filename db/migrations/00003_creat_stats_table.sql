-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE stats (
    id UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
    ip_address TEXT ,
    user_agent TEXT,
    clicked_at TIMESTAMP DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE stats; 
-- +goose StatementEnd
