package i18n

import (
	"embed"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// DefaultLocale 은 Accept-Language 를 해석할 수 없을 때 사용하는 로케일이다.
const DefaultLocale = "zh-CN"

// localeAliases 는 요청에 담겨 오는 언어 태그를 locales/ 디렉터리의
// 메시지 파일 이름(= 지원 로케일)으로 매핑한다. 키는 모두 소문자이며
// 하이픈 표기(ko-kr)를 쓴다.
var localeAliases = map[string]string{
	"zh":      "zh-CN",
	"zh-cn":   "zh-CN",
	"zh-hans": "zh-CN",
	"zh-sg":   "zh-CN",
	"en":      "en-US",
	"en-us":   "en-US",
	"en-gb":   "en-US",
	"ko":      "ko-KR",
	"ko-kr":   "ko-KR",
}

//go:embed locales/*.json
var localeFS embed.FS

var (
	bundle *i18n.Bundle
	once   sync.Once
)

func init() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.Chinese)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// 加载语言文件
		loadLocaleFiles()
	})
}

func loadLocaleFiles() {
	files, err := localeFS.ReadDir("locales")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := localeFS.ReadFile("locales/" + file.Name())
		if err != nil {
			panic(err)
		}
		bundle.MustParseMessageFileBytes(data, file.Name())
	}
}

// SupportedLocales 는 locales/ 에 메시지 파일이 있는 로케일 목록을 돌려준다.
func SupportedLocales() []string {
	files, err := localeFS.ReadDir("locales")
	if err != nil {
		return []string{DefaultLocale}
	}

	locales := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		locales = append(locales, strings.TrimSuffix(file.Name(), ".json"))
	}
	sort.Strings(locales)
	return locales
}

// NormalizeLocale 은 Accept-Language 헤더 값을 지원 로케일 하나로 환산한다.
// "ko", "ko_KR", "ko-KR,ko;q=0.9,en;q=0.8" 같은 입력을 모두 "ko-KR" 로 맞추고,
// 지원하지 않는 언어이거나 비어 있으면 DefaultLocale 을 돌려준다.
func NormalizeLocale(header string) string {
	type weighted struct {
		tag   string
		q     float64
		order int
	}

	var candidates []weighted
	for i, part := range strings.Split(header, ",") {
		tag := part
		q := 1.0
		if idx := strings.Index(tag, ";"); idx >= 0 {
			params := tag[idx+1:]
			tag = tag[:idx]
			for _, p := range strings.Split(params, ";") {
				p = strings.TrimSpace(p)
				if !strings.HasPrefix(strings.ToLower(p), "q=") {
					continue
				}
				if v, err := strconv.ParseFloat(strings.TrimSpace(p[2:]), 64); err == nil {
					q = v
				}
			}
		}
		tag = strings.TrimSpace(tag)
		if tag == "" || tag == "*" || q <= 0 {
			continue
		}
		candidates = append(candidates, weighted{tag: tag, q: q, order: i})
	}

	// q 값이 큰 순서로, 같으면 헤더에 등장한 순서를 유지한다.
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].q != candidates[j].q {
			return candidates[i].q > candidates[j].q
		}
		return candidates[i].order < candidates[j].order
	})

	for _, c := range candidates {
		tag := strings.ToLower(strings.ReplaceAll(c.tag, "_", "-"))
		if locale, ok := localeAliases[tag]; ok {
			return locale
		}
		// 등록되지 않은 지역 변형(예: ko-KP)은 언어 코드로 다시 맞춘다.
		if base := strings.SplitN(tag, "-", 2)[0]; base != "" {
			if locale, ok := localeAliases[base]; ok {
				return locale
			}
		}
	}
	return DefaultLocale
}

func Translate(locale, messageID string) string {
	if messageID == "" {
		return ""
	}
	localizer := i18n.NewLocalizer(bundle, NormalizeLocale(locale))
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID // 如果翻译失败，返回原始key
	}
	return msg
}
