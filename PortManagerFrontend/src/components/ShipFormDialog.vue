<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useShipStore } from '@/stores/ship'
import { LEVEL_RANGE, RARITY_OPTIONS, SHIP_TYPES, FACTIONS } from '@/constants/ship'

/**
 * 新建 / 编辑舰船的对话框。
 *
 * ship 传 null 表示新建，传对象表示编辑；打开时把数据灌进表单，
 * 保存成功后发 saved，由父级刷新列表。
 */
const props = defineProps({
  visible: { type: Boolean, default: false },
  ship: { type: Object, default: null }
})

const emit = defineEmits(['update:visible', 'saved'])

const shipStore = useShipStore()
const formRef = ref()
const submitLoading = ref(false)

const form = reactive({
  id: null,
  shipName: '',
  rarity: '',
  shipType: '',
  faction: '',
  level: 1
})

const rules = {
  shipName: [{ required: true, message: '请输入舰船名称', trigger: 'blur' }],
  rarity: [{ required: true, message: '请选择稀有度', trigger: 'change' }],
  shipType: [{ required: true, message: '请选择舰种', trigger: 'change' }],
  faction: [{ required: true, message: '请选择阵营', trigger: 'change' }]
}

const isEdit = computed(() => Boolean(form.id))
const title = computed(() => (isEdit.value ? '编辑舰船' : '添加舰船'))

watch(
  () => props.visible,
  (open) => {
    if (open) fillForm(props.ship)
  }
)

function fillForm(ship) {
  Object.assign(form, {
    id: ship?.id ?? null,
    shipName: ship?.shipName ?? '',
    rarity: ship?.rarity ?? '',
    shipType: ship?.shipType ?? '',
    faction: ship?.faction ?? '',
    level: ship?.level || 1
  })
  formRef.value?.clearValidate()
}

function close() {
  emit('update:visible', false)
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  const payload = {
    shipName: form.shipName,
    rarity: form.rarity,
    shipType: form.shipType,
    faction: form.faction,
    level: form.level
  }

  submitLoading.value = true
  try {
    if (isEdit.value) {
      await shipStore.updateShip(form.id, payload)
      ElMessage.success('更新成功')
    } else {
      await shipStore.createShip(payload)
      ElMessage.success('添加成功')
    }
    close()
    emit('saved')
  } catch (error) {
    // 提示由拦截器统一弹，这里只留日志
    console.error('提交失败:', error)
  } finally {
    submitLoading.value = false
  }
}
</script>

<template>
  <el-dialog :model-value="visible" :title="title" width="600px" @update:model-value="close">
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="舰船名称" prop="shipName">
        <el-input v-model="form.shipName" placeholder="请输入舰船名称" />
      </el-form-item>

      <el-form-item label="稀有度" prop="rarity">
        <el-select v-model="form.rarity" placeholder="请选择稀有度" class="full-width">
          <el-option
            v-for="item in RARITY_OPTIONS"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="舰种" prop="shipType">
        <el-select v-model="form.shipType" placeholder="请选择舰种" class="full-width">
          <el-option v-for="type in SHIP_TYPES" :key="type" :label="type" :value="type" />
        </el-select>
      </el-form-item>

      <el-form-item label="阵营" prop="faction">
        <el-select v-model="form.faction" placeholder="请选择阵营" class="full-width">
          <el-option v-for="faction in FACTIONS" :key="faction" :label="faction" :value="faction" />
        </el-select>
      </el-form-item>

      <el-form-item label="等级" prop="level">
        <el-input-number
          v-model="form.level"
          :min="LEVEL_RANGE.min"
          :max="LEVEL_RANGE.max"
          class="full-width"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="submit">确定</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.full-width {
  width: 100%;
}
</style>
