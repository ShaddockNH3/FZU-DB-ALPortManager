# FZU-DB-ALPortManager
福州大学2025年下半学期数据库实践大作业 - 碧蓝航线港口管理系统

## 项目简介
基于 Go + Vue3 + **openGauss** 的碧蓝航线港口管理系统，模拟游戏中的港口管理功能，实现舰船信息管理、装备系统、星级管理等核心功能。

### 主要功能
- 🚢 **舰船管理** - 舰船信息的增删改查，支持多条件筛选
- ⭐ **星级系统** - 舰船升星/降星，满星后解锁兵装栏
- ⚔️ **装备系统** - 277种装备（主炮、鱼雷、防空炮、舰载机、设备、兵装）
- 🎯 **装备限制** - 根据舰种和装备类型进行智能过滤
- 📊 **数据统计** - 按阵营、稀有度、舰种统计舰船数量
- 📄 **分页查询** - 支持大数据量的分页展示
- 🎨 **现代化UI** - 响应式设计，Element Plus 组件库

### 技术特点
- 前后端分离架构
- RESTful API 设计
- GORM ORM + Code Generation
- **云端 openGauss 数据库**，充分体验企业级数据库特性
- Docker 容器化部署

## 技术栈
- **后端**: Go 1.23 + Hertz (CloudWeGo) + GORM 1.24 + Gen 0.3.21
- **前端**: Vue3 + Vite + Element Plus + Pinia + Axios
- **数据库（云端部署）**: **openGauss 5.0.1 （基于华为云ECS鲲鹏ARM架构手动部署）** 
- **数据库（本地开发）**: PostgreSQL 15
- **容器化**: Docker + Docker Compose （用于本地环境）

## 核心挑战：手动部署并迁移至 openGauss
作为本项目的核心实践内容，我们完成了从本地 PostgreSQL 到华为云 openGauss 的完整迁移。这不仅仅是更换一个连接字符串，而是一次涉及底层架构、网络安全、驱动兼容性的深度实践。

### 迁移过程与挑战：
1.  **环境搭建**: 在华为云 ECS (鲲鹏 ARM 架构) 上，从零开始手动编译安装 openGauss 5.0.1 数据库，而非使用便捷的云数据库服务，以此深入理解数据库的部署与运维。
2.  **网络攻坚**: 配置并排查了云服务器的**安全组**规则，解决了从本地开发环境到云端数据库的网络连接超时问题。
3.  **权限配置**: 深入学习并修改了 openGauss 的核心认证文件 `pg_hba.conf`，从 `trust` 到 `md5` 再到最终的 `sha256`，理解了不同认证方式的安全性与适用场景。
4.  **驱动适配**: 遭遇了标准 Go PostgreSQL 驱动与 openGauss `sha256` 认证协议的兼容性难题，表现为 `received unexpected message`、`empty password` 等多种深层错误。最终通过查阅官方资料，找到了并成功集成了 openGauss 官方 Go 驱动 **`openGauss-connector-go-pq`**，完美解决了兼容性问题。
5.  **权限管理**: 在连接成功后，进一步解决了 GORM `AutoMigrate` 因数据库用户权限不足导致的 `permission denied for schema public` 问题，通过 `GRANT` 语句为应用用户授予了必要的 schema 权限。

这次迁移是一次宝贵的实战经历，它让我们深刻体会到企业级数据库在安全性、认证协议上的严谨性，以及在异构环境中解决复杂兼容性问题的重要性。

参考资料：

[pq - A pure Go openGauss driver for Go's database/sql package](https://pkg.go.dev/gitee.com/opengauss/openGauss-connector-go-pq)

## 数据库设计

本项目是数据库实践课程的大作业，数据库设计是核心部分。系统采用关系型数据库设计，包含3个主要表和完整的外键约束。

### 1. ship_info 表（舰船信息表）
舰船基本信息表，存储游戏中每艘舰船的核心属性。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键，自增 |
| ship_name | VARCHAR(255) | NOT NULL, INDEX | 舰船名称，建立索引以优化查询 |
| rarity | VARCHAR(64) | NOT NULL | 稀有度（普通/稀有/精锐/超稀有/海上传奇/最高方案/决战方案） |
| ship_type | VARCHAR(64) | NOT NULL, INDEX | 舰种（驱逐/轻巡/重巡/战列/战巡/航母/轻航/潜艇/维修） |
| faction | VARCHAR(64) | NOT NULL, INDEX | 阵营（白鹰/皇家/重樱/铁血/东煌/北方联合等） |
| level | INTEGER | DEFAULT 1 | 等级（1-125） |
| stars | INTEGER | DEFAULT 1 | 当前星级（满星后可装备兵装） |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | INDEX | 软删除时间（GORM软删除支持） |

**设计要点**：
- 使用 `ship_type` 和 `faction` 索引优化筛选查询
- `stars` 字段根据稀有度有不同的最大值（普通4星、稀有5星、超稀有6星）
- 软删除机制，删除的数据不会真正被删除，便于数据恢复

### 2. equipments 表（装备表）
装备库表，存储游戏中所有可用的装备。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键，自增 |
| name | VARCHAR(255) | UNIQUE, NOT NULL | 装备名称，唯一索引 |
| type | VARCHAR(64) | NOT NULL, INDEX | 装备类型（见下方枚举） |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | INDEX | 软删除时间 |

**装备类型枚举**：
- `SMALL_CALIBER_MAIN_GUN` - 小口径主炮（驱逐炮）
- `MEDIUM_CALIBER_MAIN_GUN` - 中口径主炮（轻巡炮）
- `LARGE_CALIBER_MAIN_GUN` - 大口径主炮（重巡炮）
- `BATTLESHIP_MAIN_GUN` - 战列炮
- `TORPEDO` - 水面鱼雷
- `SUBMARINE_TORPEDO` - 潜艇鱼雷
- `ANTI_AIR_GUN` - 防空炮
- `AUXILIARY` - 设备
- `FIGHTER` - 战斗机
- `DIVE_BOMBER` - 轰炸机
- `TORPEDO_BOMBER` - 鱼雷机
- `AUGMENT` - 兵装（满星专属）

**设计要点**：
- 装备名称唯一索引，防止重复添加
- 类型字段建立索引，优化按类型查询装备
- 系统初始化时自动导入277种装备数据

### 3. ship_equipments 表（舰船装备关联表）
多对多关系表，记录每艘舰船的装备配置。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ship_id | INTEGER | NOT NULL, FK | 外键关联 ship_info.id |
| slot_index | INTEGER | NOT NULL | 装备栏位索引（0-5，5为兵装栏） |
| equipment_id | INTEGER | NOT NULL, FK | 外键关联 equipments.id |
| created_at | TIMESTAMP | NOT NULL | 装备时间 |
| PRIMARY KEY | (ship_id, slot_index) | - | 复合主键 |

**外键约束**：
```sql
FOREIGN KEY (ship_id) REFERENCES ship_info(id) ON DELETE CASCADE
FOREIGN KEY (equipment_id) REFERENCES equipments(id) ON DELETE RESTRICT
```

**设计要点**：
- 复合主键 `(ship_id, slot_index)` 确保每个装备栏位唯一
- 级联删除：删除舰船时自动删除其装备配置
- 限制删除：装备被使用时不能删除
- `slot_index` 0-4 为常规装备栏，5 为兵装栏（需满星）

### 数据库关系图

```
┌─────────────────┐
│   ship_info     │
├─────────────────┤
│ id (PK)         │───┐
│ ship_name       │   │
│ rarity          │   │ 1:N
│ ship_type       │   │
│ faction         │   │
│ level           │   │
│ stars           │   │
│ ...             │   │
└─────────────────┘   │
                      │
                      ↓
            ┌─────────────────────┐
            │ ship_equipments     │
            ├─────────────────────┤
            │ ship_id (PK,FK) ────┼── 关联到 ship_info.id
            │ slot_index (PK)     │
            │ equipment_id (FK) ──┼── 关联到 equipments.id
            │ created_at          │
            └─────────────────────┘
                      ↑
                      │ N:1
                      │
┌─────────────────┐   │
│   equipments    │   │
├─────────────────┤   │
│ id (PK)         │───┘
│ name (UNIQUE)   │
│ type            │
│ ...             │
└─────────────────┘
```

### 数据完整性保证

1. **主键约束** - 所有表都有主键，确保数据唯一性
2. **外键约束** - 维护舰船和装备的关联关系
3. **唯一约束** - 装备名称唯一，防止重复
4. **非空约束** - 关键字段不允许为空
5. **索引优化** - 在常用查询字段建立索引
6. **软删除** - 保留历史数据，支持数据恢复

### 业务规则实现

1. **星级限制** - 应用层校验星级范围（不同稀有度有不同上限）
2. **装备类型限制** - 根据舰种过滤可装备的装备类型
3. **兵装限制** - 只有满星舰船可以装备兵装
4. **按需创建** - 装备栏记录在装备时才创建，避免空记录

## 快速开始

本项目采用 **Go 构建标签 (Build Tags)** 实现云端与本地环境的无缝切换，无需修改数据库初始化代码，只需在启动时添加一个参数即可。

### 模式一：连接云端 openGauss 数据库 (最终部署方案)

这是本项目的最终形态，应用后端连接到部署在华为云 ECS 上的 openGauss 数据库。

1.  **配置连接信息**
    *   打开 `PortManagerBackend/pkg/constants/constants.go` 文件，将其中的数据库连接信息修改为您的**云端 openGauss 配置**。

2.  **启动后端服务 (openGauss 模式)**
    *   执行以下命令启动后端。`-tags opengauss` 会告诉 Go 编译器选择使用 openGauss 专用驱动。
    ```powershell
    cd PortManagerBackend
    go run -tags opengauss .
    ```
    首次启动会自动在云端数据库中创建表结构并初始化数据。

3.  **启动前端服务**
    ```powershell
    cd PortManagerFrontend
    npm install
    npm run dev
    ```

### 模式二：使用 Docker 启动本地 PostgreSQL (用于开发测试)

如果您需要在本地进行快速开发和测试，可以使用 Docker Compose 一键启动 PostgreSQL 数据库。

1.  **配置连接信息**
    *   打开 `PortManagerBackend/pkg/constants/constants.go` 文件，将其中的数据库连接信息修改为**本地 PostgreSQL 配置**。

2.  **启动本地数据库**
    ```powershell
    docker-compose up -d
    ```

3.  **启动后端服务 (PostgreSQL 模式)**
    *   执行以下命令启动后端。`-tags postgres` 会告诉 Go 编译器选择使用标准的 PostgreSQL 驱动。
    ```powershell
    cd PortManagerBackend
    go run -tags postgres .
    ```

4.  **启动前端服务**
    ```powershell
    cd PortManagerFrontend
    npm install
    npm run dev
    ```

### 数据库配置

配置文件位于 `PortManagerBackend/pkg/constants/constants.go`

#### 云端部署 (openGauss) 配置示例
```go
const (
	PostgresHost     = "your-ip"    // 华为云 ECS 公网 IP
	PostgresPort     = 26000
	PostgresUser     = "shaddocknh3"     // 数据库使用者
	PostgresPassword = "shenmidazhi"    // 数据库密码
	PostgresDBName   = "al_port_manager" // 数据库名
	PostgresSSLMode  = "disable"
	PostgresTimeZone = "Asia/Shanghai"
)
```

#### 本地开发 (PostgreSQL) 配置示例

```go
const (
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresUser     = "commander"
	PostgresPassword = "mysecretpassword"
	PostgresDBName   = "al_port_db"
	PostgresSSLMode  = "disable"
	PostgresTimeZone = "Asia/Shanghai"
)
```

## API 接口

### 舰船管理
- `POST /api/v1/ship` - 创建舰船
- `GET /api/v1/ship/:id` - 获取舰船详情（含装备信息）
- `PUT /api/v1/ship/:id` - 更新舰船信息
- `DELETE /api/v1/ship/:id` - 删除舰船
- `GET /api/v1/ships` - 获取舰船列表（分页+筛选）
- `GET /api/v1/stats` - 获取统计数据

### 装备管理
- `GET /api/v1/equipments` - 获取装备列表
- `PUT /api/v1/ship/:id/equip` - 装备/卸载装备

### 请求示例

**创建舰船**：
```json
{
  "shipName": "企业",
  "rarity": "超稀有",
  "shipType": "航母",
  "faction": "白鹰"
}
```

**装备装备**：
```json
{
  "slotIndex": 0,
  "equipmentId": 25
}
```

## 许可证
[MIT License](LICENSE)

## 作者
[ShaddockNH3](https://github.com/ShaddockNH3)
