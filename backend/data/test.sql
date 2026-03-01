INSERT INTO user (id, username, password, app_admin, created_at, avatar, language, app_theme, last_visited_group) VALUES
    (1, 'testuser1', '$2a$12$q7Nm8q9c9g9unKbhjqcWS.Y7tQplxJvgTi8wjsWh7IOPE9ilUwNVm', 0, datetime('now', '-1 day'), '/static/assets/avatar1.jpg', 'EN', 'LIGHT', NULL),
    (2, 'testuser2', '$2a$12$Z30jTp2WrTWT1jOcnZiXvOcIcqhFNyNnKt7yS7FcUUaIHdgVPy3k2', 1, datetime('now', '-1 day'), '/static/assets/avatar2.jpg', 'FR', 'DARK', NULL),
    (3, 'testuser3', '$2a$12$flHptXw2TVYQs3b74duKJO.AkxIoaFPctDSp0AtquuTc82xte4wwy', 0, datetime('now', '-1 day'), '/static/assets/avatar3.jpg', 'EN', 'SYSTEM', NULL),
    (4, 'testuser4', '$2a$12$8dCvoylHH5QIRHlpurXJ3ORMqeGwRkfP3XzytQUVxuPjoIbzj9PWa', 0, datetime('now', '-1 day'), NULL, 'EN', 'LIGHT', NULL);

INSERT INTO user_group (id, name, image_url, created_at) VALUES
    (1, 'Family', '/static/images/family.jpg', datetime('now', '-1 day')),
    (2, 'Friends', '/static/images/friends.jpg', datetime('now', '-1 day')),
    (3, 'Work', NULL, datetime('now', '-1 day'));

-- Update users with last_visited_group
UPDATE user SET last_visited_group = 1 WHERE id = 1;
UPDATE user SET last_visited_group = 1 WHERE id = 2;
UPDATE user SET last_visited_group = 2 WHERE id = 3;

INSERT INTO member_group (user_id, user_group_id, admin, joined_at) VALUES
    (1, 1, 1, datetime('now', '-1 day')),
    (2, 1, 1, datetime('now', '-1 day')),
    (3, 1, 0, datetime('now', '-1 day')),
    (1, 2, 0, datetime('now', '-1 day')),
    (3, 2, 1, datetime('now', '-1 day'));

INSERT INTO session (id, session_id, created_at, expires_at, last_activity, ip_address, user_agent, user_id) VALUES
    (1, 'session123abc', datetime('now', '-1 day'), datetime('now', '-1 day', '+1 month'), datetime('now', '-1 day', '+5 hours'), '192.168.1.100', 'Mozilla/5.0', 1),
    (2, 'session456def', datetime('now', '-1 day'), datetime('now', '-1 day', '+1 month'), datetime('now', '-1 day', '+1 hour'), '192.168.1.101', 'Chrome/120.0', 2),
    (3, 'session789ghi', datetime('now', '-1 day'), datetime('now', '-1 day', '+1 month'), datetime('now', '-1 day', '+2 hours'), '192.168.1.102', 'Safari/17.0', 3),
    (4, 'session999xyz', datetime('now', '-1 day'), datetime('now', '-1 day', '+1 month'), datetime('now', '-1 day', '+3 hours'), '192.168.1.103', NULL, 4);

INSERT INTO unit (id, name, factor, unit_type) VALUES
    (1, 'Kilogram', 1000.0, 'WEIGHT'),
    (2, 'Gram', 1.0, 'WEIGHT'),
    (3, 'Liter', 1000.0, 'VOLUME'),
    (4, 'Milliliter', 1.0, 'VOLUME'),
    (5, 'Cup', 240.0, 'VOLUME'),
    (6, 'Tablespoon', 15.0, 'VOLUME'),
    (7, 'Teaspoon', 5.0, 'VOLUME'),
    (8, 'Piece', 1.0, 'PIECE'),
    (9, 'Bag', 1.0, 'BAG'),
    (10, 'Count', 1.0, 'NUMERIC'),
    (11, 'Undefined', 1.0, 'UNDEFINED');

INSERT INTO item (id, name, description, average_market_price, unit_type, category) VALUES
    (1, 'Flour', 'All-purpose flour', 2.50, 'WEIGHT', 'GRAINS AND PASTA'),
    (2, 'Sugar', 'White granulated sugar', 1.80, 'WEIGHT', 'BAKED GOODS'),
    (3, 'Salt', 'Table salt', 0.50, 'WEIGHT', 'SPICES AND CONDIMENTS'),
    (4, 'Eggs', 'Large eggs', 3.50, 'PIECE', 'DAIRY'),
    (5, 'Milk', 'Whole milk', 2.20, 'VOLUME', 'DAIRY'),
    (6, 'Butter', 'Unsalted butter', 4.00, 'WEIGHT', 'DAIRY'),
    (7, 'Chicken Breast', 'Boneless skinless chicken breast', 8.50, 'WEIGHT', 'MEAT'),
    (8, 'Tomatoes', 'Fresh tomatoes', 3.00, 'WEIGHT', 'VEGETABLES'),
    (9, 'Onions', 'Yellow onions', 1.50, 'WEIGHT', 'VEGETABLES'),
    (10, 'Garlic', 'Fresh garlic', 2.00, 'WEIGHT', 'VEGETABLES'),
    (11, 'Water', NULL, NULL, 'VOLUME', NULL),
    (12, 'Pepper', NULL, 1.20, 'WEIGHT', 'SPICES AND CONDIMENTS'),
    (13, 'Olive Oil', 'Extra virgin olive oil', NULL, 'VOLUME', 'SPICES AND CONDIMENTS'),
    (14, 'Potato Chips', 'Salted potato chips', 2.99, 'BAG', 'SNACKS'),
    (15, 'Canned Beans', 'Black beans', 1.50, 'NUMERIC', 'CANNED GOODS');

INSERT INTO recipe (id, name, description, image_url, original_link, preparation_time, cooking_time, servings, instructions, created_at, public, comment, user_group_id) VALUES
    (1, 'Chocolate Chip Cookies', 'Classic homemade chocolate chip cookies', '/static/recipes/cookies.jpg', 'https://example.com/cookies', 15, 12, 24, 'Mix ingredients and bake at 350F', datetime('now', '-1 day'), 1, 'Family favorite!', 1),
    (2, 'Grilled Chicken', 'Simple grilled chicken breast with herbs', '/static/recipes/chicken.jpg', NULL, 10, 20, 4, 'Season and grill until cooked through', datetime('now', '-1 day'), 1, NULL, 1),
    (3, 'Tomato Soup', 'Creamy tomato soup', '/static/recipes/soup.jpg', 'https://example.com/soup', 10, 30, 6, 'Cook tomatoes with onions and blend', datetime('now', '-1 day'), 0, 'Great for winter', 2),
    (4, 'Quick Salad', NULL, NULL, NULL, NULL, NULL, NULL, NULL, datetime('now', '-1 day'), 1, NULL, 1),
    (5, 'Secret Recipe', 'Top secret family recipe', NULL, NULL, 5, 15, 2, 'Cannot reveal instructions', datetime('now', '-1 day'), 0, 'Do not share!', 2);

INSERT INTO category (id, name, user_group_id) VALUES
    (1, 'DESSERT', 1),
    (2, 'MAIN COURSE', 1),
    (3, 'SOUP', 1),
    (4, 'VEGETARIAN', 2),
    (5, 'SALAD', 2),
    (6, 'BREAKFAST', 1),
    (7, 'VEGAN', 2),
    (8, 'GLUTEN FREE', 2);

INSERT INTO recipe_category (recipe_id, category_id) VALUES
    (1, 1),  -- Cookies -> DESSERT
    (2, 2),  -- Chicken -> MAIN COURSE
    (3, 3),  -- Soup -> SOUP
    (3, 4),  -- Soup -> VEGETARIAN
    (4, 5),  -- Salad -> SALAD
    (4, 7);  -- Salad -> VEGAN

INSERT INTO ingredient (id, quantity, item_id, unit_id, reciped_id) VALUES
    -- Chocolate Chip Cookies ingredients
    (1, 2.0, 1, 1, 1),      -- 2 cups flour
    (2, 1.0, 2, 1, 1),      -- 1 cup sugar
    (3, 0.5, 6, 1, 1),      -- 0.5 cup butter
    (4, 2.0, 4, 8, 1),      -- 2 eggs
    -- Grilled Chicken ingredients
    (5, 4.0, 7, 8, 2),      -- 4 pieces chicken breast
    (6, 2.0, 10, 2, 2),     -- 2 cloves garlic
    (7, 0.5, 3, 7, 2),      -- 0.5 tsp salt
    -- Tomato Soup ingredients
    (8, 6.0, 8, 8, 3),      -- 6 tomatoes
    (9, 1.0, 9, 8, 3),      -- 1 onion
    (10, 2.0, 10, 2, 3),     -- 2 cloves garlic
    (11, 1.0, 3, 7, 3),      -- 1 tsp salt
    -- Quick Salad ingredients
    (12, NULL, 8, 8, 4),     -- Tomatoes with NULL quantity (to taste)
    (13, 1.0, 13, NULL, 4),  -- Olive oil with NULL unit
    (14, NULL, 12, NULL, 4); -- Pepper with NULL quantity and unit

-- Dishes
INSERT INTO dish (id, portion, bought, datetime, recipe_id, user_group_id) VALUES
    (1, 4, 0, datetime('now', '-1 day', '+6 hours'), 2, 1),  -- Grilled chicken for dinner
    (2, 6, 0, datetime('now', '-1 day', '+12 hours'), 3, 1),  -- Tomato soup for lunch
    (3, 12, 1, datetime('now', '-1 day', '+15 hours'), 1, 1); -- Cookies (bought)

-- Grocery List
INSERT INTO grocery_list (id, quantity, quantity_bought, user_quantity, item_id, unit_id, user_group_id) VALUES
    (1, 2.0, 0.0, 2.0, 1, 1, 1),       -- 2 kg flour
    (2, 1.0, 0.5, 1.0, 2, 1, 1),       -- 1 kg sugar (half bought)
    (3, 12.0, 12.0, 12.0, 4, 8, 1),    -- 12 eggs (all bought)
    (4, 2.0, 0.0, 2.0, 5, 3, 1),       -- 2 liters milk
    (5, 1.0, 0.0, 1.0, 7, 1, 2),       -- 1 kg chicken breast
    (6, 4.0, 2.0, 4.0, 8, 8, 2),       -- 4 tomatoes (2 bought)
    (7, 1.0, 0.0, 1.0, 12, NULL, 1),   -- Pepper with NULL unit
    (8, 3.0, 0.0, 3.0, 11, NULL, 2);   -- Water with NULL unit
