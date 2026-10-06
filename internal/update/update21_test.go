package update

import (
	"strings"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/i18n"
)

// 기본 제공 플러그인의 설명문은 DB 에 메시지 키로 저장된다. 플러그인을 새로
// 추가하면서 키나 번역을 빠뜨리면 화면에 키 문자열이 그대로 보이므로 여기서 막는다.
func TestBuiltinPluginTextsAreMessageKeys(t *testing.T) {
	for _, p := range constants.Plugins {
		for field, value := range map[string]string{"Help": p.Help, "Introduction": p.Introduction} {
			if !strings.HasPrefix(value, i18n.BuiltinPluginPrefix) {
				t.Errorf("plugin %s: %s is %q, want a %s* message key", p.Name, field, value, i18n.BuiltinPluginPrefix)
			}
		}
	}
}

func TestBuiltinPluginTextsResolveInEveryLocale(t *testing.T) {
	for _, locale := range i18n.SupportedLocales() {
		for _, p := range constants.Plugins {
			for field, key := range map[string]string{"Help": p.Help, "Introduction": p.Introduction} {
				got := i18n.Translate(locale, key)
				if got == key {
					t.Errorf("%s: plugin %s %s key %q has no translation", locale, p.Name, field, key)
				}
			}
		}
	}
}

// 마이그레이션 표가 시드와 어긋나면, 기존 설치본의 플러그인 설명이 중국어로 남거나
// 반대로 사라진 플러그인을 계속 들고 있게 된다.
func TestMigrationTableCoversEverySeededPlugin(t *testing.T) {
	seeded := make(map[string]string, len(constants.Plugins))
	for _, p := range constants.Plugins {
		seeded[p.Hash] = p.Name
	}

	for hash, name := range seeded {
		if _, ok := builtinPluginTexts[hash]; !ok {
			t.Errorf("plugin %s (hash %s) is seeded but missing from builtinPluginTexts", name, hash)
		}
	}
	for hash := range builtinPluginTexts {
		if _, ok := seeded[hash]; !ok {
			t.Errorf("builtinPluginTexts has hash %s, which is no longer seeded", hash)
		}
	}
}

func TestTranslateStoredLeavesUserTextAlone(t *testing.T) {
	// 사용자가 가져온 플러그인의 설명은 자유 문장이므로 그대로 지나가야 한다.
	custom := "my own plugin, takes -foo"
	if got := i18n.TranslateStored("ko-KR", custom); got != custom {
		t.Errorf("TranslateStored changed user text to %q", got)
	}

	key := "plugin.builtin.nuclei.introduction"
	if got := i18n.TranslateStored("ko-KR", key); got != "취약점 스캔" {
		t.Errorf("TranslateStored(ko-KR, %s) = %q", key, got)
	}
}
