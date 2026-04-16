CREATE TABLE IF NOT EXISTS polls (
    id TEXT PRIMARY KEY,
    team_id UUID NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS poll_options (
    poll_id TEXT NOT NULL REFERENCES polls (id) ON DELETE CASCADE,
    restaurant_id TEXT NOT NULL REFERENCES restaurants (id) ON DELETE CASCADE,
    PRIMARY KEY (poll_id, restaurant_id)
);

CREATE TABLE IF NOT EXISTS poll_votes (
    poll_id TEXT NOT NULL REFERENCES polls (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    restaurant_id TEXT NOT NULL REFERENCES restaurants (id) ON DELETE CASCADE,
    voted_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (poll_id, user_id)
);