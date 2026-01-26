-- Yumsday database schema
-- user_group
CREATE TABLE IF NOT EXISTS user_group (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- USER
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    app_admin BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    avatar VARCHAR,
    language VARCHAR DEFAULT 'en' NOT NULL,
    app_theme VARCHAR DEFAULT 'light' NOT NULL,
    last_visited_group INTEGER,
    FOREIGN KEY (last_visited_group) REFERENCES user_group(id)
);

-- SESSION (manage cookies sessions)
CREATE TABLE IF NOT EXISTS session (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    session_id VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address VARCHAR NOT NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_agent VARCHAR,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- MEMBER (Many-to-Many relationship between USER and USER_GROUP)
-- Stores the admin flag to, for admins of the group
CREATE TABLE IF NOT EXISTS member_group (
    user_id INTEGER NOT NULL,
    user_group_id INTEGER NOT NULL,
    admin BOOLEAN DEFAULT FALSE NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, user_group_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- UNIT
CREATE TABLE IF NOT EXISTS unit (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    factor FLOAT NOT NULL,
    unit_type VARCHAR NOT NULL
);

-- ITEM
CREATE TABLE IF NOT EXISTS item (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    average_market_price FLOAT,
    unit_type VARCHAR NOT NULL,
    category VARCHAR
);

-- RECIPE
CREATE TABLE IF NOT EXISTS recipe (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    recipe_category VARCHAR,
    image_url VARCHAR,
    original_link VARCHAR,
    preparation_time INTEGER,
    cooking_time INTEGER,
    servings INTEGER,
    instructions VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    public BOOLEAN DEFAULT FALSE NOT NULL,
    comment VARCHAR,
    user_group_id INTEGER NOT NULL,
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- INGREDIENT
CREATE TABLE IF NOT EXISTS ingredient (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    quantity FLOAT,
    item_id INTEGER NOT NULL,
    unit_id INTEGER,
    reciped_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES item(id),
    FOREIGN KEY (unit_id) REFERENCES unit(id),
    FOREIGN KEY (reciped_id) REFERENCES recipe(id)
);

-- DISH
CREATE TABLE IF NOT EXISTS dish (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    portion INTEGER NOT NULL,
    bought BOOLEAN DEFAULT FALSE NOT NULL,
    datetime TIMESTAMP NOT NULL,
    recipe_id INTEGER NOT NULL,
    user_group_id INTEGER NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipe(id),
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- GROCERY LIST
CREATE TABLE IF NOT EXISTS grocery_list (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    quantity FLOAT,
    quantity_bought FLOAT,
    user_quantity FLOAT,
    item_id INTEGER NOT NULL,
    unit_id INTEGER,
    user_group_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES item(id),
    FOREIGN KEY (unit_id) REFERENCES unit(id),
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);