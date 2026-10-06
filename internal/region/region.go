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
	// KR 은 국내 환경이다.
	KR Region = "KR"
	// CN 은 중국 환경이다. 업스트림이 전제하던 조합이다.
	CN Region = "CN"
	// Global 은 어느 쪽 지역 자료도 쓰지 않는 조합이다.
	Global Region = "GLOBAL"
)

// Default 는 설정이 없을 때 쓰는 프로필이다. 이 포크는 국내 환경을 기준으로
// 하므로 KR 이다. 업스트림 그대로의 동작은 CN 프로필이다.
//
// 이미 쓰던 설치본은 이 값을 따라가지 않는다. Update22 가 저장된 값을 읽어
// 그대로 쓰고, 설정이 아예 없을 때만 이 기본값을 넣는다.
const Default = KR

// All 은 고를 수 있는 프로필을 설정 화면에 보여줄 순서대로 돌려준다.
func All() []Region {
	return []Region{KR, CN, Global}
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
