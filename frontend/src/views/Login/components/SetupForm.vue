<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus'
import { setupApi } from '@/api/login'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()
const emit = defineEmits<{ (event: 'complete'): void }>()
const form = reactive({ username: '', password: '', confirmPassword: '' })
const loading = ref(false)

const submit = async () => {
  if (form.username.trim().length < 3 || form.username.trim().length > 64) {
    ElMessage.warning(t('setup.usernameHint'))
    return
  }
  if (form.password.length < 12 || form.password.length > 72) {
    ElMessage.warning(t('setup.passwordHint'))
    return
  }
  if (form.password !== form.confirmPassword) {
    ElMessage.warning(t('setup.passwordMismatch'))
    return
  }
  loading.value = true
  try {
    const res = await setupApi({ username: form.username.trim(), password: form.password })
    if (res?.code === 200) {
      ElMessage.success(t('setup.complete'))
      emit('complete')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="setup-card">
    <h1>{{ t('setup.title') }}</h1>
    <p>{{ t('setup.description') }}</p>
    <ElForm label-position="top" @submit.prevent="submit">
      <ElFormItem :label="t('login.username')">
        <ElInput v-model="form.username" autocomplete="username" maxlength="64" />
      </ElFormItem>
      <ElFormItem :label="t('login.password')">
        <ElInput v-model="form.password" type="password" show-password autocomplete="new-password" />
      </ElFormItem>
      <ElFormItem :label="t('setup.confirmPassword')">
        <ElInput v-model="form.confirmPassword" type="password" show-password autocomplete="new-password" />
      </ElFormItem>
      <ElButton type="primary" native-type="submit" :loading="loading" class="setup-submit">
        {{ t('setup.create') }}
      </ElButton>
    </ElForm>
  </section>
</template>

<style scoped>
.setup-card { width: 100%; padding: 32px; background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--el-box-shadow-light); }
.setup-card h1 { margin: 0 0 8px; font-size: 24px; color: var(--el-text-color-primary); }
.setup-card p { margin: 0 0 24px; color: var(--el-text-color-secondary); line-height: 1.6; }
.setup-submit { width: 100%; height: 42px; }
</style>
