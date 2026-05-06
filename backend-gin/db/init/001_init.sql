CREATE DATABASE IF NOT EXISTS selfcook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE selfcook;



CREATE TABLE IF NOT EXISTS pickup_points (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  name VARCHAR(100) NOT NULL,

  contact_name VARCHAR(50) DEFAULT '',

  contact_phone VARCHAR(20) DEFAULT '',

  address VARCHAR(255) NOT NULL,

  business_hours VARCHAR(100) DEFAULT '',

  status VARCHAR(20) NOT NULL DEFAULT 'active',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

);



CREATE TABLE IF NOT EXISTS products (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  name VARCHAR(100) NOT NULL,

  subtitle VARCHAR(255) DEFAULT '',

  cover_image VARCHAR(255) DEFAULT '',

  category_name VARCHAR(50) DEFAULT '',

  description TEXT,

  status VARCHAR(20) NOT NULL DEFAULT 'on_sale',

  sort_order INT NOT NULL DEFAULT 0,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

);



CREATE TABLE IF NOT EXISTS product_skus (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  product_id BIGINT UNSIGNED NOT NULL,

  sku_name VARCHAR(100) NOT NULL,

  sku_code VARCHAR(50) NOT NULL,

  price DECIMAL(10,2) NOT NULL,

  original_price DECIMAL(10,2) NOT NULL DEFAULT 0,

  stock_total INT NOT NULL DEFAULT 0,

  stock_available INT NOT NULL DEFAULT 0,

  limit_per_user INT NOT NULL DEFAULT 0,

  limit_per_order INT NOT NULL DEFAULT 0,

  status VARCHAR(20) NOT NULL DEFAULT 'active',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_product_skus_product_id (product_id)

);



CREATE TABLE IF NOT EXISTS `groups` (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  title VARCHAR(150) NOT NULL,

  cover_image VARCHAR(255) DEFAULT '',

  leader_user_id BIGINT UNSIGNED NOT NULL DEFAULT 1,

  status VARCHAR(20) NOT NULL DEFAULT 'ongoing',

  start_at DATETIME NOT NULL,

  cutoff_at DATETIME NOT NULL,

  fulfillment_mode VARCHAR(20) NOT NULL DEFAULT 'mixed',

  allow_modify_before_cutoff TINYINT(1) NOT NULL DEFAULT 0,

  show_join_list TINYINT(1) NOT NULL DEFAULT 0,

  pickup_rule_desc VARCHAR(255) DEFAULT '',

  delivery_rule_desc VARCHAR(255) DEFAULT '',

  group_notice TEXT,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_groups_status (status),

  KEY idx_groups_cutoff_at (cutoff_at)

);



CREATE TABLE IF NOT EXISTS group_items (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  group_id BIGINT UNSIGNED NOT NULL,

  product_id BIGINT UNSIGNED NOT NULL,

  product_sku_id BIGINT UNSIGNED NOT NULL,

  product_name_snapshot VARCHAR(100) NOT NULL,

  sku_name_snapshot VARCHAR(100) NOT NULL,

  cover_image_snapshot VARCHAR(255) DEFAULT '',

  price_snapshot DECIMAL(10,2) NOT NULL,

  original_price_snapshot DECIMAL(10,2) NOT NULL DEFAULT 0,

  stock_total_snapshot INT NOT NULL DEFAULT 0,

  stock_available_snapshot INT NOT NULL DEFAULT 0,

  limit_per_user_snapshot INT NOT NULL DEFAULT 0,

  limit_per_order_snapshot INT NOT NULL DEFAULT 0,

  status VARCHAR(20) NOT NULL DEFAULT 'active',

  sort_order INT NOT NULL DEFAULT 0,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_group_items_group_id (group_id)

);



CREATE TABLE IF NOT EXISTS orders (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  order_no VARCHAR(40) NOT NULL UNIQUE,

  group_id BIGINT UNSIGNED NOT NULL,

  user_id BIGINT UNSIGNED NOT NULL DEFAULT 1,

  status VARCHAR(30) NOT NULL,

  fulfillment_mode VARCHAR(20) NOT NULL,

  contact_name VARCHAR(50) NOT NULL,

  contact_phone VARCHAR(20) NOT NULL,

  pickup_point_id BIGINT UNSIGNED NULL,

  address_id BIGINT UNSIGNED NULL,

  delivery_address_snapshot VARCHAR(255) DEFAULT '',

  goods_amount DECIMAL(10,2) NOT NULL,

  discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,

  delivery_fee DECIMAL(10,2) NOT NULL DEFAULT 0,

  payable_amount DECIMAL(10,2) NOT NULL,

  paid_amount DECIMAL(10,2) NOT NULL,

  remark VARCHAR(255) DEFAULT '',

  cutoff_at_snapshot DATETIME NOT NULL,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_orders_group_id (group_id),

  KEY idx_orders_status (status)

);



CREATE TABLE IF NOT EXISTS order_items (

  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  order_id BIGINT UNSIGNED NOT NULL,

  group_item_id BIGINT UNSIGNED NOT NULL,

  product_id BIGINT UNSIGNED NOT NULL,

  product_sku_id BIGINT UNSIGNED NOT NULL,

  product_name_snapshot VARCHAR(100) NOT NULL,

  sku_name_snapshot VARCHAR(100) NOT NULL,

  unit_price_snapshot DECIMAL(10,2) NOT NULL,

  quantity INT NOT NULL,

  subtotal_amount DECIMAL(10,2) NOT NULL,

  taste_remark VARCHAR(100) DEFAULT '',

  item_status VARCHAR(20) NOT NULL DEFAULT 'normal',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_order_items_order_id (order_id)

);



INSERT INTO pickup_points (id, name, contact_name, contact_phone, address, business_hours, status)

VALUES (1, 'A栋一楼前台', '管理员', '13800001111', '科技园 A栋一楼前台', '11:00-13:00', 'active')

ON DUPLICATE KEY UPDATE

  name = VALUES(name),

  contact_name = VALUES(contact_name),

  contact_phone = VALUES(contact_phone),

  address = VALUES(address),

  business_hours = VALUES(business_hours),

  status = VALUES(status);



INSERT INTO products (id, name, subtitle, cover_image, category_name, description, status, sort_order)

VALUES

(1, '番茄牛腩饭', '大块牛腩，酸甜开胃', 'https://example.com/beef-rice.jpg', '热销套餐', '含米饭、配菜、例汤', 'on_sale', 1),

(2, '黑椒鸡排饭', '现煎鸡排，香浓黑椒汁', 'https://example.com/chicken-rice.jpg', '热销套餐', '含米饭、配菜、例汤', 'on_sale', 2)

ON DUPLICATE KEY UPDATE

  name = VALUES(name),

  subtitle = VALUES(subtitle),

  cover_image = VALUES(cover_image),

  category_name = VALUES(category_name),

  description = VALUES(description),

  status = VALUES(status),

  sort_order = VALUES(sort_order);



INSERT INTO product_skus (id, product_id, sku_name, sku_code, price, original_price, stock_total, stock_available, limit_per_user, limit_per_order, status)

VALUES

(1, 1, '大份', 'BEEF_RICE_L', 22.00, 25.00, 100, 100, 2, 4, 'active'),

(2, 2, '标准份', 'CHICKEN_RICE_M', 18.00, 20.00, 100, 100, 2, 4, 'active')

ON DUPLICATE KEY UPDATE

  product_id = VALUES(product_id),

  sku_name = VALUES(sku_name),

  sku_code = VALUES(sku_code),

  price = VALUES(price),

  original_price = VALUES(original_price),

  stock_total = VALUES(stock_total),

  stock_available = VALUES(stock_available),

  limit_per_user = VALUES(limit_per_user),

  limit_per_order = VALUES(limit_per_order),

  status = VALUES(status);



INSERT INTO `groups` (id, title, cover_image, leader_user_id, status, start_at, cutoff_at, fulfillment_mode, allow_modify_before_cutoff, show_join_list, pickup_rule_desc, delivery_rule_desc, group_notice)

VALUES

(1, '3号楼午餐接龙', 'https://example.com/group-lunch.jpg', 1, 'ongoing', NOW(), DATE_ADD(NOW(), INTERVAL 6 HOUR), 'mixed', 1, 1, '11:30-12:00 到 A栋一楼前台自提', '配送范围限科技园周边', '午餐现做现配，请按时取餐')

ON DUPLICATE KEY UPDATE

  title = VALUES(title),

  cover_image = VALUES(cover_image),

  leader_user_id = VALUES(leader_user_id),

  status = VALUES(status),

  start_at = VALUES(start_at),

  cutoff_at = VALUES(cutoff_at),

  fulfillment_mode = VALUES(fulfillment_mode),

  allow_modify_before_cutoff = VALUES(allow_modify_before_cutoff),

  show_join_list = VALUES(show_join_list),

  pickup_rule_desc = VALUES(pickup_rule_desc),

  delivery_rule_desc = VALUES(delivery_rule_desc),

  group_notice = VALUES(group_notice);



INSERT INTO group_items (id, group_id, product_id, product_sku_id, product_name_snapshot, sku_name_snapshot, cover_image_snapshot, price_snapshot, original_price_snapshot, stock_total_snapshot, stock_available_snapshot, limit_per_user_snapshot, limit_per_order_snapshot, status, sort_order)

VALUES

(1, 1, 1, 1, '番茄牛腩饭', '大份', 'https://example.com/beef-rice.jpg', 22.00, 25.00, 50, 50, 2, 4, 'active', 1),

(2, 1, 2, 2, '黑椒鸡排饭', '标准份', 'https://example.com/chicken-rice.jpg', 18.00, 20.00, 50, 50, 2, 4, 'active', 2)

ON DUPLICATE KEY UPDATE

  group_id = VALUES(group_id),

  product_id = VALUES(product_id),

  product_sku_id = VALUES(product_sku_id),

  product_name_snapshot = VALUES(product_name_snapshot),

  sku_name_snapshot = VALUES(sku_name_snapshot),

  cover_image_snapshot = VALUES(cover_image_snapshot),

  price_snapshot = VALUES(price_snapshot),

  original_price_snapshot = VALUES(original_price_snapshot),

  stock_total_snapshot = VALUES(stock_total_snapshot),

  stock_available_snapshot = VALUES(stock_available_snapshot),

  limit_per_user_snapshot = VALUES(limit_per_user_snapshot),

  limit_per_order_snapshot = VALUES(limit_per_order_snapshot),

  status = VALUES(status),

  sort_order = VALUES(sort_order);

