<script setup lang="ts">
import {
  ElRow,
  ElCol,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElOption,
  ElSelect,
  ElText,
  ElDivider
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ElCard } from 'element-plus'
import { ref, reactive, onBeforeMount, onMounted } from 'vue'
import notification from './components/notification.vue'
import Deduplication from './components/Deduplication.vue'
import { getSystemConfigurationApi, saveSystemConfigurationApi } from '@/api/Configuration'
import { REGION_PROFILES, useRegionStore } from '@/store/modules/region'
import { Codemirror } from 'vue-codemirror'
import { javascript } from '@codemirror/lang-javascript'
import { oneDark } from '@codemirror/theme-one-dark'
const extensions = [javascript(), oneDark]
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
      <ElRow>
        <ElCol :span="3" style="height: 100%">
          <span>{{ t('configuration.system') }}</span>
        </ElCol>
      </ElRow>
    </template>
    <ElForm :model="form" label-width="auto" style="max-width: 600px">
      <ElFormItem :label="t('configuration.timezone')">
        <ElInput v-model="form.timezone" />
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
      <ElFormItem label="Module Config">
        <Codemirror
          v-model="form.ModulesConfig"
          :extensions="extensions"
          :autofocus="true"
          :indent-with-tab="true"
          :tab-size="2"
          :style="{ height: '550px', width: '100%' }"
        />
      </ElFormItem>
    </ElForm>

    <ElRow>
      <!-- <ElCol :span="20" :offset="12"/> -->
      <ElCol :span="12" :offset="2">
        <ElButton type="primary" @click="confirmAdd" :loading="saveLoading">Save</ElButton>
        <ElDivider direction="vertical" />
        <ElText size="small" type="danger">{{ t('configuration.threadMsg') }}</ElText>
      </ElCol>
    </ElRow>
  </ElCard>
  <notification />
  <Deduplication />
</template>

<style scoped>
.header-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>
