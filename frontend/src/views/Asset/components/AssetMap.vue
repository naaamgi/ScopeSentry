<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import { ElButton, ElEmpty, ElInput, ElMessage, ElRadioButton, ElRadioGroup } from 'element-plus'
import { getAssetApi } from '@/api/asset'
import type { AssetData } from '@/api/asset/types'
import { useI18n } from '@/hooks/web/useI18n'
import { useAppStore } from '@/store/modules/app'
import { getCssVar } from '@/utils'

const { t } = useI18n()
const appStore = useAppStore()
const canvas = ref<HTMLDivElement>()
const records = ref<AssetData[]>([])
const total = ref(0)
const loading = ref(false)
const query = ref('')
const selected = ref('')
const expandedIps = ref(new Set<string>())
const layout = ref<'top-down' | 'free'>(localStorage.getItem('asset-map-layout') === 'free' ? 'free' : 'top-down')
const visibleRecords = computed(() => records.value.filter((asset) =>
  [asset.domain, asset.host, asset.ip, asset.service, asset.title, String(asset.port)]
    .some((value) => String(value || '').toLowerCase().includes(query.value.trim().toLowerCase()))
))
let chart: echarts.ECharts | undefined
let observer: ResizeObserver | undefined

type GraphNode = {
  id: string
  name: string
  category: number
  symbolSize: number
  value: string
  x?: number
  y?: number
  groupIp?: string
}

// IPs form separate lanes. Names sharing an IP are spread across that lane;
// otherwise several domain labels land at the same coordinates.
const placeTopDown = (nodes: GraphNode[], links: { source: string; target: string }[]) => {
  const byId = new Map(nodes.map((node) => [node.id, node]))
  const children = new Map<string, string[]>()
  links.forEach(({ source, target }) => children.set(source, [...(children.get(source) || []), target]))
  const ips = nodes.filter((node) => node.category === 2).sort((a, b) => a.name.localeCompare(b.name))
  let cursor = 0
  ips.forEach((ip) => {
    const serviceRoots = (children.get(ip.id) || []).map((id) => byId.get(id))
      .filter((node): node is GraphNode => !!node && node.category === 3)
    const expanded = serviceRoots.some((node) => (children.get(node.id) || []).length > 0)
    const width = expanded ? 470 : 320
    ip.x = cursor + width / 2
    ip.y = 220
    serviceRoots.forEach((root) => {
      root.x = ip.x
      root.y = 330
      const services = (children.get(root.id) || []).map((id) => byId.get(id))
        .filter((node): node is GraphNode => !!node)
      services.forEach((service, index) => {
        service.x = ip.x! + (index % 2 === 0 ? -115 : 115)
        service.y = 425 + Math.floor(index / 2) * 56
      })
    })
    cursor += width + 100
  })
  const averageChildX = (node: GraphNode) => {
    const positions = (children.get(node.id) || []).map((id) => byId.get(id)?.x)
      .filter((position): position is number => position !== undefined)
    return positions.length ? positions.reduce((sum, position) => sum + position, 0) / positions.length : undefined
  }
  const placeDnsRow = (category: number, y: number) => {
    const groups = new Map<number, GraphNode[]>()
    nodes.filter((node) => node.category === category).sort((a, b) => a.name.localeCompare(b.name))
      .forEach((node) => {
        let anchor = averageChildX(node)
        if (anchor === undefined) { anchor = cursor + 160; cursor += 320 }
        groups.set(anchor, [...(groups.get(anchor) || []), node])
      })
    const row: GraphNode[] = []
    groups.forEach((group, anchor) => group.forEach((node, index) => {
      node.x = anchor + (index - (group.length - 1) / 2) * 210
      node.y = y
      row.push(node)
    }))
    row.sort((a, b) => a.x! - b.x!)
    row.forEach((node, index) => {
      if (index && node.x! < row[index - 1].x! + 210) node.x = row[index - 1].x! + 210
    })
  }
  placeDnsRow(1, 125)
  placeDnsRow(0, 55)
}

const render = async () => {
  await nextTick()
  if (!canvas.value) return
  chart ||= echarts.init(canvas.value)
  const nodes = new Map<string, GraphNode>()
  const edges = new Map<string, { source: string; target: string }>()
  const addNode = (id: string, name: string, category: number, size: number) => {
    if (!name) return
    const existing = nodes.get(id)
    if (existing) {
      // A DNS name can occur as both a root domain and a scanned host.
      if (category === 0) { existing.category = 0; existing.symbolSize = 35 }
      return
    }
    nodes.set(id, { id, name, category, symbolSize: size, value: name })
  }
  const addEdge = (source: string, target: string) => {
    if (source !== target && nodes.has(source) && nodes.has(target)) edges.set(`${source}|${target}`, { source, target })
  }
  visibleRecords.value.forEach((asset) => {
    const domain = String(asset.domain || '').trim().replace(/\.$/, '')
    const host = String(asset.host || domain).trim().replace(/\.$/, '')
    const ip = String(asset.ip || '').trim()
    const service = asset.port ? `${asset.port}${asset.service ? ` / ${asset.service}` : ''}` : String(asset.service || '')
    const domainId = `dns:${domain.toLowerCase()}`
    const hostId = `dns:${host.toLowerCase()}`
    if (domain) addNode(domainId, domain, 0, 35)
    if (host) addNode(hostId, host, 1, 27)
    if (ip) addNode(`ip:${ip}`, ip, 2, 25)
    if (service && ip) addNode(`service:${ip}:${service}`, service, 3, 18)
    if (domain && host && domainId !== hostId) addEdge(domainId, hostId)
    if (host && ip) addEdge(hostId, `ip:${ip}`)
    if (ip && service) addEdge(`ip:${ip}`, `service:${ip}:${service}`)
  })
  let graphNodes = [...nodes.values()]
  let graphEdges = [...edges.values()]
  const categoryNames = [t('map.domain'), t('map.host'), 'IP', t('map.service')]
  const topDown = layout.value === 'top-down'
  const palette = ['--chart-1', '--chart-2', '--chart-5', '--chart-3'].map((name) => getCssVar(name).trim())
  const textColor = getCssVar('--text-primary').trim()
  if (topDown) {
    const servicesByIp = new Map<string, GraphNode[]>()
    graphEdges.filter((edge) => edge.source.startsWith('ip:') && edge.target.startsWith('service:'))
      .forEach((edge) => {
        const service = nodes.get(edge.target)
        if (service) servicesByIp.set(edge.source, [...(servicesByIp.get(edge.source) || []), service])
      })
    graphNodes = graphNodes.filter((node) => node.category !== 3)
    graphEdges = graphEdges.filter((edge) => !edge.target.startsWith('service:'))
    servicesByIp.forEach((services, ipId) => {
      services.sort((a, b) => Number.parseInt(a.name) - Number.parseInt(b.name) || a.name.localeCompare(b.name))
      if (services.length === 1) {
        graphNodes.push(services[0])
        graphEdges.push({ source: ipId, target: services[0].id })
      } else {
        const expanded = expandedIps.value.has(ipId)
        const group: GraphNode = {
          id: `service-group:${ipId}`, groupIp: ipId, category: 3, symbolSize: 14,
          name: t('map.serviceGroup', { count: services.length }) + (expanded ? ' ▾' : ' ▸'),
          value: String(services.length)
        }
        graphNodes.push(group)
        graphEdges.push({ source: ipId, target: group.id })
        if (expanded) services.forEach((service) => {
          graphNodes.push(service)
          graphEdges.push({ source: group.id, target: service.id })
        })
      }
    })
    placeTopDown(graphNodes, graphEdges)
    graphNodes.forEach((node) => { node.symbolSize = node.category === 3 ? 12 : 15 })
  }
  chart.setOption({
    animationDuration: 500,
    tooltip: {
      renderMode: 'richText',
      backgroundColor: getCssVar('--bg-elevated').trim(),
      borderColor: getCssVar('--border-strong').trim(),
      textStyle: { color: textColor },
      formatter: (params: any) => params.dataType === 'node'
        ? `${categoryNames[params.data.category]} · ${params.data.name}${params.data.groupIp ? `\n${t('map.groupHint')}` : ''}` : ''
    },
    legend: [{ data: categoryNames, top: 8, textStyle: { color: textColor, fontFamily: 'Inter, Noto Sans KR, sans-serif', fontSize: 12 } }],
    series: [{
      type: 'graph', layout: topDown ? 'none' : 'force', roam: true, draggable: !topDown,
      data: graphNodes, links: graphEdges,
      categories: [
        { name: t('map.domain'), itemStyle: { color: palette[0] } },
        { name: t('map.host'), itemStyle: { color: palette[1] } },
        { name: 'IP', itemStyle: { color: palette[2] } },
        { name: t('map.service'), itemStyle: { color: palette[3] } }
      ],
      label: { show: true, position: 'right', color: textColor, fontFamily: 'Inter, Noto Sans KR, sans-serif', fontSize: 12, fontWeight: 500, distance: 8 },
      labelLayout: { hideOverlap: false },
      lineStyle: { color: getCssVar('--text-secondary').trim(), opacity: topDown ? .42 : .7, width: 1, curveness: topDown ? 0 : .06 },
      emphasis: { focus: 'none', lineStyle: { width: 2.4 } },
      force: layout.value === 'free' ? { repulsion: 140, edgeLength: [75, 150], gravity: .07 } : undefined
    }]
  }, true)
  chart.off('click')
  chart.on('click', (event) => {
    if (event.dataType !== 'node') { selected.value = ''; return }
    const node = event.data as GraphNode
    selected.value = node.name
    if (node.groupIp) {
      if (expandedIps.value.has(node.groupIp)) expandedIps.value.delete(node.groupIp)
      else expandedIps.value.add(node.groupIp)
      render()
    }
  })
  chart.resize()
}

const load = async () => {
  loading.value = true
  try {
    const result = await getAssetApi('', 1, 500, {})
    records.value = result.data.list || []
    total.value = Math.max(Number(result.data.total) || 0, records.value.length)
    await render()
  } catch {
    ElMessage.error(t('map.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
  if (canvas.value) {
    observer = new ResizeObserver(() => chart?.resize())
    observer.observe(canvas.value)
  }
})
onBeforeUnmount(() => { observer?.disconnect(); chart?.dispose() })
watch(() => appStore.getIsDark, () => { render() })
watch(layout, (value) => { localStorage.setItem('asset-map-layout', value); render() })
</script>

<template>
  <section class="asset-map">
    <header class="map-toolbar">
      <div><h2>{{ t('map.title') }}</h2><p>{{ t('map.description') }}</p></div>
      <div class="map-controls">
        <ElRadioGroup v-model="layout" :aria-label="t('map.layout')">
          <ElRadioButton label="top-down">{{ t('map.topDown') }}</ElRadioButton>
          <ElRadioButton label="free">{{ t('map.freeLayout') }}</ElRadioButton>
        </ElRadioGroup>
        <ElInput v-model="query" clearable :placeholder="t('map.search')" @input="render" />
        <ElButton :loading="loading" @click="load">{{ t('map.refresh') }}</ElButton>
      </div>
    </header>
    <div class="map-meta">{{ t('map.shown', { count: visibleRecords.length, total }) }} <span v-if="selected">· {{ selected }}</span></div>
    <ElEmpty v-if="!loading && records.length === 0" :description="t('map.empty')" />
    <div ref="canvas" class="map-canvas" :aria-label="t('map.title')" />
  </section>
</template>

<style scoped>
.asset-map { min-height: 680px; padding: 22px; background: var(--el-bg-color); border: 1px solid var(--el-border-color); border-radius: 12px; }
.map-toolbar { display: flex; align-items: start; justify-content: space-between; flex-wrap: wrap; gap: 16px; }
.map-toolbar h2 { margin: 0 0 5px; font-size: 20px; }
.map-toolbar p { margin: 0; color: var(--el-text-color-secondary); }
.map-controls { display: flex; flex-wrap: wrap; gap: 8px; }
.map-controls .el-input { width: 250px; }
.map-meta { min-height: 34px; padding-top: 12px; font-size: 12px; color: var(--el-text-color-secondary); }
.map-canvas { width: 100%; height: min(60vh, 620px); min-height: 420px; background: var(--bg-card); }
@media (max-width: 640px) { .map-controls, .map-controls .el-input { width: 100%; } }
</style>
