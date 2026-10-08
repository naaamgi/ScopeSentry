<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { reactive } from 'vue'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { getPluginInfoApi } from '@/api/node'
import type { pluginInfoData } from '@/api/node/types'
import { useIcon } from '@/hooks/web/useIcon'
import { BaseButton } from '@/components/Button'
import { reCheckPluginApi, reInstallPluginApi, uninstallPluginApi } from '@/api/plugins'
const { t } = useI18n()

const props = defineProps<{
  closeDialog: () => void
  name: string
}>()
const correctIcon = useIcon({ icon: 'icon-park:check-one' })
const errorIcon = useIcon({ icon: 'line-md:close-circle', color: 'var(--critical)' })

const progressColums = reactive<TableColumn[]>([
  {
    field: 'index',
    label: t('tableDemo.index'),
    type: 'index',
    minWidth: '15'
  },
  {
    field: 'name',
    label: t('plugin.name')
  },
  {
    field: 'install',
    label: 'Install',
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      return cellValue == '1' ? correctIcon : errorIcon
    }
  },
  {
    field: 'check',
    label: 'Check',
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      return cellValue == '1' ? correctIcon : errorIcon
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 200,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton type="warning" onClick={() => handleAction('reinstall', row)}>
            {t('plugin.reInstall')}
          </BaseButton>
          <BaseButton type="success" onClick={() => handleAction('recheck', row)}>
            {t('plugin.reCheck')}
          </BaseButton>
          <BaseButton type="danger" onClick={() => handleAction('uninstall', row)}>
            {t('plugin.uninstall')}
          </BaseButton>
        </>
      )
    }
  }
])
const { tableRegister, tableState } = useTable({
  fetchDataApi: async () => {
    const res = await getPluginInfoApi(props.name)
    return {
      list: res.data.list
    }
  },
  immediate: true
})
const { loading, dataList } = tableState
const handleAction = (type: string, row: pluginInfoData) => {
  switch (type) {
    case 'reinstall':
      reInstallPluginApi(props.name, row.hash, row.module)
      break
    case 'recheck':
      reCheckPluginApi(props.name, row.hash, row.module)
      break
    case 'uninstall':
      uninstallPluginApi(props.name, row.hash, row.module)
      break
  }
}
</script>
<template>
  <p class="node-plugin-note">{{ t('plugin.nodeUsageHint') }}</p>
  <Table
    @register="tableRegister"
    :columns="progressColums"
    :data="dataList"
    rowKey="hash"
    :loading="loading"
    :resizable="true"
    max-height="600"
    :tooltip-options="{
      offset: 1,
      showArrow: false,
      effect: 'dark',
      enterable: true,
      showAfter: 0,
      popperOptions: {},
      popperClass: 'test',
      placement: 'top',
      hideAfter: 0,
      disabled: false
    }"
    :style="{
      fontFamily:
        '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
    }"
  />
</template>

<style scoped>
.node-plugin-note { margin: 0 0 12px; color: var(--text-secondary); font-size: 13px; line-height: 1.5; }
</style>
