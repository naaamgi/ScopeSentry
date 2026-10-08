<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElButton, ElPopover } from 'element-plus'
import { guides } from './helpGuides'

const route = useRoute()
const open = ref(false)
const guide = computed(() => {
  const key = Object.keys(guides).sort((a, b) => b.length - a.length).find((item) => route.path.startsWith(item))
  return key ? guides[key] : undefined
})
watch(() => route.path, () => { open.value = false })
</script>

<template>
  <ElPopover v-if="guide" v-model:visible="open" placement="bottom-end" :width="500" trigger="click">
    <template #reference><ElButton text class="help-trigger" :aria-expanded="open">사용법</ElButton></template>
    <article class="help-panel">
      <h2>{{ guide.title }}</h2>
      <p class="help-summary">{{ guide.summary }}</p>
      <section v-for="section in guide.sections" :key="section.heading">
        <h3>{{ section.heading }}</h3>
        <ul><li v-for="item in section.items" :key="item">{{ item }}</li></ul>
        <pre v-if="section.example"><code>{{ section.example }}</code></pre>
      </section>
    </article>
  </ElPopover>
</template>

<style scoped>
.help-trigger { font-weight: 650; color: var(--el-color-primary); }
.help-panel { max-height: min(72vh, 680px); overflow-y: auto; padding: 4px 8px 4px 2px; line-height: 1.55; color: var(--el-text-color-regular); }
.help-panel h2 { margin: 0 0 6px; color: var(--el-text-color-primary); font-size: 18px; }
.help-summary { margin: 0 0 16px; color: var(--el-text-color-secondary); }
.help-panel section + section { border-top: 1px solid var(--el-border-color); margin-top: 14px; padding-top: 14px; }
.help-panel h3 { margin: 0 0 6px; color: var(--el-text-color-primary); font-size: 14px; }
.help-panel ul { margin: 0; padding-left: 20px; }
.help-panel li + li { margin-top: 5px; }
.help-panel pre { margin: 10px 0 0; padding: 12px; overflow-x: auto; border-radius: 8px; background: var(--el-fill-color-light); font-size: 12px; }
</style>
