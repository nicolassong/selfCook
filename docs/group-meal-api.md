# 团购客饭接龙小程序 API 文档

## 1. 约定

### 1.1 基础信息

- Base URL: `/api/v1`
- 数据格式：`application/json`
- 认证方式：小程序登录后服务端签发 token，Header 使用 `Authorization: Bearer <token>`
- 时间格式：ISO 8601，例如 `2025-01-01T10:00:00+08:00`

### 1.2 通用响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 1.3 通用错误码

| code | 含义 |
|---|---|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40002 | 登录失效 |
| 40003 | 无权限 |
| 40004 | 资源不存在 |
| 40010 | 活动已截单 |
| 40011 | 库存不足 |
| 40012 | 超出限购 |
| 40013 | 优惠券不可用 |
| 40014 | 履约信息不合法 |
| 40015 | 订单状态不可操作 |
| 40016 | 重复请求 |
| 50000 | 系统异常 |

## 2. 认证接口

## 2.1 小程序登录

### POST `/auth/wechat/login`

请求体：

```json
{
  "code": "wx-login-code",
  "nickname": "小王",
  "avatarUrl": "https://example.com/avatar.png"
}
```

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 1001,
      "nickname": "小王",
      "avatarUrl": "https://example.com/avatar.png",
      "role": "user",
      "phone": null
    }
  }
}
```

## 2.2 绑定手机号

### POST `/auth/bind-phone`

请求体：

```json
{
  "phoneCode": "wechat-phone-code"
}
```

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "phone": "13800001111"
  }
}
```

## 3. 用户端接口

## 3.1 获取团购列表

### GET `/groups`

查询参数：

| 参数 | 必填 | 说明 |
|---|---|---|
| status | 否 | ongoing/cutoff/completed |
| keyword | 否 | 搜索关键词 |
| page | 否 | 页码 |
| pageSize | 否 | 每页数量 |
| mine | 否 | 是否仅我的团 |

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 2001,
        "title": "3号楼午餐接龙",
        "leaderName": "张团长",
        "cutoffAt": "2025-01-01T10:30:00+08:00",
        "status": "ongoing",
        "fulfillmentMode": "mixed",
        "joinedUserCount": 38,
        "itemCount": 12,
        "coverImage": "https://example.com/group.jpg"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1
    }
  }
}
```

## 3.2 获取团详情

### GET `/groups/{groupId}`

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 2001,
    "title": "3号楼午餐接龙",
    "status": "ongoing",
    "leader": {
      "id": 101,
      "nickname": "张团长"
    },
    "cutoffAt": "2025-01-01T10:30:00+08:00",
    "fulfillmentMode": "mixed",
    "allowModifyBeforeCutoff": true,
    "showJoinList": true,
    "pickupPoints": [
      {
        "id": 1,
        "name": "A栋一楼前台"
      }
    ],
    "items": [
      {
        "groupItemId": 9001,
        "productId": 301,
        "productSkuId": 601,
        "productName": "番茄牛腩饭",
        "skuName": "大份",
        "price": 22,
        "originalPrice": 25,
        "stockAvailable": 18,
        "limitPerUser": 2,
        "coverImage": "https://example.com/a.jpg",
        "status": "active"
      }
    ],
    "joinPreview": [
      {
        "nicknameMasked": "小**",
        "summary": "番茄牛腩饭 x1"
      }
    ]
  }
}
```

## 3.3 创建订单

### POST `/orders`

请求 Header：

- `X-Idempotency-Key: create-order-uuid`

请求体：

```json
{
  "groupId": 2001,
  "fulfillmentMode": "pickup",
  "pickupPointId": 1,
  "pickupSlotId": 11,
  "contactName": "王女士",
  "contactPhone": "13800001111",
  "couponId": 5001,
  "remark": "少辣",
  "items": [
    {
      "groupItemId": 9001,
      "quantity": 2,
      "tasteRemark": "一份不要香菜"
    }
  ]
}
```

成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "orderId": 70001,
    "orderNo": "GM202501010001",
    "status": "joined",
    "goodsAmount": 44,
    "discountAmount": 5,
    "deliveryFee": 0,
    "payableAmount": 39
  }
}
```

库存不足响应：

```json
{
  "code": 40011,
  "message": "库存不足",
  "data": {
    "groupItemId": 9001,
    "availableQty": 1
  }
}
```

## 3.4 修改订单

### PUT `/orders/{orderId}`

说明：仅未截单且团活动允许改单时可用。

请求体：

```json
{
  "pickupPointId": 2,
  "pickupSlotId": 12,
  "remark": "改成微辣",
  "items": [
    {
      "groupItemId": 9001,
      "quantity": 1,
      "tasteRemark": "不要葱"
    }
  ]
}
```

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "orderId": 70001,
    "status": "joined",
    "payableAmount": 22
  }
}
```

## 3.5 获取订单详情

### GET `/orders/{orderId}`

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 70001,
    "orderNo": "GM202501010001",
    "status": "ready_for_pickup",
    "group": {
      "id": 2001,
      "title": "3号楼午餐接龙"
    },
    "fulfillmentMode": "pickup",
    "pickupInfo": {
      "pickupPointName": "A栋一楼前台",
      "pickupSlotName": "11:30-12:00"
    },
    "items": [
      {
        "productName": "番茄牛腩饭",
        "skuName": "大份",
        "unitPrice": 22,
        "quantity": 2,
        "subtotalAmount": 44,
        "tasteRemark": "一份不要香菜"
      }
    ],
    "amount": {
      "goodsAmount": 44,
      "discountAmount": 5,
      "deliveryFee": 0,
      "paidAmount": 39
    },
    "timeline": [
      {
        "status": "joined",
        "time": "2025-01-01T09:20:00+08:00",
        "desc": "已接龙成功"
      },
      {
        "status": "cutoff_locked",
        "time": "2025-01-01T10:30:00+08:00",
        "desc": "已截单"
      }
    ]
  }
}
```

## 3.6 我的订单列表

### GET `/me/orders`

查询参数：

| 参数 | 必填 | 说明 |
|---|---|---|
| status | 否 | joined/completed/cancelled |
| page | 否 | 页码 |
| pageSize | 否 | 每页 |

## 3.7 取消订单

### POST `/orders/{orderId}/cancel`

请求体：

```json
{
  "reason": "临时有事不吃了"
}
```

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "orderId": 70001,
    "status": "cancelled"
  }
}
```

## 3.8 地址管理

### GET `/me/addresses`
### POST `/me/addresses`
### PUT `/me/addresses/{addressId}`
### DELETE `/me/addresses/{addressId}`

新增地址请求示例：

```json
{
  "contactName": "王女士",
  "contactPhone": "13800001111",
  "province": "广东省",
  "city": "深圳市",
  "district": "南山区",
  "detailAddress": "科技园1号楼1202",
  "communityName": "科技园",
  "isDefault": true
}
```

## 3.9 用户优惠券

### GET `/me/coupons`

### POST `/coupons/{couponId}/claim`

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "userCouponId": 10001,
    "status": "unused"
  }
}
```

## 3.10 用户积分

### GET `/me/points`

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "balance": 120,
    "records": [
      {
        "id": 1,
        "changeValue": 20,
        "sourceType": "order_reward",
        "remark": "订单完成赠送积分",
        "createdAt": "2025-01-01T12:30:00+08:00"
      }
    ]
  }
}
```

## 3.11 订阅消息授权上报

### POST `/subscribe`

请求体：

```json
{
  "sceneCode": "order_created",
  "templateId": "tmpl_001",
  "accepted": true
}
```

## 4. 团长端接口

## 4.1 发起团

### POST `/leader/groups`

请求体：

```json
{
  "title": "3号楼午餐接龙",
  "coverImage": "https://example.com/g.jpg",
  "templateId": 1,
  "startAt": "2025-01-01T08:00:00+08:00",
  "cutoffAt": "2025-01-01T10:30:00+08:00",
  "fulfillmentMode": "mixed",
  "allowModifyBeforeCutoff": true,
  "showJoinList": true,
  "pickupPointIds": [1, 2],
  "items": [
    {
      "productSkuId": 601,
      "stockTotal": 30,
      "price": 22
    }
  ]
}
```

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "groupId": 2001,
    "status": "ongoing"
  }
}
```

## 4.2 我的团列表

### GET `/leader/groups`

## 4.3 团管理详情

### GET `/leader/groups/{groupId}`

## 4.4 团订单列表

### GET `/leader/groups/{groupId}/orders`

查询参数：

| 参数 | 必填 | 说明 |
|---|---|---|
| status | 否 | 订单状态 |
| fulfillmentMode | 否 | pickup/delivery |
| keyword | 否 | 用户昵称/手机号/订单号 |

## 4.5 手动截单

### POST `/leader/groups/{groupId}/cutoff`

请求 Header：

- `X-Idempotency-Key: cutoff-group-uuid`

请求体：

```json
{
  "reason": "按计划截单"
}
```

响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "groupId": 2001,
    "status": "cutoff",
    "cutoffAt": "2025-01-01T10:30:00+08:00"
  }
}
```

## 4.6 获取汇总

### GET `/leader/groups/{groupId}/summary`

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "bySku": [
      {
        "productName": "番茄牛腩饭",
        "skuName": "大份",
        "totalQty": 28
      }
    ],
    "byPickupPoint": [
      {
        "pickupPointName": "A栋一楼前台",
        "orderCount": 18,
        "totalQty": 26
      }
    ],
    "byDeliveryRegion": [
      {
        "regionName": "科技园南区",
        "orderCount": 6,
        "totalQty": 8
      }
    ]
  }
}
```

## 4.7 发送团通知

### POST `/leader/groups/{groupId}/notify`

请求体：

```json
{
  "sceneCode": "group_cutoff",
  "target": "joined_users"
}
```

## 5. 后台管理接口

## 5.1 商品管理

### GET `/admin/products`
### POST `/admin/products`
### PUT `/admin/products/{productId}`
### POST `/admin/products/{productId}/toggle-status`

新增商品示例：

```json
{
  "name": "番茄牛腩饭",
  "subtitle": "大块牛腩，酸甜开胃",
  "coverImage": "https://example.com/a.jpg",
  "categoryName": "热销套餐",
  "description": "含米饭、配菜、汤",
  "skus": [
    {
      "skuName": "大份",
      "skuCode": "BEEF_RICE_L",
      "price": 22,
      "originalPrice": 25,
      "stockTotal": 100,
      "limitPerUser": 2,
      "limitPerOrder": 4
    }
  ]
}
```

## 5.2 库存调整

### POST `/admin/inventory/adjust`

请求体：

```json
{
  "groupItemId": 9001,
  "changeQty": 10,
  "reason": "午间补货"
}
```

## 5.3 优惠券管理

### GET `/admin/coupons`
### POST `/admin/coupons`
### POST `/admin/coupons/{couponId}/issue`

发券请求示例：

```json
{
  "userIds": [1001, 1002],
  "reason": "新用户拉新活动"
}
```

## 5.4 自提点管理

### GET `/admin/pickup-points`
### POST `/admin/pickup-points`
### PUT `/admin/pickup-points/{pickupPointId}`

## 5.5 订单履约操作

### POST `/admin/orders/{orderId}/ready-for-pickup`
### POST `/admin/orders/{orderId}/start-delivery`
### POST `/admin/orders/{orderId}/complete`
### POST `/admin/orders/{orderId}/manual-adjust`

订单完成请求示例：

```json
{
  "remark": "已核销完成"
}
```

## 5.6 积分调整

### POST `/admin/points/adjust`

请求体：

```json
{
  "userId": 1001,
  "changeValue": 20,
  "remark": "活动补偿积分"
}
```

## 5.7 消息模板配置

### GET `/admin/message-templates`
### POST `/admin/message-templates`

请求体：

```json
{
  "sceneCode": "group_cutoff",
  "templateId": "tmpl_002",
  "templateName": "截单通知",
  "status": "active"
}
```

## 6. 关键业务校验规则

### 6.1 创建订单校验

- 团状态必须为 `ongoing`
- 当前时间必须早于 `cutoffAt`
- 每个订单项数量 > 0
- 数量不能超过活动库存
- 数量不能超过限购
- 自提/配送信息必须完整
- 优惠券必须属于本人且未使用

### 6.2 修改订单校验

- 团允许改单
- 订单状态为 `joined`
- 未截单
- 改单后重新校验库存、限购与优惠券

### 6.3 截单校验

- 团长或管理员权限
- 团状态必须为 `ongoing`
- 幂等键不可重复产生新动作

## 7. 建议异步事件

- `OrderCreated`
- `OrderUpdated`
- `OrderCancelled`
- `GroupCutoff`
- `OrderReadyForPickup`
- `OrderDelivering`
- `OrderCompleted`
- `CouponIssued`
- `PointsGranted`

## 8. OpenAPI 落地建议

可后续将本文拆分为：

- `auth.yaml`
- `group.yaml`
- `order.yaml`
- `leader.yaml`
- `admin.yaml`

并统一放入 `backend-gin/docs/openapi/` 目录中生成接口文档站点。
