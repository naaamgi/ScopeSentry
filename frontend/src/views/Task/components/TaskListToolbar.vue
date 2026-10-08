<script setup lang="ts">
import { ElButton, ElInput } from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { useIcon } from '@/hooks/web/useIcon'

defineProps<{
  modelValue: string
  label: string
  searchId: string
  placeholder: string
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
  (event: 'search'): void
}>()

const { t } = useI18n()
const searchIcon = useIcon({ icon: 'iconoir:search' })
</script>

<template>
  <div class="task-list-toolbar">
    <div class="task-list-search app-list-search">
      <label :for="searchId">{{ label }}:</label>
      <ElInput
        :id="searchId"
        :model-value="modelValue"
        :placeholder="placeholder"
        @update:model-value="emit('update:modelValue', $event)"
        @keyup.enter="emit('search')"
      />
      <ElButton type="primary" :icon="searchIcon" @click="emit('search')">
        {{ t('common.search') }}
      </ElButton>
    </div>
    <div class="task-list-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.task-list-toolbar {
  display: grid;
  gap: 12px;
  margin-bottom: 14px;
}

.task-list-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.task-list-actions :deep(.el-button) {
  margin-left: 0;
  white-space: nowrap;
}

</style>
