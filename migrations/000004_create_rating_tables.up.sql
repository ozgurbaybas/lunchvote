CREATE TABLE IF NOT EXISTS ratings (
    id TEXT PRIMARY KEY,
    restaurant_id TEXT NOT NULL REFERENCES restaurants (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    score INTEGER NOT NULL CHECK (
        score >= 1
        AND score <= 5
    ),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (restaurant_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_restaurant_id ON ratings (restaurant_id);