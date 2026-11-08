<script setup>
import { computed, h } from 'vue'
import { NButton, NCard, NDataTable, NIcon, NPagination, NTag } from 'naive-ui'
import { PhMinus, PhPencilSimple, PhPlus, PhSword, PhTrash } from '@phosphor-icons/vue'
import { initialStars, maxStars, PAGE_SIZES, rarityTagType } from '@/constants/ship'
import { useBreakpoints } from '@/composables/useBreakpoints'

/**
 * 舰船表格 + 分页。纯展示，所有动作都往上抛事件，状态在父级。
 */
const props = defineProps({
  ships: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 }
})

const emit = defineEmits(['edit', 'delete', 'equip', 'star', 'page-change', 'size-change'])

// 稀有度的标签色：领域里叫 danger，Naive 里叫 error
const TAG_TYPES = { info: 'default', success: 'success', warning: 'warning', danger: 'error' }

const pageSizes = PAGE_SIZES.map((size) => ({ label: `${size} 条/页`, value: size }))

const { isCompact, isNarrow } = useBreakpoints()

function starCell(row) {
  const minusDisabled = row.stars <= initialStars(row.rarity)
  const plusDisabled = row.stars >= maxStars(row.rarity)

  return h('div', { style: 'display:flex;align-items:center;gap:6px' }, [
    h(NButton, {
      size: 'tiny',
      circle: true,
      quaternary: true,
      disabled: minusDisabled,
      onClick: () => emit('star', row, -1)
    }, { icon: () => h(NIcon, null, { default: () => h(PhMinus) }) }),
    h('span', { style: 'min-width:52px;text-align:center;font-variant-numeric:tabular-nums' },
      `${row.stars}★${plusDisabled ? ' 满' : ''}`),
    h(NButton, {
      size: 'tiny',
      circle: true,
      quaternary: true,
      disabled: plusDisabled,
      onClick: () => emit('star', row, 1)
    }, { icon: () => h(NIcon, null, { default: () => h(PhPlus) }) })
  ])
}

function actionCell(row) {
  // 窄屏只留图标：三个带文字的按钮要 200px，图标版 100px 就够
  const iconOnly = isCompact.value
  const label = (text) => (iconOnly ? {} : { default: () => text })

  return h('div', { style: 'display:flex;gap:6px' }, [
    h(NButton, {
      size: 'tiny',
      secondary: true,
      type: 'primary',
      circle: iconOnly,
      title: '编辑',
      onClick: () => emit('edit', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhPencilSimple) }), ...label('编辑') }),
    h(NButton, {
      size: 'tiny',
      secondary: true,
      circle: iconOnly,
      title: '装备',
      onClick: () => emit('equip', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhSword) }), ...label('装备') }),
    h(NButton, {
      size: 'tiny',
      quaternary: true,
      type: 'error',
      circle: iconOnly,
      title: '删除',
      onClick: () => emit('delete', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhTrash) }), ...label('删除') })
  ])
}

// 列的取舍顺序：先丢创建时间，再丢 ID / 舰种 / 阵营 / 等级，名称和操作永远留着
const baseColumns = [
  { title: 'ID', key: 'id', width: 70, sorter: (a, b) => a.id - b.id, hideOn: 'narrow' },
  { title: '舰船名称', key: 'shipName', minWidth: 140 },
  {
    title: '稀有度',
    key: 'rarity',
    width: 108,
    render: (row) =>
      h(NTag, {
        size: 'small',
        bordered: false,
        type: TAG_TYPES[rarityTagType(row.rarity)] || 'default'
      }, { default: () => row.rarity })
  },
  { title: '舰种', key: 'shipType', width: 100, hideOn: 'narrow' },
  { title: '阵营', key: 'faction', width: 110, hideOn: 'narrow' },
  { title: '星级', key: 'stars', width: 150, render: starCell },
  {
    title: '等级',
    key: 'level',
    width: 88,
    sorter: (a, b) => a.level - b.level,
    hideOn: 'compact'
  },
  { title: '创建时间', key: 'createdAt', width: 186, className: 'nowrap', hideOn: 'compact' },
  {
    title: '操作',
    key: 'actions',
    width: isCompact.value ? 116 : 214,
    render: actionCell
  }
]

const columns = computed(() =>
  baseColumns.filter((column) => {
    if (column.hideOn === 'compact' && isCompact.value) return false
    if (column.hideOn === 'narrow' && isNarrow.value) return false
    return true
  })
)

// 横向滚动的宽度跟着留几列走，不然窄屏会滚出一大片空白
const scrollX = computed(() => {
  if (isNarrow.value) return 520
  if (isCompact.value) return 800
  return 1180
})
</script>

<template>
  <n-card size="small">
    <n-data-table
      :columns="columns"
      :data="props.ships"
      :loading="props.loading"
      :bordered="false"
      :row-key="(row) => row.id"
      :scroll-x="scrollX"
      size="medium"
    />

    <div class="pager">
      <n-pagination
        :page="props.page"
        :page-size="props.pageSize"
        :item-count="props.total"
        :page-sizes="pageSizes"
        :simple="isNarrow"
        :show-size-picker="!isNarrow"
        :show-quick-jumper="!isNarrow"
        :page-slot="isCompact ? 5 : 9"
        :disabled="props.loading"
        @update:page="(value) => emit('page-change', value)"
        @update:page-size="(value) => emit('size-change', value)"
      />
    </div>
  </n-card>
</template>

<style lang="scss" scoped>
.pager {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}
</style>
