# 团购客饭接龙小程序数据库设计

## 1. 设计原则

- 优先满足团购活动、接龙订单、库存追踪、履约与会员营销的核心需求。
- 采用关系型数据库设计，推荐 MySQL 8.0 或 PostgreSQL 15+。
- 历史订单、商品价格、团内商品使用快照表，避免后续商品变更影响历史数据。
- 所有关键状态流转需保留操作人与操作时间。
- 对库存、订单创建、截单等高并发操作预留幂等与事务能力。

## 2. 枚举定义建议

### 2.1 用户角色

- `user` 普通用户
- `leader` 团长
- `admin` 商家管理员

### 2.2 团状态

- `draft` 草稿
- `published` 已发布
- `ongoing` 进行中
- `cutoff` 已截单
- `fulfilling` 履约中
- `completed` 已完成
- `closed` 已关闭

### 2.3 订单状态

- `pending_confirm` 待确认
- `joined` 已接龙
- `cutoff_locked` 截单锁定
- `ready_for_pickup` 待取餐
- `delivering` 配送中
- `completed` 已完成
- `cancelled` 已取消
- `closed` 已关闭

### 2.4 履约方式

- `pickup` 自提
- `delivery` 配送

### 2.5 库存流水类型

- `reserve` 占用
- `release` 释放
- `deduct` 扣减
- `adjust_increase` 增加
- `adjust_decrease` 减少

## 3. 核心表结构

## 3.1 users

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 用户 ID |
| openid | varchar(64) unique | 微信 openid |
| unionid | varchar(64) null | 微信 unionid |
| nickname | varchar(64) | 昵称 |
| avatar_url | varchar(255) | 头像 |
| phone | varchar(20) null | 手机号 |
| role | varchar(20) | user/leader/admin |
| status | varchar(20) | active/disabled |
| points_balance | int | 当前积分 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- unique(openid)
- index(role)
- index(phone)

## 3.2 addresses

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 地址 ID |
| user_id | bigint | 用户 ID |
| contact_name | varchar(50) | 联系人 |
| contact_phone | varchar(20) | 联系电话 |
| province | varchar(50) | 省 |
| city | varchar(50) | 市 |
| district | varchar(50) | 区 |
| detail_address | varchar(255) | 详细地址 |
| community_name | varchar(100) | 小区/楼宇 |
| latitude | decimal(10,7) null | 纬度 |
| longitude | decimal(10,7) null | 经度 |
| is_default | tinyint | 是否默认 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- index(user_id)

## 3.3 pickup_points

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 自提点 ID |
| name | varchar(100) | 自提点名称 |
| contact_name | varchar(50) | 联系人 |
| contact_phone | varchar(20) | 联系电话 |
| address | varchar(255) | 地址 |
| business_hours | varchar(100) | 营业时间 |
| status | varchar(20) | active/inactive |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 3.4 pickup_point_slots

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 时间段 ID |
| pickup_point_id | bigint | 自提点 ID |
| slot_name | varchar(50) | 如 11:30-12:00 |
| start_time | time | 开始时间 |
| end_time | time | 结束时间 |
| sort_order | int | 排序 |
| status | varchar(20) | active/inactive |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 3.5 products

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 商品 ID |
| name | varchar(100) | 商品名称 |
| subtitle | varchar(255) | 副标题 |
| cover_image | varchar(255) | 封面图 |
| category_name | varchar(50) | 类目 |
| description | text | 商品说明 |
| status | varchar(20) | on_sale/off_sale |
| sort_order | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 3.6 product_skus

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | SKU ID |
| product_id | bigint | 商品 ID |
| sku_name | varchar(100) | 规格名 |
| sku_code | varchar(50) | 规格编码 |
| price | decimal(10,2) | 售价 |
| original_price | decimal(10,2) | 划线价 |
| stock_total | int | 总库存 |
| stock_available | int | 可售库存 |
| limit_per_user | int | 每人限购 |
| limit_per_order | int | 每单限购 |
| weight | decimal(10,2) null | 重量 |
| status | varchar(20) | active/inactive |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- index(product_id)
- index(status)

## 3.7 group_templates

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 模板 ID |
| name | varchar(100) | 模板名称 |
| description | varchar(255) | 描述 |
| default_fulfillment_mode | varchar(20) | 默认履约方式 |
| status | varchar(20) | active/inactive |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 3.8 groups

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 团活动 ID |
| title | varchar(150) | 团标题 |
| cover_image | varchar(255) | 封面图 |
| leader_user_id | bigint | 团长 ID |
| template_id | bigint null | 模板 ID |
| status | varchar(20) | 团状态 |
| start_at | datetime | 开团时间 |
| cutoff_at | datetime | 截单时间 |
| fulfillment_mode | varchar(20) | pickup/delivery/mixed |
| allow_modify_before_cutoff | tinyint | 是否允许改单 |
| show_join_list | tinyint | 是否展示接龙明细 |
| pickup_rule_desc | varchar(255) null | 自提说明 |
| delivery_rule_desc | varchar(255) null | 配送说明 |
| group_notice | text null | 团说明 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- index(leader_user_id)
- index(status)
- index(cutoff_at)
- index(start_at)

## 3.9 group_items

说明：存储团内商品快照，避免商品信息变更影响历史团活动。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 记录 ID |
| group_id | bigint | 团活动 ID |
| product_id | bigint | 商品 ID |
| product_sku_id | bigint | SKU ID |
| product_name_snapshot | varchar(100) | 商品名快照 |
| sku_name_snapshot | varchar(100) | SKU 快照 |
| cover_image_snapshot | varchar(255) | 图片快照 |
| price_snapshot | decimal(10,2) | 售价快照 |
| original_price_snapshot | decimal(10,2) | 原价快照 |
| stock_total_snapshot | int | 初始库存快照 |
| stock_available_snapshot | int | 当前活动可售库存 |
| limit_per_user_snapshot | int | 用户限购快照 |
| limit_per_order_snapshot | int | 订单限购快照 |
| status | varchar(20) | active/inactive/sold_out |
| sort_order | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- unique(group_id, product_sku_id)
- index(group_id)

## 3.10 group_pickup_points

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 记录 ID |
| group_id | bigint | 团活动 ID |
| pickup_point_id | bigint | 自提点 ID |
| is_default | tinyint | 是否默认 |
| created_at | datetime | 创建时间 |

## 3.11 group_delivery_regions

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 记录 ID |
| group_id | bigint | 团活动 ID |
| region_name | varchar(100) | 区域名称 |
| address_keyword | varchar(100) null | 地址关键字 |
| delivery_fee | decimal(10,2) | 运费 |
| min_amount | decimal(10,2) | 起送金额 |
| created_at | datetime | 创建时间 |

## 3.12 orders

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 订单 ID |
| order_no | varchar(40) unique | 订单号 |
| group_id | bigint | 团活动 ID |
| user_id | bigint | 用户 ID |
| status | varchar(30) | 订单状态 |
| fulfillment_mode | varchar(20) | pickup/delivery |
| contact_name | varchar(50) | 联系人 |
| contact_phone | varchar(20) | 联系电话 |
| pickup_point_id | bigint null | 自提点 |
| pickup_slot_id | bigint null | 自提时间段 |
| address_id | bigint null | 地址 ID |
| delivery_address_snapshot | varchar(255) null | 配送地址快照 |
| coupon_id | bigint null | 使用优惠券 ID |
| goods_amount | decimal(10,2) | 商品总价 |
| discount_amount | decimal(10,2) | 优惠金额 |
| delivery_fee | decimal(10,2) | 运费 |
| payable_amount | decimal(10,2) | 应付金额 |
| paid_amount | decimal(10,2) | 实付金额 |
| remark | varchar(255) null | 用户备注 |
| cutoff_at_snapshot | datetime | 截单时间快照 |
| source | varchar(30) | miniprogram |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| cancelled_at | datetime null | 取消时间 |
| completed_at | datetime null | 完成时间 |

索引建议：

- unique(order_no)
- index(group_id)
- index(user_id)
- index(status)
- index(created_at)

## 3.13 order_items

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 订单项 ID |
| order_id | bigint | 订单 ID |
| group_item_id | bigint | 团商品快照 ID |
| product_id | bigint | 商品 ID |
| product_sku_id | bigint | SKU ID |
| product_name_snapshot | varchar(100) | 商品名称快照 |
| sku_name_snapshot | varchar(100) | SKU 名称快照 |
| unit_price_snapshot | decimal(10,2) | 单价快照 |
| quantity | int | 数量 |
| subtotal_amount | decimal(10,2) | 小计 |
| taste_remark | varchar(100) null | 口味备注 |
| item_status | varchar(20) | normal/out_of_stock/replaced |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- index(order_id)
- index(product_sku_id)

## 3.14 inventory_logs

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 日志 ID |
| group_id | bigint | 团活动 ID |
| group_item_id | bigint | 团商品 ID |
| product_sku_id | bigint | SKU ID |
| order_id | bigint null | 关联订单 |
| change_type | varchar(30) | reserve/release/deduct/adjust |
| change_qty | int | 变更数量，正负均可 |
| before_stock | int | 变更前库存 |
| after_stock | int | 变更后库存 |
| operator_id | bigint null | 操作人 |
| operator_role | varchar(20) null | 操作角色 |
| remark | varchar(255) null | 备注 |
| created_at | datetime | 创建时间 |

索引建议：

- index(group_item_id)
- index(order_id)
- index(created_at)

## 3.15 coupons

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 优惠券 ID |
| name | varchar(100) | 优惠券名称 |
| coupon_type | varchar(20) | full_reduction/discount/fixed |
| amount | decimal(10,2) | 面额/优惠值 |
| threshold_amount | decimal(10,2) | 使用门槛 |
| applicable_scope | varchar(20) | all/group/product |
| status | varchar(20) | active/inactive |
| valid_from | datetime | 生效时间 |
| valid_to | datetime | 失效时间 |
| total_count | int | 总发放量 |
| per_user_limit | int | 每人限领 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 3.16 user_coupons

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 用户券 ID |
| coupon_id | bigint | 优惠券 ID |
| user_id | bigint | 用户 ID |
| status | varchar(20) | unused/used/expired |
| acquired_at | datetime | 领取时间 |
| used_at | datetime null | 使用时间 |
| order_id | bigint null | 使用订单 |
| valid_from | datetime | 生效时间 |
| valid_to | datetime | 失效时间 |

索引建议：

- index(user_id, status)
- index(coupon_id)

## 3.17 points_ledger

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 积分流水 ID |
| user_id | bigint | 用户 ID |
| change_value | int | 积分变化 |
| balance_after | int | 变更后余额 |
| source_type | varchar(30) | order_reward/manual_adjust/exchange |
| source_id | bigint null | 来源记录 ID |
| remark | varchar(255) null | 备注 |
| created_at | datetime | 创建时间 |

索引建议：

- index(user_id)
- index(source_type, source_id)

## 3.18 notifications

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 通知记录 ID |
| user_id | bigint | 用户 ID |
| order_id | bigint null | 关联订单 |
| group_id | bigint null | 关联团活动 |
| scene_code | varchar(50) | 通知场景 |
| template_id | varchar(100) | 模板 ID |
| send_status | varchar(20) | pending/sent/failed |
| request_payload | json | 请求参数 |
| response_payload | json null | 响应结果 |
| fail_reason | varchar(255) null | 失败原因 |
| sent_at | datetime null | 发送时间 |
| created_at | datetime | 创建时间 |

## 3.19 idempotent_requests

用于下单、截单、库存调整等操作幂等控制。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigint PK | 主键 |
| biz_type | varchar(50) | 业务类型 |
| idem_key | varchar(100) | 幂等键 |
| request_hash | varchar(128) | 请求摘要 |
| response_snapshot | json null | 响应快照 |
| status | varchar(20) | processing/success/failed |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

索引建议：

- unique(biz_type, idem_key)

## 4. 状态机与事务设计

## 4.1 订单状态机

```text
pending_confirm -> joined -> cutoff_locked -> ready_for_pickup -> completed
pending_confirm -> joined -> cutoff_locked -> delivering -> completed
pending_confirm -> cancelled
joined -> cancelled
pending_confirm/joined -> closed
```

状态含义：

- `pending_confirm`：订单创建中，适合保留为系统短暂处理中间态。
- `joined`：订单已有效接龙，库存已占用。
- `cutoff_locked`：团已截单，订单不允许用户修改。
- `ready_for_pickup`：已备餐，可自提。
- `delivering`：已出配送。
- `completed`：用户完成履约。
- `cancelled`：用户或后台取消。
- `closed`：系统关闭，如异常超时关闭。

## 4.2 库存状态机

```text
available -> reserve -> release
reserve -> deduct
available -> adjust_increase/adjust_decrease
```

建议规则：

- 创建订单时扣减 `group_items.stock_available_snapshot`，并写入 `reserve` 日志。
- 取消订单时回补活动库存，并写入 `release` 日志。
- 截单后是否额外写 `deduct` 取决于库存管理策略；MVP 可将“占用即销售”视作已锁定库存。
- 人工调整库存只能通过后台入口执行，必须记录操作人。

## 4.3 创建订单事务

建议在单事务中完成：

1. 校验团状态为 `ongoing`。
2. 校验当前时间未超过 `cutoff_at`。
3. 查询团商品与库存并加行锁。
4. 校验限购、库存、券使用条件。
5. 写入 `orders`。
6. 写入 `order_items`。
7. 扣减 `group_items.stock_available_snapshot`。
8. 写入 `inventory_logs`。
9. 标记 `user_coupons` 为已使用（如使用券）。
10. 提交事务。

## 4.4 截单事务

建议在事务中完成：

1. 校验团状态为 `ongoing`。
2. 更新团状态为 `cutoff`。
3. 批量将 `joined` 订单更新为 `cutoff_locked`。
4. 写入操作日志。
5. 异步触发汇总缓存刷新与通知。

## 5. 并发与幂等策略

### 5.1 下单并发控制

- 使用数据库事务 + 行级锁控制 `group_items` 库存。
- 每个 SKU 下单前使用 `SELECT ... FOR UPDATE`。
- 若一次下单包含多个 SKU，则固定顺序加锁，避免死锁。

### 5.2 幂等策略

适用接口：

- 创建订单
- 截单
- 库存调整
- 发券
- 手工完成订单

建议：

- 客户端生成 `X-Idempotency-Key`
- 服务端结合 `biz_type + idem_key` 做唯一约束
- 成功后返回相同响应快照

### 5.3 防止重复通知

- `notifications` 表保存唯一业务场景记录
- 可约束 `(user_id, order_id, scene_code)` 或 `(user_id, group_id, scene_code)`

## 6. 汇总查询建议

### 6.1 商品维度汇总

按 `group_id + product_sku_id` 聚合 `order_items.quantity`

### 6.2 自提点汇总

按 `orders.pickup_point_id` 分组统计订单数与商品份数

### 6.3 配送区域汇总

按区域规则或地址关键字进行聚合展示

## 7. 审计与扩展建议

- 建议补充 `operation_logs` 表记录后台关键操作。
- 若后续接入支付，可新增 `payments` 与 `refunds` 表。
- 若后续支持团长佣金，可新增 `leader_settlements` 表。
- 若后续支持司机配送，可新增 `delivery_tasks` 表。

## 8. MVP 最小落表范围

MVP 必须具备：

- `users`
- `addresses`
- `pickup_points`
- `pickup_point_slots`
- `products`
- `product_skus`
- `groups`
- `group_items`
- `orders`
- `order_items`
- `inventory_logs`
- `coupons`
- `user_coupons`
- `points_ledger`
- `notifications`
- `idempotent_requests`
