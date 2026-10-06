package i18n

import "strings"

// BuiltinPluginPrefix 는 기본 제공 플러그인의 설명문이 DB에 키로 저장될 때 쓰는
// 접두사다. 사용자가 가져온 플러그인은 설명을 자유 문장으로 쓰므로, 저장된 값이
// 번역할 키인지 그대로 보여줄 문장인지 이 접두사로 구분한다.
const BuiltinPluginPrefix = "plugin.builtin."

// TranslateStored 는 저장된 값이 메시지 키일 때만 번역하고, 그 밖에는 값을
// 그대로 돌려준다. 기본 제공 플러그인과 사용자 플러그인이 같은 필드를 공유하는
// 곳에서 쓴다.
func TranslateStored(locale, value string) string {
	if !strings.HasPrefix(value, BuiltinPluginPrefix) {
		return value
	}
	return Translate(locale, value)
}
