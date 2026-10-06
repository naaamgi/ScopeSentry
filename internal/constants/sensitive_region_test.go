package constants

import (
	"encoding/json"
	"regexp"
	"testing"
)

// 규칙은 스캐너(Go, RE2)에서 컴파일되므로 여기서 같은 엔진으로 검증한다.
// 깨진 정규식은 스캔 중에야 드러나고, 그때는 결과가 그냥 비어 보인다.
func TestKoreanSensitiveRulesCompile(t *testing.T) {
	rules, err := SensitiveRulesKR()
	if err != nil {
		t.Fatalf("failed to load the Korean rules: %v", err)
	}
	if len(rules) == 0 {
		t.Fatal("no Korean rules were loaded")
	}

	for _, r := range rules {
		if r.Name == "" {
			t.Error("a rule has no name")
		}
		if _, err := regexp.Compile(r.Regular); err != nil {
			t.Errorf("%s: regex does not compile: %v", r.Name, err)
		}
		switch r.Color {
		case "red", "orange", "yellow", "green", "cyan":
		default:
			t.Errorf("%s: unexpected color %q", r.Name, r.Color)
		}
	}
}

func TestKoreanSensitiveRulesMatchWhatTheyShould(t *testing.T) {
	rules, err := SensitiveRulesKR()
	if err != nil {
		t.Fatalf("failed to load the Korean rules: %v", err)
	}

	byName := make(map[string]string, len(rules))
	for _, r := range rules {
		byName[r.Name] = r.Regular
	}

	// 기존 규칙들과 같은 방식으로 앞뒤 한 글자를 경계로 쓰기 때문에, 본문 중간에
	// 끼어 있는 형태로 확인한다.
	cases := []struct {
		rule    string
		matches []string
		rejects []string
	}{
		{
			rule:    "Korean Resident Registration Number",
			matches: []string{"가입자 900101-1234567 확인", "id=9001011234567&", "x=001231-4567890;"},
			rejects: []string{
				"x 901301-1234567 y", // 13월
				"x 900132-1234567 y", // 32일
				"x 900101-9234567 y", // 성별 자리가 1~4 밖
				"x 900101-5234567 y", // 외국인등록번호 쪽
			},
		},
		{
			rule:    "Korean Foreign Registration Number",
			matches: []string{"x 900101-5234567 y", "x 900101-8234567 y"},
			rejects: []string{"x 900101-1234567 y", "x 900101-9234567 y"},
		},
		{
			rule: "Korean Mobile Phone Number",
			matches: []string{
				"tel: 010-1234-5678 .",
				"tel: 01012345678 .",
				"tel: +82 10-1234-5678 .",
				"tel: 019-123-4567 .",
			},
			rejects: []string{
				"tel: 02-1234-5678 .",  // 지역번호
				"tel: 015-1234-5678 .", // 쓰이지 않는 식별번호
			},
		},
		{
			rule:    "Korean Business Registration Number",
			matches: []string{"사업자 123-45-67890 입니다"},
			rejects: []string{"x 1234-56-78901 y", "x 123-456-7890 y"},
		},
		{
			rule:    "Korean Driver License Number",
			matches: []string{"면허 11-22-333333-44 끝"},
			rejects: []string{"x 11-22-33333-44 y"},
		},
		{
			rule:    "Korean Passport Number",
			matches: []string{"passport=M12345678&", "passport=S87654321&"},
			rejects: []string{"passport=M1234567&", "passport=Z12345678&"},
		},
		{
			rule:    "Toss Payments Secret Key",
			matches: []string{"test_sk_abcdefghij1234567890ab", "live_ck_ABCDEFGHIJ1234567890ab"},
			rejects: []string{"test_sk_tooshort", "demo_sk_abcdefghij1234567890ab"},
		},
		{
			rule:    "Kakao REST API Key",
			matches: []string{"Authorization: KakaoAK 0123456789abcdef0123456789abcdef"},
			rejects: []string{"Authorization: KakaoAK 0123456789abcdef", "KakaoAKx 0123456789abcdef0123456789abcdef"},
		},
	}

	if len(cases) != len(rules) {
		t.Errorf("%d rules are shipped but %d are covered here", len(rules), len(cases))
	}

	for _, c := range cases {
		pattern, ok := byName[c.rule]
		if !ok {
			t.Errorf("no rule named %q is shipped", c.rule)
			continue
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Errorf("%s: regex does not compile: %v", c.rule, err)
			continue
		}
		for _, s := range c.matches {
			if !re.MatchString(s) {
				t.Errorf("%s: should match %q", c.rule, s)
			}
		}
		for _, s := range c.rejects {
			if re.MatchString(s) {
				t.Errorf("%s: should not match %q", c.rule, s)
			}
		}
	}
}

// 중국 전용 규칙 이름은 기본 규칙 파일을 보고 찾으므로, 파일과 어긋나면
// 지역을 바꿔도 규칙이 켜지거나 꺼지지 않는다.
func TestChineseRuleNamesExistInTheBaseFile(t *testing.T) {
	var base []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(SensData), &base); err != nil {
		t.Fatalf("failed to parse the base rules: %v", err)
	}

	present := make(map[string]bool, len(base))
	for _, r := range base {
		present[r.Name] = true
	}

	for _, name := range SensitiveRuleNamesCN {
		if !present[name] {
			t.Errorf("SensitiveRuleNamesCN lists %q, which the base rule file does not contain", name)
		}
	}
}

// 국내 규칙 이름이 기본 파일의 이름과 겹치면 지역 전환이 서로의 규칙을 건드린다.
func TestKoreanRuleNamesDoNotCollideWithTheBaseFile(t *testing.T) {
	var base []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(SensData), &base); err != nil {
		t.Fatalf("failed to parse the base rules: %v", err)
	}
	present := make(map[string]bool, len(base))
	for _, r := range base {
		present[r.Name] = true
	}

	names, err := SensitiveRuleNamesKR()
	if err != nil {
		t.Fatalf("failed to read the Korean rule names: %v", err)
	}
	for _, name := range names {
		if present[name] {
			t.Errorf("Korean rule %q has the same name as a base rule", name)
		}
	}
}
