USE selfcook;
SET NAMES utf8mb4;

INSERT INTO products (id, name, subtitle, cover_image, category_name, description, status, sort_order)
VALUES
(101, '红烧鸡腿饭', '经典下饭，适合工作餐', '', '单人套餐', '鸡腿、米饭、时蔬，口味咸香，出餐稳定。', 'on_sale', 101),
(102, '黑椒鸡排饭', '黑椒香浓，方便分餐', '', '单人套餐', '鸡排、米饭、配菜，适合办公室午餐。', 'on_sale', 102),
(103, '番茄牛腩饭', '酸甜开胃，牛腩软烂', '', '单人套餐', '番茄牛腩、米饭、配菜，适合不想吃辣的用户。', 'on_sale', 103),
(104, '鱼香肉丝饭', '微辣下饭，人气套餐', '', '单人套餐', '鱼香肉丝、米饭、配菜，偏微辣。', 'on_sale', 104),
(105, '宫保鸡丁饭', '轻微辣，坚果香', '', '单人套餐', '宫保鸡丁、米饭、配菜，口味浓郁。', 'on_sale', 105),
(106, '香菇滑鸡饭', '清爽少油，老人小孩友好', '', '单人套餐', '香菇滑鸡、米饭、时蔬，口味清淡。', 'on_sale', 106),
(107, '咖喱土豆鸡饭', '微甜咖喱，饱腹感强', '', '单人套餐', '咖喱鸡块、土豆、米饭，适合午餐。', 'on_sale', 107),
(108, '梅菜扣肉饭', '浓香下饭，适合重口味', '', '单人套餐', '梅菜扣肉、米饭、青菜，油脂较高。', 'on_sale', 108),
(109, '清炒时蔬盒', '素食补充，清爽解腻', '', '配菜小食', '当季绿叶菜，适合搭配多人套餐。', 'on_sale', 109),
(110, '蒜蓉西兰花', '低脂蔬菜，适合控油', '', '配菜小食', '西兰花蒜蓉清炒，适合清淡需求。', 'on_sale', 110),
(111, '番茄炒蛋', '家常不辣，儿童友好', '', '家庭分享菜', '番茄炒蛋，适合老人小孩和不吃辣人群。', 'on_sale', 111),
(112, '土豆炖牛腩', '家庭分享装，适合多人', '', '家庭分享菜', '牛腩、土豆、胡萝卜，适合3-4人分享。', 'on_sale', 112),
(113, '糖醋里脊', '酸甜口，聚餐人气菜', '', '家庭分享菜', '酸甜口味，适合多人分食。', 'on_sale', 113),
(114, '小炒黄牛肉', '香辣下饭，适合重口味', '', '家庭分享菜', '黄牛肉、小米椒，偏辣。', 'on_sale', 114),
(115, '麻婆豆腐', '川味热菜，经济实惠', '', '家庭分享菜', '豆腐、肉末、豆瓣酱，麻辣口。', 'on_sale', 115),
(116, '菌菇鸡汤', '清淡热汤，适合搭配', '', '汤品饮品', '鸡肉、菌菇，清淡暖胃。', 'on_sale', 116),
(117, '冬瓜排骨汤', '少油清汤，适合家庭餐', '', '汤品饮品', '冬瓜、排骨，清爽不腻。', 'on_sale', 117),
(118, '酸梅汤', '解腻饮品，适合多人', '', '汤品饮品', '冷饮，适合搭配油香菜品。', 'on_sale', 118),
(119, '儿童虾仁蛋炒饭', '少盐少油，儿童友好', '', '儿童轻食', '虾仁、鸡蛋、米饭，少油少盐。', 'on_sale', 119),
(120, '低脂鸡胸沙拉', '轻食控卡，适合晚餐', '', '轻食套餐', '鸡胸肉、蔬菜、鸡蛋，低脂清爽。', 'on_sale', 120)
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
(101, 101, '标准份', 'RICE_CHICKEN_LEG_M', 26.00, 30.00, 100, 80, 0, 10, 'active'),
(102, 102, '标准份', 'RICE_BLACK_PEPPER_CHICKEN_M', 24.00, 28.00, 100, 75, 0, 10, 'active'),
(103, 103, '标准份', 'RICE_TOMATO_BEEF_M', 32.00, 36.00, 80, 50, 0, 8, 'active'),
(104, 104, '标准份', 'RICE_YUXIANG_PORK_M', 25.00, 29.00, 80, 55, 0, 8, 'active'),
(105, 105, '标准份', 'RICE_KUNGPAO_CHICKEN_M', 25.00, 29.00, 80, 55, 0, 8, 'active'),
(106, 106, '标准份', 'RICE_MUSHROOM_CHICKEN_M', 26.00, 30.00, 90, 65, 0, 8, 'active'),
(107, 107, '标准份', 'RICE_CURRY_CHICKEN_M', 24.00, 28.00, 90, 60, 0, 8, 'active'),
(108, 108, '标准份', 'RICE_PORK_MUSTARD_M', 29.00, 34.00, 60, 35, 0, 6, 'active'),
(109, 109, '分享份', 'VEG_SEASONAL_SHARE', 18.00, 22.00, 60, 40, 0, 6, 'active'),
(110, 110, '分享份', 'VEG_BROCCOLI_SHARE', 20.00, 24.00, 60, 40, 0, 6, 'active'),
(111, 111, '分享份', 'EGG_TOMATO_SHARE', 22.00, 26.00, 80, 60, 0, 6, 'active'),
(112, 112, '家庭份', 'BEEF_POTATO_FAMILY', 68.00, 78.00, 40, 28, 0, 4, 'active'),
(113, 113, '分享份', 'PORK_SWEET_SOUR_SHARE', 48.00, 56.00, 50, 32, 0, 5, 'active'),
(114, 114, '分享份', 'BEEF_STIR_FRIED_SPICY_SHARE', 58.00, 68.00, 40, 26, 0, 4, 'active'),
(115, 115, '分享份', 'TOFU_MAPO_SHARE', 24.00, 28.00, 70, 45, 0, 6, 'active'),
(116, 116, '2-3人份', 'SOUP_MUSHROOM_CHICKEN', 32.00, 38.00, 50, 30, 0, 5, 'active'),
(117, 117, '2-3人份', 'SOUP_WINTER_MELON_RIB', 35.00, 42.00, 50, 30, 0, 5, 'active'),
(118, 118, '1L装', 'DRINK_SOUR_PLUM_1L', 16.00, 20.00, 80, 60, 0, 8, 'active'),
(119, 119, '儿童份', 'KIDS_SHRIMP_FRIED_RICE', 22.00, 26.00, 70, 45, 0, 6, 'active'),
(120, 120, '轻食份', 'SALAD_CHICKEN_BREAST', 28.00, 32.00, 70, 45, 0, 6, 'active')
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
(101, '办公室午餐快取团', '', 1, 'ongoing', NOW(), DATE_ADD(NOW(), INTERVAL 8 HOUR), 'mixed', 1, 1, '11:30-12:30 到前台自提', '科技园3公里内可配送', '适合办公室拼餐，优先推荐出餐快、方便分餐的套餐。'),
(102, '家庭晚餐分享团', '', 1, 'ongoing', NOW(), DATE_ADD(NOW(), INTERVAL 12 HOUR), 'mixed', 1, 1, '17:30-19:00 社区门口自提', '周边社区可配送', '适合3-5人家庭晚餐，包含分享菜、汤品和儿童友好菜。'),
(103, '轻食少油健康团', '', 1, 'ongoing', NOW(), DATE_ADD(NOW(), INTERVAL 10 HOUR), 'mixed', 1, 1, '12:00-13:00 健身房门口自提', '科技园和社区可配送', '适合少油、清淡、控卡和不吃辣的人群。')
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
(101, 101, 101, 101, '红烧鸡腿饭', '标准份', '', 26.00, 30.00, 80, 80, 0, 10, 'active', 1),
(102, 101, 102, 102, '黑椒鸡排饭', '标准份', '', 24.00, 28.00, 75, 75, 0, 10, 'active', 2),
(103, 101, 103, 103, '番茄牛腩饭', '标准份', '', 32.00, 36.00, 50, 50, 0, 8, 'active', 3),
(104, 101, 104, 104, '鱼香肉丝饭', '标准份', '', 25.00, 29.00, 55, 55, 0, 8, 'active', 4),
(105, 101, 105, 105, '宫保鸡丁饭', '标准份', '', 25.00, 29.00, 55, 55, 0, 8, 'active', 5),
(106, 101, 106, 106, '香菇滑鸡饭', '标准份', '', 26.00, 30.00, 65, 65, 0, 8, 'active', 6),
(107, 101, 107, 107, '咖喱土豆鸡饭', '标准份', '', 24.00, 28.00, 60, 60, 0, 8, 'active', 7),
(108, 101, 118, 118, '酸梅汤', '1L装', '', 16.00, 20.00, 60, 60, 0, 8, 'active', 8),
(109, 102, 111, 111, '番茄炒蛋', '分享份', '', 22.00, 26.00, 60, 60, 0, 6, 'active', 1),
(110, 102, 112, 112, '土豆炖牛腩', '家庭份', '', 68.00, 78.00, 28, 28, 0, 4, 'active', 2),
(111, 102, 113, 113, '糖醋里脊', '分享份', '', 48.00, 56.00, 32, 32, 0, 5, 'active', 3),
(112, 102, 114, 114, '小炒黄牛肉', '分享份', '', 58.00, 68.00, 26, 26, 0, 4, 'active', 4),
(113, 102, 115, 115, '麻婆豆腐', '分享份', '', 24.00, 28.00, 45, 45, 0, 6, 'active', 5),
(114, 102, 116, 116, '菌菇鸡汤', '2-3人份', '', 32.00, 38.00, 30, 30, 0, 5, 'active', 6),
(115, 102, 117, 117, '冬瓜排骨汤', '2-3人份', '', 35.00, 42.00, 30, 30, 0, 5, 'active', 7),
(116, 102, 119, 119, '儿童虾仁蛋炒饭', '儿童份', '', 22.00, 26.00, 45, 45, 0, 6, 'active', 8),
(117, 103, 106, 106, '香菇滑鸡饭', '标准份', '', 26.00, 30.00, 65, 65, 0, 8, 'active', 1),
(118, 103, 109, 109, '清炒时蔬盒', '分享份', '', 18.00, 22.00, 40, 40, 0, 6, 'active', 2),
(119, 103, 110, 110, '蒜蓉西兰花', '分享份', '', 20.00, 24.00, 40, 40, 0, 6, 'active', 3),
(120, 103, 116, 116, '菌菇鸡汤', '2-3人份', '', 32.00, 38.00, 30, 30, 0, 5, 'active', 4),
(121, 103, 117, 117, '冬瓜排骨汤', '2-3人份', '', 35.00, 42.00, 30, 30, 0, 5, 'active', 5),
(122, 103, 119, 119, '儿童虾仁蛋炒饭', '儿童份', '', 22.00, 26.00, 45, 45, 0, 6, 'active', 6),
(123, 103, 120, 120, '低脂鸡胸沙拉', '轻食份', '', 28.00, 32.00, 45, 45, 0, 6, 'active', 7)
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
