<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import {
  ElButton,
  ElTabPane,
  ElTabs,
  ElInput,
  ElSwitch,
  ElDropdown,
  ElDropdownMenu,
  ElDropdownItem,
  ElIcon,
  ElMessageBox
} from 'element-plus'
import ProjectList from './components/ProjectList.vue'
import AddProject from './components/AddProject.vue'
import { useI18n } from '@/hooks/web/useI18n'
import { h, reactive, ref } from 'vue'
import { Dialog } from '@/components/Dialog'
import { deleteProjectApi, getProjectDataApi } from '@/api/project'
import { useIcon } from '@/hooks/web/useIcon'
const { t } = useI18n()
let allProjectData = reactive({})
let tabNames = ref<string[]>([])
let tagNum = reactive({})
const projectListLoading = ref(false)
const getProjectTag = async (pageIndex: number, pageSize: number) => {
  if (pageIndex === 0) {
    pageIndex = currentPage.value
    pageSize = currentpageSize.value
  } else {
    currentPage.value = pageIndex
    currentpageSize.value = pageSize
  }
  try {
    const res = await getProjectDataApi(search.value, pageIndex, pageSize)
    // 更新响应式对象
    Object.assign(allProjectData, res.data.result)
    tabNames.value = Object.keys(res.data.tag)
    Object.assign(tagNum, res.data.tag)
    const index = tabNames.value.indexOf('All')
    if (index !== -1) {
      tabNames.value.splice(index, 1)
    }
  } catch (error) {
    console.error('An error occurred:', error)
  }
}
const dialogVisible = ref(false)
const addProject = async () => {
  dialogVisible.value = true
}
const closeDialog = () => {
  dialogVisible.value = false
}
const search = ref('')
const searchicon = useIcon({ icon: 'iconoir:search' })
const currentPage = ref(1)
const currentpageSize = ref(50)

const loading = ref(false)
const handleSearch = async () => {
  loading.value = true
  projectListLoading.value = true
  await getProjectTag(currentPage.value, currentpageSize.value)
  loading.value = false
  projectListLoading.value = false
}
handleSearch()
const multipleSelection = ref(false)
const deleteicon = useIcon({ icon: 'openmoji:delete' })
const delSelect = async () => {
  const deleteAsset = ref<boolean>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAsset.value,
          'onUpdate:modelValue': (val: boolean) => {
            deleteAsset.value = val
          }
        })
      ])
  }).then(async () => {
    await deleteProjectApi(selectedRowIds.value, deleteAsset.value)
    getProjectTag(currentPage.value, currentpageSize.value)
  })
}
const elDropdownicon = useIcon({ icon: 'ri:arrow-drop-down-line' })
const selectedRowIds = ref([])
</script>

<template>
  <ContentWrap>
    <div class="project-toolbar">
      <div class="app-list-search">
        <label for="project-search">{{ t('form.input') }}:</label>
        <ElInput id="project-search" v-model="search" :placeholder="t('common.inputText')" @keyup.enter="handleSearch" />
        <ElButton :loading="loading" type="primary" :icon="searchicon" @click="handleSearch">
          {{ t('common.search') }}
        </ElButton>
      </div>
      <div class="project-actions">
        <ElButton type="primary" @click="addProject">{{ t('project.addProject') }}</ElButton>
        <label class="project-multi-select">
          {{ t('common.multipleSelection') }}
          <ElSwitch v-model="multipleSelection" inline-prompt active-text="Yes" inactive-text="No" />
        </label>
        <ElDropdown v-if="multipleSelection" trigger="click">
          <ElButton plain class="custom-button align-bottom">
            {{ t('common.operation') }}
            <ElIcon class="el-icon--right"><elDropdownicon /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem :icon="deleteicon" @click="delSelect">{{
                t('common.delete')
              }}</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </div>
    </div>
    <ElTabs class="demo-tabs" v-loading="projectListLoading">
      <ElTabPane :label="`All (${tagNum['All']})`"
        ><ProjectList
          :tableDataList="allProjectData['All']"
          :getProjectTag="getProjectTag"
          :total="tagNum['All']"
          :multipleSelection="multipleSelection"
          v-model:selectedRows="selectedRowIds"
      /></ElTabPane>
      <ElTabPane
        v-for="tagName in tabNames"
        :label="`${tagName} (${tagNum[tagName]})`"
        :key="tagName"
        ><ProjectList
          :tableDataList="allProjectData[tagName]"
          :getProjectTag="getProjectTag"
          :total="tagNum[tagName]"
          :multipleSelection="multipleSelection"
          v-model:selectedRows="selectedRowIds"
      /></ElTabPane>
    </ElTabs>
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="t('project.addProject')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <AddProject
      :closeDialog="closeDialog"
      projectid=""
      :getProjectData="getProjectTag"
      :schedule="false"
    />
  </Dialog>
</template>
<style scoped>
.project-toolbar { display: grid; gap: 12px; margin-bottom: 14px; }
.project-actions { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; }
.project-actions :deep(.el-button) { margin-left: 0; }
.project-multi-select { display: inline-flex; align-items: center; gap: 8px; color: var(--text-secondary); font-size: 14px; white-space: nowrap; }
.demo-tabs :deep(.el-tabs__content) { padding: 16px 0 0; color: var(--text-secondary); font-size: 14px; font-weight: 400; }
</style>
