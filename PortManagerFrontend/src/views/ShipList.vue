<template>
  <div class="ship-list-container">
    <el-card class="filter-card" shadow="hover">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="舰船名称">
          <el-input
            v-model="filterForm.shipName"
            placeholder="请输入舰船名称"
            clearable
            @clear="handleSearch"
          />
        </el-form-item>
        <el-form-item label="稀有度">
          <el-select
            v-model="filterForm.rarity"
            placeholder="请选择稀有度"
            clearable
            @change="handleSearch"
            style="width: 180px"
          >
            <el-option label="普通" value="普通" />
            <el-option label="稀有" value="稀有" />
            <el-option label="精锐" value="精锐" />
            <el-option label="超稀有" value="超稀有" />
            <el-option label="海上传奇" value="海上传奇" />
            <el-option label="最高方案 (科研)" value="最高方案" />
            <el-option label="决战方案 (科研)" value="决战方案" />
          </el-select>
        </el-form-item>
        <el-form-item label="舰种">
          <el-select
            v-model="filterForm.shipType"
            placeholder="请选择舰种"
            clearable
            @change="handleSearch"
            style="width: 150px"
          >
            <el-option label="驱逐" value="驱逐" />
            <el-option label="轻巡" value="轻巡" />
            <el-option label="重巡" value="重巡" />
            <el-option label="战列" value="战列" />
            <el-option label="战巡" value="战巡" />
            <el-option label="航母" value="航母" />
            <el-option label="轻航" value="轻航" />
            <el-option label="潜艇" value="潜艇" />
            <el-option label="维修" value="维修" />
          </el-select>
        </el-form-item>
        <el-form-item label="阵营">
          <el-select
            v-model="filterForm.faction"
            placeholder="请选择阵营"
            clearable
            @change="handleSearch"
            style="width: 150px"
          >
            <el-option label="白鹰" value="白鹰" />
            <el-option label="皇家" value="皇家" />
            <el-option label="重樱" value="重樱" />
            <el-option label="铁血" value="铁血" />
            <el-option label="东煌" value="东煌" />
            <el-option label="北方联合" value="北方联合" />
            <el-option label="自由鸢尾" value="自由鸢尾" />
            <el-option label="维希教廷" value="维希教廷" />
            <el-option label="鸢尾教国" value="鸢尾教国" />
            <el-option label="撒丁帝国" value="撒丁帝国" />
            <el-option label="郁金王国" value="郁金王国" />
            <el-option label="飓风" value="飓风" />
            <el-option label="META" value="META" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :icon="Search">
            搜索
          </el-button>
          <el-button @click="handleReset" :icon="RefreshRight">
            重置
          </el-button>
          <el-button type="success" @click="handleAdd" :icon="Plus">
            添加舰船
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card" shadow="hover">
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        style="width: 100%"
        :default-sort="{ prop: 'id', order: 'descending' }"
      >
        <el-table-column prop="id" label="ID" width="80" sortable />
        <el-table-column prop="shipName" label="舰船名称" min-width="150" />
        <el-table-column prop="rarity" label="稀有度" width="150">
          <template #default="{ row }">
            <el-tag :type="getRarityType(row.rarity)">{{ row.rarity }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="shipType" label="舰种" width="120" />
        <el-table-column prop="faction" label="阵营" width="150" />
        <el-table-column label="星级" width="180">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-button 
                size="small" 
                :icon="Minus" 
                circle 
                @click="handleStarChange(row, -1)"
                :disabled="row.stars <= getInitialStarsByRarity(row.rarity)"
              />
              <span style="min-width: 80px; text-align: center;">
                {{ row.stars }}★ {{ row.stars >= getMaxStarsByRarity(row.rarity) ? '(满)' : '' }}
              </span>
              <el-button 
                size="small" 
                :icon="Plus" 
                circle 
                @click="handleStarChange(row, 1)"
                :disabled="row.stars >= getMaxStarsByRarity(row.rarity)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="等级" width="100" sortable />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="success"
              size="small"
              @click="handleEquipment(row)"
            >
              装备
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="舰船名称" prop="shipName">
          <el-input v-model="form.shipName" placeholder="请输入舰船名称" />
        </el-form-item>
        <el-form-item label="稀有度" prop="rarity">
          <el-select v-model="form.rarity" placeholder="请选择稀有度" style="width: 100%">
            <el-option label="普通" value="普通" />
            <el-option label="稀有" value="稀有" />
            <el-option label="精锐" value="精锐" />
            <el-option label="超稀有" value="超稀有" />
            <el-option label="海上传奇" value="海上传奇" />
            <el-option label="最高方案 (科研)" value="最高方案" />
            <el-option label="决战方案 (科研)" value="决战方案" />
          </el-select>
        </el-form-item>
        <el-form-item label="舰种" prop="shipType">
          <el-select v-model="form.shipType" placeholder="请选择舰种" style="width: 100%">
            <el-option label="驱逐" value="驱逐" />
            <el-option label="轻巡" value="轻巡" />
            <el-option label="重巡" value="重巡" />
            <el-option label="战列" value="战列" />
            <el-option label="战巡" value="战巡" />
            <el-option label="航母" value="航母" />
            <el-option label="轻航" value="轻航" />
            <el-option label="潜艇" value="潜艇" />
            <el-option label="维修" value="维修" />
          </el-select>
        </el-form-item>
        <el-form-item label="阵营" prop="faction">
          <el-select v-model="form.faction" placeholder="请选择阵营" style="width: 100%">
            <el-option label="白鹰" value="白鹰" />
            <el-option label="皇家" value="皇家" />
            <el-option label="重樱" value="重樱" />
            <el-option label="铁血" value="铁血" />
            <el-option label="东煌" value="东煌" />
            <el-option label="北方联合" value="北方联合" />
            <el-option label="自由鸢尾" value="自由鸢尾" />
            <el-option label="维希教廷" value="维希教廷" />
            <el-option label="鸢尾教国" value="鸢尾教国" />
            <el-option label="撒丁帝国" value="撒丁帝国" />
            <el-option label="郁金王国" value="郁金王国" />
            <el-option label="飓风" value="飓风" />
            <el-option label="META" value="META" />
          </el-select>
        </el-form-item>
        <el-form-item label="等级" prop="level">
          <el-input-number
            v-model="form.level"
            :min="1"
            :max="125"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 装备对话框 -->
    <el-dialog
      v-model="equipDialogVisible"
      :title="`装备管理 - ${currentShip?.shipName}`"
      width="800px"
      @close="handleEquipDialogClose"
    >
      <div v-if="currentShip" class="equipment-container">
        <div v-for="(slot, index) in currentShip.equipmentSlots" :key="index" class="equipment-slot">
          <div class="slot-label">装备栏 {{ index + 1 }}</div>
          <div class="slot-content">
            <el-tag v-if="slot.equippedItem" type="success" closable @close="handleUnequip(index)">
              {{ slot.equippedItem.name }}
            </el-tag>
            <el-button v-else type="primary" size="small" @click="handleSelectEquipment(index, slot.acceptableTypes)">
              选择装备
            </el-button>
          </div>
          <div class="acceptable-types">
            可装备类型: {{ formatAcceptableTypes(slot.acceptableTypes) }}
          </div>
        </div>
        
        <div v-if="currentShip.augmentSlot" class="equipment-slot augment-slot">
          <div class="slot-label">兵装栏 <span class="tip">(需满星)</span></div>
          <div class="slot-content">
            <el-tag v-if="currentShip.augmentSlot.equippedItem" type="warning" closable @close="handleUnequip(5)">
              {{ currentShip.augmentSlot.equippedItem.name }}
            </el-tag>
            <el-button 
              v-else 
              type="primary" 
              size="small" 
              @click="handleSelectEquipment(5, currentShip.augmentSlot.acceptableTypes)"
              :disabled="!isMaxStars(currentShip)"
            >
              {{ isMaxStars(currentShip) ? '选择兵装' : '未满星' }}
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 选择装备对话框 -->
    <el-dialog
      v-model="selectEquipDialogVisible"
      title="选择装备"
      width="600px"
    >
      <el-table
        :data="filteredEquipments"
        @row-click="handleEquipmentSelected"
        highlight-current-row
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="装备名称" min-width="200" />
        <el-table-column prop="type" label="类型" width="150">
          <template #default="{ row }">
            {{ formatEquipmentType(row.type) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, RefreshRight, Plus, Edit, Delete, Minus } from '@element-plus/icons-vue'
import { useShipStore } from '@/stores/ship'
import { getShipById, equipShip, getEquipmentList } from '@/api/ship'

const shipStore = useShipStore()

// 筛选表单
const filterForm = reactive({
  shipName: '',
  rarity: '',
  shipType: '',
  faction: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表格数据
const tableData = ref([])
const loading = ref(false)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => (isEdit.value ? '编辑舰船' : '添加舰船'))
const isEdit = ref(false)
const submitLoading = ref(false)

// 表单
const formRef = ref()
const form = reactive({
  id: null,
  shipName: '',
  rarity: '',
  shipType: '',
  faction: '',
  level: 1,
  stars: 1
})

// 表单验证规则
const rules = {
  shipName: [
    { required: true, message: '请输入舰船名称', trigger: 'blur' }
  ],
  rarity: [
    { required: true, message: '请选择稀有度', trigger: 'change' }
  ],
  shipType: [
    { required: true, message: '请选择舰种', trigger: 'change' }
  ],
  faction: [
    { required: true, message: '请选择阵营', trigger: 'change' }
  ]
}

// 获取稀有度对应的初始星级
const getInitialStarsByRarity = (rarity) => {
  const initialStarsMap = {
    '普通': 1,
    '稀有': 2,
    '精锐': 2,
    '超稀有': 3,
    '海上传奇': 3,
    '最高方案': 3,
    '决战方案': 3
  }
  return initialStarsMap[rarity] || 1
}

// 获取稀有度对应的最大星级
const getMaxStarsByRarity = (rarity) => {
  const maxStarsMap = {
    '普通': 4,
    '稀有': 5,
    '精锐': 5,
    '超稀有': 6,
    '海上传奇': 6,
    '最高方案': 6,
    '决战方案': 6
  }
  return maxStarsMap[rarity] || 3
}

// 检查是否满星
const isMaxStars = (ship) => {
  if (!ship || !ship.stars || !ship.rarity) return false
  const maxStars = getMaxStarsByRarity(ship.rarity)
  return ship.stars >= maxStars
}

// 根据舰种过滤兵装
const filterAugmentsByShipType = (augments, shipType) => {
  const augmentMap = {
    '驱逐': ['双剑', '单手锤'],
    '轻巡': ['铁剑', '手弩'],
    '重巡': ['大剑', '骑枪'],
    '战列': ['指挥刀', '轻弩'],
    '战巡': ['指挥刀', '轻弩'],
    '航母': ['猎弓', '权杖'],
    '轻航': ['猎弓', '权杖'],
    '维修': ['维修手弩'],
    '潜艇': ['短剑', '若无']
  }
  
  const allowedNames = augmentMap[shipType] || []
  return augments.filter(aug => allowedNames.includes(aug.name))
}

// 升星/降星
const handleStarChange = async (row, delta) => {
  const newStars = row.stars + delta
  const minStars = getInitialStarsByRarity(row.rarity)
  const maxStars = getMaxStarsByRarity(row.rarity)
  
  if (newStars < minStars || newStars > maxStars) {
    return
  }
  
  try {
    // 传递完整的舰船信息，只更新 stars 字段
    await shipStore.updateShip(row.id, {
      shipName: row.shipName,
      rarity: row.rarity,
      shipType: row.shipType,
      faction: row.faction,
      level: row.level,
      stars: newStars
    })
    ElMessage.success(delta > 0 ? '升星成功' : '降星成功')
    loadData()
  } catch (error) {
    console.error('星级调整失败:', error)
    ElMessage.error('星级调整失败')
  }
}

// 获取稀有度标签类型
const getRarityType = (rarity) => {
  const typeMap = {
    '普通': 'info',
    '稀有': 'success',
    '精锐': 'warning',
    '超稀有': 'danger',
    '海上传奇': 'danger',
    '最高方案': 'danger',
    '决战方案': 'danger'
  }
  return typeMap[rarity] || 'info'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...filterForm
    }
    // 移除空值
    Object.keys(params).forEach(key => {
      if (!params[key]) delete params[key]
    })
    
    const data = await shipStore.fetchShipList(params)
    tableData.value = data.ships || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  Object.assign(filterForm, {
    shipName: '',
    rarity: '',
    shipType: '',
    faction: ''
  })
  handleSearch()
}

// 添加
const handleAdd = () => {
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    shipName: row.shipName,
    rarity: row.rarity,
    shipType: row.shipType,
    faction: row.faction,
    level: row.level || 1
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除舰船"${row.shipName}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await shipStore.deleteShip(row.id)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 取消删除
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitLoading.value = true
    try {
      const data = {
        shipName: form.shipName,
        rarity: form.rarity,
        shipType: form.shipType,
        faction: form.faction,
        level: form.level
      }
      
      if (isEdit.value) {
        await shipStore.updateShip(form.id, data)
        ElMessage.success('更新成功')
      } else {
        await shipStore.createShip(data)
        ElMessage.success('添加成功')
      }
      
      dialogVisible.value = false
      loadData()
    } catch (error) {
      console.error('提交失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 关闭对话框
const handleDialogClose = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    id: null,
    shipName: '',
    rarity: '',
    shipType: '',
    faction: '',
    level: 1
  })
}

// 分页大小改变
const handleSizeChange = () => {
  pagination.page = 1
  loadData()
}

// 页码改变
const handlePageChange = () => {
  loadData()
}

// 装备管理
const equipDialogVisible = ref(false)
const selectEquipDialogVisible = ref(false)
const currentShip = ref(null)
const currentSlotIndex = ref(null)
const allEquipments = ref([])
const filteredEquipments = ref([])

// 打开装备对话框
const handleEquipment = async (row) => {
  try {
    // 重新获取最新的舰船数据（包含装备信息）
    currentShip.value = await getShipById(row.id)
    console.log('舰船装备信息:', currentShip.value)
    equipDialogVisible.value = true
  } catch (error) {
    console.error('获取舰船装备信息失败:', error)
    ElMessage.error('获取舰船装备信息失败')
  }
}

// 关闭装备对话框
const handleEquipDialogClose = () => {
  currentShip.value = null
  loadData() // 重新加载列表
}

// 选择装备
const handleSelectEquipment = async (slotIndex, acceptableTypes) => {
  currentSlotIndex.value = slotIndex
  
  console.log('选择装备 - 栏位:', slotIndex, '可接受类型:', acceptableTypes)
  
  // 如果是兵装栏（slotIndex === 5），需要根据舰种过滤
  const isAugmentSlot = slotIndex === 5
  
  // 获取所有装备
  if (allEquipments.value.length === 0) {
    try {
      const data = await getEquipmentList()
      allEquipments.value = data.equipments || []
      console.log('所有装备:', allEquipments.value)
    } catch (error) {
      console.error('获取装备列表失败:', error)
      ElMessage.error('获取装备列表失败')
      return
    }
  }
  
  // 根据可接受的类型过滤装备
  let filtered = allEquipments.value.filter(equip => 
    acceptableTypes.includes(equip.type)
  )
  
  // 如果是兵装栏，进一步根据舰种过滤
  if (isAugmentSlot) {
    filtered = filterAugmentsByShipType(filtered, currentShip.value.shipType)
  }
  
  filteredEquipments.value = filtered
  
  console.log('过滤后的装备:', filteredEquipments.value)
  
  if (filteredEquipments.value.length === 0) {
    ElMessage.warning('该栏位暂无可用装备')
    return
  }
  
  selectEquipDialogVisible.value = true
}

// 装备被选中
const handleEquipmentSelected = async (equipment) => {
  try {
    console.log('装备舰船 - 舰船ID:', currentShip.value.id, '栏位:', currentSlotIndex.value, '装备ID:', equipment.id)
    
    const result = await equipShip(currentShip.value.id, {
      slotIndex: currentSlotIndex.value,
      equipmentId: equipment.id
    })
    
    console.log('装备结果:', result)
    
    ElMessage.success('装备成功')
    selectEquipDialogVisible.value = false
    
    // 重新获取舰船信息
    currentShip.value = await getShipById(currentShip.value.id)
    console.log('更新后的舰船信息:', currentShip.value)
  } catch (error) {
    console.error('装备失败:', error)
    ElMessage.error(`装备失败: ${error.message || error}`)
  }
}

// 卸下装备
const handleUnequip = async (slotIndex) => {
  try {
    await equipShip(currentShip.value.id, {
      slotIndex: slotIndex,
      equipmentId: 0 // 0表示卸下装备
    })
    
    ElMessage.success('卸下装备成功')
    
    // 重新获取舰船信息
    currentShip.value = await getShipById(currentShip.value.id)
  } catch (error) {
    console.error('卸下装备失败:', error)
    ElMessage.error('卸下装备失败')
  }
}

// 格式化装备类型
const formatEquipmentType = (type) => {
  // type 可能是数字或字符串
  const typeMap = {
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
    'SMALL_CALIBER_MAIN_GUN': '驱逐主炮',
    'MEDIUM_CALIBER_MAIN_GUN': '轻巡主炮',
    'LARGE_CALIBER_MAIN_GUN': '重巡主炮',
    'BATTLESHIP_MAIN_GUN': '战列主炮',
    'TORPEDO': '鱼雷',
    'SUBMARINE_TORPEDO': '潜艇鱼雷',
    'ANTI_AIR_GUN': '防空炮',
    'AUXILIARY': '设备',
    'DIVE_BOMBER': '轰炸机',
    'TORPEDO_BOMBER': '鱼雷机',
    'FIGHTER': '战斗机',
    'AUGMENT': '兵装'
  }
  return typeMap[type] || type
}

// 格式化可接受的装备类型
const formatAcceptableTypes = (types) => {
  return types.map(t => formatEquipmentType(t)).join(', ')
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.ship-list-container {
  .filter-card {
    margin-bottom: 20px;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
  }

  .table-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.equipment-container {
  .equipment-slot {
    padding: 15px;
    margin-bottom: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    background: #f9f9f9;
    
    &.augment-slot {
      background: #fff9e6;
      border-color: #ffd700;
    }
    
    .slot-label {
      font-weight: bold;
      color: #333;
      margin-bottom: 8px;
      font-size: 14px;
      
      .tip {
        font-size: 12px;
        color: #999;
        font-weight: normal;
      }
    }
    
    .slot-content {
      margin-bottom: 8px;
      min-height: 32px;
      display: flex;
      align-items: center;
    }
    
    .acceptable-types {
      font-size: 12px;
      color: #666;
      font-style: italic;
    }
  }
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>
