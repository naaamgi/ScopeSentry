package region

import "github.com/Autumn-27/ScopeSentry/internal/constants"

// SensitiveRuleStates 는 지역 전용 민감정보 규칙이 프로필 target 에서 어떤 상태여야
// 하는지를 규칙 이름별로 돌려준다.
//
// 고른 지역의 규칙은 배포 기본값대로 켜고(기본값이 꺼짐인 규칙은 그대로 꺼진다),
// 다른 지역의 규칙은 끈다. 어느 지역에도 속하지 않는 규칙은 결과에 들어가지
// 않으므로, 호출부가 이 표에 있는 이름만 건드리면 된다.
//
// 설치 시 시드, 기존 설치본 마이그레이션, 설정 화면에서의 프로필 변경이 모두 이
// 함수를 쓴다. 세 곳이 서로 다른 판단을 하면 DB 상태와 설정이 어긋난다.
func SensitiveRuleStates(target Region) (map[string]bool, error) {
	defaults, err := constants.SensitiveRuleDefaultState()
	if err != nil {
		return nil, err
	}

	krNames, err := constants.SensitiveRuleNamesKR()
	if err != nil {
		return nil, err
	}

	owners := []struct {
		owner Region
		names []string
	}{
		{CN, constants.SensitiveRuleNamesCN},
		{KR, krNames},
	}

	states := make(map[string]bool)
	for _, group := range owners {
		for _, name := range group.names {
			states[name] = group.owner == target && defaults[name]
		}
	}
	return states, nil
}
