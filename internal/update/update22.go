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
func Update22() {
	ctx := context.Background()

	profile, firstRun := ensureRegionConfig(ctx)
	ensureKoreanSensitiveRules(ctx, profile)

	// 프로필 설정이 방금 처음 들어간 설치본만 규칙 상태를 맞춘다. 이미 설정이
	// 있었다면 운영자가 고른 상태이므로, 업그레이드가 손으로 켠 규칙을
	// 되돌려서는 안 된다.
	if firstRun {
		alignRegionSensitiveRules(ctx, profile)
	}
}

// ensureRegionConfig 는 지역 프로필 설정을 확인하고, 없으면 기본값으로 넣는다.
// 돌려주는 bool 은 이번 호출에서 처음 넣었는지를 뜻한다.
func ensureRegionConfig(ctx context.Context) (region.Region, bool) {
	coll := mongodb.DB.Collection("config")
	filter := bson.M{"type": "system", "name": "region"}

	var stored struct {
		Value string `bson:"value"`
	}
	err := coll.FindOne(ctx, filter).Decode(&stored)
	if err == nil {
		return region.Parse(stored.Value), false
	}

	if _, err := coll.InsertOne(ctx, bson.M{
		"type":  "system",
		"name":  "region",
		"value": region.Default.String(),
	}); err != nil {
		logger.Error("failed to insert the region config", zap.Error(err))
		return region.Default, false
	}

	logger.Info("set the region profile", zap.String("region", region.Default.String()))
	return region.Default, true
}

// ensureKoreanSensitiveRules 는 국내 전용 규칙이 없으면 넣는다.
// 상태는 프로필에 맞춘 값으로 들어간다.
func ensureKoreanSensitiveRules(ctx context.Context, profile region.Region) {
	rules, err := constants.SensitiveRulesKR()
	if err != nil {
		logger.Error("failed to read the Korean sensitive rules", zap.Error(err))
		return
	}
	states, err := region.SensitiveRuleStates(profile)
	if err != nil {
		logger.Error("failed to resolve the region sensitive rule states", zap.Error(err))
		return
	}

	coll := mongodb.DB.Collection("SensitiveRule")
	added := 0
	for _, r := range rules {
		existing, err := coll.CountDocuments(ctx, bson.M{"name": r.Name})
		if err != nil {
			logger.Error("failed to look up a sensitive rule", zap.String("name", r.Name), zap.Error(err))
			continue
		}
		if existing > 0 {
			continue
		}
		if _, err := coll.InsertOne(ctx, bson.M{
			"name":    r.Name,
			"regular": r.Regular,
			"color":   r.Color,
			"state":   states[r.Name],
		}); err != nil {
			logger.Error("failed to insert a sensitive rule", zap.String("name", r.Name), zap.Error(err))
			continue
		}
		added++
	}

	if added > 0 {
		logger.Info("added Korean sensitive rules", zap.Int("count", added))
	}
}

// alignRegionSensitiveRules 는 지역 전용 규칙의 사용 여부를 프로필에 맞춘다.
// 설정 화면에서 프로필을 바꿀 때와 같은 판단(region.SensitiveRuleStates)을 쓴다.
func alignRegionSensitiveRules(ctx context.Context, profile region.Region) {
	states, err := region.SensitiveRuleStates(profile)
	if err != nil {
		logger.Error("failed to resolve the region sensitive rule states", zap.Error(err))
		return
	}

	grouped := map[bool][]string{}
	for name, state := range states {
		grouped[state] = append(grouped[state], name)
	}

	coll := mongodb.DB.Collection("SensitiveRule")
	for state, names := range grouped {
		if len(names) == 0 {
			continue
		}
		if _, err := coll.UpdateMany(ctx,
			bson.M{"name": bson.M{"$in": names}},
			bson.M{"$set": bson.M{"state": state}},
		); err != nil {
			logger.Error("failed to align the region sensitive rules", zap.Error(err))
		}
	}
}
