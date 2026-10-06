package i18n

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNormalizeLocale(t *testing.T) {
	cases := []struct {
		header string
		want   string
	}{
		{"", DefaultLocale},
		{"*", DefaultLocale},
		{"ja", DefaultLocale},
		{"zh-CN", "zh-CN"},
		{"zh", "zh-CN"},
		{"zh-Hans", "zh-CN"},
		{"en", "en-US"},
		{"en-GB", "en-US"},
		{"ko", "ko-KR"},
		{"ko-KR", "ko-KR"},
		{"ko_KR", "ko-KR"},
		{"KO-kr", "ko-KR"},
		{"ko-KP", "ko-KR"},
		{"ko-KR,ko;q=0.9,en;q=0.8", "ko-KR"},
		{"en;q=0.8,ko-KR;q=0.9", "ko-KR"},
		{"ja,ko;q=0.5", "ko-KR"},
		{"ko;q=0", DefaultLocale},
	}

	for _, c := range cases {
		if got := NormalizeLocale(c.header); got != c.want {
			t.Errorf("NormalizeLocale(%q) = %q, want %q", c.header, got, c.want)
		}
	}
}

func TestLocaleFilesShareSameKeys(t *testing.T) {
	locales := SupportedLocales()

	messages := make(map[string]map[string]string, len(locales))
	for _, locale := range locales {
		data, err := localeFS.ReadFile("locales/" + locale + ".json")
		if err != nil {
			t.Fatalf("missing message file for %s: %v", locale, err)
		}
		var parsed map[string]string
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("invalid JSON in %s.json: %v", locale, err)
		}
		if len(parsed) == 0 {
			t.Fatalf("%s.json has no messages", locale)
		}
		messages[locale] = parsed
	}

	reference := locales[0]
	for _, locale := range locales[1:] {
		for key := range messages[reference] {
			if _, ok := messages[locale][key]; !ok {
				t.Errorf("%s.json is missing key %q present in %s.json", locale, key, reference)
			}
		}
		for key := range messages[locale] {
			if _, ok := messages[reference][key]; !ok {
				t.Errorf("%s.json has extra key %q absent from %s.json", locale, key, reference)
			}
		}
	}
}

func TestTranslateResolvesEveryLocale(t *testing.T) {
	cases := []struct {
		locale string
		want   string
	}{
		{"zh-CN", "操作成功"},
		{"en-US", "Operation completed successfully"},
		{"ko-KR", "요청을 처리했습니다"},
		{"ko", "요청을 처리했습니다"},
	}

	for _, c := range cases {
		if got := Translate(c.locale, "api.success"); got != c.want {
			t.Errorf("Translate(%q, api.success) = %q, want %q", c.locale, got, c.want)
		}
	}

	if got := Translate("ko-KR", ""); got != "" {
		t.Errorf("Translate with empty message id = %q, want empty string", got)
	}
	if got := Translate("ko-KR", "api.does.not.exist"); got != "api.does.not.exist" {
		t.Errorf("unknown message id should fall back to the key, got %q", got)
	}
}

func TestErrorLocalizesItsKey(t *testing.T) {
	plain := NewError("api.plugin.already_exists", nil)
	if got := plain.Error(); got != "api.plugin.already_exists" {
		t.Errorf("Error() = %q, want the bare key", got)
	}
	if got := plain.Localize("ko-KR"); got != "이미 등록된 플러그인입니다" {
		t.Errorf("Localize(ko-KR) = %q", got)
	}
	if got := plain.Localize("en-US"); got != "The plugin already exists" {
		t.Errorf("Localize(en-US) = %q", got)
	}
}

func TestErrorFillsTemplateData(t *testing.T) {
	err := NewError("api.plugin.module_invalid", map[string]interface{}{"Module": "NoSuchModule"})

	cases := map[string]string{
		"ko-KR": "올바르지 않은 모듈입니다: NoSuchModule",
		"en-US": "Invalid module: NoSuchModule",
	}
	for locale, want := range cases {
		if got := err.Localize(locale); got != want {
			t.Errorf("Localize(%s) = %q, want %q", locale, got, want)
		}
	}
}

func TestWrapErrorKeepsTheCause(t *testing.T) {
	cause := errors.New("unexpected EOF")
	err := WrapError("api.plugin.zip.read_failed", cause)

	if !errors.Is(err, cause) {
		t.Error("errors.Is should reach the wrapped cause")
	}

	want := "zip 파일을 읽지 못했습니다: unexpected EOF"
	if got := err.Localize("ko-KR"); got != want {
		t.Errorf("Localize(ko-KR) = %q, want %q", got, want)
	}

	var target *Error
	if !errors.As(error(err), &target) {
		t.Error("errors.As should recognise *i18n.Error")
	}
}
