<script setup lang="ts">
import { ElButton } from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ElCard } from 'element-plus'
import { ref, onBeforeMount } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { getSubfinderConfigurationApi, saveSubfinderConfigurationApi } from '@/api/Configuration'
const { t } = useI18n()
const content = ref(``)
onBeforeMount(async () => {
  try {
    const res = await getSubfinderConfigurationApi()

    if (res.code === 200) {
      content.value = res.data.content
    } else {
      console.error(`API request failed with status code ${res.code}`)
    }
  } catch (error) {
    console.error('An error occurred while fetching the subfinder config:', error)
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
    await saveSubfinderConfigurationApi(content.value)
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
        <span>{{ t('configuration.subfinder') }}</span>
        <ElButton type="primary" @click="confirmAdd" :loading="saveLoading">{{ t('common.save') }}</ElButton>
      </div>
    </template>
    <p class="config-intro">서브도메인 수집 정보원의 API 키 목록입니다. <code>[]</code>는 키가 없다는 뜻이며, 모든 소스를 채울 필요는 없습니다. 계정이 있는 소스에만 키를 넣으세요. 예: <code>github:</code> 아래에 <code>- YOUR_TOKEN</code>을 한 줄 추가합니다. 자세한 형식은 오른쪽 위 사용법에 있습니다.</p>
    <codemirror
      v-model="content"
      :style="{ height: 'min(50vh, 500px)' }"
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
