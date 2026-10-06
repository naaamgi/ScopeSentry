package project

import "testing"

func TestRootDomainForTarget(t *testing.T) {
	cases := []struct {
		name   string
		target string
		want   string
	}{
		{"plain domain", "www.example.com", "example.com"},
		{"domain with path", "https://a.b.example.com/x", "example.com"},
		{"company name", "CMP:Example Corp", "CMP:Example Corp"},
		{"app name", "APP:Example", "APP:Example"},
		{"app package", "APP-ID:com.example.app", "APP-ID:com.example.app"},

		// ICP 번호는 하위 등록번호를 떼어 같은 등록번호끼리 모은다.
		{"icp without suffix", "ICP:京ICP证1234号", "ICP:京ICP证1234号"},
		{"icp with suffix", "ICP:京ICP证1234号-1", "ICP:京ICP证1234号"},
		{"icp with second suffix", "ICP:京ICP证1234号-12", "ICP:京ICP证1234号"},

		// 사업자등록번호는 번호 자체에 하이픈이 있으므로 잘라서는 안 된다.
		{"business registration number", "BRN:123-45-67890", "BRN:123-45-67890"},
		{"business registration number without hyphens", "BRN:1234567890", "BRN:1234567890"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := rootDomainForTarget(c.target); got != c.want {
				t.Errorf("rootDomainForTarget(%q) = %q, want %q", c.target, got, c.want)
			}
		})
	}
}

func TestHasTargetPrefix(t *testing.T) {
	withPrefix := []string{
		"CMP:Example Corp",
		"ICP:京ICP证1234号",
		"BRN:123-45-67890",
		"APP:Example",
		"APP-ID:com.example.app",
	}
	for _, target := range withPrefix {
		if !hasTargetPrefix(target) {
			t.Errorf("hasTargetPrefix(%q) = false, want true", target)
		}
	}

	withoutPrefix := []string{
		"example.com",
		"192.168.1.1",
		"192.168.0.0/18",
		"CIDR:192.168.0.0/18",
	}
	for _, target := range withoutPrefix {
		if hasTargetPrefix(target) {
			t.Errorf("hasTargetPrefix(%q) = true, want false", target)
		}
	}
}
