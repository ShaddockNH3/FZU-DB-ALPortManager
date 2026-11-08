<script setup>
import { computed, h } from 'vue'
import { NButton, NCard, NDataTable, NIcon, NPagination, NTag } from 'naive-ui'
import { PhMinus, PhPencilSimple, PhPlus, PhSword, PhTrash } from '@phosphor-icons/vue'
import { initialStars, maxStars, PAGE_SIZES, rarityTagType } from '@/constants/ship'

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
  return h('div', { style: 'display:flex;gap:6px' }, [
    h(NButton, {
      size: 'tiny',
      secondary: true,
      type: 'primary',
      onClick: () => emit('edit', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhPencilSimple) }), default: () => '编辑' }),
    h(NButton, {
      size: 'tiny',
      secondary: true,
      onClick: () => emit('equip', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhSword) }), default: () => '装备' }),
    h(NButton, {
      size: 'tiny',
      quaternary: true,
      type: 'error',
      onClick: () => emit('delete', row)
    }, { icon: () => h(NIcon, null, { default: () => h(PhTrash) }), default: () => '删除' })
  ])
}

const columns = computed(() => [
  { title: 'ID', key: 'id', width: 70, sorter: (a, b) => a.id - b.id },
  { title: '舰船名称', key: 'shipName', minWidth: 150 },
  {
    title: '稀有度',
    key: 'rarity',
    width: 120,
    render: (row) =>
      h(NTag, {
        size: 'small',
        bordered: false,
        type: TAG_TYPES[rarityTagType(row.rarity)] || 'default'
      }, { default: () => row.rarity })
  },
  { title: '舰种', key: 'shipType', width: 100 },
  { title: '阵营', key: 'faction', width: 110 },
  { title: '星级', key: 'stars', width: 160, render: starCell },
  { title: '等级', key: 'level', width: 90, sorter: (a, b) => a.level - b.level },
  { title: '创建时间', key: 'createdAt', width: 190, className: 'nowrap' },
  { title: '操作', key: 'actions', width: 240, render: actionCell }
])
</script>

<template>
  <n-card size="small">
    <n-data-table
      :columns="columns"
      :data="props.ships"
      :loading="props.loading"
      :bordered="false"
      :row-key="(row) => row.id"
      :scroll-x="1100"
      size="medium"
    />

    <div class="pager">
      <n-pagination
        :page="props.page"
        :page-size="props.pageSize"
        :item-count="props.total"
        :page-sizes="pageSizes"
        show-size-picker
        show-quick-jumper
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
  margin-top: 16px;
}
</style>
