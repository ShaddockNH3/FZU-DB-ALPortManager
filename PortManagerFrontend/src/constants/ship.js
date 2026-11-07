// 稀有度、舰种、阵营、装备类型这些枚举只在这里定义一次：
// 筛选条、编辑表单、装备对话框、统计页都从这里取。
// 星级规则和后端 biz/domain/ship/stars.go 保持一致，两边改动要同步。

export const RARITY_OPTIONS = [
  { label: '普通', value: '普通' },
  { label: '稀有', value: '稀有' },
  { label: '精锐', value: '精锐' },
  { label: '超稀有', value: '超稀有' },
  { label: '海上传奇', value: '海上传奇' },
  { label: '最高方案 (科研)', value: '最高方案' },
  { label: '决战方案 (科研)', value: '决战方案' }
]

export const SHIP_TYPES = Object.freeze([
  '驱逐', '轻巡', '重巡', '战列', '战巡', '航母', '轻航', '潜艇', '维修'
])

export const FACTIONS = Object.freeze([
  '白鹰', '皇家', '重樱', '铁血', '东煌', '北方联合', '自由鸢尾',
  '维希教廷', '鸢尾教国', '撒丁帝国', '郁金王国', '飓风', 'META'
])

const STAR_RULES = {
  普通: { initial: 1, max: 4, tag: 'info', color: '#909399' },
  稀有: { initial: 2, max: 5, tag: 'success', color: '#67C23A' },
  精锐: { initial: 2, max: 5, tag: 'warning', color: '#E6A23C' },
  超稀有: { initial: 3, max: 6, tag: 'danger', color: '#F56C6C' },
  海上传奇: { initial: 3, max: 6, tag: 'danger', color: '#C71585' },
  最高方案: { initial: 3, max: 6, tag: 'danger', color: '#8B008B' },
  决战方案: { initial: 3, max: 6, tag: 'danger', color: '#FF1493' }
}

const FALLBACK = { initial: 1, max: 3, tag: 'info', color: '#409eff' }

function ruleOf(rarity) {
  return STAR_RULES[rarity] || FALLBACK
}

export const initialStars = (rarity) => ruleOf(rarity).initial
export const maxStars = (rarity) => ruleOf(rarity).max
export const rarityTagType = (rarity) => ruleOf(rarity).tag
export const rarityColor = (rarity) => ruleOf(rarity).color

/** 满星判断：星级到顶才解锁兵装槽。 */
export function isMaxStars(ship) {
  if (!ship || !ship.stars || !ship.rarity) return false
  return ship.stars >= maxStars(ship.rarity)
}

// 后端返回的装备类型既可能是数字枚举，也可能是字符串，两个都认。
export const EQUIPMENT_TYPE_LABELS = Object.freeze({
  0: '驱逐主炮',
  1: '轻巡主炮',
  2: '重巡主炮',
  3: '战列主炮',
  4: '鱼雷',
  5: '潜艇鱼雷',
  6: '防空炮',
  7: '设备',
  8: '轰炸机',
  9: '鱼雷机',
  10: '战斗机',
  11: '兵装',
  SMALL_CALIBER_MAIN_GUN: '驱逐主炮',
  MEDIUM_CALIBER_MAIN_GUN: '轻巡主炮',
  LARGE_CALIBER_MAIN_GUN: '重巡主炮',
  BATTLESHIP_MAIN_GUN: '战列主炮',
  TORPEDO: '鱼雷',
  SUBMARINE_TORPEDO: '潜艇鱼雷',
  ANTI_AIR_GUN: '防空炮',
  AUXILIARY: '设备',
  DIVE_BOMBER: '轰炸机',
  TORPEDO_BOMBER: '鱼雷机',
  FIGHTER: '战斗机',
  AUGMENT: '兵装'
})

export const equipmentTypeLabel = (type) => EQUIPMENT_TYPE_LABELS[type] ?? type

export const acceptableTypesLabel = (types) =>
  (types || []).map(equipmentTypeLabel).join(', ')

// 兵装按舰种限定可用名单（后端只校验到 AUGMENT 类型，具体名字由前端筛）。
export const AUGMENTS_BY_SHIP_TYPE = Object.freeze({
  驱逐: ['双剑', '单手锤'],
  轻巡: ['铁剑', '手弩'],
  重巡: ['大剑', '骑枪'],
  战列: ['指挥刀', '轻弩'],
  战巡: ['指挥刀', '轻弩'],
  航母: ['猎弓', '权杖'],
  轻航: ['猎弓', '权杖'],
  维修: ['维修手弩'],
  潜艇: ['短剑', '若无']
})

export function augmentsForShipType(augments, shipType) {
  const allowed = AUGMENTS_BY_SHIP_TYPE[shipType] || []
  return augments.filter((augment) => allowed.includes(augment.name))
}

// 阵营配色：统计页的进度条用，和 FACTIONS 一一对应。
export const FACTION_COLORS = Object.freeze({
  白鹰: '#4169E1',
  皇家: '#FFD700',
  重樱: '#FF1493',
  铁血: '#696969',
  东煌: '#FF4500',
  北方联合: '#DC143C',
  自由鸢尾: '#9370DB',
  维希教廷: '#8B4513',
  鸢尾教国: '#BA55D3',
  撒丁帝国: '#228B22',
  郁金王国: '#FFB6C1',
  飓风: '#00CED1',
  META: '#8B008B'
})

export const factionColor = (faction) => FACTION_COLORS[faction] || '#409eff'
