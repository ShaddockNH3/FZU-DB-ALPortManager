# FZU-DB-ALPortManager

福州大学 2025 数据库实践大作业：碧蓝航线风格的港口管理系统，Go（Hertz + GORM Gen）+ Vue 3（Naive UI），数据跑在华为云 openGauss 上。

## 快速开始

```bash
cp .env.example .env                    # 按需改数据库地址；不建也能跑，默认值和 docker-compose 对得上
docker compose up -d                    # 本地 PostgreSQL，开发用

cd PortManagerBackend
go build -tags postgres -o output/bin/port-manager .   # Go 1.27+
./output/bin/port-manager               # 默认 :8888

cd ../PortManagerFrontend && npm install && npm run dev # :5173，/api 代理到 8888
```

数据库有两套初始化实现，靠构建标签选：`-tags postgres` 走本地 PostgreSQL，`-tags opengauss` 走华为云 openGauss（认证协议不同，驱动也不同）。连接串、建表清单、装备种子都在 `db/schema.go`，两个驱动文件只负责开连接。

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v1/ships` | 舰船列表，支持分页与名称/稀有度/舰种/阵营筛选 |
| POST | `/api/v1/ship` | 新建舰船，等级固定 1，星级按稀有度给 |
| GET | `/api/v1/ship/:id` | 舰船详情，含 5 个装备槽 + 兵装槽 |
| PUT | `/api/v1/ship/:id` | 更新舰船，空字段不覆盖 |
| DELETE | `/api/v1/ship/:id` | 删除舰船并清掉装备记录 |
| PUT | `/api/v1/ship/:id/equip` | 装/卸装备，`equipmentId` 传 0 表示卸下 |
| GET | `/api/v1/equipments` | 装备列表，可按类型过滤 |
| GET | `/api/v1/stats` | 按阵营/稀有度/舰种统计数量 |

```json
{ "shipName": "企业", "rarity": "超稀有", "shipType": "航母", "faction": "白鹰" }
```

## 配置

| 变量 | 默认 | 说明 |
| --- | --- | --- |
| `POSTGRES_HOST` | `localhost` | 数据库地址 |
| `POSTGRES_PORT` | `5432` | 端口 |
| `POSTGRES_USER` | `commander` | 用户名 |
| `POSTGRES_PASSWORD` | `mysecretpassword` | 密码 |
| `POSTGRES_DB` | `al_port_db` | 库名 |
| `POSTGRES_SSLMODE` | `disable` | openGauss 需要按服务端配置调整 |
| `POSTGRES_TIMEZONE` | `Asia/Shanghai` | 会话时区 |

## 项目结构

```
PortManagerBackend/
├── biz/domain/ship/    纯规则：装备槽匹配、星级、分页（13 个单测）
├── biz/service/        用例：取数 → 过规则 → 存回去
├── biz/convert/        数据库行 ↔ API 结构
├── biz/handler/        HTTP 入口，错误码映射在 errors.go
├── query/ model/       gorm gen 生成，手不改
└── db/
    ├── schema.go        共用：连接串、建表清单、装备种子
    ├── init_postgres.go 本地 PostgreSQL（-tags postgres）
    └── init_opengauss.go 华为云 openGauss（-tags opengauss）
PortManagerFrontend/src/
├── theme.js            浅色调色板与 Naive UI themeOverrides，颜色只在这里改
├── utils/feedback.js   全局 message / dialog（组件外也能弹提示）
├── constants/ship.js   稀有度、舰种、阵营、装备类型、兵装规则只在这里定义
├── composables/        列表状态（useShipList）、装备对话框用例（useEquipmentDialog）
├── components/         筛选条、表格、编辑对话框、装备对话框、装备槽位、统计卡片
├── api/ stores/        请求（错误提示统一在拦截器）与列表状态
└── views/              页面编排（列表 95 行，统计 95 行）
```

## 测试

```bash
cd PortManagerBackend

go test ./biz/domain/... ./biz/convert/...
```

覆盖装备槽规则（9 种舰种 × 6 个槽）、星级规则、分页边界，以及数据库行到接口结构的转换。

## 已知问题

- openGauss 分支只做到「能编译」（`-tags opengauss`）。手上没有可用实例，端到端没跑过。
- 前端没有自动化测试，改动只能靠手点。
- 装备列表只支持按单一类型过滤，多选要改 thrift 定义再生成。
- 兵装槽满星才能装，后端会拦；前端只把按钮禁掉，理由写在标题旁边的括号里。

## 更多

- 表结构、约束与 openGauss 迁移过程：[docs/database.md](docs/database.md)
- 分层与领域规则：[docs/DESIGN.md](docs/DESIGN.md)

## 致谢

[Hertz](https://github.com/cloudwego/hertz)、[GORM Gen](https://gorm.io/gen/)、[Element Plus](https://element-plus.org/)、[openGauss](https://opengauss.org/)。

## License

MIT，见 [LICENSE](LICENSE)。
