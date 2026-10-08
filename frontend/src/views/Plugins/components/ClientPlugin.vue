<script setup lang="tsx">
import { moduleBadgeStyle } from '@/styles/moduleColors'
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, onMounted, h, inject, computed, type Ref } from 'vue'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElMessageBox,
  ElTag,
  ElTooltip,
  ElScrollbar,
  UploadInstance,
  UploadProps,
  UploadRawFile,
  ElUpload,
  ElMessage,
  ElDropdownItem,
  ElDropdownMenu,
  ElDropdown,
  ElIcon,
  ElBadge,
  ElDrawer,
  ElSpace,
  ElSwitch
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import {
  cleanAllPluginLogApi,
  cleanPluginLogApi,
  deletePluginDataApi,
  getPluginDataApi,
  getPluginLogApi,
  reCheckPluginApi,
  reInstallPluginApi,
  uninstallPluginApi
} from '@/api/plugins'
import detail from './detail.vue'
import { useUserStore } from '@/store/modules/user'

const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  getList()
}

// 从父组件注入插件市场相关方法
const openMarketDialog = inject<() => void>('openMarketDialog', () => {})
const pendingPluginsCount = inject<Ref<number>>('pendingPluginsCount', ref(0))
const pendingPluginsCountValue = computed(() => pendingPluginsCount.value)

const taskColums = reactive<TableColumn[]>([
  {
    field: 'index',
    label: t('tableDemo.index'),
    type: 'index',
    minWidth: '15'
  },
  {
    field: 'selection',
    type: 'selection',
    minWidth: 55,
    selectable: (row) => !row.isSystem
  },
  {
    field: 'name',
    label: t('plugin.name'),
    formatter: (row, __: TableColumn, value: string) => {
      return (
        <a
          href={`https://plugin.scope-sentry.top/plugin/${row.hash}`}
          style="color: var(--accent); text-decoration: none;"
          target="_blank"
        >
          {value}
        </a>
      )
    }
  },
  {
    field: 'module',
    label: t('plugin.module'),
    formatter: (_, __: TableColumn, value: string) => {
      return <ElTag style={moduleBadgeStyle(value)}>{value}</ElTag>
    }
  },
  {
    field: 'version',
    label: t('plugin.version'),
    minWidth: 100
  },
  {
    field: 'parameter',
    label: t('plugin.parameter'),
    formatter: (row, __: TableColumn, value: string) => {
      return (
        <ElTooltip content={row.help} placement="top" effect="light">
          <span style="cursor: pointer;">{value}</span>
        </ElTooltip>
      )
    }
  },
  {
    field: 'introduction',
    label: t('plugin.introduction'),
    minWidth: 200
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    width: 450,
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      const handleCommand = (command) => {
        switch (command) {
          case 'reinstall':
            reInstallPluginApi('all', row.hash, row.module)
            break
          case 'recheck':
            reCheckPluginApi('all', row.hash, row.module)
            break
          case 'uninstall':
            uninstallPluginApi('all', row.hash, row.module)
            break
        }
      }
      const retestAndDeleteDropdown = h(
        ElDropdown,
        {
          onCommand: handleCommand
        },
        {
          default: () =>
            h(
              ElButton,
              {
                style: { outline: 'none', boxShadow: 'none' }
              },
              () => [
                t('common.operation'), // 下拉菜单触发按钮文字
                h(
                  ElIcon,
                  {},
                  () => h(ArrowDown) // 向下箭头图标
                )
              ]
            ),
          dropdown: () =>
            h(ElDropdownMenu, null, () => {
              return [
                h(ElDropdownItem, { command: 'reinstall' }, () => t('plugin.reInstall')),
                h(ElDropdownItem, { command: 'recheck' }, () => t('plugin.reCheck')),
                h(ElDropdownItem, { command: 'uninstall' }, () => t('plugin.uninstall'))
              ]
            })
        }
      )
      return (
        <div class="plugin-row-actions">
          {retestAndDeleteDropdown}
          <BaseButton type="warning" onClick={() => openLogDialogVisible(row)}>
            {t('common.log')}
          </BaseButton>
          <BaseButton type="info" onClick={() => confirmCleanLog(row.hash, row.module)}>
            {t('common.cleanLog')}
          </BaseButton>
          <BaseButton type="success" onClick={() => editPlugin(row.id)}>
            {t('common.edit')}
          </BaseButton>
          <BaseButton
            type="danger"
            onClick={() => confirmDelete(row.hash, row.module)}
            disabled={row.isSystem}
          >
            {t('common.delete')}
          </BaseButton>
        </div>
      )
    }
  }
])

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    // 不传递 type 参数，默认为扫描端插件
    const res = await getPluginDataApi(search.value, currentPage.value, pageSize.value, 'scan')
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: true
})
const { loading, dataList, total, currentPage, pageSize } = tableState
pageSize.value = 20
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const dialogVisible = ref(false)

let DialogTitle = t('plugin.new')
const closeDialog = () => {
  dialogVisible.value = false
}

const confirmDeleteSelect = async () => {
  ElMessageBox({
    title: 'Delete',
    draggable: true
  }).then(async () => {
    await delSelect()
  })
}

const confirmDelete = async (hash: string, module: string) => {
  ElMessageBox({
    title: 'Delete',
    draggable: true
  }).then(async () => {
    await del(hash, module)
  })
}

const confirmCleanLog = async (hash: string, module: string) => {
  ElMessageBox({
    title: 'Clean Log',
    message: 'Are you sure you want to clean the logs?',
    draggable: true
  }).then(async () => {
    await cleanPluginLogApi(module, hash)
  })
}

const confirmCleanAllLog = async () => {
  ElMessageBox({
    title: 'Clean All Plugin Logs',
    message: 'Are you sure you want to clean all plugin logs?',
    draggable: true
  }).then(async () => {
    await cleanAllPluginLogApi()
  })
}
const delLoading = ref(false)
const del = async (hash: string, module: string) => {
  delLoading.value = true
  try {
    const res = await deletePluginDataApi([{ hash, module }])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  const deleteData = selectedRows.map((row) => ({
    hash: row.hash,
    module: row.module
  }))

  delLoading.value = true
  try {
    const res = await deletePluginDataApi(deleteData)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}

const addPlugin = async () => {
  id.value = ''
  dialogVisible.value = true
}

const id = ref('')

const editPlugin = async (data) => {
  id.value = data
  DialogTitle = t('common.edit')
  dialogVisible.value = true
}
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})

const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.7
}
const logDialogVisible = ref(false)
const closeLogDialogVisible = () => {
  logDialogVisible.value = false
  logSearchText.value = ''
}
const logContent = ref('')
const logSearchText = ref('')
const logScrollbarRef = ref<InstanceType<typeof ElScrollbar>>()
const autoScroll = ref(true)

const logModule = ref('')
const logHash = ref('')
const openLogDialogVisible = async (data) => {
  logModule.value = data.module
  logHash.value = data.hash
  await refreshLog()
  logDialogVisible.value = true
  // 等待DOM更新后滚动到底部
  setTimeout(() => {
    scrollToBottom()
  }, 100)
}

const refreshLog = async () => {
  const res = await getPluginLogApi(logModule.value, logHash.value)
  logContent.value = res.data.data
  if (autoScroll.value) {
    setTimeout(() => {
      scrollToBottom()
    }, 50)
  }
}

const cleanLog = async () => {
  await cleanPluginLogApi(logModule.value, logHash.value)
  logContent.value = ''
  ElMessage.success(t('common.cleanLog') + ' ' + t('common.success'))
}

const scrollToBottom = () => {
  if (logScrollbarRef.value) {
    const scrollbar = logScrollbarRef.value
    const scrollbarEl = scrollbar.$el
    const wrapEl = scrollbarEl?.querySelector('.el-scrollbar__wrap')
    if (wrapEl) {
      wrapEl.scrollTop = wrapEl.scrollHeight
    }
  }
}

const copyLog = async () => {
  if (logContent.value) {
    try {
      await navigator.clipboard.writeText(logContent.value)
      ElMessage.success(t('common.copySuccess'))
    } catch (error) {
      ElMessage.error(t('common.copyFailed'))
    }
  }
}

// 过滤后的日志内容（用于搜索高亮）
const filteredLogContent = computed(() => {
  if (!logSearchText.value || !logContent.value) {
    return logContent.value
  }
  const searchText = logSearchText.value
  const lines = logContent.value.split('\n')
  return lines.filter((line) => line.toLowerCase().includes(searchText.toLowerCase())).join('\n')
})

// 高亮搜索关键词
const highlightLog = (text: string) => {
  if (!logSearchText.value || !text) {
    return text
  }
  const searchText = logSearchText.value
  const regex = new RegExp(`(${searchText})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}
const userStore = useUserStore()
const uploadHeaders = { Authorization: `${userStore.getToken}` }
const upload = ref<UploadInstance>()
const handleExceed: UploadProps['onExceed'] = (files) => {
  upload.value!.clearFiles()
  const file = files[0] as UploadRawFile
  upload.value!.handleStart(file)
}

const handleUploadSuccess = (response) => {
  console.log(response)
  if (response.code === 200) {
    ElMessage.success(t('common.uploadSuccess'))
  } else {
    ElMessage.error(response.message)
  }
  getList()
  upload.value?.clearFiles()
}
const handleFileChange = (_file, fileList) => {
  if (fileList.length > 0) {
    upload.value!.submit()
  }
}

// 暴露刷新列表方法给父组件
defineExpose({
  refreshList: getList
})
</script>

<template>
  <ContentWrap>
    <div class="plugin-search-toolbar app-list-search">
      <label for="client-plugin-search">{{ t('plugin.name') }}</label>
      <ElInput id="client-plugin-search" v-model="search" :placeholder="t('common.inputText')" @keyup.enter="handleSearch" />
      <ElButton type="primary" :icon="searchicon" @click="handleSearch">{{ t('common.search') }}</ElButton>
    </div>
    <ElRow :gutter="16" class="mt-4">
      <ElCol :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <div class="flex flex-wrap gap-3 items-center">
          <BaseButton type="primary" @click="addPlugin">
            {{ t('plugin.new') }}
          </BaseButton>

          <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
            {{ t('plugin.delete') }}
          </BaseButton>

          <ElBadge
            :value="pendingPluginsCountValue"
            :hidden="pendingPluginsCountValue === 0"
            :max="99"
          >
            <BaseButton type="success" @click="openMarketDialog">
              {{ t('plugin.market') }}
            </BaseButton>
          </ElBadge>

          <BaseButton type="warning" @click="confirmCleanAllLog">
            {{ t('common.cleanAllLog') }}
          </BaseButton>

          <ElUpload
            ref="upload"
            class="flex items-center"
            action="/api/plugin/import"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :limit="1"
            :on-exceed="handleExceed"
            :auto-upload="false"
            @change="handleFileChange"
          >
            <template #trigger>
              <BaseButton>
                <template #icon>
                  <Icon icon="iconoir:upload" />
                </template>
                {{ t('plugin.import') }}
              </BaseButton>
            </template>
          </ElUpload>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="taskColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
        :resizable="true"
        :pagination="{
          total: total,
          pageSizes: [20, 30, 50, 100, 200, 500, 1000]
        }"
        @register="tableRegister"
        :headerCellStyle="tableHeaderColor"
      />
    </div>
  </ContentWrap>
  <ElDrawer v-model="dialogVisible" size="50%" direction="rtl">
    <template #title>
      <div style="display: flex; align-items: center; gap: 16px; width: 100%">
        <span style="font-weight: 500; white-space: nowrap">{{ DialogTitle }}</span>
        <span style="color: var(--critical); font-size: 12px; font-weight: normal; line-height: 1.4">
          {{ t('plugin.parameterConfigTip') }}
        </span>
      </div>
    </template>
    <detail :closeDialog="closeDialog" :getList="getList" :id="id" tp="scan" />
  </ElDrawer>
  <ElDrawer
    v-model="logDialogVisible"
    :title="t('node.log')"
    size="80%"
    direction="rtl"
    :with-header="true"
  >
    <template #header>
      <div style="display: flex; align-items: center; justify-content: space-between; width: 100%">
        <span style="font-weight: 500">{{ t('node.log') }}</span>
        <ElSpace>
          <ElInput
            v-model="logSearchText"
            :placeholder="t('common.search')"
            clearable
            style="width: 200px"
          >
            <template #prefix>
              <ElIcon><Search /></ElIcon>
            </template>
          </ElInput>
          <ElSwitch v-model="autoScroll" :active-text="t('common.autoScroll')" inactive-text="" />
        </ElSpace>
      </div>
    </template>
    <div style="display: flex; flex-direction: column; height: 100%">
      <ElScrollbar ref="logScrollbarRef" style="flex: 1; height: 0">
        <div
          style="
            padding: 16px;
            background: var(--bg-elevated);
            color: var(--text-secondary);
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            line-height: 1.6;
            min-height: 100%;
            white-space: pre-wrap;
            word-wrap: break-word;
          "
        >
          <div v-if="logSearchText && filteredLogContent">
            <div
              v-for="(line, index) in filteredLogContent.split('\n')"
              :key="index"
              v-html="highlightLog(line)"
            ></div>
          </div>
          <div v-else-if="logContent">
            {{ logContent }}
          </div>
          <div v-else style="color: var(--text-muted); text-align: center; padding: 40px">
            {{ t('common.noData') }}
          </div>
        </div>
      </ElScrollbar>
      <div
        style="
          padding: 16px;
          border-top: 1px solid var(--el-border-color);
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <div style="color: var(--el-text-color-secondary); font-size: 12px">
          <span v-if="logSearchText">
            {{ t('common.searchResult') }}:
            {{ filteredLogContent.split('\n').filter((l) => l).length }}
            {{ t('common.lines') }}
          </span>
          <span v-else-if="logContent">
            {{ logContent.split('\n').filter((l) => l).length }} {{ t('common.lines') }}
          </span>
        </div>
        <ElSpace>
          <BaseButton @click="refreshLog" type="primary">
            {{ t('common.refresh') }}
          </BaseButton>
          <BaseButton @click="copyLog" type="info">
            {{ t('common.copy') }}
          </BaseButton>
          <BaseButton @click="cleanLog" type="danger">{{ t('common.cleanLog') }}</BaseButton>
          <BaseButton @click="closeLogDialogVisible">{{ t('common.off') }}</BaseButton>
        </ElSpace>
      </div>
    </div>
  </ElDrawer>
</template>

<style scoped lang="less">
.plugin-row-actions { display: flex; align-items: center; flex-wrap: nowrap; gap: 8px; }
.plugin-row-actions :deep(.el-button) { margin-left: 0 !important; }
</style>
