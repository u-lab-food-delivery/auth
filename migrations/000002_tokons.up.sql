CREATE TABLE tokens(
    user_id UUID REFERENCES users(id),
    access_token TEXT,
    refresh_token TEXT,
    is_revoked BOOLEAN 'FALSE',
    expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)