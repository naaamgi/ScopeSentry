package constants

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// 두 설정은 문자열 리터럴로 보관되어 설치 시 DB 에 들어가고, 사용자가 설정 화면의
// 편집기에서 고친다. 주석만 손봤더라도 리터럴을 건드린 뒤에는 모양이 그대로인지
// 확인해 둔다.
func TestModulesConfigIsValidYAML(t *testing.T) {
	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(ModulesConfig), &parsed); err != nil {
		t.Fatalf("ModulesConfig does not parse as YAML: %v", err)
	}

	// 모듈별 동시 실행 수는 스캐너가 모듈 이름으로 찾아 쓴다. 키가 바뀌면 설정이
	// 조용히 무시되므로 키 집합을 고정한다.
	want := []string{
		"maxGoroutineCount",
		"subdomainScan", "subdomainSecurity", "assetMapping", "assetHandle",
		"portScanPreparation", "portScan", "portFingerprint",
		"URLScan", "URLSecurity", "webCrawler", "dirScan", "vulnerabilityScan",
	}
	for _, key := range want {
		if _, ok := parsed[key]; !ok {
			t.Errorf("ModulesConfig is missing the key %q", key)
		}
	}
	if len(parsed) != len(want) {
		t.Errorf("ModulesConfig has %d keys, want %d", len(parsed), len(want))
	}
}

// RadConfig 는 업스트림에서 이미 들여쓰기가 흐트러져 있다. 세 개의
// *_filter_config 블록의 하위 키가 최상위에 붙어 있어 키가 겹치고, 엄격한 YAML
// 파서는 이를 거부한다. 고치면 rad 플러그인에 전달되는 설정이 달라지므로
// 별도 변경으로 다뤄야 한다. 여기서는 그 상태를 기록해 두고, 줄과 키가
// 그대로인지만 확인한다.
func TestRadConfigShapeIsUnchanged(t *testing.T) {
	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(RadConfig), &parsed); err == nil {
		t.Error("RadConfig now parses as strict YAML; the upstream indentation defect looks fixed, so update this test")
	} else if !strings.Contains(err.Error(), "already defined") {
		t.Errorf("RadConfig fails to parse for an unexpected reason: %v", err)
	}

	lines := strings.Split(RadConfig, "\n")
	if len(lines) != 64 {
		t.Errorf("RadConfig has %d lines, want 64", len(lines))
	}

	// 줄마다 주석 앞의 코드 부분에서 키 이름만 뽑아 센다.
	keys := 0
	for _, line := range lines {
		code := strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		code = strings.TrimPrefix(code, "- ")
		if strings.Contains(code, ":") {
			keys++
		}
	}
	if keys != 64 {
		t.Errorf("RadConfig has %d keys, want 64", keys)
	}
}

// 설정은 사용자가 읽는 쪽이라, 한쪽만 영문으로 남아 섞이지 않게 한다.
func TestDefaultConfigCommentsAreNotChinese(t *testing.T) {
	for name, content := range map[string]string{"ModulesConfig": ModulesConfig, "RadConfig": RadConfig} {
		for i, line := range strings.Split(content, "\n") {
			idx := strings.Index(line, "#")
			if idx < 0 {
				continue
			}
			for _, r := range line[idx:] {
				if r >= 0x4e00 && r <= 0x9fff {
					t.Errorf("%s line %d still has a Chinese comment: %s", name, i+1, strings.TrimSpace(line))
					break
				}
			}
		}
	}
}
