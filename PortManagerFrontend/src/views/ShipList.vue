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
        <el-table-column prop="level" label="等级" width="100" sortable />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, RefreshRight, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { useShipStore } from '@/stores/ship'

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
  level: 1
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
</style>
