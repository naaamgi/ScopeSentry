<script setup lang="tsx">
import { moduleBackgroundMap, moduleBadgeStyle } from '@/styles/moduleColors'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, computed } from 'vue'
import {
  ElDrawer,
  ElRow,
  ElCol,
  ElInput,
  ElSelect,
  ElOption,
  ElCard,
  ElTag,
  ElDialog,
  ElMessageBox,
  ElRadioGroup,
  ElRadioButton
} from 'element-plus'
import { BaseButton } from '@/components/Button'
import { Icon } from '@/components/Icon'
import type { RemotePluginData } from '@/api/plugins/types'

interface PluginMarketProps {
  visible: boolean
  remotePluginList: RemotePluginData[]
  marketLoading: boolean
  onClose: () => void
  onRefresh: () => void
  onInstall: (plugin: RemotePluginData, token?: string) => Promise<void>
}

const props = defineProps<PluginMarketProps>()

const { t } = useI18n()

// 筛选条件
const marketSearch = ref('')
const selectedModule = ref('')
const selectedPriceStatus = ref<number | string>('')
const selectedInstallStatus = ref<string>('')
const selectedType = ref<string>('') // '' 表示全部, 'scan' 表示扫描, 'server' 表示服务端

// 模块选项
const moduleOptions = [
  { label: 'TargetHandler', value: 'TargetHandler' },
  { label: 'SubdomainScan', value: 'SubdomainScan' },
  { label: 'SubdomainSecurity', value: 'SubdomainSecurity' },
  { label: 'PortScanPreparation', value: 'PortScanPreparation' },
  { label: 'PortScan', value: 'PortScan' },
  { label: 'AssetMapping', value: 'AssetMapping' },
  { label: 'URLScan', value: 'URLScan' },
  { label: 'WebCrawler', value: 'WebCrawler' },
  { label: 'DirScan', value: 'DirScan' },
  { label: 'VulnerabilityScan', value: 'VulnerabilityScan' },
  { label: 'AssetHandle', value: 'AssetHandle' },
  { label: 'PortFingerprint', value: 'PortFingerprint' },
  { label: 'URLSecurity', value: 'URLSecurity' },
  { label: 'PassiveScan', value: 'PassiveScan' }
]

const resetFilters = () => {
  marketSearch.value = ''
  selectedModule.value = ''
  selectedPriceStatus.value = ''
  selectedInstallStatus.value = ''
  selectedType.value = ''
}

const filteredRemotePlugins = computed(() => {
  let filtered = props.remotePluginList

  // 搜索筛选
  if (marketSearch.value) {
    const searchLower = marketSearch.value.toLowerCase()
    filtered = filtered.filter(
      (plugin) =>
        plugin.name.toLowerCase().includes(searchLower) ||
        plugin.module.toLowerCase().includes(searchLower) ||
        plugin.introduction.toLowerCase().includes(searchLower) ||
        plugin.username.toLowerCase().includes(searchLower)
    )
  }

  // 模块筛选
  if (selectedModule.value) {
    filtered = filtered.filter((plugin) => plugin.module === selectedModule.value)
  }

  // 价格状态筛选
  if (selectedPriceStatus.value !== '') {
    if (selectedPriceStatus.value === 'free') {
      filtered = filtered.filter((plugin) => plugin.priceStatus === 0)
    } else if (selectedPriceStatus.value === 'paid') {
      filtered = filtered.filter((plugin) => plugin.priceStatus !== 0)
    }
  }

  // 安装状态筛选
  if (selectedInstallStatus.value) {
    if (selectedInstallStatus.value === 'installed') {
      filtered = filtered.filter((plugin) => plugin.isInstalled && !plugin.needUpdate)
    } else if (selectedInstallStatus.value === 'needUpdate') {
      filtered = filtered.filter((plugin) => plugin.isInstalled && plugin.needUpdate)
    } else if (selectedInstallStatus.value === 'notInstalled') {
      filtered = filtered.filter((plugin) => !plugin.isInstalled)
    }
  }

  // 类型筛选
  if (selectedType.value) {
    filtered = filtered.filter((plugin) => {
      const pluginType = plugin.type || 'scan' // 空字符串视为 'scan'
      return pluginType === selectedType.value
    })
  }

  return filtered
})

// 跳转到插件详情页
const handleCardClick = (plugin: RemotePluginData) => {
  if (plugin.hash) {
    window.open(`https://plugin.scope-sentry.top/plugin/${plugin.hash}`, '_blank')
  }
}

// Token 输入对话框
const tokenDialogVisible = ref(false)
const pluginToken = ref('')
const currentPlugin = ref<RemotePluginData | null>(null)

// 处理安装/更新按钮点击
const handleInstallClick = (plugin: RemotePluginData) => {
  // 如果是收费插件，先弹出 token 输入框
  if (plugin.priceStatus !== 0) {
    currentPlugin.value = plugin
    pluginToken.value = ''
    tokenDialogVisible.value = true
  } else {
    // 免费插件直接安装
    props.onInstall(plugin)
  }
}

// 确认 token 并安装
const confirmTokenAndInstall = async () => {
  if (!pluginToken.value.trim()) {
    ElMessageBox.alert(t('plugin.tokenRequired'), t('common.reminder'), {
      type: 'warning'
    })
    return
  }

  if (currentPlugin.value) {
    tokenDialogVisible.value = false
    // 调用父组件的安装函数，传递 token
    await props.onInstall(currentPlugin.value, pluginToken.value.trim())
    pluginToken.value = ''
    currentPlugin.value = null
  }
}
</script>

<template>
  <ElDrawer
    :model-value="visible"
    :title="t('plugin.market')"
    direction="rtl"
    size="90%"
    :close-on-click-modal="false"
    @update:model-value="(val) => !val && onClose()"
  >
    <div class="flex flex-col gap-4">
      <ElRow :gutter="16" style="margin-bottom: 16px">
        <ElCol :span="24">
          <ElRadioGroup v-model="selectedType" size="default">
            <ElRadioButton value="">{{ t('common.all') }}</ElRadioButton>
            <ElRadioButton value="scan">{{ t('plugin.scanPlugin') }}</ElRadioButton>
            <ElRadioButton value="server">{{ t('plugin.serverPlugin') }}</ElRadioButton>
          </ElRadioGroup>
        </ElCol>
      </ElRow>
      <ElRow :gutter="16">
        <ElCol :span="6">
          <ElInput
            v-model="marketSearch"
            :placeholder="t('common.inputText')"
            clearable
            style="height: 38px"
          >
            <template #prefix>
              <Icon icon="iconoir:search" />
            </template>
          </ElInput>
        </ElCol>
        <ElCol :span="4">
          <ElSelect
            v-model="selectedModule"
            :placeholder="t('plugin.module')"
            clearable
            style="width: 100%; height: 38px"
          >
            <ElOption
              v-for="option in moduleOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </ElSelect>
        </ElCol>
        <ElCol :span="4">
          <ElSelect
            v-model="selectedPriceStatus"
            :placeholder="t('plugin.priceStatus')"
            clearable
            style="width: 100%; height: 38px"
          >
            <ElOption :label="t('plugin.free')" value="free" />
            <ElOption :label="t('plugin.paid')" value="paid" />
          </ElSelect>
        </ElCol>
        <ElCol :span="4">
          <ElSelect
            v-model="selectedInstallStatus"
            :placeholder="t('plugin.installStatus')"
            clearable
            style="width: 100%; height: 38px"
          >
            <ElOption :label="t('plugin.notInstalled')" value="notInstalled" />
            <ElOption :label="t('plugin.installed')" value="installed" />
            <ElOption :label="t('plugin.needUpdate')" value="needUpdate" />
          </ElSelect>
        </ElCol>
        <ElCol :span="3">
          <BaseButton type="primary" @click="onRefresh" :loading="marketLoading">
            {{ t('plugin.refresh') }}
          </BaseButton>
        </ElCol>
        <ElCol :span="3">
          <BaseButton @click="resetFilters">
            {{ t('common.reset') }}
          </BaseButton>
        </ElCol>
      </ElRow>
      <div v-loading="marketLoading" class="plugin-market-container">
        <ElRow :gutter="20">
          <ElCol
            v-for="plugin in filteredRemotePlugins"
            :key="plugin.id"
            :xs="24"
            :sm="12"
            :md="8"
            :lg="6"
            :xl="6"
          >
            <ElCard
              class="plugin-market-card"
              :style="{
                width: '100%',
                position: 'relative',
                marginBottom: '20px',
                cursor: 'pointer'
              }"
              shadow="hover"
              :body-style="{ padding: '0' }"
              @click="handleCardClick(plugin)"
            >
              <div class="plugin-card-cover">
                <div
                  :style="{
                    height: '150px',
                    background: moduleBackgroundMap[plugin.module] || 'var(--cat-neutral-bg)',
                    borderBottom: '1px solid var(--border)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    position: 'relative'
                  }"
                >
                  <div
                    :style="{
                      fontSize: '18px',
                      fontWeight: 'bold',
                      color: 'var(--text-primary)',
                      letterSpacing: '0.5px',
                      textAlign: 'center'
                    }"
                  >
                    {{ plugin.name }}
                  </div>
                  <ElTag
                    :type="plugin.priceStatus === 0 ? 'success' : 'danger'"
                    :style="{ position: 'absolute', bottom: '8px', left: '8px' }"
                  >
                    {{ plugin.priceStatus === 0 ? t('plugin.free') : t('plugin.paid') }}
                  </ElTag>
                  <ElTag
                    :style="{ ...moduleBadgeStyle(plugin.module), position: 'absolute', bottom: '8px', right: '8px' }"
                    v-if="plugin.type !== 'server'"
                  >
                    {{
                      moduleOptions.find((option) => option.value === plugin.module)?.label ||
                      t('plugin.unknownModule')
                    }}
                  </ElTag>
                </div>
              </div>
              <div class="plugin-card-content">
                <div class="plugin-info-item">
                  <span class="label">{{ t('plugin.version') }}：</span>
                  <span class="value">{{ plugin.version || 'N/A' }}</span>
                </div>
                <div class="plugin-info-item">
                  <span class="label">{{ t('plugin.author') }}：</span>
                  <span class="value">{{ plugin.username }}</span>
                </div>
                <div class="plugin-info-item">
                  <span class="label">{{ t('plugin.createTime') }}：</span>
                  <span class="value">{{ plugin.createTime }}</span>
                </div>
                <div class="plugin-info-item introduction">
                  <span class="label">{{ t('plugin.introduction') }}：</span>
                  <span class="value">{{ plugin.introduction || t('plugin.noIntroduction') }}</span>
                </div>
                <div class="plugin-status">
                  <ElTag
                    v-if="plugin.isInstalled && plugin.needUpdate"
                    type="warning"
                    style="margin-bottom: 10px"
                  >
                    {{ t('plugin.needUpdate') }}
                  </ElTag>
                  <ElTag v-else-if="plugin.isInstalled" type="success" style="margin-bottom: 10px">
                    {{ t('plugin.installed') }}
                  </ElTag>
                  <ElTag v-else style="margin-bottom: 10px">{{ t('plugin.notInstalled') }}</ElTag>
                </div>
                <div class="plugin-actions" @click.stop>
                  <BaseButton
                    v-if="plugin.isInstalled && plugin.needUpdate"
                    type="warning"
                    style="width: 100%"
                    @click="handleInstallClick(plugin)"
                  >
                    {{ t('plugin.update') }}
                  </BaseButton>
                  <BaseButton
                    v-else-if="plugin.isInstalled"
                    type="info"
                    disabled
                    style="width: 100%"
                  >
                    {{ t('plugin.installed') }}
                  </BaseButton>
                  <BaseButton
                    v-else
                    type="primary"
                    style="width: 100%"
                    @click="handleInstallClick(plugin)"
                  >
                    {{ t('plugin.install') }}
                  </BaseButton>
                </div>
              </div>
            </ElCard>
          </ElCol>
        </ElRow>
        <div v-if="filteredRemotePlugins.length === 0" class="empty-state">
          <p>{{ t('plugin.noPluginData') }}</p>
        </div>
      </div>
    </div>
  </ElDrawer>
  <!-- Token 输入对话框 -->
  <ElDialog
    v-model="tokenDialogVisible"
    :title="t('plugin.token')"
    width="400px"
    :close-on-click-modal="false"
  >
    <div class="flex flex-col gap-4">
      <ElInput
        v-model="pluginToken"
        :placeholder="t('plugin.tokenPlaceholder')"
        type="password"
        show-password
        @keyup.enter="confirmTokenAndInstall"
      />
    </div>
    <template #footer>
      <BaseButton @click="tokenDialogVisible = false">{{ t('common.cancel') }}</BaseButton>
      <BaseButton type="primary" @click="confirmTokenAndInstall">
        {{ t('common.ok') }}
      </BaseButton>
    </template>
  </ElDialog>
</template>

<style scoped lang="less">
.plugin-market-container {
  min-height: 400px;
  padding: 10px 0;
}

.plugin-market-card {
  border-radius: 12px;
  overflow: hidden;

  :deep(.el-card__body) {
    border-radius: 12px;
    overflow: hidden;
    padding: 0;
  }
}

.plugin-card-cover {
  width: 100%;
  overflow: hidden;
  border-radius: 12px 12px 0 0;
}

.plugin-card-content {
  padding: 15px;
}

.plugin-info-item {
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 1.6;

  .label {
    font-weight: 600;
    color: var(--text-secondary);
    margin-right: 5px;
  }

  .value {
    color: var(--text-primary);
  }

  &.introduction {
    .value {
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
      color: var(--text-muted);
    }
  }
}

.plugin-status {
  margin: 15px 0;
}

.plugin-actions {
  margin-top: 15px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
  font-size: 16px;
}
</style>
