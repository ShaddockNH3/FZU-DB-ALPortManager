<script setup>
import { computed, reactive, ref, watch } from 'vue'
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect
} from 'naive-ui'
import { useShipStore } from '@/stores/ship'
import { message } from '@/utils/feedback'
import { FACTIONS, LEVEL_RANGE, RARITY_OPTIONS, SHIP_TYPES } from '@/constants/ship'

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
  shipName: { required: true, message: '请输入舰船名称', trigger: ['input', 'blur'] },
  rarity: { required: true, message: '请选择稀有度', trigger: ['change', 'blur'] },
  shipType: { required: true, message: '请选择舰种', trigger: ['change', 'blur'] },
  faction: { required: true, message: '请选择阵营', trigger: ['change', 'blur'] }
}

const rarityOptions = RARITY_OPTIONS.map((item) => ({ label: item.label, value: item.value }))
const shipTypeOptions = SHIP_TYPES.map((type) => ({ label: type, value: type }))
const factionOptions = FACTIONS.map((faction) => ({ label: faction, value: faction }))

const isEdit = computed(() => Boolean(form.id))
const title = computed(() => (isEdit.value ? '编辑舰船' : '添加舰船'))

const show = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

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
  formRef.value?.restoreValidation()
}

function close() {
  show.value = false
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
      message.success('更新成功')
    } else {
      await shipStore.createShip(payload)
      message.success('添加成功')
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
  <n-modal
    v-model:show="show"
    preset="card"
    :title="title"
    style="width: 560px; max-width: 92vw"
  >
    <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="88">
      <n-form-item label="舰船名称" path="shipName">
        <n-input v-model:value="form.shipName" placeholder="例如：企业" />
      </n-form-item>

      <n-form-item label="稀有度" path="rarity">
        <n-select v-model:value="form.rarity" :options="rarityOptions" placeholder="请选择" />
      </n-form-item>

      <n-form-item label="舰种" path="shipType">
        <n-select v-model:value="form.shipType" :options="shipTypeOptions" placeholder="请选择" />
      </n-form-item>

      <n-form-item label="阵营" path="faction">
        <n-select v-model:value="form.faction" :options="factionOptions" placeholder="请选择" />
      </n-form-item>

      <n-form-item label="等级" path="level">
        <n-input-number
          v-model:value="form.level"
          :min="LEVEL_RANGE.min"
          :max="LEVEL_RANGE.max"
          style="width: 100%"
        />
      </n-form-item>
    </n-form>

    <template #footer>
      <div class="modal-footer">
        <n-button quaternary @click="close">取消</n-button>
        <n-button type="primary" :loading="submitLoading" @click="submit">确定</n-button>
      </div>
    </template>
  </n-modal>
</template>

<style lang="scss" scoped>
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
