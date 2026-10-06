// Package region 은 배포 지역 프로필을 다룬다.
//
// ScopeSentry 는 중국 환경을 전제로 만들어져 있다. 자산 필드 ICP, 미니프로그램
// 자산 타입, 중국 주민번호·휴대전화 민감정보 규칙이 그렇다. 지역 프로필은 그중
// 어떤 조합을 쓸지 고르는 설정이며, 언어 설정(locale)과는 별개다. 영어 화면에서
// 국내 규칙을 쓰거나 한국어 화면에서 중국 규칙을 쓰는 조합도 가능해야 한다.
package region

import "strings"

// Region 은 지역 프로필이다.
type Region string

const (
	// CN 은 중국 환경이다. 업스트림이 전제하던 기본값이라 이 값을 default 로 둔다.
	CN Region = "CN"
	// KR 은 국내 환경이다.
	KR Region = "KR"
	// Global 은 어느 쪽 지역 자료도 쓰지 않는 조합이다.
	Global Region = "GLOBAL"
)

// Default 는 설정이 없을 때 쓰는 프로필이다. 기존 설치본의 동작을 그대로
// 유지하려고 CN 으로 둔다.
const Default = CN

// All 은 고를 수 있는 프로필을 설정 화면에 보여줄 순서대로 돌려준다.
func All() []Region {
	return []Region{CN, KR, Global}
}

// Parse 는 저장된 값을 Region 으로 읽는다. 대소문자는 가리지 않으며,
// 비어 있거나 모르는 값이면 Default 를 돌려준다.
func Parse(value string) Region {
	switch Region(strings.ToUpper(strings.TrimSpace(value))) {
	case CN:
		return CN
	case KR:
		return KR
	case Global:
		return Global
	default:
		return Default
	}
}

// String 으로 저장·전송에 쓰는 표기를 얻는다.
func (r Region) String() string {
	return string(r)
}
