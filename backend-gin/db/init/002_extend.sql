USE selfcook;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  open_id VARCHAR(64) NOT NULL UNIQUE,
  union_id VARCHAR(64) DEFAULT '',
  nickname VARCHAR(64) NOT NULL,
  avatar_url VARCHAR(255) DEFAULT '',
  phone VARCHAR(20) DEFAULT '',
  role VARCHAR(20) NOT NULL DEFAULT 'user',
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  points_balance INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS addresses (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  contact_name VARCHAR(50) NOT NULL,
  contact_phone VARCHAR(20) NOT NULL,
  province VARCHAR(50) DEFAULT '',
  city VARCHAR(50) DEFAULT '',
  district VARCHAR(50) DEFAULT '',
  detail_address VARCHAR(255) NOT NULL,
  community_name VARCHAR(100) DEFAULT '',
  latitude DECIMAL(10,7) DEFAULT 0,
  longitude DECIMAL(10,7) DEFAULT 0,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_addresses_user_id (user_id)
);

CREATE TABLE IF NOT EXISTS inventory_logs (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  group_id BIGINT UNSIGNED NOT NULL,
  group_item_id BIGINT UNSIGNED NOT NULL,
  product_sku_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NULL,
  change_type VARCHAR(30) NOT NULL,
  change_qty INT NOT NULL,
  before_stock INT NOT NULL,
  after_stock INT NOT NULL,
  operator_id BIGINT UNSIGNED NULL,
  operator_role VARCHAR(20) DEFAULT '',
  remark VARCHAR(255) DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_inventory_logs_group_item_id (group_item_id),
  KEY idx_inventory_logs_order_id (order_id)
);

CREATE TABLE IF NOT EXISTS coupons (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  coupon_type VARCHAR(20) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  threshold_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  applicable_scope VARCHAR(20) NOT NULL DEFAULT 'all',
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  valid_from DATETIME NOT NULL,
  valid_to DATETIME NOT NULL,
  total_count INT NOT NULL DEFAULT 0,
  per_user_limit INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_coupons (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  coupon_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'unused',
  acquired_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  used_at DATETIME NULL,
  order_id BIGINT UNSIGNED NULL,
  valid_from DATETIME NOT NULL,
  valid_to DATETIME NOT NULL,
  KEY idx_user_coupons_user_id (user_id),
  KEY idx_user_coupons_coupon_id (coupon_id)
);

CREATE TABLE IF NOT EXISTS points_ledgers (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  change_value INT NOT NULL,
  balance_after INT NOT NULL,
  source_type VARCHAR(30) NOT NULL,
  source_id BIGINT UNSIGNED NULL,
  remark VARCHAR(255) DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_points_ledgers_user_id (user_id)
);

CREATE TABLE IF NOT EXISTS notifications (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NULL,
  group_id BIGINT UNSIGNED NULL,
  scene_code VARCHAR(50) NOT NULL,
  template_id VARCHAR(100) NOT NULL,
  send_status VARCHAR(20) NOT NULL DEFAULT 'pending',
  request_payload TEXT,
  response_payload TEXT,
  fail_reason VARCHAR(255) DEFAULT '',
  sent_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_notifications_user_id (user_id),
  KEY idx_notifications_order_id (order_id),
  KEY idx_notifications_group_id (group_id)
);

ALTER TABLE orders ADD COLUMN cancelled_at DATETIME NULL;
ALTER TABLE orders ADD COLUMN completed_at DATETIME NULL;

INSERT INTO users (id, open_id, nickname, phone, role, status, points_balance)
VALUES (1, 'demo-openid-1', '演示用户', '13800001111', 'user', 'active', 0)
ON DUPLICATE KEY UPDATE nickname = VALUES(nickname);

INSERT INTO addresses (id, user_id, contact_name, contact_phone, province, city, district, detail_address, community_name, is_default)
VALUES (1, 1, '演示用户', '13800001111', '广东省', '深圳市', '南山区', '科技园1号楼1202', '科技园', 1)
ON DUPLICATE KEY UPDATE detail_address = VALUES(detail_address);

INSERT INTO coupons (id, name, coupon_type, amount, threshold_amount, applicable_scope, status, valid_from, valid_to, total_count, per_user_limit)
VALUES (1, '新人满20减5', 'full_reduction', 5.00, 20.00, 'all', 'active', NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1000, 1)
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO user_coupons (id, coupon_id, user_id, status, acquired_at, valid_from, valid_to)
VALUES (1, 1, 1, 'unused', NOW(), NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY))
ON DUPLICATE KEY UPDATE status = VALUES(status);
