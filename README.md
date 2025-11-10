# FZU-DB-ALPortManager
福州大学2025年下半学期数据库实践大作业 - 碧蓝航线港口管理系统

## 项目简介
基于 Go + Vue3 + PostgreSQL 的碧蓝航线港口管理系统，用于管理游戏中的舰船信息，包括舰船的基本属性、稀有度、舰种、阵营等信息。

### 主要功能
- 🚢 舰船信息管理（增删改查）
- 🔍 多条件筛选查询（按舰船名称、稀有度、舰种、阵营）
- 📊 数据统计展示
- 📄 分页展示
- 🎨 现代化的 UI 界面

### 技术特点
- 前后端分离架构
- RESTful API 设计
- GORM ORM 框架，支持自动迁移
- 响应式前端界面
- Docker 容器化部署

## 技术栈
- **后端**: Go 1.23 + Hertz (CloudWeGo) + GORM 1.24 + Gen 0.3.21
- **前端**: Vue3 + Vite + Element Plus + Pinia + Axios
- **数据库**: PostgreSQL 15 → **openGauss (华为云平台)** ⚠️ 待迁移
- **容器化**: Docker + Docker Compose

> **⚠️ 重要提示**: 根据学校课程要求，本项目后续将迁移至华为云平台的 openGauss 数据库。当前使用 PostgreSQL 仅用于本地开发测试。

## 数据库结构
### ship_info 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键，自增 |
| ship_name | varchar(255) | 舰船名称 |
| rarity | varchar(64) | 稀有度（普通/稀有/精锐/超稀有/海上传奇/最高方案/决战方案） |
| ship_type | varchar(64) | 舰种（驱逐/轻巡/重巡/战列/战巡/航母/轻航/潜艇/维修） |
| faction | varchar(64) | 阵营（白鹰/皇家/重樱/铁血/东煌/北方联合等） |
| level | int | 等级（1-125） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 软删除时间 |

## 快速开始

1. **启动数据库**
   ```powershell
   docker-compose up -d
   ```

2. **启动后端**
   ```powershell
   cd PortManagerBackend
   go run .
   ```

3. **启动前端**
   ```powershell
   cd PortManagerFrontend
   npm install
   npm run dev
   ```

4. 访问 `http://localhost:5173` 使用系统

## 数据库配置
数据库配置位于 `PortManagerBackend/pkg/constants/constants.go`

默认配置：
- Host: localhost
- Port: 5432
- User: commander
- Password: mysecretpassword
- Database: al_port_db

## API 接口
### 舰船管理
- `POST /api/v1/ship` - 创建舰船
- `GET /api/v1/ship/:id` - 获取单个舰船信息
- `PUT /api/v1/ship/:id` - 更新舰船信息
- `DELETE /api/v1/ship/:id` - 删除舰船
- `GET /api/v1/ships` - 获取舰船列表（支持分页和筛选）
- `GET /api/v1/stats` - 获取统计信息
- `GET /ping` - 健康检查

### 请求参数示例
```json
{
  "shipName": "企业",
  "rarity": "超稀有",
  "shipType": "航母",
  "faction": "白鹰",
  "level": 120
}
```

## 项目结构
```
FZU-DB-ALPortManager/
├── PortManagerBackend/          # 后端代码
│   ├── biz/                     # 业务逻辑层
│   │   ├── handler/             # HTTP 处理器
│   │   ├── model/               # API 模型
│   │   ├── router/              # 路由注册
│   │   └── service/             # 业务服务
│   ├── db/                      # 数据库初始化
│   ├── model/                   # 数据模型
│   ├── query/                   # GORM Gen 生成的查询代码
│   ├── pkg/                     # 工具包
│   │   └── constants/           # 常量配置
│   ├── idl/                     # 接口定义
│   ├── main.go                  # 入口文件
│   └── generate.go              # 代码生成器
├── PortManagerFrontend/         # 前端代码
│   ├── src/
│   │   ├── views/               # 页面组件
│   │   ├── components/          # 公共组件
│   │   ├── router/              # 路由配置
│   │   ├── stores/              # 状态管理
│   │   ├── api/                 # API 请求封装
│   │   └── assets/              # 静态资源
│   ├── index.html
│   └── vite.config.js           # Vite 配置
├── docker-compose.yml           # Docker Compose 配置
└── README.md                    # 项目文档
```

## 未来计划
### 迁移至华为云 openGauss 数据库
> **学校课程要求**: 本项目后续将迁移至华为云平台的 openGauss 数据库

openGauss 是华为开源的企业级关系型数据库，基于 PostgreSQL 开发，兼容 PostgreSQL 协议。由于两者的高度兼容性，迁移工作主要包括：
- 修改数据库连接配置
- 测试 GORM 驱动兼容性
- 数据导入导出

参考资源：
- [openGauss 官方文档](https://docs.opengauss.org/)
- [华为云 GaussDB](https://www.huaweicloud.com/product/gaussdb.html)


## 许可证
本项目基于 [MIT License](LICENSE) 开源。

## 作者
[ShaddockNH3](https://github.com/ShaddockNH3)
