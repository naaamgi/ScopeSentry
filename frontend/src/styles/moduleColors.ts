// Module categories use the supplied tonal chip palette across all screens.
const moduleCategoryMap: Record<string, number> = {
  TargetHandler: 1,
  SubdomainScan: 2,
  SubdomainSecurity: 3,
  PortScanPreparation: 4,
  PortScan: 5,
  AssetMapping: 6,
  URLScan: 7,
  WebCrawler: 8,
  DirScan: 1,
  VulnerabilityScan: 2,
  AssetHandle: 3,
  PortFingerprint: 4,
  URLSecurity: 5,
  PassiveScan: 6
}

export const moduleColorMap: Record<string, string> = Object.fromEntries(
  Object.entries(moduleCategoryMap).map(([module, category]) => [module, `var(--cat-${category})`])
)

export const moduleBackgroundMap: Record<string, string> = Object.fromEntries(
  Object.entries(moduleCategoryMap).map(([module, category]) => [module, `var(--cat-${category}-bg)`])
)

export const moduleBadgeStyle = (module: string) => ({
  backgroundColor: moduleBackgroundMap[module] || 'var(--cat-neutral-bg)',
  borderColor: 'transparent',
  color: moduleColorMap[module] || 'var(--cat-neutral)'
})
