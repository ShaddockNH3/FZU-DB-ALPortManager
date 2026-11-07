// 稀有度只在这里定义一次：下拉选项、初始星级、满星星级、标签样式都从这里取。
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
