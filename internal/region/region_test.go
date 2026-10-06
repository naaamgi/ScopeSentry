package region

import "testing"

func TestParse(t *testing.T) {
	cases := map[string]Region{
		"CN":       CN,
		"cn":       CN,
		" kr ":     KR,
		"KR":       KR,
		"GLOBAL":   Global,
		"global":   Global,
		"":         Default,
		"JP":       Default,
		"nonsense": Default,
	}

	for input, want := range cases {
		if got := Parse(input); got != want {
			t.Errorf("Parse(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestAllIsParseable(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("All returned nothing")
	}

	seen := make(map[Region]bool, len(all))
	for _, r := range all {
		if seen[r] {
			t.Errorf("All lists %q twice", r)
		}
		seen[r] = true

		// 설정 화면이 보여준 값을 그대로 저장하므로, 왕복이 깨지면 안 된다.
		if got := Parse(r.String()); got != r {
			t.Errorf("Parse(%q) = %q, want %q", r.String(), got, r)
		}
	}

	if !seen[Default] {
		t.Errorf("All does not include the default profile %q", Default)
	}
}
