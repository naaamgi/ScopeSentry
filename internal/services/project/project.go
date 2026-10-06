package project

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"sort"
	"strings"
	"sync"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"golang.org/x/sync/errgroup"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	repo "github.com/Autumn-27/ScopeSentry/internal/repositories/project"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetProjectsByTag(ctx *gin.Context) ([]models.TagGroup, error)
	GetTarget(ctx *gin.Context, id string) (string, error)
	GetTargets(ctx context.Context, ids []string) (string, error)

	GetProjectsData(ctx *gin.Context, search string, pageIndex, pageSize int) (models.ProjectListResponse, error)
	GetProjectContent(ctx *gin.Context, id string) (*models.ProjectContentResponse, error)
	AddProject(ctx *gin.Context, p *models.Project) error
	DeleteProjects(ctx *gin.Context, ids []string, delA bool) error
	UpdateProject(ctx *gin.Context, p *models.UpdateProject) error
	UpdateAssetsProject(ctx *gin.Context, id string) error
}

type service struct {
	projectRepo repo.Repository
}

func NewService() Service {
	return &service{
		projectRepo: repo.NewRepository(),
	}
}

func (s *service) GetProjectsByTag(ctx *gin.Context) ([]models.TagGroup, error) {
	// 在service层定义需要的字段
	projection := bson.M{
		"_id":  1,
		"name": 1,
		"tag":  1,
	}

	projects, err := s.projectRepo.GetProjectsByTag(ctx, projection)
	if err != nil {
		return nil, err
	}

	// 按tag分组
	tagMap := make(map[string][]models.TagProject)
	for _, p := range projects {
		tagProject := models.TagProject{
			Value: p.ID.Hex(),
			Label: p.Name,
		}
		tagMap[p.Tag] = append(tagMap[p.Tag], tagProject)
	}

	// 转换为响应格式并按tag排序
	var result []models.TagGroup
	for tag, projects := range tagMap {
		tagGroup := models.TagGroup{
			Label:    tag,
			Value:    "",
			Children: projects,
		}
		result = append(result, tagGroup)
	}

	// 按tag排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result, nil
}

func (s *service) GetTarget(ctx *gin.Context, id string) (string, error) {
	return s.projectRepo.GetTarget(ctx, id)
}

func (s *service) GetTargets(ctx context.Context, ids []string) (string, error) {
	return s.projectRepo.GetTargets(ctx, ids)
}

func (s *service) GetProjectsData(ctx *gin.Context, search string, pageIndex, pageSize int) (models.ProjectListResponse, error) {
	// tag 聚合统计
	tagCounts, err := s.projectRepo.AggregateTagCounts(ctx, search)
	if err != nil {
		return models.ProjectListResponse{}, err
	}
	// 计算 All 总数
	all := 0
	for _, v := range tagCounts {
		all += v
	}
	tagCounts["All"] = all

	result := make(map[string][]models.ProjectBrief)
	var mu sync.Mutex
	base := ctx.Request.Context()
	g, _ := errgroup.WithContext(base)
	for tag := range tagCounts {
		t := tag
		g.Go(func() error {
			items, err := s.projectRepo.ListProjectsByTag(ctx, search, t, pageIndex, pageSize)
			if err != nil {
				return err
			}
			mu.Lock()
			result[t] = items
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return models.ProjectListResponse{}, err
	}
	return models.ProjectListResponse{Result: result, Tag: tagCounts}, nil
}

func (s *service) GetProjectContent(ctx *gin.Context, id string) (*models.ProjectContentResponse, error) {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	target, err := s.projectRepo.GetTarget(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := &models.ProjectContentResponse{
		Name:           p.Name,
		Tag:            p.Tag,
		Target:         target,
		Node:           p.Node,
		Logo:           p.Logo,
		ScheduledTasks: p.ScheduledTasks,
		Hour:           p.Hour,
		AllNode:        p.AllNode,
		Duplicates:     p.Duplicates,
		Template:       p.Template,
		Ignore:         p.Ignore,
	}
	return resp, nil
}

func (s *service) AddProject(ctx *gin.Context, p *models.Project) error {
	// 名称唯一
	exists, err := s.projectRepo.ExistsByName(ctx, p.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("name already exists")
	}
	p.Tp = "project"
	p.Target = strings.TrimSpace(p.Target)
	targetList, err := helper.GetTargetList(p.Target, p.Ignore)
	if err != nil {
		return err
	}
	rootDomains := []string{}
	for _, tg := range targetList {
		rootDomain := rootDomainForTarget(tg)

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}
	p.RootDomains = rootDomains
	rawTarget := p.Target
	p.Target = ""
	projectId, err := s.projectRepo.InsertProject(ctx, p)
	if err != nil {
		return err
	}
	err = s.projectRepo.UpsertProjectTarget(ctx, projectId, rawTarget)
	if err != nil {
		return err
	}
	go func() {
		err = s.projectRepo.UpdateAssetsProject(context.Background(), rootDomains, projectId, false)
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	//id, err := s.projectRepo.InsertProject(ctx, p)
	//if err != nil {
	//	return err
	//}
	//if err := s.projectRepo.UpsertProjectTarget(ctx, id, p.Target); err != nil {
	//	return err
	//}
	// 调度
	//if p.ScheduledTasks {
	//	if err := s.projectRepo.CreateOrUpdateProjectSchedule(ctx, id, p.Hour, true, p.Name); err != nil {
	//		return err
	//	}
	//}
	// runNow 可后续接入任务系统
	return nil
}

func (s *service) UpdateProject(ctx *gin.Context, p *models.UpdateProject) error {
	p.Tp = "project"
	p.Target = strings.TrimSpace(p.Target)
	targetList, err := helper.GetTargetList(p.Target, p.Ignore)
	if err != nil {
		return err
	}
	rootDomains := []string{}
	for _, tg := range targetList {
		rootDomain := rootDomainForTarget(tg)

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}
	p.RootDomains = rootDomains
	rawTarget := p.Target
	p.Target = ""
	if err := s.projectRepo.UpdateProject(ctx, p); err != nil {
		return err
	}

	if err = s.projectRepo.UpsertProjectTarget(ctx, p.ID, rawTarget); err != nil {
		return err
	}
	go func() {
		err = s.projectRepo.UpdateAssetsProject(context.Background(), rootDomains, p.ID, true)
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	return nil
}

func (s *service) DeleteProjects(ctx *gin.Context, ids []string, delA bool) error {

	// 删除项目
	if err := s.projectRepo.DeleteProjects(ctx, ids); err != nil {
		return err
	}
	if err := s.projectRepo.DeleteProjectTargets(ctx, ids); err != nil {
		return err
	}
	if delA {
		go func() {
			s.projectRepo.DeleteProjectAsset(ctx, ids)
		}()

	}
	return nil
}

// targetPrefixes 는 도메인이 아닌 타깃 앞에 붙는 식별자다.
//
//	CMP:    회사명
//	ICP:    중국 ICP 등록번호
//	BRN:    국내 사업자등록번호
//	APP:    앱 이름
//	APP-ID: 앱 패키지명
//
// 서버는 이 타깃을 파싱해서 보관할 뿐이고, 실제 수집은 스캐너 플러그인이 한다.
var targetPrefixes = []string{"CMP:", "ICP:", "BRN:", "APP:", "APP-ID:"}

func hasTargetPrefix(target string) bool {
	for _, prefix := range targetPrefixes {
		if strings.Contains(target, prefix) {
			return true
		}
	}
	return false
}

// rootDomainForTarget 은 프로젝트 타깃 한 줄을 루트 도메인 목록에 넣을 값으로
// 바꾼다. 프리픽스가 붙은 타깃은 도메인이 아니라서 그대로 쓴다.
//
// ICP 번호만 예외로 뒤에 붙는 하위 등록번호(-1, -2 …)를 떼어, 같은 등록번호의
// 타깃이 하나로 모이게 한다. BRN 은 번호 자체가 하이픈을 품고 있어서
// (123-45-67890) 같은 처리를 하면 번호가 잘린다.
func rootDomainForTarget(target string) string {
	if !hasTargetPrefix(target) {
		rootDomain, _ := helper.GetRootDomain(target)
		return rootDomain
	}
	if strings.Contains(target, "ICP:") {
		return "ICP:" + getBeforeLastDash(strings.ReplaceAll(target, "ICP:", ""))
	}
	return target
}

func getBeforeLastDash(s string) string {
	index := strings.LastIndex(s, "-") // 查找最后一个 '-'
	if index != -1 {
		return s[:index]
	}
	return s
}

func (s *service) UpdateAssetsProject(ctx *gin.Context, id string) error {
	// 获取项目信息
	project, err := s.projectRepo.FindByID(ctx.Request.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	// 获取项目目标
	target, err := s.projectRepo.GetTarget(ctx.Request.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to get project target: %w", err)
	}

	// 解析目标获取根域名列表
	targetList, err := helper.GetTargetList(target, project.Ignore)
	if err != nil {
		return fmt.Errorf("failed to parse target list: %w", err)
	}

	rootDomains := []string{}
	for _, tg := range targetList {
		rootDomain := rootDomainForTarget(tg)

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}

	// 更新资产项目关联
	err = s.projectRepo.UpdateAssetsProject(ctx.Request.Context(), rootDomains, id, true)
	if err != nil {
		return fmt.Errorf("failed to update assets project: %w", err)
	}

	return nil
}
