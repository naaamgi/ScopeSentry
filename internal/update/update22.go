package update

import (
	"context"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/region"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// Update22 는 지역 프로필 설정과 국내 전용 민감정보 규칙을 기존 설치본에 채운다.
// 멱등하므로 매번 호출해도 된다.
//
// 규칙은 꺼진 상태로 들어간다. 지역을 KR 로 바꾸는 시점에 설정 화면이 켜주며,
// 여기서 미리 켜면 중국 환경으로 쓰던 설치본에 갑자기 국내 PII 규칙이 돌아간다.
func Update22() {
	ctx := context.Background()

	// 1) 지역 프로필 기본값
	configColl := mongodb.DB.Collection("config")
	filter := bson.M{"type": "system", "name": "region"}
	count, err := configColl.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("failed to look up the region config", zap.Error(err))
	} else if count == 0 {
		_, err = configColl.InsertOne(ctx, bson.M{
			"type":  "system",
			"name":  "region",
			"value": region.Default.String(),
		})
		if err != nil {
			logger.Error("failed to insert the region config", zap.Error(err))
		}
	}

	// 2) 국내 전용 민감정보 규칙
	rules, err := constants.SensitiveRulesKR()
	if err != nil {
		logger.Error("failed to read the Korean sensitive rules", zap.Error(err))
		return
	}

	sensitiveColl := mongodb.DB.Collection("SensitiveRule")
	added := 0
	for _, r := range rules {
		existing, err := sensitiveColl.CountDocuments(ctx, bson.M{"name": r.Name})
		if err != nil {
			logger.Error("failed to look up a sensitive rule", zap.String("name", r.Name), zap.Error(err))
			continue
		}
		if existing > 0 {
			continue
		}
		_, err = sensitiveColl.InsertOne(ctx, bson.M{
			"name":    r.Name,
			"regular": r.Regular,
			"color":   r.Color,
			"state":   false,
		})
		if err != nil {
			logger.Error("failed to insert a sensitive rule", zap.String("name", r.Name), zap.Error(err))
			continue
		}
		added++
	}

	if added > 0 {
		logger.Info("added Korean sensitive rules", zap.Int("count", added))
	}
}
