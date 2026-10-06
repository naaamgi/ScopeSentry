package constants

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/models"
)

// 지역 전용 민감정보 규칙.
//
// 기본 규칙 파일(assets/ScopeSentry.SensitiveRule.json)에는 중국 주민번호·
// 휴대전화·은행카드 규칙이 섞여 있다. 국내 환경에서는 쓸 일이 없고 오탐만
// 늘어나므로 지역 프로필에 따라 켜고 끈다. 국내 규칙은 별도 파일로 둔다.

//go:embed assets/SensitiveRule.KR.json
var sensDataKR string

// SensitiveRuleNamesCN 은 기본 규칙 파일에 들어 있는 중국 전용 규칙 이름이다.
// 이 이름으로 규칙을 찾으므로 기본 파일의 name 과 정확히 같아야 한다.
var SensitiveRuleNamesCN = []string{
	"Chinese IDCard",
	"Chinese Mobile Number",
	"Chinese Bank Card ID",
}

// SensitiveRulesKR 은 국내 전용 민감정보 규칙을 돌려준다.
func SensitiveRulesKR() ([]models.SensitiveRuleItem, error) {
	var rules []models.SensitiveRuleItem
	if err := json.Unmarshal([]byte(sensDataKR), &rules); err != nil {
		return nil, fmt.Errorf("failed to parse the Korean sensitive rules: %w", err)
	}
	return rules, nil
}

// SensitiveRuleNamesKR 은 국내 전용 규칙 이름만 돌려준다.
func SensitiveRuleNamesKR() ([]string, error) {
	rules, err := SensitiveRulesKR()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(rules))
	for _, r := range rules {
		names = append(names, r.Name)
	}
	return names, nil
}

// SensitiveRuleDefaultState 는 규칙 이름별 배포 기본값(켜짐 여부)을 돌려준다.
// 지역을 바꿀 때 그 지역 규칙을 "원래 켜져 있던 것만" 켜는 데 쓴다.
func SensitiveRuleDefaultState() (map[string]bool, error) {
	var base []models.SensitiveRuleItem
	if err := json.Unmarshal([]byte(SensData), &base); err != nil {
		return nil, fmt.Errorf("failed to parse the base sensitive rules: %w", err)
	}

	kr, err := SensitiveRulesKR()
	if err != nil {
		return nil, err
	}

	states := make(map[string]bool, len(base)+len(kr))
	for _, r := range base {
		states[r.Name] = r.State
	}
	for _, r := range kr {
		states[r.Name] = r.State
	}
	return states, nil
}
