<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h, onMounted } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import {
  ElButton,
  ElInput,
  ElProgress,
  ElTag,
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElDropdown,
  ElDropdownMenu,
  ElDropdownItem,
  ElIcon,
  ElRadioGroup,
  ElRadioButton,
  ElSelect,
  ElOption,
  ElDialog
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import {
  getTaskDataApi,
  deleteTaskApi,
  retestTaskApi,
  stopTaskApi,
  starTaskApi,
  syancProjectApi
} from '@/api/task'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import AddTask from './components/AddTask.vue'
import ProgressInfo from './components/ProgressInfo.vue'
import TaskListToolbar from './components/TaskListToolbar.vue'
import { useRouter } from 'vue-router'
import { getProjectAllApi } from '@/api/project'
const { push } = useRouter()
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  getList()
}
const taskColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    minWidth: 55
  },
  {
    field: 'name',
    label: t('task.taskName'),
    minWidth: 100
  },
  {
    field: 'taskNum',
    label: t('task.taskCount'),
    minWidth: 70,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          type: 'info'
        },
        () => cellValue
      )
    }
  },
  {
    field: 'progress',
    label: t('task.taskProgress'),
    minWidth: 200,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(ElProgress, {
        percentage: cellValue,
        type: 'line',
        striped: true,
        status: cellValue < 100 ? '' : 'success',
        stripedFlow: cellValue < 100 ? true : false
      })
    }
  },
  {
    field: 'status',
    label: t('common.state'),
    minWidth: 200,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      // 1为运行中，2为暂停，3为运行完成
      let tagType, tagText
      switch (cellValue) {
        case 1:
          tagType = 'info'
          tagText = t('task.running') // 运行中
          break
        case 2:
          tagType = 'warning'
          tagText = t('task.stop') // 暂停
          break
        case 3:
          tagType = 'success'
          tagText = t('task.finish') // 运行完成
          break
        default:
          tagType = 'default'
          tagText = '' // 未知状态
      }
      return h(ElTag, { type: tagType }, () => tagText)
    }
  },
  {
    field: 'creatTime',
    minWidth: 200,
    label: t('task.createTime')
  },
  {
    field: 'endTime',
    label: t('task.endTime'),
    minWidth: 200,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      if (cellValue == '') {
        return '-'
      }
      return cellValue
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 420,
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      const handleCommand = (command) => {
        ids.value = []
        switch (command) {
          case 'retest':
            confirmRetest(row)
            break
          case 'delete':
            confirmDelete(row)
            break
          case 'stop':
            ids.value.push(row.id)
            stopTask(ids.value)
            break
          case 'start':
            ids.value.push(row.id)
            startTask(ids.value)
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
              // 根据 row.status 渲染不同的菜单项
              if (row.status === 3) {
                return [
                  h(ElDropdownItem, { command: 'retest' }, () => t('task.retest')),
                  h(ElDropdownItem, { command: 'delete' }, () => t('common.delete'))
                ]
              } else {
                return [
                  h(ElDropdownItem, { command: 'start' }, () => t('task.start')),
                  h(ElDropdownItem, { command: 'stop' }, () => t('task.stop')), // 如果是运行中，显示“停止”按钮
                  h(ElDropdownItem, { command: 'retest' }, () => t('task.retest')),
                  h(ElDropdownItem, { command: 'delete' }, () => t('common.delete'))
                ]
              }
            })
        }
      )
      return (
        <div class="task-row-actions">
          {retestAndDeleteDropdown}
          <BaseButton type="primary" onClick={() => getTaskResult(row.name)}>
            {t('task.result')}
          </BaseButton>
          <BaseButton type="success" onClick={() => getTaskContent(row)}>
            {t('common.view')}
          </BaseButton>
          <ElButton type="warning" onClick={() => getProgressInfo(row.id)}>
            {t('task.taskProgress')}
          </ElButton>
        </div>
      )
    }
  }
])

const progressDialogVisible = ref(false)
let getProgressInfoID = ''
const getProgressInfo = async (id) => {
  getProgressInfoID = id
  progressDialogVisible.value = true
}

const getTaskResult = async (id) => {
  push(`/asset-information/index?task=${id}`)
}

const stopTask = async (ids) => {
  console.log('begin stop')
  await stopTaskApi(ids)
}

const startTask = async (ids) => {
  console.log('begin start')
  await starTaskApi(ids)
}

const progresscloseDialog = () => {
  progressDialogVisible.value = false
}
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getTaskDataApi(search.value, currentPage.value, pageSize.value)
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
const addTask = async () => {
  taskid.value = ''
  DialogTitle = t('task.addTask')
  Create.value = true
  dialogVisible.value = true
}

let DialogTitle = t('task.addTask')
const closeDialog = () => {
  dialogVisible.value = false
}

let Create = ref(true)
const taskid = ref('')
const getTaskContent = async (data) => {
  taskid.value = data.id
  dialogVisible.value = true
  Create.value = false
  DialogTitle = t('common.view')
}
const confirmDeleteSelect = async () => {
  const deleteAssetS = ref<boolean | string | number>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAssetS.value,
          'onUpdate:modelValue': (val: boolean | string | number) => {
            deleteAssetS.value = val
          }
        })
      ])
  }).then(async () => {
    await delSelect(deleteAssetS.value)
  })
}

const confirmStopSelect = async () => {
  ElMessageBox({
    title: 'Stop Task',
    draggable: true
  }).then(async () => {
    await stopTaskSelect()
  })
}
const stopTaskSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    await stopTask(ids.value)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error Stop data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmStartSelect = async () => {
  ElMessageBox({
    title: 'Start Task',
    draggable: true
  }).then(async () => {
    await startTaskSelect()
  })
}
const startTaskSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    await startTask(ids.value)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error Stop data:', error)
    delLoading.value = false
    getList()
  }
}
interface ProjectOption { value: string; label: string }
const projectOptions = ref<ProjectOption[]>([])
const syncDialogVisible = ref(false)
const syncMode = ref<'existing' | 'new'>('existing')
const selectedProjectId = ref('')
const newProjectName = ref('')
const newProjectTag = ref('')
const syncTaskIds = ref<string[]>([])
const syncLoading = ref(false)

const confirmSyncToProjectSelect = async () => {
  const table = await getElTableExpose()
  syncTaskIds.value = (table?.getSelectionRows() || []).map((row) => row.id)
  if (syncTaskIds.value.length === 0) {
    ElMessage.warning('프로젝트에 동기화할 스캔 작업을 먼저 선택하세요.')
    return
  }
  syncMode.value = 'existing'
  selectedProjectId.value = ''
  newProjectName.value = ''
  newProjectTag.value = ''
  try {
    const res = await getProjectAllApi()
    projectOptions.value = (res.data?.list || []).flatMap((group: { label: string; children?: ProjectOption[] }) =>
      (group.children || []).map((project) => ({
        value: project.value,
        label: `${group.label} / ${project.label} · ${project.value.slice(-6)}`
      }))
    )
    syncDialogVisible.value = true
  } catch (error) {
    ElMessage.error('프로젝트 목록을 불러오지 못했습니다.')
  }
}

const syncToProject = async () => {
  if (syncMode.value === 'existing' && !selectedProjectId.value) {
    ElMessage.warning('동기화할 프로젝트를 선택하세요.')
    return
  }
  if (syncMode.value === 'new' && (!newProjectName.value.trim() || !newProjectTag.value.trim())) {
    ElMessage.warning('새 프로젝트 이름과 태그를 입력하세요.')
    return
  }
  syncLoading.value = true
  try {
    const res = await syancProjectApi(
      syncTaskIds.value, syncMode.value, selectedProjectId.value,
      newProjectTag.value.trim(), newProjectName.value.trim()
    )
    if (res.code === 200) {
      ElMessage.success('선택한 작업 대상을 프로젝트에 동기화했습니다.')
      syncDialogVisible.value = false
    } else {
      ElMessage.error(res.message || '프로젝트 동기화에 실패했습니다.')
    }
  } catch (error) {
    ElMessage.error('프로젝트 동기화에 실패했습니다. 프로젝트 설정을 확인하세요.')
  } finally {
    syncLoading.value = false
  }
}

const confirmDelete = async (data) => {
  const deleteAsset = ref<boolean | string | number>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAsset.value,
          'onUpdate:modelValue': (val: boolean | string | number) => {
            deleteAsset.value = val
          }
        })
      ])
  }).then(async () => {
    await del(data, deleteAsset.value)
  })
}
const delLoading = ref(false)
const del = async (data, delA) => {
  delLoading.value = true
  try {
    const res = await deleteTaskApi([data.id], delA)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const ids = ref<string[]>([])
const delSelect = async (delA) => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deleteTaskApi(ids.value, delA)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmRetest = async (data) => {
  const confirmed = window.confirm('Are you sure you want to retest?')
  if (confirmed) {
    await retestTask(data)
  }
}
const retestTask = async (data) => {
  try {
    await retestTaskApi(data.id)
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    getList()
  }
}
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})
const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.75
}
</script>

<template>
  <ContentWrap>
    <TaskListToolbar
      v-model="search"
      search-id="scan-task-search"
      :label="t('task.taskName')"
      :placeholder="t('common.inputText')"
      @search="handleSearch"
    >
      <template #actions>
        <BaseButton type="primary" @click="addTask">{{ t('task.addTask') }}</BaseButton>
        <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">{{ t('task.delTask') }}</BaseButton>
        <BaseButton type="warning" :loading="delLoading" @click="confirmStopSelect">{{ t('task.stop') }}</BaseButton>
        <BaseButton type="success" :loading="delLoading" @click="confirmStartSelect">{{ t('task.start') }}</BaseButton>
        <BaseButton type="info" :loading="delLoading" @click="confirmSyncToProjectSelect">{{ t('task.syncToProject') }}</BaseButton>
      </template>
    </TaskListToolbar>
    <div>
      <Table
        :tooltip-options="{
          offset: 1,
          showArrow: false,
          effect: 'dark',
          enterable: false,
          showAfter: 0,
          popperOptions: {},
          popperClass: 'test',
          placement: 'bottom',
          hideAfter: 0,
          disabled: true
        }"
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="taskColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
        :max-height="maxHeight"
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
  <ElDialog v-model="syncDialogVisible" :title="t('task.syncToProject')" width="520px" :close-on-click-modal="false">
    <div class="project-sync-form">
      <p>선택한 스캔 작업 {{ syncTaskIds.length }}개의 대상을 프로젝트에 추가합니다. 기존 대상은 유지됩니다.</p>
      <ElRadioGroup v-model="syncMode">
        <ElRadioButton label="existing">{{ t('task.syncToExisting') }}</ElRadioButton>
        <ElRadioButton label="new">{{ t('task.createNewProject') }}</ElRadioButton>
      </ElRadioGroup>
      <label v-if="syncMode === 'existing'">
        <span>프로젝트</span>
        <ElSelect v-model="selectedProjectId" filterable placeholder="태그 / 프로젝트명으로 검색" style="width: 100%">
          <ElOption v-for="item in projectOptions" :key="item.value" :value="item.value" :label="item.label" />
        </ElSelect>
      </label>
      <template v-else>
        <label><span>프로젝트 이름</span><ElInput v-model="newProjectName" :placeholder="t('project.msgProject')" /></label>
        <label><span>태그</span><ElInput v-model="newProjectTag" :placeholder="t('project.msgProjectTag')" /></label>
      </template>
    </div>
    <template #footer>
      <ElButton @click="syncDialogVisible = false">{{ t('common.cancel') }}</ElButton>
      <ElButton type="primary" :loading="syncLoading" @click="syncToProject">{{ t('common.confirmed') }}</ElButton>
    </template>
  </ElDialog>
  <Dialog
    v-model="dialogVisible"
    :title="DialogTitle"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <AddTask
      :closeDialog="closeDialog"
      :getList="getList"
      :create="Create"
      :taskid="taskid"
      :schedule="false"
      tp="scan"
      :targetIds="[]"
    />
  </Dialog>
  <Dialog
    v-model="progressDialogVisible"
    :title="t('task.taskProgress')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    width="70%"
    max-height="700"
  >
    <ProgressInfo
      :closeDialog="progresscloseDialog"
      :getProgressInfoID="getProgressInfoID"
      getProgressInfotype="scan"
      getProgressInforunnerid=""
  /></Dialog>
</template>

<style scoped>
.project-sync-form { display: grid; gap: 16px; }
.project-sync-form p { margin: 0; color: var(--text-secondary); line-height: 1.5; }
.project-sync-form label { display: grid; gap: 6px; min-width: 0; }
.project-sync-form label span { color: var(--text-primary); font-size: 13px; font-weight: 600; }
</style>

<style scoped>
:deep(.task-row-actions) { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; }
:deep(.task-row-actions .el-button) { margin-left: 0 !important; }
</style>
