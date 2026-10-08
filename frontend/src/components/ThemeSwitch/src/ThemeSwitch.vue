<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { ElSwitch } from 'element-plus'
import { useIcon } from '@/hooks/web/useIcon'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('theme-switch')

const Sun = useIcon({ icon: 'emojione-monotone:sun', color: 'var(--medium)' })

const CrescentMoon = useIcon({ icon: 'emojione-monotone:crescent-moon', color: 'var(--medium)' })

const appStore = useAppStore()

// 初始化获取是否是暗黑主题
const isDark = computed({
  get() {
    return appStore.getIsDark
  },
  set(val: boolean) {
    appStore.setIsDark(val)
  }
})

// 设置switch的背景颜色
const switchColor = 'var(--border-strong)'
</script>

<template>
  <ElSwitch
    :class="prefixCls"
    v-model="isDark"
    inline-prompt
    :border-color="switchColor"
    :inactive-color="switchColor"
    :active-color="switchColor"
    :active-icon="Sun"
    :inactive-icon="CrescentMoon"
  />
</template>

<style lang="less" scoped>
:deep(.el-switch__core .el-switch__inner .is-icon) {
  overflow: visible;
}
</style>
