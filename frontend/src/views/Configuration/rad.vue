<script setup lang="ts">
import { ElButton } from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ElCard } from 'element-plus'
import { ref, onBeforeMount } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { getRadConfigurationApi, saveRadConfigurationApi } from '@/api/Configuration'
const { t } = useI18n()
const content = ref(``)
onBeforeMount(async () => {
  try {
    const res = await getRadConfigurationApi()

    if (res.code === 200) {
      content.value = res.data.content
    } else {
      console.error(`API request failed with status code ${res.code}`)
    }
  } catch (error) {
    console.error('An error occurred while fetching the rad config:', error)
  }
})
const confirmAdd = async () => {
  const confirmed = window.confirm(t('configuration.saveConfirm'))
  if (confirmed) {
    await save()
  }
}
const save = async () => {
  saveLoading.value = true
  try {
    await saveRadConfigurationApi(content.value)
  } finally {
    saveLoading.value = false
  }
}
const saveLoading = ref(false)
</script>

<template>
  <ElCard shadow="never" class="mb-20px config-panel">
    <template #header>
      <div class="config-header">
        <span>{{ t('configuration.rad') }}</span>
        <ElButton type="primary" @click="confirmAdd" :loading="saveLoading">{{ t('common.save') }}</ElButton>
      </div>
    </template>
    <p class="config-intro">웹 크롤러의 방문 범위와 동시 페이지 수를 조절합니다. 소형 노드는 <code>max_depth: 2</code>, <code>max_page_concurrent: 2</code>, <code>max_page_visit_per_site: 100</code>부터 시험하세요. 기존 YAML에서 숫자만 바꾸고, 범위 제한은 오른쪽 위 사용법을 확인하세요.</p>
    <codemirror
      v-model="content"
      :style="{ height: 'min(55vh, 560px)' }"
      :indent-with-tab="true"
      :tab-size="2"
    />
  </ElCard>
</template>

<style scoped>
.config-panel { max-width: 1200px; margin-inline: auto; }
.config-header { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.config-intro { margin: 0 0 14px; color: var(--el-text-color-regular); line-height: 1.6; }
.config-intro code { font-family: 'JetBrains Mono', monospace; }
:deep(.cm-editor) { font-family: 'JetBrains Mono', monospace; font-size: 12px; background: var(--bg-subtle); color: var(--text-primary); border: 1px solid var(--border-strong); border-radius: var(--radius); overflow: hidden; }
:deep(.cm-gutters) { background: var(--bg-subtle); color: var(--text-muted); border-right: 1px solid var(--border); }
:deep(.cm-activeLine), :deep(.cm-activeLineGutter) { background: var(--accent-bg); }
:deep(.cm-cursor) { border-left-color: var(--text-primary); }
</style>
