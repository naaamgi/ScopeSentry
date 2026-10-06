package mcp

// listAssetsToolDesc 资产查询说明
//
// MCP 도구 설명은 LLM 클라이언트가 프롬프트로 읽는 텍스트라서 요청마다 로케일을
// 바꿀 수 없다. 이 계층은 영어 하나로 고정한다.
const listAssetsToolDesc = `Query the ScopeSentry asset list.

[Endpoint and request body]
Matches the per-asset POST endpoints (for example /api/assets/asset). The body is
models.SearchRequest:
- pageIndex, pageSize: paging
- search: search expression string (the web UI's Csearch searchParams)
- filter: exact-match filter object (the web UI's project/task dropdowns, statistics
  sidebar clicks and table column filters, merged together)
- sort: sorting (ignored by most types)
- sid: SensitiveResult only, a sensitive rule name

The MCP asset_type maps to the backend Index (asset -> "asset"), which
helper.GetSearchQuery then turns into a MongoDB query.

[How the web UI splits the parameters (Csearch.vue and the asset pages)]
- search: the DSL the user typed, such as domain=example && port==443
- filter.project: the ElTreeSelect value, an array of project ObjectIDs (not names)
- filter.task: the ElSelect value, an array of task names (not task IDs). A dynamic
  tag of the form task=<name> also writes into filter
- filter can also come from statistics sidebar clicks (port/service/app/icon) and
  table column filters (statuscode/level/type/status and so on)
- Note: the asset pages' searchKeywordsData offers a project hint, but the backend
  SearchToMongoDB does not register a project keyword, so project only works
  through filter

asset_type values:
- asset: web and port assets
- RootDomain: root domains
- subdomain: subdomains
- app: mobile apps
- mp: mini-programs
- UrlScan: URL scan results
- SensitiveResult: sensitive information
- DirScanResult: directory scan
- crawler: crawler
- vulnerability: vulnerabilities
- PageMonitoring: page monitoring
- IPAsset: aggregated IP assets
- SubdomainTakerResult: subdomain takeover

[search expression] (not SQL; a custom DSL parsed by helper.SearchToMongoDB)
- field=value : fuzzy match (case-insensitive)
- field=="value" : exact match; quote a value containing spaces
- field!="value" : exclude
- expr1 && expr2 : AND
- expr1 || expr2 : OR
- (expr) : grouping
- Fields available on every type (injected by SearchToMongoDB): tag->tags,
  task->taskName (the task name), rootDomain
- project is not available in search. Use filter.project, whose values are project
  ObjectIDs from list_projects or list_projects_data

search keywords per asset_type -> MongoDB field:
- asset: domain->host, ip, port, service, app->technologies, title, statuscode, icon->faviconmmh3, banner->metadata, type, body, header->rawheaders
- RootDomain: domain, icp (brn is a read-only alias for the same field), company
- subdomain: domain->host, ip, type, value
- app: name, icp (brn is a read-only alias for the same field), company, category, description, url, apk
- mp: name, icp (brn is a read-only alias for the same field), company, category, description, url
- UrlScan: url->output, input, source, resultId, type->outputtype (no statuscode in search; the HTTP status code is filter.status only)
- SensitiveResult: url, sname->sid, body, info->match, md5
- DirScanResult: url, statuscode->status, redirect->msg, length
- vulnerability: url, vulname, matched, request, response, level
- crawler: url, method, body, resultId
- PageMonitoring: url, hash, diff, response
- IPAsset: ip, domain->ports.server.domain, port->ports.port, service->ports.server.service, webServer->ports.server.webServer, app->ports.server.technologies
- SubdomainTakerResult: domain->input, value, type->cname, response

search examples:
- domain=example && port==443
- port==443 && service=nginx
- title="admin login" || body=admin

[filter] (combinable with search; helper.GetSearchQuery appends it as $and conditions)
A JSON object whose keys are filter dimensions and whose values are string arrays.
Several values under one key are OR; different keys are AND. Only keys defined in
filterKeyCache take effect; the rest are ignored.

filter notes:
- project: filter only. An array of project ObjectIDs (the database stores IDs; some
  endpoints convert them to names for display)
- task: either search (task=="<task name>") or filter.task (an array of task names).
  The value is the name from list_tasks, not the task ID
- Several values under one key are OR, different keys are AND. A key that is not in
  filterKeyCache is ignored

global filter key -> MongoDB field:
- project->project, port->port, service->service, app->technologies
- icon->faviconmmh3, statuscode->statuscode, status->status
- level->level, type->type, color->color, tags->tags
- task->taskName, sname->sid

filter keys per asset_type (a key not listed does not apply to that collection):
- asset: project, port, service, app, icon, statuscode, type, task, tags
- RootDomain: project, tags
- subdomain: project, type, task, tags
- app / mp: project, tags
- UrlScan: status (the HTTP status code, not statuscode), tags
- SensitiveResult: status (1 unprocessed / 2 processing / 3 ignored / 4 suspected / 5 confirmed / 6 processed), color, sname, tags
- DirScanResult: status (the HTTP status code), tags
- crawler: project, task, tags
- vulnerability: project, level (critical/high/medium/low/info/unknown), status (1-6), task, tags
- PageMonitoring: tags
- IPAsset: project, port, service, app (nested under ports, resolved with an aggregation)
- SubdomainTakerResult: tags

filter examples (project must be an ObjectID, never a project name):
- asset: {"project":["<project ObjectID>"],"port":["443"]}
- subdomain: {"project":["<project ObjectID>"],"type":["A"]}
- vulnerability: {"project":["<project ObjectID>"],"level":["high"]}

search and filter combined:
- search: domain=example && port==443, filter: {"project":["<project ObjectID>"]}

Notes:
- Do not write project=... or project=="..." in search; it does not work and errors
  when combined with &&
- For DirScanResult the HTTP status code is only available in search, as
  statuscode==200
- UrlScan has no statuscode in search; filter by filter.status
- To filter SensitiveResult by rule name use sname=<rule name> in search, or
  filter.sname

[sort]
- UrlScan / DirScanResult: {"length":"ascending"} sorts ascending; any other value
  (including "descending" and "-1") sorts descending
- Other types: the server sorts by time or _id, and sort is usually ignored

[sid]
SensitiveResult only: a sensitive rule name, used to expand the matches for that
rule (what clicking a rule name does in the web UI)`

// createScanTemplateToolDesc 扫描模板创建说明
const createScanTemplateToolDesc = `Create a scan template.

A scan template is made of modules (pipeline stages), each carrying a number of
plugins. The module fields reference a plugin's hash, not its id and not its name.

Suggested flow:
1. list_plugin_modules for every module name
2. list_plugins (optionally filtered by module) for the hash and default parameter of
   each plugin in a module
3. Pass modules as module -> list of plugin hashes (plugins within a module run in
   array order)
4. This tool fills Parameters in from each plugin's default parameter; pass
   parameters to override them

Parameters:
- name: the template name, required
- modules: {module name: [plugin hash, ...]}, for example
  {"SubdomainScan":["d60ba73c..."],"PortScan":["..."]}
- parameters: optional, {module name: {plugin hash: parameter string}}, overriding the
  defaults
- vullist: optional, a list of nuclei POC template IDs (when VulnerabilityScan runs
  nuclei)
- template_json: optional, a complete ScanTemplate JSON. It takes precedence over
  modules and is there for advanced cases

On success the template id is returned, ready to pass as create_scan_task's template.

Common modules: TargetHandler, SubdomainScan, SubdomainSecurity, PortScanPreparation,
PortScan, PortFingerprint, AssetMapping, AssetHandle, URLScan, WebCrawler,
URLSecurity, DirScan, VulnerabilityScan, PassiveScan`

// createScanTaskToolDesc 扫描任务创建说明（与 Web 端 /api/task/add 及 common.Insert 逻辑一致）
const createScanTaskToolDesc = `Create a scan task. name and node are required;
template is a scan template ObjectID from list_scan_templates.

[targetSource] decides how the task's targets are resolved (see
internal/services/task/common/common.go):
- general: use the target field directly (domains, IPs or URLs, one per line or comma
  separated)
- project: read the targets from the linked projects. Pass project, an array of
  project ObjectIDs; target is not needed
- asset: select from the web asset store. Requires search; project, filter and
  targetNumber are optional
- RootDomain: select from the root domain store. Requires search; project, filter and
  targetNumber are optional
- subdomain: select from the subdomain store. Requires search; project, filter and
  targetNumber are optional
- UrlScan: select from the URL scan results. Requires search; project, filter and
  targetNumber are optional
- assetSource / RootDomainSource / subdomainSource / UrlScanSource: created from the
  matching asset page
  - targetTp=search: select targets with search + filter + project + targetNumber
  - targetTp=select: pass targetIds, a list of asset ObjectIDs

[search] A search expression with the same syntax as list_assets, such as
task=="<some task name>" or domain=^example.com.
To continue from subdomains: targetSource=subdomain, search=task=="<the subdomain
collection task name>"

[filter] An exact-match filter JSON, combinable with search. filter.project holds
project ObjectIDs.
[targetNumber] The cap on how many targets a search mode selects. 0 means no cap.
[targetIds] The asset ObjectIDs selected in select mode.

[Other parameters]
- allNode: automatically include every online node
- ignore / duplicates: targets to skip, and the deduplication strategy
- bindProject: the project the results belong to
- scheduledTasks together with cycleType/hour/minute/day/week: a scheduled task

For a full sweep of a root domain, two stages work better: first scan the root
domains with general and only SubdomainScan/SubdomainSecurity, then, once that
finishes, create a task for the remaining modules with subdomain and
search=task=="<the first task's name>".`

// countAssetsToolDesc 资产数量统计（对应 POST /api/assets/common/total）
const countAssetsToolDesc = `Count the assets matching a condition. This is the same
endpoint the web UI's pagination total comes from, /api/assets/common/total.

It takes the same search and filter as list_assets but does not page, and returns
only total. asset_type takes the same values as list_assets (asset, RootDomain,
subdomain, vulnerability and the rest).

Example, counting the subdomains in a project:
{"asset_type":"subdomain","filter":{"project":["<project ObjectID>"]}}

Example, counting the web assets a task produced:
{"asset_type":"asset","search":"task==\"<some task name>\""}`
