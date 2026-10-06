package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/i18n"
	"github.com/gin-gonic/gin"
)

func newContextWithAcceptLanguage(header string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if header != "" {
		c.Request.Header.Set("Accept-Language", header)
	}
	return c, recorder
}

func TestGetLocaleNormalizesAcceptLanguage(t *testing.T) {
	cases := []struct {
		header string
		want   string
	}{
		{"", "zh-CN"},
		{"ko", "ko-KR"},
		{"ko-KR,ko;q=0.9,en;q=0.8", "ko-KR"},
		{"en", "en-US"},
		{"fr", "zh-CN"},
	}

	for _, c := range cases {
		ctx, _ := newContextWithAcceptLanguage(c.header)
		if got := getLocale(ctx); got != c.want {
			t.Errorf("getLocale with Accept-Language %q = %q, want %q", c.header, got, c.want)
		}
	}
}

func TestSuccessTranslatesMessageForRequestLocale(t *testing.T) {
	cases := []struct {
		header string
		want   string
	}{
		{"ko-KR", "요청을 처리했습니다"},
		{"en-US", "Operation completed successfully"},
		{"", "操作成功"},
	}

	for _, c := range cases {
		ctx, recorder := newContextWithAcceptLanguage(c.header)
		Success(ctx, nil, "api.success")

		var body Response
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatalf("unexpected response body for %q: %v", c.header, err)
		}
		if body.Message != c.want {
			t.Errorf("Success message for %q = %q, want %q", c.header, body.Message, c.want)
		}
	}
}

func TestSuccessKeepsEmptyMessageEmpty(t *testing.T) {
	ctx, recorder := newContextWithAcceptLanguage("ko-KR")
	Success(ctx, nil, "")

	var body Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected response body: %v", err)
	}
	if body.Message != "" {
		t.Errorf("message = %q, want empty string", body.Message)
	}
}

func TestErrorDetailLocalizesI18nErrors(t *testing.T) {
	cases := []struct {
		header string
		want   string
	}{
		{"ko-KR", "올바르지 않은 모듈입니다: NoSuchModule"},
		{"en-US", "Invalid module: NoSuchModule"},
	}

	for _, c := range cases {
		ctx, recorder := newContextWithAcceptLanguage(c.header)
		err := i18n.NewError("api.plugin.module_invalid", map[string]interface{}{"Module": "NoSuchModule"})
		InternalServerError(ctx, "api.plugin.key.error", err)

		var body Response
		if jsonErr := json.Unmarshal(recorder.Body.Bytes(), &body); jsonErr != nil {
			t.Fatalf("unexpected response body: %v", jsonErr)
		}
		if body.Data != c.want {
			t.Errorf("data for %q = %v, want %q", c.header, body.Data, c.want)
		}
	}
}

func TestErrorDetailPassesThroughPlainErrors(t *testing.T) {
	ctx, recorder := newContextWithAcceptLanguage("ko-KR")
	BadRequest(ctx, "api.bad_request", errors.New("ids is required"))

	var body Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected response body: %v", err)
	}
	if body.Data != "ids is required" {
		t.Errorf("data = %v, want the plain error text", body.Data)
	}
}

// 에러 없이 호출하면 data 는 비어야 한다. 예전에는 fmt.Sprintf("%v", nil) 때문에
// "<nil>" 이 들어가 화면의 토스트에 그대로 붙었다.
func TestErrorDetailIsEmptyWithoutAnError(t *testing.T) {
	for _, call := range []struct {
		name string
		fn   func(*gin.Context)
	}{
		{"Unauthorized", func(c *gin.Context) { Unauthorized(c, "api.unauthorized.header_required", nil) }},
		{"BadRequest", func(c *gin.Context) { BadRequest(c, "api.bad_request", nil) }},
		{"InternalServerError", func(c *gin.Context) { InternalServerError(c, "api.error", nil) }},
	} {
		ctx, recorder := newContextWithAcceptLanguage("ko-KR")
		call.fn(ctx)

		var body Response
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: unexpected response body: %v", call.name, err)
		}
		if body.Data != nil && body.Data != "" {
			t.Errorf("%s: data = %v, want empty", call.name, body.Data)
		}
	}
}
