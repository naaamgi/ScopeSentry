<script setup lang="ts">
import {
  ElFormItem,
  ElForm,
  ElButton,
  ElMessage
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { Codemirror } from 'vue-codemirror'
import { ref } from 'vue'
import { toRefs } from '@vueuse/core'
import { updatePocDataApi, addPocDataApi } from '@/api/poc'
const { t } = useI18n()
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  pocForm: {
    id: string
    name: string
    content: string
    level: string
    tags: string[]
  }
}>()
const { pocForm } = toRefs(props)
const localForm = ref({ ...pocForm.value })

const exampleTemplate = `id: example-response-marker

info:
  name: Example response marker
  author: your-name
  severity: info
  description: Replace the marker with a response value you expect.
  tags: example

http:
  - method: GET
    path:
      - "{{BaseURL}}/"
    matchers:
      - type: word
        part: body
        words:
          - "REPLACE_WITH_UNIQUE_MARKER"
`
const saveLoading = ref(false)
const submitForm = async () => {
  if (!localForm.value.content.trim()) {
    ElMessage.warning(t('poc.contentMsg'))
    return
  }
  saveLoading.value = true
  try {
    const res = localForm.value.id
      ? await updatePocDataApi(localForm.value.id, localForm.value.name, localForm.value.content, localForm.value.level, localForm.value.tags)
      : await addPocDataApi(localForm.value.name, localForm.value.content, localForm.value.level, localForm.value.tags)
    if (res.code === 200) { props.getList(); props.closeDialog() }
  } catch (error) {
    console.error('POC save failed:', error)
  } finally {
    saveLoading.value = false
  }
}
</script>
<template>
  <ElForm :model="localForm" label-position="top" status-icon class="poc-editor-form">
    <div class="poc-editor-intro">
      <strong>{{ t('poc.editorFormat') }}</strong>
      <p>{{ t('poc.editorHelp') }}</p>
      <ElButton size="small" @click="localForm.content = exampleTemplate">{{ t('poc.insertExample') }}</ElButton>
    </div>
    <ElFormItem :label="t('poc.content')" prop="content">
      <codemirror
        v-model="localForm.content"
        class="poc-code-editor"
        :style="{ height: 'min(48vh, 430px)', width: '100%' }"
        :indent-with-tab="true"
        :tab-size="2"
      />
    </ElFormItem>
    <div class="poc-editor-actions">
      <ElButton @click="props.closeDialog()">{{ t('common.cancel') }}</ElButton>
      <ElButton type="primary" @click="submitForm" :loading="saveLoading">{{ t('task.save') }}</ElButton>
    </div>
  </ElForm>
</template>

<style scoped>
.poc-editor-intro { padding: 14px 16px; margin-bottom: 16px; background: var(--bg-subtle); border: 1px solid var(--border); border-radius: var(--radius); }
.poc-editor-intro strong { color: var(--text-primary); }
.poc-editor-intro p { margin: 6px 0 10px; color: var(--text-secondary); line-height: 1.5; }
.poc-editor-actions { display: flex; justify-content: flex-end; gap: 8px; padding-top: 14px; border-top: 1px solid var(--border); }
.poc-editor-actions .el-button { margin-left: 0; }
:deep(.poc-code-editor .cm-editor) { font-family: 'JetBrains Mono', monospace; font-size: 12px; background: var(--bg-subtle); color: var(--text-primary); border: 1px solid var(--border-strong); border-radius: var(--radius); overflow: hidden; }
:deep(.poc-code-editor .cm-gutters) { background: var(--bg-subtle); color: var(--text-muted); border-right: 1px solid var(--border); }
:deep(.poc-code-editor .cm-activeLine), :deep(.poc-code-editor .cm-activeLineGutter) { background: var(--accent-bg); }
:deep(.poc-code-editor .cm-cursor) { border-left-color: var(--text-primary); }
:deep(.poc-code-editor .cm-selectionBackground) { background: var(--accent-bg) !important; }
</style>
