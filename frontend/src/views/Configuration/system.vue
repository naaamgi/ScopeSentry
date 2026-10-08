<script setup lang="ts">
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElOption,
  ElSelect,
  ElText
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ElCard } from 'element-plus'
import { ref, reactive, onBeforeMount } from 'vue'
import notification from './components/notification.vue'
import Deduplication from './components/Deduplication.vue'
import { getSystemConfigurationApi, saveSystemConfigurationApi } from '@/api/Configuration'
import { REGION_PROFILES, useRegionStore } from '@/store/modules/region'
import { Codemirror } from 'vue-codemirror'
const { t } = useI18n()
const regionStore = useRegionStore()
const form = reactive({
  timezone: '',
  ModulesConfig: '',
  region: regionStore.getRegion
})
onBeforeMount(async () => {
  try {
    const res = await getSystemConfigurationApi()

    if (res.code == 200) {
      form.timezone = res.data.timezone
      form.ModulesConfig = res.data.ModulesConfig
      regionStore.setRegion(res.data.region)
      form.region = regionStore.getRegion
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
  const res = await saveSystemConfigurationApi(form.timezone, form.ModulesConfig, form.region)
  saveLoading.value = false
  if (res.code == 200) {
    // 지역이 바뀌면 자산 표의 컬럼과 탭이 달라지므로 화면을 다시 읽는다.
    const changed = regionStore.getRegion !== form.region
    regionStore.setRegion(form.region)
    if (changed) {
      window.location.reload()
    }
  }
}
const saveLoading = ref(false)
</script>

<template>
  <ElCard shadow="never" class="mb-20px">
    <template #header>
      <span>{{ t('configuration.system') }}</span>
    </template>
    <ElForm :model="form" label-position="top" class="system-form">
      <div class="system-basic-grid">
        <ElFormItem :label="t('configuration.timezone')">
          <ElInput v-model="form.timezone" placeholder="Asia/Seoul" />
        </ElFormItem>
        <ElFormItem :label="t('configuration.region')">
          <ElSelect v-model="form.region">
            <ElOption
              v-for="profile in REGION_PROFILES"
              :key="profile"
              :label="t(`configuration.regionProfile.${profile}`)"
              :value="profile"
            />
          </ElSelect>
          <ElText size="small" type="info">{{ t('configuration.regionMsg') }}</ElText>
        </ElFormItem>
      </div>
      <ElFormItem label="Module Config (YAML)" class="module-config-item">
        <div class="module-config-field">
          <p>스캔 단계별 동시 처리 수입니다. 기존 YAML을 유지하고 필요한 숫자만 바꾸세요. 소형 노드(2 vCPU·4GB)는 maxGoroutineCount 2, portScan·webCrawler·vulnerabilityScan 1부터 시험하세요. 자세한 항목 설명은 오른쪽 위 사용법에 있습니다.</p>
          <Codemirror
            v-model="form.ModulesConfig"
            :autofocus="true"
            :indent-with-tab="true"
            :tab-size="2"
            :style="{ height: 'min(55vh, 550px)', width: '100%' }"
          />
        </div>
      </ElFormItem>
      <div class="system-save-actions">
        <ElButton type="primary" @click="confirmAdd" :loading="saveLoading">{{ t('common.save') }}</ElButton>
        <ElText size="small" type="danger">{{ t('configuration.threadMsg') }}</ElText>
      </div>
    </ElForm>
  </ElCard>
  <notification />
  <Deduplication />
</template>

<style scoped>
.system-form { width: 100%; }
.system-basic-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 0 24px; }
.system-basic-grid .el-select { width: 100%; }
.module-config-field { width: 100%; min-width: 0; }
.module-config-field p { margin: 0 0 10px; color: var(--el-text-color-secondary); line-height: 1.55; }
.system-save-actions { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; }
:deep(.module-config-item .el-form-item__content) { display: block; min-width: 0; }
:deep(.module-config-field .cm-editor) { font-family: 'JetBrains Mono', monospace; font-size: 12px; background: var(--bg-subtle); color: var(--text-primary); border: 1px solid var(--border-strong); border-radius: var(--radius); overflow: hidden; }
:deep(.module-config-field .cm-gutters) { background: var(--bg-subtle); color: var(--text-muted); border-right: 1px solid var(--border); }
:deep(.module-config-field .cm-activeLine), :deep(.module-config-field .cm-activeLineGutter) { background: var(--accent-bg); }
:deep(.module-config-field .cm-cursor) { border-left-color: var(--text-primary); }
@media (max-width: 800px) { .system-basic-grid { grid-template-columns: 1fr; } }
</style>
