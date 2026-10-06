package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/app"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/asset"
	assetCommon "github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/crawler"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/dirscan"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/ip"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/mp"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/page_monitoring"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/root_domain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/sensitive"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/subdomain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/url"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/vulnerability"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/services/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	taskCommon "github.com/Autumn-27/ScopeSentry/internal/services/task/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/task"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/template"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type deps struct {
	projectService    project.Service
	taskService       task.Service
	taskCommonService taskCommon.Service
	templateService   template.Service
	assetService      asset.Service
	rootDomainService root_domain.Service
	subdomainService  subdomain.Service
	appService        app.Service
	mpService         mp.Service
	urlService        url.Service
	sensitiveService  sensitive.Service
	dirscanService    dirscan.Service
	crawlerService    crawler.Service
	vulnService       vulnerability.Service
	ipService         ip.Service
	pageMonService    page_monitoring.Service
	commonService     assetCommon.Service
	nodeService       node.Service
	pluginService     plugin.Service
}

var d = &deps{
	projectService:    project.NewService(),
	taskService:       task.NewService(),
	taskCommonService: taskCommon.NewService(),
	templateService:   template.NewService(),
	assetService:      asset.NewService(),
	rootDomainService: root_domain.NewService(),
	subdomainService:  subdomain.NewService(),
	appService:        app.NewService(),
	mpService:         mp.NewService(),
	urlService:        url.NewService(),
	sensitiveService:  sensitive.NewService(),
	dirscanService:    dirscan.NewService(),
	crawlerService:    crawler.NewService(),
	vulnService:       vulnerability.NewService(),
	ipService:         ip.NewService(),
	pageMonService:    page_monitoring.NewService(),
	commonService:     assetCommon.NewService(),
	nodeService:       node.NewService(),
	pluginService:     plugin.NewService(),
}

func registerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects",
		Description: "List projects grouped by tag",
	}, listProjects)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects_data",
		Description: "List projects, paged, with search",
	}, listProjectsData)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project",
		Description: "Get a project by its ID",
	}, getProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_project",
		Description: "Create a project. name and target are required; tag, template, node and others are optional",
	}, createProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tasks",
		Description: "List scan tasks, paged",
	}, listTasks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_task",
		Description: "Get a task by its ID",
	}, getTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_scan_templates",
		Description: "List scan templates, paged",
	}, listScanTemplates)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_scan_template",
		Description: "Get a scan template by its ID",
	}, getScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugin_modules",
		Description: "List every scan template module name, the stages of the pipeline. A scan template is made of these modules, each carrying a number of plugins.",
	}, listPluginModules)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugins",
		Description: "List the available scan plugins, optionally filtered by module. Returns each plugin hash, name, module and default parameter. A scan template module field references the plugin hash.",
	}, listPlugins)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_template",
		Description: createScanTemplateToolDesc,
	}, createScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_task",
		Description: createScanTaskToolDesc,
	}, createScanTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_assets",
		Description: listAssetsToolDesc,
	}, listAssets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "count_assets",
		Description: countAssetsToolDesc,
	}, countAssets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_asset_detail",
		Description: "Get an asset; asset_type accepts asset or vulnerability",
	}, getAssetDetail)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_asset_tag",
		Description: "Add a tag to an asset",
	}, addAssetTag)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_nodes",
		Description: "List the scan nodes; online_only=true returns only the online ones",
	}, listNodes)
}

type listProjectsDataInput struct {
	Search    string `json:"search,omitempty" jsonschema:"a keyword for fuzzy search on the project name"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"the page number, starting at 1; 1 by default"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"items per page; 20 by default"`
}

func listProjects(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsByTag(c)
	if err != nil {
		return errorResult("failed to retrieve the project list", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listProjectsData(ctx context.Context, _ *mcp.CallToolRequest, input listProjectsDataInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsData(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("failed to retrieve the project data", err)
	}
	return jsonToolResult(result)
}

type getProjectInput struct {
	ID string `json:"id" jsonschema:"the project MongoDB ObjectID"`
}

func getProject(ctx context.Context, _ *mcp.CallToolRequest, input getProjectInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id must not be empty", nil)
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectContent(c, input.ID)
	if err != nil {
		return errorResult("failed to retrieve the project details", err)
	}
	if result == nil {
		return errorResult("the project does not exist", nil)
	}
	return jsonToolResult(result)
}

type createProjectInput struct {
	Name           string   `json:"name" jsonschema:"the project name; required"`
	Tag            string   `json:"tag,omitempty" jsonschema:"the project tag, used for grouping"`
	Target         string   `json:"target" jsonschema:"the scan targets; required, one per line or comma separated, accepting domains, IPs, URLs and so on"`
	Template       string   `json:"template,omitempty" jsonschema:"the scan template ID to link"`
	Node           []string `json:"node,omitempty" jsonschema:"the names of the scan nodes to use"`
	AllNode        bool     `json:"allNode,omitempty" jsonschema:"whether to use every node"`
	Ignore         string   `json:"ignore,omitempty" jsonschema:"the targets to skip, in the same format as target"`
	Duplicates     string   `json:"duplicates,omitempty" jsonschema:"the deduplication strategy"`
	ScheduledTasks bool     `json:"scheduledTasks,omitempty" jsonschema:"whether to enable scheduled scanning"`
	Hour           int      `json:"hour,omitempty" jsonschema:"the scheduled scan interval in hours; only used when scheduledTasks is on"`
}

func createProject(ctx context.Context, _ *mcp.CallToolRequest, input createProjectInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || input.Target == "" {
		return errorResult("name and target must not be empty", nil)
	}
	p := &models.Project{
		Name:           input.Name,
		Tag:            input.Tag,
		Target:         input.Target,
		Template:       input.Template,
		Node:           input.Node,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Tp:             "project",
	}
	c := ginContext(ctx)
	if err := d.projectService.AddProject(c, p); err != nil {
		return errorResult("failed to create the project", err)
	}
	return jsonToolResult(ginH{"success": true, "message": "the project was created"})
}

type listTasksInput struct {
	Search    string `json:"search,omitempty" jsonschema:"fuzzy search on the task name"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"the page number, starting at 1; 1 by default"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"items per page; 20 by default"`
}

func listTasks(ctx context.Context, _ *mcp.CallToolRequest, input listTasksInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	tasks, total, err := d.taskService.List(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("failed to retrieve the task list", err)
	}
	return jsonToolResult(ginH{"list": tasks, "total": total})
}

type getTaskInput struct {
	ID string `json:"id" jsonschema:"the task MongoDB ObjectID"`
}

func getTask(ctx context.Context, _ *mcp.CallToolRequest, input getTaskInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id must not be empty", nil)
	}
	c := ginContext(ctx)
	result, err := d.taskService.GetTaskDetail(c, input.ID)
	if err != nil {
		return errorResult("failed to retrieve the task details", err)
	}
	if result == nil {
		return errorResult("the task does not exist", nil)
	}
	return jsonToolResult(result)
}

type listScanTemplatesInput struct {
	Query     string `json:"query,omitempty" jsonschema:"fuzzy search on the template name"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"the page number, starting at 1; 1 by default"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"items per page; 20 by default"`
}

func listScanTemplates(ctx context.Context, _ *mcp.CallToolRequest, input listScanTemplatesInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	result, err := d.templateService.List(ctx, input.PageIndex, input.PageSize, input.Query)
	if err != nil {
		return errorResult("failed to retrieve the template list", err)
	}
	return jsonToolResult(result)
}

type getScanTemplateInput struct {
	ID string `json:"id" jsonschema:"the scan template MongoDB ObjectID"`
}

func getScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input getScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id must not be empty", nil)
	}
	c := ginContext(ctx)
	result, err := d.templateService.Detail(c, input.ID)
	if err != nil {
		return errorResult("failed to retrieve the template details", err)
	}
	return jsonToolResult(result)
}

type createScanTemplateInput struct {
	Name         string                       `json:"name" jsonschema:"the template name; required"`
	Modules      map[string][]string          `json:"modules,omitempty" jsonschema:"a map from module name to plugin hashes. The key is a module name (from list_plugin_modules) and the value is the hashes of the plugins to enable in it (from list_plugins), which run in array order. For example {\"SubdomainScan\":[\"d60ba73c...\"]}"`
	Parameters   map[string]map[string]string `json:"parameters,omitempty" jsonschema:"optional; overrides the plugin run parameters, shaped as module name -> plugin hash -> parameter string. Left out, each plugin default is used"`
	VulList      []string                     `json:"vullist,omitempty" jsonschema:"optional; a list of nuclei POC template IDs, used only when VulnerabilityScan runs nuclei"`
	TemplateJSON string                       `json:"template_json,omitempty" jsonschema:"optional; a complete ScanTemplate JSON. It takes precedence over modules and is there for advanced cases"`
}

func createScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input createScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" && input.TemplateJSON == "" {
		return errorResult("provide at least one of name or template_json", nil)
	}

	var tmpl *models.ScanTemplate

	switch {
	case input.TemplateJSON != "":
		tmpl = &models.ScanTemplate{}
		if err := json.Unmarshal([]byte(input.TemplateJSON), tmpl); err != nil {
			return errorResult("template_json is malformed", err)
		}
		if tmpl.Name == "" {
			tmpl.Name = input.Name
		}
	case len(input.Modules) > 0:
		built, err := buildTemplateFromModules(ctx, input)
		if err != nil {
			return errorResult(err.Error(), nil)
		}
		tmpl = built
	default:
		tmpl = &models.ScanTemplate{Name: input.Name}
	}

	if tmpl.Name == "" {
		return errorResult("the template name must not be empty", nil)
	}

	id, err := d.templateService.Save(ctx, "", tmpl)
	if err != nil {
		return errorResult("failed to create the template", err)
	}
	return jsonToolResult(ginH{"success": true, "id": id, "message": "the template was created and can be passed as the template parameter of create_scan_task"})
}

// buildTemplateFromModules 根据「模块->插件hash列表」组装扫描模板，并自动回填插件默认参数
func buildTemplateFromModules(ctx context.Context, input createScanTemplateInput) (*models.ScanTemplate, error) {
	validModule := make(map[string]bool, len(constants.PLUGINSMODULES))
	for _, m := range constants.PLUGINSMODULES {
		validModule[m] = true
	}

	templateMap := map[string]any{
		"name": input.Name,
	}
	paramMap := map[string]map[string]string{}

	for module, hashes := range input.Modules {
		if !validModule[module] {
			return nil, fmt.Errorf("invalid module name: %s (call list_plugin_modules for the valid ones)", module)
		}
		if len(hashes) == 0 {
			continue
		}
		templateMap[module] = hashes

		modParams := map[string]string{}
		for _, hash := range hashes {
			// 优先使用调用方提供的参数覆盖
			if input.Parameters != nil {
				if mp, ok := input.Parameters[module]; ok {
					if v, ok := mp[hash]; ok {
						modParams[hash] = v
						continue
					}
				}
			}
			// 否则回填插件默认参数
			plg, err := d.pluginService.GetPluginByHash(ctx, hash)
			if err != nil || plg == nil {
				return nil, fmt.Errorf("no plugin with hash %s (module %s)", hash, module)
			}
			if plg.Module != module {
				return nil, fmt.Errorf("plugin %s (hash=%s) belongs to module %s and cannot be placed in module %s", plg.Name, hash, plg.Module, module)
			}
			modParams[hash] = plg.Parameter
		}
		if len(modParams) > 0 {
			paramMap[module] = modParams
		}
	}

	if len(paramMap) > 0 {
		templateMap["Parameters"] = paramMap
	}
	if len(input.VulList) > 0 {
		templateMap["vullist"] = input.VulList
	}

	data, err := json.Marshal(templateMap)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble the template: %w", err)
	}
	tmpl := &models.ScanTemplate{}
	if err := json.Unmarshal(data, tmpl); err != nil {
		return nil, fmt.Errorf("failed to assemble the template: %w", err)
	}
	return tmpl, nil
}

type createScanTaskInput struct {
	Name           string              `json:"name" jsonschema:"the task name; required and must be unique"`
	Target         string              `json:"target,omitempty" jsonschema:"the scan targets; required when targetSource is general, one per line or comma separated"`
	Node           []string            `json:"node" jsonschema:"the names of the nodes that run the scan; required"`
	Template       string              `json:"template,omitempty" jsonschema:"the scan template ObjectID; a scan cannot run without it"`
	AllNode        bool                `json:"allNode,omitempty" jsonschema:"whether to automatically include every online node"`
	Ignore         string              `json:"ignore,omitempty" jsonschema:"the targets to skip, in the same format as target"`
	Duplicates     string              `json:"duplicates,omitempty" jsonschema:"the deduplication strategy, for example None"`
	Project        []string            `json:"project,omitempty" jsonschema:"the project ObjectIDs to link; required when targetSource is project"`
	TargetSource   string              `json:"targetSource,omitempty" jsonschema:"where the targets come from: general, project, asset, RootDomain, subdomain, UrlScan, or any of those with the Source suffix. general by default; see this tool description"`
	TargetTp       string              `json:"targetTp,omitempty" jsonschema:"how a Source origin selects targets: search or select"`
	Search         string              `json:"search,omitempty" jsonschema:"the search expression that selects targets from the asset store; same syntax as list_assets, see this tool description"`
	Filter         map[string][]string `json:"filter,omitempty" jsonschema:"the exact-match filter, combinable with search; filter.project holds project ObjectIDs"`
	TargetNumber   int                 `json:"targetNumber,omitempty" jsonschema:"the cap on targets in search mode; 0 means no cap"`
	TargetIds      []string            `json:"targetIds,omitempty" jsonschema:"the asset ObjectIDs selected in select mode"`
	BindProject    string              `json:"bindProject,omitempty" jsonschema:"the project ObjectID to bind, which the scan results belong to"`
	ScheduledTasks bool                `json:"scheduledTasks,omitempty" jsonschema:"whether to create this as a scheduled task"`
	Hour           int                 `json:"hour,omitempty" jsonschema:"the scheduled task interval, hours"`
	Minute         int                 `json:"minute,omitempty" jsonschema:"the scheduled task interval, minutes"`
	Day            int                 `json:"day,omitempty" jsonschema:"the scheduled task interval, days"`
	Week           int                 `json:"week,omitempty" jsonschema:"the scheduled task interval, weeks (the weekly cycle)"`
	CycleType      string              `json:"cycleType,omitempty" jsonschema:"the cycle type, for example nhours, daily or weekly"`
}

func createScanTask(ctx context.Context, _ *mcp.CallToolRequest, input createScanTaskInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || len(input.Node) == 0 {
		return errorResult("name and node must not be empty", nil)
	}
	c := ginContext(ctx)
	exists, err := d.taskService.CheckTaskNameExists(c, input.Name)
	if err != nil {
		return errorResult("failed to check the task name", err)
	}
	if exists {
		return errorResult("a task with that name already exists", nil)
	}

	targetSource := input.TargetSource
	if targetSource == "" {
		targetSource = "general"
	}

	taskModel := &models.Task{
		Name:           input.Name,
		Target:         input.Target,
		Node:           input.Node,
		Template:       input.Template,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		Project:        input.Project,
		TargetSource:   targetSource,
		TargetTp:       input.TargetTp,
		Search:         input.Search,
		TargetNumber:   input.TargetNumber,
		TargetIds:      input.TargetIds,
		BindProject:    input.BindProject,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Minute:         input.Minute,
		Day:            input.Day,
		Week:           input.Week,
		CycleType:      input.CycleType,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		taskModel.Filter = filter
	}
	taskID, err := d.taskCommonService.Insert(ctx, taskModel)
	if err != nil {
		return errorResult("failed to create the scan task", err)
	}
	return jsonToolResult(ginH{"success": true, "id": taskID})
}

type listAssetsInput struct {
	AssetType        string              `json:"asset_type" jsonschema:"the asset type; required. One of asset, RootDomain, subdomain, app, mp, UrlScan, SensitiveResult, DirScanResult, crawler, vulnerability, PageMonitoring, IPAsset or SubdomainTakerResult"`
	PageIndex        int                 `json:"pageIndex,omitempty" jsonschema:"the page number, starting at 1; 1 by default"`
	PageSize         int                 `json:"pageSize,omitempty" jsonschema:"items per page; 20 by default"`
	SearchExpression string              `json:"search,omitempty" jsonschema:"the search expression; not SQL, see this tool description for the syntax"`
	Filter           map[string][]string `json:"filter,omitempty" jsonschema:"the exact-match filter JSON, combinable with search. filter.project holds project ObjectIDs (from list_projects, not project names) and filter.task holds task names. See this tool description"`
	Sort             map[string]string   `json:"sort,omitempty" jsonschema:"sorting. Only UrlScan and DirScanResult support length: ascending sorts ascending, any other value descending"`
	Sid              string              `json:"sid,omitempty" jsonschema:"a sensitive information rule name; only used when asset_type is SensitiveResult"`
}

type countAssetsInput struct {
	AssetType        string              `json:"asset_type" jsonschema:"the asset type; required, taking the same values as list_assets"`
	SearchExpression string              `json:"search,omitempty" jsonschema:"the search expression; same syntax as list_assets"`
	Filter           map[string][]string `json:"filter,omitempty" jsonschema:"the exact-match filter JSON; same syntax as list_assets, and filter.project holds project ObjectIDs"`
	Sid              string              `json:"sid,omitempty" jsonschema:"a sensitive information rule name; only used when asset_type is SensitiveResult"`
}

func countAssets(ctx context.Context, _ *mcp.CallToolRequest, input countAssetsInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}

	query := models.SearchRequest{
		Index:            index,
		SearchExpression: input.SearchExpression,
		Sid:              input.Sid,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		query.Filter = filter
	}

	total, err := d.commonService.TotalData(ctx, &query)
	if err != nil {
		return errorResult("failed to count the assets", err)
	}
	return jsonToolResult(ginH{"total": total})
}

func listAssets(ctx context.Context, _ *mcp.CallToolRequest, input listAssetsInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	query := models.SearchRequest{
		PageIndex:        input.PageIndex,
		PageSize:         input.PageSize,
		Index:            index,
		SearchExpression: input.SearchExpression,
		Sort:             input.Sort,
		Sid:              input.Sid,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		query.Filter = filter
	}

	c := ginContext(ctx)
	data, err := queryAssets(ctx, c, index, query)
	if err != nil {
		return errorResult("failed to query the assets", err)
	}
	return jsonToolResult(data)
}

type getAssetDetailInput struct {
	AssetType string `json:"asset_type" jsonschema:"the asset type; a detail lookup accepts asset or vulnerability"`
	ID        string `json:"id" jsonschema:"the asset MongoDB ObjectID; for the vulnerability type pass the hash instead"`
}

func getAssetDetail(ctx context.Context, _ *mcp.CallToolRequest, input getAssetDetailInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id must not be empty", nil)
	}
	c := ginContext(ctx)
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}

	switch index {
	case "asset":
		result, err := d.assetService.GetAssetByID(c, input.ID)
		if err != nil {
			return errorResult("failed to retrieve the asset details", err)
		}
		if result == nil {
			return errorResult("the asset does not exist", nil)
		}
		return jsonToolResult(result)
	case "vulnerability":
		result, err := d.vulnService.GetVulnerabilityDetailByHash(c, input.ID)
		if err != nil {
			return errorResult("failed to retrieve the vulnerability details", err)
		}
		return jsonToolResult(result)
	default:
		return errorResult("asset_type must be asset or vulnerability for a detail lookup", nil)
	}
}

type addAssetTagInput struct {
	AssetType string `json:"asset_type" jsonschema:"the asset type, as in list_assets"`
	ID        string `json:"id" jsonschema:"the asset MongoDB ObjectID"`
	Tag       string `json:"tag" jsonschema:"the tag name to add"`
}

func addAssetTag(ctx context.Context, _ *mcp.CallToolRequest, input addAssetTagInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.ID == "" || input.Tag == "" {
		return errorResult("id and tag must not be empty", nil)
	}
	c := ginContext(ctx)
	req := &models.TagRequest{Type: index, ID: input.ID, Tag: input.Tag}
	if err := d.commonService.AddTag(c, req); err != nil {
		return errorResult("failed to add the tag", err)
	}
	return jsonToolResult(ginH{"success": true})
}

type listNodesInput struct {
	OnlineOnly bool `json:"online_only,omitempty" jsonschema:"true returns only online nodes; the default false returns all of them"`
}

func listNodes(ctx context.Context, _ *mcp.CallToolRequest, input listNodesInput) (*mcp.CallToolResult, any, error) {
	result, err := d.nodeService.GetNodeData(ctx, input.OnlineOnly)
	if err != nil {
		return errorResult("failed to retrieve the node list", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listPluginModules(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	return jsonToolResult(ginH{"modules": constants.PLUGINSMODULES})
}

type listPluginsInput struct {
	Module string `json:"module,omitempty" jsonschema:"filter by module name; left empty, every scan plugin is returned. Module names come from list_plugin_modules"`
	Search string `json:"search,omitempty" jsonschema:"fuzzy search on the plugin name; only applies when module is left out"`
}

type pluginBrief struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Module    string `json:"module"`
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	Status    bool   `json:"status"`
}

func toPluginBriefs(plugins []models.Plugin) []pluginBrief {
	briefs := make([]pluginBrief, 0, len(plugins))
	for _, p := range plugins {
		// 跳过服务端插件，它们不参与扫描流水线
		if p.Type == "server" {
			continue
		}
		briefs = append(briefs, pluginBrief{
			Hash:      p.Hash,
			Name:      p.Name,
			Module:    p.Module,
			Parameter: p.Parameter,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	return briefs
}

func listPlugins(ctx context.Context, _ *mcp.CallToolRequest, input listPluginsInput) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)

	if input.Module != "" {
		plugins, err := d.pluginService.ListByModule(c, input.Module)
		if err != nil {
			return errorResult("failed to retrieve the plugin list", err)
		}
		return jsonToolResult(ginH{"list": toPluginBriefs(plugins)})
	}

	resp, err := d.pluginService.List(c, &models.PluginListRequest{
		PageIndex: 1,
		PageSize:  500,
		Search:    input.Search,
	})
	if err != nil {
		return errorResult("failed to retrieve the plugin list", err)
	}
	return jsonToolResult(ginH{"list": toPluginBriefs(resp.List), "total": resp.Total})
}

func queryAssets(ctx context.Context, c *gin.Context, index string, query models.SearchRequest) (any, error) {
	switch index {
	case "asset":
		list, err := d.assetService.GetAssets(c, query)
		return ginH{"list": list}, err
	case "RootDomain":
		return d.rootDomainService.GetRootDomainData(c, query)
	case "subdomain":
		list, err := d.subdomainService.GetSubdomains(ctx, query)
		return ginH{"list": list}, err
	case "app":
		return d.appService.GetAppData(c, query)
	case "mp":
		return d.mpService.GetMPData(c, query)
	case "UrlScan":
		list, err := d.urlService.GetURLs(ctx, query)
		return ginH{"list": list}, err
	case "SensitiveResult":
		list, err := d.sensitiveService.GetSensitiveInfo(ctx, query)
		return ginH{"list": list}, err
	case "DirScanResult":
		list, err := d.dirscanService.List(ctx, query)
		return ginH{"list": list}, err
	case "crawler":
		list, err := d.crawlerService.GetCrawlers(ctx, query)
		return ginH{"list": list}, err
	case "vulnerability":
		return d.vulnService.GetVulnerabilities(ctx, query)
	case "PageMonitoring":
		return d.pageMonService.GetResult(c, query)
	case "IPAsset":
		return d.ipService.GetIPAssets(c, query)
	case "SubdomainTakerResult":
		return querySubdomainTaker(ctx, query)
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", index)
	}
}

func querySubdomainTaker(ctx context.Context, query models.SearchRequest) (any, error) {
	takerService := subdomain.NewTakerService()
	c := ginContext(ctx)
	return takerService.GetSubdomainTakerData(c, query)
}

func normalizeAssetIndex(assetType string) (string, error) {
	assetType = strings.TrimSpace(assetType)
	if assetType == "" {
		return "", fmt.Errorf("asset_type must not be empty")
	}
	aliases := map[string]string{
		"asset":                "asset",
		"web":                  "asset",
		"rootdomain":           "RootDomain",
		"root_domain":          "RootDomain",
		"root-domain":          "RootDomain",
		"subdomain":            "subdomain",
		"app":                  "app",
		"mp":                   "mp",
		"miniprogram":          "mp",
		"mini_program":         "mp",
		"url":                  "UrlScan",
		"urlscan":              "UrlScan",
		"sensitive":            "SensitiveResult",
		"sensitiveresult":      "SensitiveResult",
		"sensitive_result":     "SensitiveResult",
		"dirscan":              "DirScanResult",
		"dir_scan":             "DirScanResult",
		"dirscanresult":        "DirScanResult",
		"directory":            "DirScanResult",
		"crawler":              "crawler",
		"vulnerability":        "vulnerability",
		"vuln":                 "vulnerability",
		"pagemonitoring":       "PageMonitoring",
		"page_monitoring":      "PageMonitoring",
		"ip":                   "IPAsset",
		"ipasset":              "IPAsset",
		"subdomaintaker":       "SubdomainTakerResult",
		"subdomain_taker":      "SubdomainTakerResult",
		"subdomaintakerresult": "SubdomainTakerResult",
	}
	key := strings.ToLower(assetType)
	if v, ok := aliases[key]; ok {
		return v, nil
	}
	if _, ok := aliases[strings.ReplaceAll(key, "-", "_")]; ok {
		return aliases[strings.ReplaceAll(key, "-", "_")], nil
	}
	// 允许直接传 MongoDB 集合名
	valid := []string{"asset", "RootDomain", "subdomain", "app", "mp", "UrlScan",
		"SensitiveResult", "DirScanResult", "crawler", "vulnerability",
		"PageMonitoring", "IPAsset", "SubdomainTakerResult"}
	for _, v := range valid {
		if v == assetType {
			return v, nil
		}
	}
	return "", fmt.Errorf("unsupported asset_type: %s", assetType)
}

type ginH map[string]any
