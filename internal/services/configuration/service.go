// configuration-------------------------------------
// @file      : service.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/29 21:20
// -------------------------------------------

package configuration

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/region"
	assetCommon "github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	nservice "github.com/Autumn-27/ScopeSentry/internal/services/node"

	"github.com/Autumn-27/ScopeSentry/internal/repositories/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	repo         common.Repository
	nodeService  nservice.Service
	dedupService assetCommon.DedupService
}

func NewService() *Service {
	return &Service{
		repo:         common.NewRepository(),
		nodeService:  nservice.NewService(),
		dedupService: assetCommon.NewDedupService(),
	}
}

const (
	collConfig        = "config"
	collNotification  = "notification"
	collSensitiveRule = "SensitiveRule"

	// regionConfigName 은 config 컬렉션에서 지역 프로필을 담는 키 이름이다.
	regionConfigName = "region"
)

// GetSubfinderContent 读取 SubfinderApiConfig
func (s *Service) GetSubfinderContent(ctx *gin.Context) (string, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "SubfinderApiConfig"}, bson.M{"_id": 0})
	if err != nil {
		return "", err
	}
	if v, ok := doc["value"].(string); ok {
		return v, nil
	}
	return "", nil
}

// SaveSubfinderContent 保存 SubfinderApiConfig 并通知节点刷新
func (s *Service) SaveSubfinderContent(ctx *gin.Context, content string) error {
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "SubfinderApiConfig"}, bson.M{"name": "SubfinderApiConfig", "value": content}); err != nil {
		return err
	}
	// 通知所有节点刷新 subfinder 配置
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "subfinder"})
	return nil
}

// GetRadContent 读取 RadConfig
func (s *Service) GetRadContent(ctx *gin.Context) (string, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "RadConfig"}, bson.M{"_id": 0})
	if err != nil {
		return "", err
	}
	if v, ok := doc["value"].(string); ok {
		return v, nil
	}
	return "", nil
}

// SaveRadContent 保存 RadConfig 并通知节点刷新
func (s *Service) SaveRadContent(ctx *gin.Context, content string) error {
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "RadConfig"}, bson.M{"name": "RadConfig", "value": content}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "rad"})
	return nil
}

// GetSystemData 返回 type=system 的所有键值
func (s *Service) GetSystemData(ctx *gin.Context) (map[string]interface{}, error) {
	list, err := s.repo.FindMany(ctx.Request.Context(), collConfig, bson.M{"type": "system"}, nil)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	for _, m := range list {
		name, _ := m["name"].(string)
		if name == "" {
			continue
		}
		out[name] = m["value"]
	}

	// 지역 프로필은 나중에 들어온 설정이라 기존 설치본에는 없다.
	// 화면이 항상 하나를 고른 상태로 보이게 기본값을 채워 돌려준다.
	current, _ := out[regionConfigName].(string)
	out[regionConfigName] = region.Parse(current).String()

	return out, nil
}

// SaveSystemData 保存 system 配置，并通知节点
func (s *Service) SaveSystemData(ctx *gin.Context, kv map[string]interface{}) error {
	// 지역이 바뀌었는지 저장 전에 확인한다. 저장할 때마다 규칙 상태를 다시 맞추면
	// 사용자가 개별 규칙을 켠 것까지 되돌려 버린다.
	var regionChange *region.Region
	if raw, ok := kv[regionConfigName]; ok {
		next := region.Parse(fmt.Sprintf("%v", raw))
		kv[regionConfigName] = next.String()
		if next != s.GetRegion(ctx.Request.Context()) {
			regionChange = &next
		}
	}

	for k, v := range kv {
		doc := bson.M{"type": "system", "name": k, "value": v}
		if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"type": "system", "name": k}, doc); err != nil {
			return err
		}
	}
	timezone, _ := kv["timezone"].(string)
	modulesConfig := fmt.Sprintf("%v", kv["ModulesConfig"]) // 兼容字符串或其他类型
	msg := timezone + "[*]" + modulesConfig
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "system", Content: msg})

	if regionChange != nil {
		if err := s.applyRegionSensitiveRules(ctx.Request.Context(), *regionChange); err != nil {
			// 설정 자체는 저장됐으므로 요청을 실패로 만들지는 않는다.
			logger.Error("failed to apply the region sensitive rules", zap.Error(err))
		}
	}

	return nil
}

// GetRegion 은 저장된 지역 프로필을 읽는다. 값이 없거나 모르는 값이면 기본값이다.
func (s *Service) GetRegion(ctx context.Context) region.Region {
	doc, err := s.repo.FindOne(ctx, collConfig, bson.M{"type": "system", "name": regionConfigName}, bson.M{"_id": 0})
	if err != nil {
		return region.Default
	}
	value, _ := doc["value"].(string)
	return region.Parse(value)
}

// applyRegionSensitiveRules 는 지역 전용 민감정보 규칙의 사용 여부를 프로필에
// 맞춘다. 어떤 규칙이 켜지고 꺼지는지는 region.SensitiveRuleStates 가 정하며,
// 시드와 마이그레이션도 같은 함수를 쓴다.
func (s *Service) applyRegionSensitiveRules(ctx context.Context, target region.Region) error {
	states, err := region.SensitiveRuleStates(target)
	if err != nil {
		return err
	}

	// 같은 상태로 갈 규칙을 묶어 업데이트 두 번으로 끝낸다.
	grouped := map[bool][]string{}
	for name, state := range states {
		grouped[state] = append(grouped[state], name)
	}

	for state, names := range grouped {
		if len(names) == 0 {
			continue
		}
		_, err := s.repo.UpdateMany(ctx,
			collSensitiveRule,
			bson.M{"name": bson.M{"$in": names}},
			bson.M{"$set": bson.M{"state": state}},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDeduplicationConfig 查询去重配置与下一次运行时间（若有）
func (s *Service) GetDeduplicationConfig(ctx *gin.Context) (map[string]interface{}, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "deduplication"}, bson.M{"_id": 0})
	if err != nil {
		// 若不存在，返回空 map
		return map[string]interface{}{}, nil
	}
	return doc, nil
}

// SaveDeduplicationConfig 保存去重配置；如需要立即运行，透传到节点或调度器（此处仅保存与刷新）
func (s *Service) SaveDeduplicationConfig(ctx *gin.Context, cfg models.DepConfig) error {
	cfg.Name = "deduplication"
	runNow := cfg.RunNow
	cfg.RunNow = false
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "deduplication"}, cfg); err != nil {
		return err
	}
	// 这里按需可扩展调用内部调度器或消息下发；当前版本仅保存
	if runNow {
		go func() {
			err := s.dedupService.DoAssetDeduplication()
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}()
	}
	return nil
}

// GetNotificationList 获取通知列表
func (s *Service) GetNotificationList(ctx *gin.Context) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.repo.Find(ctx.Request.Context(), collNotification, bson.M{}, nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// AddNotification 新增通知并刷新节点
func (s *Service) AddNotification(ctx *gin.Context, data models.Notification) error {
	_, err := s.repo.InsertOne(ctx.Request.Context(), collNotification, data)
	if err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}

// UpdateNotification 更新通知并刷新节点
func (s *Service) UpdateNotification(ctx *gin.Context, data models.UpdateNotification) error {
	if data.ID == "" {
		return fmt.Errorf("missing id")
	}
	objID, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}
	if err := s.repo.UpdateOne(ctx.Request.Context(), collNotification, bson.M{"_id": objID}, bson.M{"$set": data}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	if err := mongodb.FindAll("notification", bson.M{"state": true}, bson.M{"_id": 0, "method": 1, "url": 1, "contentType": 1, "data": 1, "state": 1}, &config.NotificationApi); err != nil {
		return err
	}
	return nil
}

// DeleteNotifications 批量删除通知并刷新节点
func (s *Service) DeleteNotifications(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// 将 string 转换为 ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid ObjectID: %s, error: %w", id, err)
		}
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	if _, err := s.repo.DeleteMany(ctx.Request.Context(), collNotification, filter); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}

// GetNotificationConfig 读取通知配置
func (s *Service) GetNotificationConfig(ctx *gin.Context) (map[string]interface{}, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "notification"}, bson.M{})
	if err != nil {
		return map[string]interface{}{}, nil
	}
	delete(doc, "_id")
	delete(doc, "type")
	delete(doc, "name")
	return doc, nil
}

// UpdateNotificationConfig 更新通知配置并刷新节点
func (s *Service) UpdateNotificationConfig(ctx *gin.Context, data map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if err := s.repo.UpdateOne(ctx.Request.Context(), collConfig, bson.M{"name": "notification"}, bson.M{"$set": data}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}
