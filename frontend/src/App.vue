<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { ConfigGlobal } from '@/components/ConfigGlobal'
import { useDesign } from '@/hooks/web/useDesign'
import { useStorage } from '@/hooks/web/useStorage'
import { setCssVar } from './utils'
import { isDark } from '@/utils/is'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('app')

const appStore = useAppStore()

const currentSize = computed(() => appStore.getCurrentSize)

const greyMode = computed(() => appStore.getGreyMode)

const { getStorage } = useStorage('localStorage')

// 根据浏览器当前主题设置系统主题色
const setDefaultTheme = () => {
  if (getStorage('isDark') !== null) {
    // 如果用户已经手动设置过主题，使用用户设置
    const isDarkValue = getStorage('isDark')
    appStore.setIsDark(isDarkValue)
  } else {
    // 如果用户没有设置过，则根据系统主题自动设置
    const isDarkTheme = isDark()
    appStore.setIsDark(isDarkTheme)
  }

  // Replace legacy stored theme colors with the approved palette on every launch.
  appStore.setTheme({
    elColorPrimary: 'var(--accent)',
    leftMenuBorderColor: 'var(--border)',
    leftMenuBgColor: 'var(--bg-card)',
    leftMenuBgLightColor: 'var(--bg-subtle)',
    leftMenuBgActiveColor: 'var(--accent-bg)',
    leftMenuCollapseBgActiveColor: 'var(--accent-bg)',
    leftMenuTextColor: 'var(--text-secondary)',
    leftMenuTextActiveColor: 'var(--accent)',
    logoTitleTextColor: 'var(--text-primary)',
    logoBorderColor: 'var(--border)',
    topHeaderBgColor: 'var(--bg-card)',
    topHeaderTextColor: 'var(--text-primary)',
    topHeaderHoverColor: 'var(--bg-subtle)',
    topToolBorderColor: 'var(--border)'
  } as ThemeTypes)
  appStore.setCssVarTheme()

  // 初始化标签页高度
  const tagsViewEnabled = appStore.getTagsView
  setCssVar('--tags-view-height', tagsViewEnabled ? '35px' : '0px')
}

setDefaultTheme()
</script>

<template>
  <ConfigGlobal :size="currentSize">
    <RouterView :class="greyMode ? `${prefixCls}-grey-mode` : ''" />
  </ConfigGlobal>
</template>

<style lang="less">
@prefix-cls: ~'@{namespace}-app';

.size {
  width: 100%;
  height: 100%;
}

html,
body {
  padding: 0 !important;
  margin: 0;
  overflow: hidden;
  .size;

  #app {
    .size;
  }
}

.@{prefix-cls}-grey-mode {
  filter: grayscale(100%);
}
</style>
