-- Yumsday database schema

-- Arbitrary limit of 1000 characters to ensure no excessively long names are added to the DB
-- that could cause performance issues or be used for malicious purposes.
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS unit_systems (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL
);

-- Populate unit_systems table
INSERT INTO unit_systems (id, name) VALUES
(1, 'METRIC'),
(2, 'US'),
(3, 'UK');

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    app_admin BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    avatar VARCHAR NOT NULL,
    language VARCHAR NOT NULL,
    app_theme VARCHAR NOT NULL,
    last_visited_group_id INTEGER,
    unit_system_id INTEGER NOT NULL,
    FOREIGN KEY (last_visited_group_id) REFERENCES groups(id)
    FOREIGN KEY (unit_system_id) REFERENCES unit_systems(id)
);

INSERT INTO users (username, password, app_admin, created_at, language, app_theme, unit_system_id) VALUES (
    "admin",
    "$2a$12$L4zK2tkbTZFR37/jFJvbgObzhyqoogNuLaLUatMfGH3QGRKBnLrNS",
    true,
    (unixepoch()),
    "/static/assets/avatar1.jpg",
    "EN",
    "SYSTEM",
    1
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
    id VARCHAR PRIMARY KEY NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL,
    ip_address VARCHAR,
    user_agent VARCHAR,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_session_id ON sessions (id);

CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    factor FLOAT NOT NULL,
    unit_type VARCHAR NOT NULL,
    abbreviation VARCHAR NOT NULL
);

-- Populate units table
INSERT INTO units (id, name, factor, unit_type, abbreviation) VALUES
-- WEIGHT
(1, 'milligram', 0.001, 'WEIGHT', 'mg'),
(2, 'gram', 1, 'WEIGHT', 'g'),
(3, 'kilogram', 1000, 'WEIGHT', 'kg'),
(4, 'ounce', 28.3495, 'WEIGHT', 'oz'),
(5, 'pound', 453.592, 'WEIGHT', 'lb'),

-- VOLUME
(6, 'milliliter', 1, 'VOLUME', 'mL'),
(7, 'centiliter', 10, 'VOLUME', 'cL'),
(8, 'deciliter', 100, 'VOLUME', 'dL'),
(9, 'liter', 1000, 'VOLUME', 'L'),
(10, 'teaspoon', 4.92892, 'VOLUME', 'tsp'),
(11, 'tablespoon', 14.7868, 'VOLUME', 'tbsp'),
(12, 'fluid ounce', 29.5735, 'VOLUME', 'fl oz'),
(13, 'cup', 236.588, 'VOLUME', 'cup'),
(14, 'pint', 473.176, 'VOLUME', 'pt'),
(15, 'quart', 946.353, 'VOLUME', 'qt'),
(16, 'imperial fluid ounce', 28.4131, 'VOLUME', 'fl oz'),
(17, 'imperial cup', 284.131, 'VOLUME', 'cup'),
(18, 'imperial pint', 568.261, 'VOLUME', 'pt'),
(19, 'imperial quart', 1136.52, 'VOLUME', 'qt'),

-- NUMERIC
(20, 'unit', 1, 'NUMERIC', 'u'),
(21, 'dozen', 12, 'NUMERIC', 'dz'),

-- PIECE
(22, 'piece', 1, 'PIECE', 'pc'),
(23, 'clove', 1, 'PIECE', 'clove'),
(24, 'slice', 1, 'PIECE', 'slice'),
(25, 'bag', 1, 'PIECE', 'bag'),
(26, 'sachet', 1, 'PIECE', 'sachet'),
(27, 'can', 1, 'PIECE', 'can'),
(28, 'jar', 1, 'PIECE', 'jar');

CREATE TABLE IF NOT EXISTS units_systems_junction (
    unit_id INTEGER NOT NULL,
    system_id INTEGER NOT NULL,
    PRIMARY KEY (unit_id, system_id),
    FOREIGN KEY (unit_id) REFERENCES units(id),
    FOREIGN KEY (system_id) REFERENCES unit_systems(id)
);

INSERT INTO units_systems_junction (unit_id, system_id) VALUES
-- Populate the units_systems_junction joint table
-- WEIGHT
(1, 1),   -- milligram -> METRIC
(2, 1),   -- gram -> METRIC
(3, 1),   -- kilogram -> METRIC
(4, 2),   -- ounce -> US
(4, 3),   -- ounce -> UK
(5, 2),   -- pound -> US
(5, 3),   -- pound -> UK

-- VOLUME (metric)
(6, 1),   -- milliliter -> METRIC
(7, 1),   -- centiliter -> METRIC
(8, 1),   -- deciliter -> METRIC
(9, 1),   -- liter -> METRIC

-- VOLUME (US)
(10, 2),  -- teaspoon -> US
(11, 2),  -- tablespoon -> US
(12, 2),  -- fluid ounce -> US
(13, 2),  -- cup -> US
(14, 2),  -- pint -> US
(15, 2),  -- quart -> US

-- VOLUME (UK/imperial)
(16, 3),  -- imperial fluid ounce -> UK
(17, 3),  -- imperial cup -> UK
(18, 3),  -- imperial pint -> UK
(19, 3),  -- imperial quart -> UK

-- NUMERIC (universal)
(20, 1), (20, 2), (20, 3),  -- unit -> METRIC, US, UK
(21, 1), (21, 2), (21, 3),  -- dozen -> METRIC, US, UK

-- PIECE (universal)
(22, 1), (22, 2), (22, 3),  -- piece -> METRIC, US, UK
(23, 1), (23, 2), (23, 3),  -- clove -> METRIC, US, UK
(24, 1), (24, 2), (24, 3),  -- slice -> METRIC, US, UK
(25, 1), (25, 2), (25, 3),  -- bag -> METRIC, US, UK
(26, 1), (26, 2), (26, 3),  -- sachet -> METRIC, US, UK
(27, 1), (27, 2), (27, 3),  -- can -> METRIC, US, UK
(28, 1), (28, 2), (28, 3);  -- jar -> METRIC, US, UK

CREATE TABLE IF NOT EXISTS item_categories (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    average_market_price FLOAT,
    unit_type VARCHAR NOT NULL,
    item_category_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (item_category_id) REFERENCES item_categories(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS recipes (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    image_url VARCHAR,
    original_link VARCHAR,
    preparation_time_min INTEGER,
    cooking_time_min INTEGER,
    servings INTEGER,
    instructions VARCHAR,
    created_at TIMESTAMP NOT NULL,
    public BOOLEAN NOT NULL,
    comment VARCHAR,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS recipe_categories (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
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
    unit_id INTEGER,
    recipe_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (unit_id) REFERENCES units(id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id)
);

CREATE TABLE IF NOT EXISTS dishes (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    portion INTEGER NOT NULL,
    bought BOOLEAN NOT NULL,
    datetime TIMESTAMP NOT NULL,
    group_id INTEGER NOT NULL,
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
    unit_id INTEGER,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (unit_id) REFERENCES units(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
