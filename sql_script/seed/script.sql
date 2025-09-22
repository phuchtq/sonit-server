-- Categories --
CREATE TABLE IF NOT EXISTS categories (
	id character varying(100) PRIMARY KEY,
	name character varying(100) NOT NULL,
	description text,
	active_status bool NOT NULL DEFAULT TRUE,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
); 

-- Collections --
CREATE TABLE IF NOT EXISTS collections (
	id character varying(100) PRIMARY KEY,
	name character varying(100) NOT NULL,
	description text,
	active_status bool NOT NULL,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
); 

-- Products --
CREATE TABLE IF NOT EXISTS products (
	id character varying(100) PRIMARY KEY,
	category_id character varying(100) NOT NULL,
	collection_id character varying(100),
	name character varying(100) UNIQUE NOT NULL,
	description text,
	image character varying(100),
	size text NOT NULL,
	color text NOT NULL,
	price money,
	currency character varying(20) NOT NULL,
	active_status bool NOT NULL DEFAULT TRUE,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_product_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
	CONSTRAINT fk_product_collection FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
); 

-- Roles --
CREATE TABLE IF NOT EXISTS roles (
	id character varying(100) PRIMARY KEY,
	name character varying(100) NOT NULL,
	active_status bool NOT NULL,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
); 

-- Users --
CREATE TABLE IF NOT EXISTS users
(
    id character varying(100) PRIMARY KEY,
	role_id character varying(100) NOT NULL,
	full_name character varying(100) NOT NULL,
	email character varying(100) NOT NULL,
	password character varying(100) NOT NULL,
	profile_avatar character varying(100) NOT NULL,
	gender character varying(100) NOT NULL,
	is_vip bool DEFAULT false,
	vip_code character varying(30),
	is_active bool NOT NULL,
	is_activated bool NOT NULL,
	is_have_to_reset_password bool,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- User Securities -- 
CREATE TABLE IF NOT EXISTS user_securities 
(
	id character varying(100) PRIMARY KEY,
	access_token character varying(100) NULL,
	refresh_token character varying(100) NULL,
	action_token character varying(100) NULL,
	fail_access int DEFAULT 0,
	last_fail timestamp without time zone NULL,
	CONSTRAINT fk_security_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);


-- Orders -- 
CREATE TABLE IF NOT EXISTS orders 
(
	id character varying(100) PRIMARY KEY,
	user_id character varying(100) NOT NULL,
	items character varying(500) NOT NULL,
	total_amount money NOT NULL,
	currency character varying(100) NOT NULL,
	status character varying(100) NOT NULL,
	note character varying(100) NOT NULL,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_order_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Payments --
CREATE TABLE IF NOT EXISTS payments 
(
	id character varying(100) PRIMARY KEY,
	user_id character varying(100) NOT NULL,
	order_id character varying(100) NOT NULL,
	transaction_id character varying(100) NOT NULL,
	amount money NOT NULL,
	currency character varying(100) NOT NULL,
	status character varying(100) NOT NULL,
	method character varying(100) NOT NULL,
	created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_payment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Vouchers --
CREATE TABLE IF NOT EXISTS vouchers (
    id character varying(100) PRIMARY KEY,
    code VARCHAR(255) UNIQUE NOT NULL,
    discount NUMERIC(5, 2) NOT NULL, -- Assuming precision 5, scale 2 (e.g., 100.00)
    amount BIGINT NOT NULL DEFAULT 0,
    description TEXT,
    active_status BOOLEAN NOT NULL DEFAULT TRUE,
    allowed_category_ids TEXT[],
    allowed_product_ids TEXT[],
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_vouchers_code ON vouchers (code);
CREATE INDEX idx_vouchers_expired_at ON vouchers (expired_at);
CREATE INDEX idx_vouchers_active_status ON vouchers (active_status);
CREATE INDEX idx_vouchers_allowed_category_ids ON vouchers USING GIN (allowed_category_ids);
CREATE INDEX idx_vouchers_allowed_product_ids ON vouchers USING GIN (allowed_product_ids);

-- Feedbacks -- 
CREATE TABLE IF NOT EXISTS feedbacks (
    id character varying(100) PRIMARY KEY,
    user_id character varying(100) NOT NULL,
	order_id character varying(100),
	staff_id character varying(100),
	rating int,
    content TEXT NOT NULL,
    attachment character varying(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Carts --
CREATE TABLE IF NOT EXISTS carts (
    id character varying(100) PRIMARY KEY,
    items character varying(500) NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_cart_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);

-- Product Inventories -- 
CREATE TABLE IF NOT EXISTS product_inventories (
    id character varying(100) PRIMARY KEY,
	current_quantity BIGINT NOT NULL,
	CONSTRAINT fk_inventory_product FOREIGN KEY (id) REFERENCES products(id) ON DELETE CASCADE
);

-- Product Inventory Transactions --
CREATE TABLE IF NOT EXISTS product_inventory_transactions (
    id character varying(100) PRIMARY KEY,
	product_id character varying(100) NOT NULL,
	amount BIGINT NOT NULL,
	action character varying(100), 
	note character varying(100),
	date TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_inventoryTx_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Shippings --
CREATE TABLE IF NOT EXISTS shippings (
    id character varying(100) PRIMARY KEY,
	delivery_code character varying(100),
	shipping_unit character varying(100),
	shipping_detail character varying(100), 
	delivered_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_shipping_order FOREIGN KEY (id) REFERENCES orders(id) ON DELETE CASCADE
);


--- Seeding dadta ---

-- Roles
INSERT INTO roles (id, name, active_status) VALUES ('R001', 'Admin', True);
INSERT INTO roles (id, name, active_status) VALUES ('R003', 'Customer', True);

-- Categories
INSERT INTO categories (id, name, description) VALUES 
('cat1', 'Cues', 'Professional and casual billiard cues'),
('cat2', 'Chalks', 'Cue chalks for better grip'),
('cat3', 'Gloves', 'Billiard gloves for smoother stroke'),
('cat4', 'Cue Cases', 'Cases to protect your cues'),
('cat5', 'Accessories', 'Other essential accessories');

-- Collections
INSERT INTO collections (id, name, description, active_status) VALUES 
('col1', '2025 Premium Line', 'Top-tier billiard equipment for 2025', True),
('col2', 'Beginner Series', 'Perfect for those just starting', True),
('col3', 'Limited Editions', 'Special edition billiard gear', True),
('col4', 'Pro Player Picks', 'Gear recommended by professionals', True),
('col5', 'Seasonal Specials', 'Limited-time products for each season', True);

-- Products
INSERT INTO products (id, category_id, collection_id, name, description, image, size, color, price, currency) VALUES 
('prod1', 'cat1', 'col1', 'Predator Revo Cue', 'High-performance carbon fiber cue', 'cue1.png', '12.4mm', 'Black', 699.99, 'VND'),
('prod2', 'cat2', 'col2', 'Kamui Chalk Beta', 'Premium chalk for precision shots', 'chalk1.png', 'Standard', 'Blue', 24.99, 'VND'),
('prod3', 'cat3', 'col3', 'Molinari Glove', 'Breathable and durable glove', 'glove1.png', 'L', 'Black', 29.99, 'VND'),
('prod4', 'cat4', 'col4', 'CueTec Case 2x4', 'Strong hard case for 2 butts and 4 shafts', 'case1.png', 'Medium', 'Silver', 89.99, 'VND'),
('prod5', 'cat5', 'col5', 'Tip Shaper Tool', 'Multi-tool for shaping cue tips', 'tool1.png', 'Compact', 'Red', 19.99, 'VND');

-- Product Inventory
INSERT INTO product_inventories (id, current_quantity) VALUES
('prod1', 100), ('prod2', 100), ('prod3', 100), ('prod4', 100), ('prod5', 100);

-- Product Inventory Transactions
INSERT INTO product_inventory_transactions (id, product_id, amount, action, note, date) VALUES
('tx1', 'prod1', 100, 'import', 'Initial stock', CURRENT_TIMESTAMP),
('tx2', 'prod2', 100, 'import', 'Initial stock', CURRENT_TIMESTAMP),
('tx3', 'prod3', 100, 'import', 'Initial stock', CURRENT_TIMESTAMP),
('tx4', 'prod4', 100, 'import', 'Initial stock', CURRENT_TIMESTAMP),
('tx5', 'prod5', 100, 'import', 'Initial stock', CURRENT_TIMESTAMP);

-- Users
-- Password shared: hashedpass1
INSERT INTO users (id, role_id, full_name, email, password, profile_avatar, gender, is_active, is_activated) VALUES 
('user1', 'R001', 'Alice CueMaster', 'alice@cuebiz.com', '$2b$10$7qmVXNnx25I3W9OLvvxZjO9e3/5CSZRF70Fr2eCG5Yl8f2bgPW08y', 'ava1.png', 'female', True, True),
('user2', 'R003', 'Bob Breaker', 'bob@cuebiz.com', '$2b$10$7qmVXNnx25I3W9OLvvxZjO9e3/5CSZRF70Fr2eCG5Yl8f2bgPW08y', 'ava2.png', 'male', True, True),
('user3', 'R003', 'Charlie Chalk', 'charlie@cuebiz.com', '$2b$10$7qmVXNnx25I3W9OLvvxZjO9e3/5CSZRF70Fr2eCG5Yl8f2bgPW08y', 'ava3.png', 'male', True, True),
('user4', 'R003', 'Diana Defense', 'diana@cuebiz.com', '$2b$10$7qmVXNnx25I3W9OLvvxZjO9e3/5CSZRF70Fr2eCG5Yl8f2bgPW08y', 'ava4.png', 'female', True, True),
('user5', 'R003', 'Evan Eightball', 'evan@cuebiz.com', '$2b$10$7qmVXNnx25I3W9OLvvxZjO9e3/5CSZRF70Fr2eCG5Yl8f2bgPW08y', 'ava5.png', 'male', True, True);

-- User Securities
INSERT INTO user_securities (id, access_token, refresh_token, action_token, fail_access, last_fail) VALUES 
('user1', 'user1_at', NULL, NULL, 0, NULL),
('user2', 'user2_at', NULL, NULL, 0, NULL),
('user3', 'user3_at', NULL, NULL, 0, NULL),
('user4', 'user4_at', NULL, NULL, 0, NULL),
('user5', 'user5_at', NULL, NULL, 0, NULL);

-- Orders
INSERT INTO orders (id, user_id, items, total_amount, currency, status, note) VALUES 
('ord1', 'user2', 'prod1', 699.99, 'VND', 'Pending', 'First custom cue'),
('ord2', 'user3', 'prod2,prod5', 44.98, 'VND', 'Completed', 'Gift set'),
('ord3', 'user4', 'prod3', 29.99, 'VND', 'Processing', 'New glove'),
('ord4', 'user5', 'prod4', 89.99, 'VND', 'Delivered', 'Cue case delivered'),
('ord5', 'user2', 'prod2,prod3', 54.98, 'VND', 'Shipped', 'Accessories combo');

-- Payments
INSERT INTO payments (id, user_id, order_id, transaction_id, amount, currency, status, method) VALUES 
('pay1', 'user2', 'ord1', 'txn1', 699.99, 'VND', 'Success', 'Credit Card'),
('pay2', 'user3', 'ord2', 'txn2', 44.98, 'VND', 'Success', 'Credit Card'),
('pay3', 'user4', 'ord3', 'txn3', 29.99, 'VND', 'Success', 'Credit Card'),
('pay4', 'user5', 'ord4', 'txn4', 89.99, 'VND', 'Success', 'Credit Card'),
('pay5', 'user2', 'ord5', 'txn5', 54.98, 'VND', 'Success', 'Credit Card');

-- Feedbacks
INSERT INTO feedbacks (id, user_id, order_id, staff_id, rating, content, attachment) VALUES 
('fb1', 'user2', 'ord1', NULL, 5, 'Feedback on prod1', NULL),
('fb2', 'user3', 'ord2', NULL, 4, 'Feedback on prod2,prod5', NULL),
('fb3', 'user4', 'ord3', NULL, 3, 'Feedback on prod3', NULL),
('fb4', 'user5', 'ord4', NULL, 5, 'Feedback on prod4', NULL),
('fb5', 'user2', 'ord5', NULL, 4, 'Feedback on prod2,prod3', NULL);

-- Vouchers
INSERT INTO vouchers (id, code, discount, amount, description, allowed_category_ids, allowed_product_ids, expired_at) VALUES 
('vch1', 'CUE10', 10.0, 50, '10% off all cues', '{cat1}', '{prod1}', NOW() + INTERVAL '30 days'),
('vch2', 'CHALKFREE', 100.0, 100, 'Free chalk with any order', '{cat2}', '{prod2}', NOW() + INTERVAL '15 days'),
('vch3', 'GLOVE15', 15.0, 20, '15% off gloves', '{cat3}', '{prod3}', NOW() + INTERVAL '45 days'),
('vch4', 'CASE20', 20.0, 30, '20% off cases', '{cat4}', '{prod4}', NOW() + INTERVAL '60 days'),
('vch5', 'TOOL5', 5.0, 200, '5 VND off accessories', '{cat5}', '{prod5}', NOW() + INTERVAL '20 days');

-- Shippings
INSERT INTO shippings (id, delivery_code, shipping_unit, shipping_detail, delivered_at) VALUES 
('ord1', 'SHIP_1000', 'DHL', 'Fast shipping', NOW()),
('ord2', 'SHIP_1001', 'DHL', 'Fast shipping', NOW() + INTERVAL '1 day'),
('ord3', 'SHIP_1002', 'DHL', 'Fast shipping', NOW() + INTERVAL '2 day'),
('ord4', 'SHIP_1003', 'DHL', 'Fast shipping', NOW() + INTERVAL '3 day'),
('ord5', 'SHIP_1004', 'DHL', 'Fast shipping', NOW() + INTERVAL '4 day');
