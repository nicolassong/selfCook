USE selfcook;

CREATE TABLE IF NOT EXISTS daily_menus (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  menu_date DATE NOT NULL,
  title VARCHAR(100) NOT NULL DEFAULT '',
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  remark VARCHAR(255) DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_daily_menus_menu_date (menu_date)
);

CREATE TABLE IF NOT EXISTS daily_menu_items (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  daily_menu_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NOT NULL,
  product_sku_id BIGINT UNSIGNED NOT NULL,
  stock_total INT NOT NULL DEFAULT 0,
  stock_available INT NOT NULL DEFAULT 0,
  price DECIMAL(10,2) NOT NULL,
  original_price DECIMAL(10,2) NOT NULL DEFAULT 0,
  limit_per_user INT NOT NULL DEFAULT 0,
  limit_per_order INT NOT NULL DEFAULT 0,
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_daily_menu_items_menu_id (daily_menu_id),
  KEY idx_daily_menu_items_product_sku_id (product_sku_id),
  CONSTRAINT fk_daily_menu_items_menu_id FOREIGN KEY (daily_menu_id) REFERENCES daily_menus(id)
    ON DELETE CASCADE ON UPDATE CASCADE
);
