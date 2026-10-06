package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
