-- Yumsday database schema

-- Arbitrary limit of 1000 characters to ensure no excessively long names are added to the DB
-- that could cause performance issues or be used for malicious purposes.
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    username VARCHAR(1000) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    app_admin BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    avatar VARCHAR,
    language VARCHAR NOT NULL,
    app_theme VARCHAR NOT NULL,
    last_visited_group_id INTEGER,
    FOREIGN KEY (last_visited_group_id) REFERENCES groups(id)
);

INSERT INTO users (username, password, app_admin, created_at, language, app_theme) VALUES (
    "admin",
    "$2a$12$L4zK2tkbTZFR37/jFJvbgObzhyqoogNuLaLUatMfGH3QGRKBnLrNS",
    true,
    (unixepoch()),
    "en",
    "system"
);

-- Many-to-Many relationship between users and groups
CREATE TABLE IF NOT EXISTS group_members (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    admin BOOLEAN NOT NULL,
    joined_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL,
    ip_address VARCHAR,
    user_agent VARCHAR,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_session_id ON sessions (id);

CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    factor FLOAT NOT NULL,
    unit_type VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS item_categories (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    description VARCHAR(1000000),
    average_market_price FLOAT,
    unit_type VARCHAR NOT NULL,
    group_id INTEGER NOT NULL,
    item_category_id INTEGER NOT NULL,
    FOREIGN KEY (item_category_id) REFERENCES item_categories(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS recipes (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    description VARCHAR(1000000),
    image_url VARCHAR,
    original_link VARCHAR,
    preparation_time_min INTEGER,
    cooking_time_min INTEGER,
    servings INTEGER,
    instructions VARCHAR(1000000),
    created_at TIMESTAMP NOT NULL,
    public BOOLEAN NOT NULL,
    comment VARCHAR(1000000),
    group_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS recipe_categories (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Many-to-Many relationship between recipes and recipe_categories
CREATE TABLE IF NOT EXISTS recipes_categories_junction (
    recipe_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, category_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id),
    FOREIGN KEY (category_id) REFERENCES recipe_categories(id)
);

CREATE TABLE IF NOT EXISTS ingredients (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    quantity FLOAT,
    item_id INTEGER NOT NULL,
    unit_id INTEGER NOT NULL,
    recipe_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (unit_id) REFERENCES units(id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id)
);

CREATE TABLE IF NOT EXISTS dishes (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    portion INTEGER NOT NULL,
    bought BOOLEAN DEFAULT FALSE NOT NULL,
    datetime TIMESTAMP NOT NULL,
    recipe_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Many-to-Many relationship between recipes and dishes
CREATE TABLE IF NOT EXISTS recipes_dishes_junction (
    recipe_id INTEGER NOT NULL,
    dish_id INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, dish_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id),
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

CREATE TABLE IF NOT EXISTS groceries (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    quantity_bought FLOAT NOT NULL,
    user_quantity FLOAT NOT NULL,
    item_id INTEGER NOT NULL,
    unit_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (unit_id) REFERENCES units(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
