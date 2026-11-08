# 代码结构

重构没有换技术栈，还是 Hertz + GORM Gen + Vue 3。动的是分层和规则摆放的位置。

## 原来什么样

| 位置 | 问题 |
| --- | --- |
| `biz/convert/converter.go` | 装备类型枚举转字符串写了 12 个 `case`，反方向在 service 里又写了一遍，生成代码里其实已经有 `String()` 和 `EquipmentTypeFromString()` |
| `biz/service/ship_service.go` | 同一份映射写第三次；统计三段循环复制粘贴；分页 `page=0` 会算出负 offset |
| `biz/service` 的 `EquipShip` | 直接把前端传来的槽位和装备写进库，不校验槽位范围、不校验舰种能不能装、满星才解锁的兵装槽也不拦 |
| `biz/handler` | 所有错误一律 500 + `err.Error()`，前端分不清"没找到"和"服务器炸了" |
| 前端 `ShipList.vue` | 稀有度下拉写了两遍，星级规则、标签颜色各一份映射，和后面的 Go 代码再各一份 |

## 现在怎么分

```
handler（HTTP）→ service（用例）→ domain/ship（规则）
                      ↑
                  convert（形状转换）
```

- `biz/domain/ship`：纯函数，不碰数据库、不碰 HTTP。装备类型名与枚举互转、舰种对应的槽位规则、星级规则、分页规范化都在这里。
- `biz/service`：取数、过规则、存回去。`EquipShip` 现在先过三道检查再写库。
- `biz/convert`：只做数据库行和 API 结构之间的形状转换。
- `biz/handler`：入口。`errors.go` 里一张小表把领域错误翻成 400/404/500。

## 装备槽规则

| 舰种 | 0 | 1 | 2 | 3 | 4 |
| --- | --- | --- | --- | --- | --- |
| 驱逐 | 小口径主炮 | 鱼雷 | 防空炮 | 设备 | 设备 |
| 轻巡 | 小/中口径主炮 | 鱼雷、小口径主炮 | 防空炮 | 设备 | 设备 |
| 重巡 | 中/大口径主炮 | 鱼雷、小/中口径主炮 | 防空炮 | 设备 | 设备 |
| 战列 / 战巡 | 战列主炮 | 小/中/大口径主炮 | 防空炮 | 设备 | 设备 |
| 航母 | 战斗机 | 轰炸机 | 鱼雷机 | 设备 | 设备 |
| 轻航 | 战斗机、轰炸机 | 鱼雷机 | 防空炮 | 设备 | 设备 |
| 潜艇 | 潜艇鱼雷 | 潜艇鱼雷 | 小口径主炮 | 设备 | 设备 |
| 维修 | 防空炮 | 防空炮 | 防空炮 | 设备 | 设备 |

兵装槽（下标 5）任何舰种都只收兵装，而且要求星级已满。

## 星级规则

| 稀有度 | 初始 | 上限 |
| --- | --- | --- |
| 普通 | 1 | 4 |
| 稀有 / 精锐 | 2 | 5 |
| 超稀有 / 海上传奇 / 最高方案 / 决战方案 | 3 | 6 |

后端 `stars.go` 和前端 `constants/ship.js` 各存一份，改的时候两边一起改——前端要提前显示星级，绕不开。

## 错误码

| 领域错误 | HTTP |
| --- | --- |
| 舰船不存在、装备不存在 | 404 |
| 槽位越界、装备类型不匹配、兵装槽未解锁、名字为空 | 400 |
| 其它 | 500 |

## 前端

`ShipList.vue` 原来 794 行：模板 307 行、脚本 418 行，筛选、表格、翻页、两个对话框、装备增删全在一个 `setup` 里。

拆完：

| 文件 | 职责 | 行数 |
| --- | --- | --- |
| `views/ShipList.vue` | 只做编排：谁触发什么、结果交给谁 | 86 |
| `composables/useShipList.js` | 列表用例：加载 / 筛选 / 分页 / 升星 / 删除 | 117 |
| `composables/useEquipmentDialog.js` | 装备用例：选、装、卸、刷新 | 116 |
| `components/ShipFilterBar.vue` | 筛选条 | 99 |
| `components/ShipTable.vue` | 表格 + 分页，动作全部往上抛事件 | 113 |
| `components/ShipFormDialog.vue` | 新建 / 编辑表单，自己校验、自己存 | 149 |
| `components/EquipmentDialog.vue` | 版面和装备选择弹窗 | 87 |
| `components/EquipmentSlot.vue` | 单个槽位（普通槽与兵装槽共用） | 87 |
| `components/StatCard.vue` | 统计页的指标卡（原来四张卡各写一遍） | 80 |
| `components/DistributionCard.vue` | 统计页的分布列表（阵营/稀有度/舰种三块共用） | 90 |
| `stores/ship.js` | 列表数据的唯一来源（ships / total / loading） | 47 |

统计页 303 行 → 95 行，四张指标卡和三块分布列表各自收敛成一个组件；装备对话框 225 行 → 87 行，槽位卡片和用例分别拿出去。

列表数据只存一份在 store 里，`useShipList` 用 `storeToRefs` 取。之前 composable 自己又存了一套 `ships/total/loading`，靠手写赋值同步，两边飘了不会有任何报错。

错误提示只在 `api/index.js` 的响应拦截器里弹一次。之前拦截器弹一次、`catch` 里再弹一次，后端挂掉时页面连着蹦两条一样的红条。

导航菜单和图标从路由 `meta` 生成（`router/index.js` 一处），不再在 `Layout.vue` 里手抄一遍。

### 界面

这一版把 Element Plus 换成了 **Naive UI + Phosphor 图标**，原因有三个：

1. 之前是紫色渐变整页背景 + 半透明磨砂卡片 + 白色文字，属于「深色/花底」那一路，跟这套后台该有的浅色底正相反。
2. 卡片、按钮、表格的样式有相当一部分是自己拿 `div` 加 CSS 拼的，组件库的能力没用上。
3. Element Plus 的产物 1147 KB，Naive UI 按需引入之后主包 281 KB。

现在：底色 `#f4f6fa`、卡片纯白、主色深蓝 `#1e40af`、强调金 `#c8963e`，全部写在 `theme.js` 的 `themeOverrides` 里，组件内不写死颜色。布局用 `NLayout` / `NGrid` / `NCard`，表格 `NDataTable`，指标 `NStatistic`，分布 `NProgress`，弹窗 `NModal`，图标统一 Phosphor。

提示走 `utils/feedback.js` 里的离散 API（`createDiscreteApi`），因为 `useMessage` 只能在组件 setup 里用，而拦截器和 composable 都要弹。

顺带修了两处：创建时间列原先被挤成两行；装备名过长会撑出槽位卡片（标签现在换行）。

### 响应式

页面最窄到 375 也不该出现半截元素。两个断点，逻辑在 `composables/useBreakpoints.js`：

| 宽度 | 变化 |
| --- | --- |
| < 1100 | 表格藏掉创建时间、等级；操作列改成图标按钮（`title` 兜住含义） |
| < 720 | 筛选条改成一行一个字段；分页退成简版（失去每页条数和跳页）；页面留白收到 12px |

表格的 `scroll-x` 跟着模式走，不然丢完列还会留一大片横向滚动。

顺带清掉三处重复：舰种和阵营的 `<el-option>` 名单（筛选和表单各写一遍）、装备类型数字/字符串两份映射、兵装按舰种过滤的映射。全部收进 `constants/ship.js`。

筛选条直接改父级传进来的 `model` 对象（一个 reactive），改完发 `search`。这样不用为四个字段写四组 `update:modelValue`，代价是组件和父级共享状态——在页内组件里可以接受，跨页面复用就得改成 `v-model`。

## 数据库层

`db/` 原来两个文件 374 / 396 行，其中各有 304 行是同一份装备种子，逐字相同。现在：

| 文件 | 行数 | 干什么 |
| --- | --- | --- |
| `db/schema.go` | 388 | 连接串、建表清单、274 条装备种子、按名幂等补数据 |
| `db/init_postgres.go` | 20 | 用 gorm 官方驱动开连接 |
| `db/init_opengauss.go` | 28 | 用华为官方驱动开连接，表在就不动结构 |

连接参数不再写死：`pkg/constants` 改成读环境变量，默认值对齐 `docker-compose.yml`。之前 `.env.example` 里列着七个变量，代码里一个都没读——线上那个华为云 IP 和密码是直接写在常量里的。

行为上有一处变化：openGauss 原来只在「建表那一次」写种子，老库加了新装备不会补；现在两种驱动都走幂等补插，缺哪条补哪条。

## 构建与依赖

`go 1.27.0`。之前卡在 Go 1.23，是因为 `sonic v1.14` 用了 Go 1.27 已经删掉的 `GoMapIterator`。这轮把链上几个包一起升了：

| 包 | 前 | 后 |
| --- | --- | --- |
| hertz | 0.10.3 | 0.10.6 |
| sonic | 1.14.0 | 1.15.5 |
| gorm | 1.24.3 | 1.31.2 |
| gorm/gen | 0.3.21 | 0.3.29 |
| gorm/driver/postgres | 1.4.8 | 1.6.3 |
| pgx | 5.5.5 | 5.10.0 |

`query/` 里的生成代码没重新生成，对着新 gorm 直接编译通过；`go vet -tags postgres ./biz/... ./db/...` 干净，13 个单测照旧全过。`-tags opengauss` 也能编出来，只是没有实例可跑。

## 代价

- 多了一个包，改装备规则要从 `service` 跳到 `domain`。
- 生成代码 `biz/handler/api/port_manager_service.go` 里 `sendError` 改成了调 `mapError`，下次 `hz` 重新生成会被覆盖，改动只在两行内。
- `query/` 和 `biz/model/` 是生成物，没碰。
