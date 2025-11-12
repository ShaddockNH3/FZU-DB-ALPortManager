-- 初始化装备数据
-- 每个装备包含：ID, 名称, 类型, 品质(通过名称体现)

-- 驱逐主炮 (SMALL_CALIBER_MAIN_GUN)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(1, NOW(), NOW(), '金B-13双联装130mm主炮B-2LM', 'SMALL_CALIBER_MAIN_GUN')
ON CONFLICT (name) DO NOTHING;

-- 轻巡主炮 (MEDIUM_CALIBER_MAIN_GUN)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(2, NOW(), NOW(), '彩试作型三联装152mm主炮', 'MEDIUM_CALIBER_MAIN_GUN')
ON CONFLICT (name) DO NOTHING;

-- 重巡主炮 (LARGE_CALIBER_MAIN_GUN)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(3, NOW(), NOW(), '金双联装203mm主炮SKC', 'LARGE_CALIBER_MAIN_GUN')
ON CONFLICT (name) DO NOTHING;

-- 鱼雷 (TORPEDO)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(4, NOW(), NOW(), '金四联装533mm磁性鱼雷', 'TORPEDO')
ON CONFLICT (name) DO NOTHING;

-- 潜艇鱼雷 (SUBMARINE_TORPEDO)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(5, NOW(), NOW(), '金G7e声导鱼雷', 'SUBMARINE_TORPEDO')
ON CONFLICT (name) DO NOTHING;

-- 防空炮 (ANTI_AIR_GUN)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(6, NOW(), NOW(), '金双联装113mm高射炮', 'ANTI_AIR_GUN')
ON CONFLICT (name) DO NOTHING;

-- 设备 (AUXILIARY)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(7, NOW(), NOW(), '紫应急修理装置', 'AUXILIARY')
ON CONFLICT (name) DO NOTHING;

-- 战斗机 (FIGHTER)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(8, NOW(), NOW(), '金Bf-109G战斗机', 'FIGHTER')
ON CONFLICT (name) DO NOTHING;

-- 轰炸机 (DIVE_BOMBER)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(9, NOW(), NOW(), '彩流星', 'DIVE_BOMBER')
ON CONFLICT (name) DO NOTHING;

-- 鱼雷机 (TORPEDO_BOMBER)
INSERT INTO equipments (id, created_at, updated_at, name, type) VALUES
(10, NOW(), NOW(), '金Barracuda', 'TORPEDO_BOMBER')
ON CONFLICT (name) DO NOTHING;

-- 重置序列（如果使用了自增ID）
SELECT setval('equipments_id_seq', (SELECT MAX(id) FROM equipments));
