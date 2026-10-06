package middleware

import (
	"github.com/Autumn-27/ScopeSentry/internal/i18n"
	"github.com/gin-gonic/gin"
)

// I18nMiddleware 国际化中间件
//
// 응답 본문의 메시지를 어떤 로케일로 만들었는지 Content-Language 로 알린다.
// 본문을 만드는 쪽(response 패키지)과 같은 환산 규칙을 써야 하므로
// i18n.NormalizeLocale 을 공유한다. 예전에는 language.Parse 로 따로 해석해서,
// "ko-KR,ko;q=0.9,en;q=0.8" 처럼 목록으로 오는 헤더를 파싱에 실패하고
// Content-Language: zh 를 내보내면서 본문과 어긋났다.
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := i18n.NormalizeLocale(c.GetHeader("Accept-Language"))

		c.Set("language", locale)
		c.Header("Content-Language", locale)

		c.Next()
	}
}
