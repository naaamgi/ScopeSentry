package region

import (
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
)

func TestSensitiveRuleStatesEnablesOnlyTheTargetRegion(t *testing.T) {
	krNames, err := constants.SensitiveRuleNamesKR()
	if err != nil {
		t.Fatalf("failed to read the Korean rule names: %v", err)
	}
	defaults, err := constants.SensitiveRuleDefaultState()
	if err != nil {
		t.Fatalf("failed to read the default states: %v", err)
	}

	owners := map[Region][]string{
		KR: krNames,
		CN: constants.SensitiveRuleNamesCN,
	}

	for _, target := range All() {
		states, err := SensitiveRuleStates(target)
		if err != nil {
			t.Fatalf("%s: %v", target, err)
		}

		for owner, names := range owners {
			for _, name := range names {
				got, ok := states[name]
				if !ok {
					t.Errorf("%s: %q is missing from the result", target, name)
					continue
				}
				// 고른 지역의 규칙은 배포 기본값대로, 나머지는 꺼짐.
				want := owner == target && defaults[name]
				if got != want {
					t.Errorf("%s: %q = %v, want %v", target, name, got, want)
				}
			}
		}

		// 지역에 속하지 않는 규칙은 손대지 않으므로 결과에 없어야 한다.
		if _, ok := states["JSON Web Token"]; ok {
			t.Errorf("%s: a non-region rule leaked into the result", target)
		}
	}
}

// 기본값이 꺼짐인 규칙(여권번호, 중국 은행카드)은 그 지역을 골라도 꺼진 채여야 한다.
func TestSensitiveRuleStatesRespectsRulesShippedDisabled(t *testing.T) {
	states, err := SensitiveRuleStates(KR)
	if err != nil {
		t.Fatalf("SensitiveRuleStates(KR): %v", err)
	}
	if states["Korean Passport Number"] {
		t.Error("the passport rule ships disabled and should stay disabled under KR")
	}
	if !states["Korean Resident Registration Number"] {
		t.Error("the resident registration rule should be enabled under KR")
	}

	states, err = SensitiveRuleStates(CN)
	if err != nil {
		t.Fatalf("SensitiveRuleStates(CN): %v", err)
	}
	if states["Chinese Bank Card ID"] {
		t.Error("the bank card rule ships disabled and should stay disabled under CN")
	}
	if !states["Chinese IDCard"] {
		t.Error("the Chinese ID card rule should be enabled under CN")
	}
	if states["Korean Resident Registration Number"] {
		t.Error("Korean rules should be disabled under CN")
	}
}

// 기본 프로필에서 실제로 어떤 규칙이 켜지는지 고정해 둔다. 기본값을 바꾸면
// 설치 직후 돌아가는 PII 규칙이 달라지므로, 의도한 변경인지 여기서 걸린다.
func TestDefaultProfileEnablesTheKoreanRules(t *testing.T) {
	if Default != KR {
		t.Fatalf("Default = %q; update this test if that is intended", Default)
	}

	states, err := SensitiveRuleStates(Default)
	if err != nil {
		t.Fatalf("SensitiveRuleStates(Default): %v", err)
	}
	for _, name := range constants.SensitiveRuleNamesCN {
		if states[name] {
			t.Errorf("%q is enabled under the default profile", name)
		}
	}
}
