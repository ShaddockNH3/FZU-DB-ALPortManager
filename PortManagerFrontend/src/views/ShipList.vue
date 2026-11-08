<script setup>
import { ref } from 'vue'
import { getShipById } from '@/api/ship'
import { useShipList } from '@/composables/useShipList'
import ShipFilterBar from '@/components/ShipFilterBar.vue'
import ShipTable from '@/components/ShipTable.vue'
import ShipFormDialog from '@/components/ShipFormDialog.vue'
import EquipmentDialog from '@/components/EquipmentDialog.vue'

// 列表状态与用例都在 composable 里，这里只做编排：谁触发什么、结果交给谁。
const {
  filterForm,
  pagination,
  ships,
  total,
  loading,
  load,
  search,
  reset,
  changePage,
  changePageSize,
  changeStars,
  remove
} = useShipList()

const formVisible = ref(false)
const editingShip = ref(null)
const equipmentVisible = ref(false)
const equipmentShip = ref(null)

function openCreate() {
  editingShip.value = null
  formVisible.value = true
}

function openEdit(ship) {
  editingShip.value = { ...ship }
  formVisible.value = true
}

/** 装备对话框要完整数据（含已装装备），列表里的行没有，单独查一次。 */
async function openEquipment(ship) {
  try {
    equipmentShip.value = await getShipById(ship.id)
    equipmentVisible.value = true
  } catch (error) {
    console.error('获取舰船装备信息失败:', error)
  }
}
</script>

<template>
  <div class="page">
    <header class="page-head">
      <h2>舰船管理</h2>
      <p>共 {{ total }} 艘在册，点星级两侧的 ± 升降星</p>
    </header>

    <ship-filter-bar
      :model="filterForm"
      :loading="loading"
      @search="search"
      @reset="reset"
      @add="openCreate"
    />

    <ship-table
      :ships="ships"
      :loading="loading"
      :page="pagination.page"
      :page-size="pagination.pageSize"
      :total="total"
      @edit="openEdit"
      @delete="remove"
      @equip="openEquipment"
      @star="changeStars"
      @page-change="changePage"
      @size-change="changePageSize"
    />

    <ship-form-dialog v-model:visible="formVisible" :ship="editingShip" @saved="load" />

    <equipment-dialog v-model:visible="equipmentVisible" :ship="equipmentShip" @changed="load" />
  </div>
</template>

<style lang="scss" scoped>
.page-head {
  margin-bottom: 18px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
  }

  p {
    margin-top: 4px;
    font-size: 13px;
    color: #8b95a7;
  }
}
</style>
