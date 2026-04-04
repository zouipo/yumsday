INSERT INTO users (username, password, app_admin, created_at, avatar, language, app_theme, last_visited_group_id) VALUES
    ('testuser1', '$2a$12$q7Nm8q9c9g9unKbhjqcWS.Y7tQplxJvgTi8wjsWh7IOPE9ilUwNVm', 0, datetime('now', '-1 day'), '/static/assets/avatar1.jpg', 'EN', 'LIGHT', NULL),
    ('testuser2', '$2a$12$Z30jTp2WrTWT1jOcnZiXvOcIcqhFNyNnKt7yS7FcUUaIHdgVPy3k2', 1, datetime('now', '-1 day'), '/static/assets/avatar2.jpg', 'FR', 'DARK', NULL),
    ('testuser3', '$2a$12$flHptXw2TVYQs3b74duKJO.AkxIoaFPctDSp0AtquuTc82xte4wwy', 0, datetime('now', '-1 day'), '/static/assets/avatar3.jpg', 'EN', 'SYSTEM', NULL),
    ('testuser4', '$2a$12$8dCvoylHH5QIRHlpurXJ3ORMqeGwRkfP3XzytQUVxuPjoIbzj9PWa', 0, datetime('now', '-1 day'), NULL, 'EN', 'LIGHT', NULL);

INSERT INTO groups (name, image_url, created_at) VALUES
    ('Family', '/static/images/family.jpg', datetime('now', '-1 day')),
    ('Friends', '/static/images/friends.jpg', datetime('now', '-1 day')),
    ('Work', NULL, datetime('now', '-1 day'));

-- Update users with last_visited_group_id
UPDATE users
SET last_visited_group_id = (SELECT id FROM groups WHERE name = 'Family')
WHERE username IN ('testuser1', 'testuser2');

UPDATE users
SET last_visited_group_id = (SELECT id FROM groups WHERE name = 'Friends')
WHERE username = 'testuser3';

INSERT INTO group_members (user_id, group_id, admin, joined_at) VALUES
    ((SELECT id FROM users WHERE username = 'testuser1'), (SELECT id FROM groups WHERE name = 'Family'), 1, datetime('now', '-1 day')),
    ((SELECT id FROM users WHERE username = 'testuser2'), (SELECT id FROM groups WHERE name = 'Family'), 1, datetime('now', '-1 day')),
    ((SELECT id FROM users WHERE username = 'testuser3'), (SELECT id FROM groups WHERE name = 'Family'), 0, datetime('now', '-1 day')),
    ((SELECT id FROM users WHERE username = 'testuser1'), (SELECT id FROM groups WHERE name = 'Friends'), 0, datetime('now', '-1 day')),
    ((SELECT id FROM users WHERE username = 'testuser3'), (SELECT id FROM groups WHERE name = 'Friends'), 1, datetime('now', '-1 day'));

INSERT INTO sessions (id, created_at, last_activity, ip_address, user_agent, user_id) VALUES
    ('session123abc', datetime('now', '-1 day'), datetime('now', '-1 day', '+5 hours'), '192.168.1.100', 'Mozilla/5.0', (SELECT id FROM users WHERE username = 'testuser1')),
    ('session456def', datetime('now', '-1 day'), datetime('now', '-1 day', '+1 hour'), '192.168.1.101', 'Chrome/120.0', (SELECT id FROM users WHERE username = 'testuser2')),
    ('session789ghi', datetime('now', '-1 day'), datetime('now', '-1 day', '+2 hours'), '192.168.1.102', 'Safari/17.0', (SELECT id FROM users WHERE username = 'testuser3')),
    ('session999xyz', datetime('now', '-1 day'), datetime('now', '-1 day', '+3 hours'), '192.168.1.103', NULL, (SELECT id FROM users WHERE username = 'testuser4'));

INSERT INTO units (name, factor, unit_type) VALUES
    ('Kilogram', 1000.0, 'WEIGHT'),
    ('Gram', 1.0, 'WEIGHT'),
    ('Liter', 1000.0, 'VOLUME'),
    ('Milliliter', 1.0, 'VOLUME'),
    ('Cup', 240.0, 'VOLUME'),
    ('Tablespoon', 15.0, 'VOLUME'),
    ('Teaspoon', 5.0, 'VOLUME'),
    ('Piece', 1.0, 'PIECE'),
    ('Bag', 1.0, 'BAG'),
    ('Count', 1.0, 'NUMERIC'),
    ('Undefined', 1.0, 'UNDEFINED');

INSERT INTO item_categories (name, group_id) VALUES
    ('GRAINS AND PASTA', (SELECT id FROM groups WHERE name = 'Family')),
    ('BAKED GOODS', (SELECT id FROM groups WHERE name = 'Family')),
    ('SPICES AND CONDIMENTS', (SELECT id FROM groups WHERE name = 'Family')),
    ('DAIRY', (SELECT id FROM groups WHERE name = 'Family')),
    ('MEAT', (SELECT id FROM groups WHERE name = 'Friends')),
    ('VEGETABLES', (SELECT id FROM groups WHERE name = 'Friends')),
    ('SNACKS', (SELECT id FROM groups WHERE name = 'Friends')),
    ('CANNED GOODS', (SELECT id FROM groups WHERE name = 'Friends'));

INSERT INTO items (name, description, average_market_price, unit_type, item_category_id, group_id) VALUES
    ('Flour', 'All-purpose flour', 2.50, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'GRAINS AND PASTA'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Sugar', 'White granulated sugar', 1.80, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'BAKED GOODS'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Salt', 'Table salt', 0.50, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'SPICES AND CONDIMENTS'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Eggs', 'Large eggs', 3.50, 'PIECE', (SELECT id FROM item_categories WHERE name = 'DAIRY'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Milk', 'Whole milk', 2.20, 'VOLUME', (SELECT id FROM item_categories WHERE name = 'DAIRY'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Butter', 'Unsalted butter', 4.00, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'DAIRY'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Chicken Breast', 'Boneless skinless chicken breast', 8.50, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'MEAT'), (SELECT id FROM groups WHERE name = 'Friends')),
    ('Tomatoes', 'Fresh tomatoes', 3.00, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'VEGETABLES'), (SELECT id FROM groups WHERE name = 'Friends')),
    ('Onions', 'Yellow onions', 1.50, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'VEGETABLES'), (SELECT id FROM groups WHERE name = 'Friends')),
    ('Garlic', 'Fresh garlic', 2.00, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'VEGETABLES'), (SELECT id FROM groups WHERE name = 'Friends')),
    ('Water', NULL, NULL, 'VOLUME', NULL, (SELECT id FROM groups WHERE name = 'Friends')),
    ('Pepper', NULL, 1.20, 'WEIGHT', (SELECT id FROM item_categories WHERE name = 'SPICES AND CONDIMENTS'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Olive Oil', 'Extra virgin olive oil', NULL, 'VOLUME', (SELECT id FROM item_categories WHERE name = 'SPICES AND CONDIMENTS'), (SELECT id FROM groups WHERE name = 'Family')),
    ('Potato Chips', 'Salted potato chips', 2.99, 'BAG', (SELECT id FROM item_categories WHERE name = 'SNACKS'), (SELECT id FROM groups WHERE name = 'Friends')),
    ('Canned Beans', 'Black beans', 1.50, 'NUMERIC', (SELECT id FROM item_categories WHERE name = 'CANNED GOODS'), (SELECT id FROM groups WHERE name = 'Friends'));

INSERT INTO recipes (name, description, image_url, original_link, preparation_time_min, cooking_time_min, servings, instructions, created_at, public, comment, group_id) VALUES
    ('Chocolate Chip Cookies', 'Classic homemade chocolate chip cookies', '/static/recipes/cookies.jpg', 'https://example.com/cookies', 15, 12, 24, 'Mix ingredients and bake at 350F', datetime('now', '-1 day'), 1, 'Family favorite!', 1),
    ('Grilled Chicken', 'Simple grilled chicken breast with herbs', '/static/recipes/chicken.jpg', NULL, 10, 20, 4, 'Season and grill until cooked through', datetime('now', '-1 day'), 1, NULL, 1),
    ('Tomato Soup', 'Creamy tomato soup', '/static/recipes/soup.jpg', 'https://example.com/soup', 10, 30, 6, 'Cook tomatoes with onions and blend', datetime('now', '-1 day'), 0, 'Great for winter', 2),
    ('Quick Salad', NULL, NULL, NULL, NULL, NULL, NULL, NULL, datetime('now', '-1 day'), 1, NULL, 1),
    ('Secret Recipe', 'Top secret family recipe', NULL, NULL, 5, 15, 2, 'Cannot reveal instructions', datetime('now', '-1 day'), 0, 'Do not share!', 2);

INSERT INTO recipe_categories (name, group_id) VALUES
    ('DESSERT', 1),
    ('MAIN COURSE', 1),
    ('SOUP', 1),
    ('VEGETARIAN', 2),
    ('SALAD', 2),
    ('BREAKFAST', 1),
    ('VEGAN', 2),
    ('GLUTEN FREE', 2);

INSERT INTO recipes_categories_junction (recipe_id, category_id) VALUES
    (1, 1),  -- Cookies -> DESSERT
    (2, 2),  -- Chicken -> MAIN COURSE
    (3, 3),  -- Soup -> SOUP
    (3, 4),  -- Soup -> VEGETARIAN
    (4, 5),  -- Salad -> SALAD
    (4, 7);  -- Salad -> VEGAN

INSERT INTO ingredients (quantity, item_id, unit_id, recipe_id) VALUES
    -- Chocolate Chip Cookies ingredients
    (2.0, 1, 1, 1),      -- 2 cups flour
    (1.0, 2, 1, 1),      -- 1 cup sugar
    (0.5, 6, 1, 1),      -- 0.5 cup butter
    (2.0, 4, 8, 1),      -- 2 eggs
    -- Grilled Chicken ingredients
    (4.0, 7, 8, 2),      -- 4 pieces chicken breast
    (2.0, 10, 2, 2),     -- 2 cloves garlic
    (0.5, 3, 7, 2),      -- 0.5 tsp salt
    -- Tomato Soup ingredients
    (6.0, 8, 8, 3),      -- 6 tomatoes
    (1.0, 9, 8, 3),      -- 1 onion
    (2.0, 10, 2, 3),     -- 2 cloves garlic
    (1.0, 3, 7, 3),      -- 1 tsp salt
    -- Quick Salad ingredients
    (NULL, 8, 8, 4),     -- Tomatoes with NULL quantity (to taste)
    (1.0, 13, 11, 4),    -- Olive oil with undefined unit
    (NULL, 12, 11, 4);   -- Pepper with NULL quantity and undefined unit

-- Dishes
INSERT INTO dishes (portion, bought, datetime, recipe_id, group_id) VALUES
    (4, 0, datetime('now', '-1 day', '+6 hours'), 2, 1),  -- Grilled chicken for dinner
    (6, 0, datetime('now', '-1 day', '+12 hours'), 3, 1),  -- Tomato soup for lunch
    (12, 1, datetime('now', '-1 day', '+15 hours'), 1, 1); -- Cookies (bought)

-- Grocery List
INSERT INTO groceries (quantity_bought, user_quantity, item_id, unit_id, group_id) VALUES
    (0.0, 2.0, 1, 1, 1),    -- 2 kg flour
    (0.5, 1.0, 2, 1, 1),    -- 1 kg sugar (half bought)
    (12.0, 12.0, 4, 8, 1),  -- 12 eggs (all bought)
    (0.0, 2.0, 5, 3, 1),    -- 2 liters milk
    (0.0, 1.0, 7, 1, 2),    -- 1 kg chicken breast
    (2.0, 4.0, 8, 8, 2),    -- 4 tomatoes (2 bought)
    (0.0, 1.0, 12, 11, 1),  -- Pepper with undefined unit
    (0.0, 3.0, 11, 11, 2);  -- Water with undefined unit
