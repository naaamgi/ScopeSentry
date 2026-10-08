<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive } from 'vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElUpload,
  ElTooltip,
  ElMessage,
  UploadProps,
  UploadRawFile,
  UploadInstance,
  ElTag
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { Icon } from '@/components/Icon'
import { useIcon } from '@/hooks/web/useIcon'
import { BaseButton } from '@/components/Button'
import { getPocDataApi, deletePocDataApi, getPocDetailApi } from '@/api/poc'
import Detail from './components/Detail.vue'
import { Dialog } from '@/components/Dialog'
import { useUserStore } from '@/store/modules/user'
const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const dialogVisible = ref(false)
const search = ref('')
const handleSearch = () => {
  getList()
}
const nodeColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'name',
    label: t('poc.pocName'),
    minWidth: 70
  },
  {
    field: 'level',
    label: t('poc.level'),
    minWidth: 50,
    columnKey: 'level',
    formatter: (_: Recordable, __: TableColumn, levelValue: string) => {
      if (levelValue == null) {
        return <div></div>
      }
      let color = ''
      let flag = ''
      if (levelValue === 'critical') {
        color = 'var(--critical)'
        flag = t('poc.critical')
      } else if (levelValue === 'high') {
        color = 'var(--high)'
        flag = t('poc.high')
      } else if (levelValue === 'medium') {
        color = 'var(--medium)'
        flag = t('poc.medium')
      } else if (levelValue === 'low') {
        color = 'var(--info)'
        flag = t('poc.low')
      } else if (levelValue === 'info') {
        color = 'var(--success)'
        flag = t('poc.info')
      } else if (levelValue === 'unknown') {
        color = 'var(--text-muted)'
        flag = t('poc.unknown')
      }
      return (
        <ElRow gutter={20} style="width: 80%">
          <ElCol span={1}>
            <Icon icon="clarity:circle-solid" color={color} size={10} />
          </ElCol>
          <ElCol span={5}>
            <ElText type="info">{flag}</ElText>
          </ElCol>
        </ElRow>
      )
    },
    filters: [
      { text: t('poc.critical'), value: 'critical' },
      { text: t('poc.high'), value: 'high' },
      { text: t('poc.medium'), value: 'medium' },
      { text: t('poc.low'), value: 'low' },
      { text: t('poc.info'), value: 'info' },
      { text: t('poc.unknown'), value: 'unknown' }
    ]
  },
  {
    field: 'tags',
    label: 'TAG',
    fit: 'true',
    formatter: (_row: Recordable, __: TableColumn, tags: string[]) => {
      if (tags.length != 0) {
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
            {tags.map((product) => (
              <div
                key={product}
                onClick={() => changeTags('app', product)}
                style={{ cursor: 'pointer' }}
              >
                <ElTag type={'success'}>{product}</ElTag>
              </div>
            ))}
          </div>
        )
        // if (ProductsValue.length > 1) {
        //   let contentTool = ''
        //   if (Array.isArray(ProductsValue)) {
        //     // It's an array, you can use forEach
        //     ProductsValue.forEach((item, _) => {
        //       contentTool += `<div>${item}</div>`
        //     })
        //   } else {
        //     console.error('ProductsValue is not an array')
        //   }
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light" round>
        //         {ProductsValue[0]}
        //       </ElTag>
        //       <ElTooltip
        //         class="box-item"
        //         effect="dark"
        //         placement="top-start"
        //         content={contentTool}
        //         popper-class="tagtooltip"
        //         rawContent
        //       >
        //         <ElTag type="info" effect="plain" round style={'left:3px; position:relative'}>
        //           {t('asset.total')} {ProductsValue.length} {t('asset.p')}
        //         </ElTag>
        //       </ElTooltip>
        //     </div>
        //   )
        // } else {
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light">
        //         {ProductsValue[0]}
        //       </ElTag>
        //     </div>
        //   )
        // }
      }
    }
  },
  {
    field: 'time',
    label: t('node.createTime'),
    minWidth: 50
  },

  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 160,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <div class="poc-row-actions">
          <BaseButton type="primary" onClick={() => edit(row)}>
            {t('common.edit')}
          </BaseButton>
          <BaseButton type="danger" onClick={() => del(row)}>
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
    const res = await getPocDataApi(search.value, currentPage.value, pageSize.value, filter)
    return {
      list: res.data.list,
      total: res.data.total
    }
  }
})
const { loading, dataList, total, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}

const changeTags = (type, value) => {
  const key = `${type}=${value}`
  console.log(key)
  // dynamicTags.value = [...dynamicTags.value, key]
}

let pocForm = reactive({
  id: '',
  name: '',
  level: '',
  content: '',
  tags: []
})
const addPoc = async () => {
  pocForm.id = ''
  pocForm.name = ''
  pocForm.level = ''
  pocForm.content = ''
  pocForm.tags = []
  dialogVisible.value = true
}
const edit = async (data) => {
  pocForm.id = data.id
  pocForm.name = data.name
  pocForm.level = data.level
  pocForm.tags = data.tags
  const res = await getPocDetailApi(pocForm.id)
  pocForm.content = res.data.content
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await deletePocDataApi([data.id])
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
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deletePocDataApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmDelete = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}
const userStore = useUserStore()
const uploadHeaders = ref({ Authorization: `${userStore.getToken}` })

const upload = ref<UploadInstance>()
const handleExceed: UploadProps['onExceed'] = (files) => {
  upload.value!.clearFiles()
  const file = files[0] as UploadRawFile
  upload.value!.handleStart(file)
}

const handleUploadSuccess = (response) => {
  console.log(response)
  if (response.code === 200) {
    ElMessage.success('Upload succes')
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
const filter = reactive<{ [key: string]: any }>({})
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
</script>

<template>
  <ContentWrap>
    <div class="poc-toolbar">
      <div class="poc-search-toolbar app-list-search">
        <label for="poc-search">{{ t('poc.pocName') }}</label>
        <ElInput id="poc-search" v-model="search" :placeholder="t('common.inputText')" @keyup.enter="handleSearch" />
        <ElButton type="primary" :icon="searchicon" @click="handleSearch">{{ t('common.search') }}</ElButton>
      </div>
      <div class="poc-actions">
        <ElButton type="primary" @click="addPoc">{{ t('common.new') }}</ElButton>
        <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">{{ t('common.delete') }}</BaseButton>
        <ElTooltip :content="t('common.uploadMsg')" placement="top">
          <ElUpload
            ref="upload"
            action="/api/poc/data/import"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :limit="1"
            :on-exceed="handleExceed"
            :auto-upload="false"
            @change="handleFileChange"
          >
            <template #trigger><BaseButton>{{ t('plugin.import') }}</BaseButton></template>
          </ElUpload>
        </ElTooltip>
      </div>
    </div>
    <Table
      v-model:pageSize="pageSize"
      v-model:currentPage="currentPage"
      @filter-change="filterChange"
      :columns="nodeColums"
      :data="dataList"
      stripe
      :border="true"
      :loading="loading"
      :resizable="true"
      :pagination="{
        total: total,
        pageSizes: [10, 20, 50, 100, 200, 500, 1000]
      }"
      @register="tableRegister"
      :headerCellStyle="tableHeaderColor"
    />
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="pocForm.id ? $t('common.edit') : $t('common.new')"
    class="poc-dialog"
    width="min(860px, calc(100vw - 32px))"
    maxHeight="min(70vh, 650px)"
  >
    <Detail :closeDialog="closeDialog" :pocForm="pocForm" :getList="getList" />
  </Dialog>
</template>

<style scoped>
.poc-toolbar { display: flex; flex-direction: column; align-items: flex-start; gap: 10px; margin-bottom: 16px; }
.poc-actions, :deep(.poc-row-actions) { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; }
.poc-actions .el-button, :deep(.poc-row-actions .el-button) { margin-left: 0 !important; }
:deep(.poc-actions .el-upload) { display: flex; align-items: center; }
@media (max-width: 720px) { .poc-search-toolbar, .poc-actions { width: 100%; } }
</style>
