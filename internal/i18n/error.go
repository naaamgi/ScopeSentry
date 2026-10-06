package i18n

import "fmt"

// Error 는 번역할 메시지 키를 들고 다니는 에러다.
//
// 서비스 계층은 요청의 Accept-Language 를 알지 못하므로 문장을 만들지 않고 키만
// 돌려주고, 번역은 응답을 쓰는 지점에서 한다. 평범한 error 로도 동작하므로
// 호출부가 errors.Is / errors.As 로 원인을 계속 따라갈 수 있다.
type Error struct {
	// Key 는 locales/*.json 의 메시지 키다.
	Key string
	// Data 는 메시지 안의 Go 템플릿 자리({{.Module}} 등)를 채우는 값이다.
	Data map[string]interface{}
	// Err 는 감싼 원인이며, 없을 수도 있다.
	Err error
}

// NewError 는 원인 없이 메시지 키만 담은 에러를 만든다.
func NewError(key string, data map[string]interface{}) *Error {
	return &Error{Key: key, Data: data}
}

// WrapError 는 원인을 감싸면서 메시지 키를 붙인다.
func WrapError(key string, cause error) *Error {
	return &Error{Key: key, Err: cause}
}

// Error 는 번역하지 않은 표현을 돌려준다. 로그와 테스트에서 키를 그대로 보는 쪽이
// 추적하기 쉬우므로 일부러 키를 남긴다.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Key, e.Err)
	}
	return e.Key
}

// Unwrap 으로 errors.Is / errors.As 가 원인까지 내려갈 수 있게 한다.
func (e *Error) Unwrap() error {
	return e.Err
}

// Localize 는 주어진 로케일로 메시지를 만든다. 감싼 원인이 있으면 뒤에 덧붙이는데,
// 원인은 보통 라이브러리가 만든 영문 문장이라 번역하지 않는다.
func (e *Error) Localize(locale string) string {
	msg := TranslateWithData(locale, e.Key, e.Data)
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return msg
}
