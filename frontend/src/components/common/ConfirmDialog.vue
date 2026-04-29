<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  visible: boolean
  title: string
  message: string
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <el-dialog :model-value="visible" :title="title" width="400px" :close-on-click-modal="false" @close="emit('cancel')">
    <p style="color: var(--text-secondary); font-size: 14px;">{{ message }}</p>
    <template #footer>
      <el-button @click="emit('cancel')">{{ t('common.cancel') }}</el-button>
      <el-button type="danger" @click="emit('confirm')">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>
