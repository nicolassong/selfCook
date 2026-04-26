# selfCook

本目录当前包含以下核心项目：

- `backend-gin/`：Gin 后端（REST API + MySQL + Docker）
- `admin-lite/`：轻量 React 管理后台（Vite + Ant Design）
- `docs/`：产品、数据库、API、通知、后台设计文档

说明：现有 `admin-web/` 目录内容与本项目后台不匹配，当前不作为本项目管理台使用；实际联调请使用 `admin-lite/`。

## 一键 Docker 启动

在根目录执行：

```bash
docker compose up --build
```

启动后：

- 后端：`http://localhost:8080/api/ping`
- 管理后台：`http://localhost:5173`
- MySQL：`localhost:3306`
  - database: `selfcook`
  - user: `root`
  - password: `root`

## 本地启动

### backend-gin

```bash
cd backend-gin
go mod tidy
go run ./cmd/server
```

### admin-lite

```bash
cd admin-lite
npm install
npm run dev
```

## 当前已实现能力

### 后端

- 团购列表/详情
- 商品列表
- 自提点列表
- 下单/取消订单/我的订单
- 团长截单/汇总
- 后台订单履约流转
- 后台商品创建/团创建/订单列表/商品列表/团列表
- MySQL 初始化脚本与 Docker 部署

### 管理后台

- 仪表盘
- 商品管理列表
- 团活动管理列表
- 订单管理列表
- 对接后端基础管理 API

## 设计文档

- `docs/group-meal-prd.md`
- `docs/group-meal-database.md`
- `docs/group-meal-api.md`
- `docs/group-meal-notifications.md`
- `docs/group-meal-admin.md`
