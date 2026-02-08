-- Yumsday database schema

-- user_group
-- Arbitrary limit of 1000 characters to ensure no excessively long names are added to the DB
-- that could cause performance issues or be used for malicious purposes.
CREATE TABLE IF NOT EXISTS user_group (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP NOT NULL
);

-- USER
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    username VARCHAR(1000) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    app_admin BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    avatar VARCHAR,
    language VARCHAR NOT NULL,
    app_theme VARCHAR NOT NULL,
    last_visited_group INTEGER,
    FOREIGN KEY (last_visited_group) REFERENCES user_group(id)
);

-- MEMBER (Many-to-Many relationship between USER and USER_GROUP)
-- Stores the admin flag for admins of the group
CREATE TABLE IF NOT EXISTS member_group (
    user_id INTEGER NOT NULL,
    user_group_id INTEGER NOT NULL,
    admin BOOLEAN NOT NULL,
    joined_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, user_group_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- SESSION (manage cookies sessions)
CREATE TABLE IF NOT EXISTS session (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    session_id VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL,
    ip_address VARCHAR NOT NULL,
    user_agent VARCHAR,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- UNIT
CREATE TABLE IF NOT EXISTS unit (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    factor FLOAT NOT NULL,
    unit_type VARCHAR NOT NULL
);

-- ITEM
CREATE TABLE IF NOT EXISTS item (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    description VARCHAR(1000000),
    average_market_price FLOAT,
    unit_type VARCHAR NOT NULL,
    category VARCHAR
);

-- RECIPE
CREATE TABLE IF NOT EXISTS recipe (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    description VARCHAR(1000000),
    image_url VARCHAR,
    original_link VARCHAR,
    preparation_time INTEGER,
    cooking_time INTEGER,
    servings INTEGER,
    instructions VARCHAR(1000000),
    created_at TIMESTAMP NOT NULL,
    public BOOLEAN NOT NULL,
    comment VARCHAR(1000000),
    user_group_id INTEGER NOT NULL,
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- CATEGORY
CREATE TABLE IF NOT EXISTS category (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(1000) NOT NULL,
    user_group_id INTEGER NOT NULL,
    FOREIGN KEY (user_group_id) REFERENCES user_group(id)
);

-- RECIPE_CATEGORY (Many-to-Many relationship between RECIPE and CATEGORY)
CREATE TABLE IF NOT EXISTS recipe_category (
    recipe_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, category_id),
    FOREIGN KEY (recipe_id) REFERENCES recipe(id),
    FOREIGN KEY (category_id) REFERENCES category(id)
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