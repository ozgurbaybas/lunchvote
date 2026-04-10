CREATE TABLE IF NOT EXISTS restaurants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    district TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS restaurant_meal_cards (
    restaurant_id TEXT NOT NULL REFERENCES restaurants (id) ON DELETE CASCADE,
    meal_card TEXT NOT NULL,
    PRIMARY KEY (restaurant_id, meal_card)
);