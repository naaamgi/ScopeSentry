<div align=center>
	<img src="docs/images/favicon.ico"/>
</div>

한국어 | [English](./README.md) | [中文](./README_CN.md)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Autumn-27/ScopeSentry-Scan)

## 소개

Scope Sentry는 자산 매핑, 서브도메인 열거, 정보 노출 탐지, 취약점 스캔, 디렉터리 스캔, 서브도메인 테이크오버 점검, 크롤러, 페이지 모니터링 기능을 갖춘 도구입니다. 노드를 여러 개 띄워 두고 스캔 작업을 실행할 노드를 골라 쓸 수 있습니다. 새 취약점이 나왔을 때 관심 자산에 해당 구성요소가 있는지 빠르게 확인하는 데 쓸 수 있습니다.

분산 구성 참고 글: [https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA](https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA)

## Discord

[https://discord.gg/GWVwSBBm48](https://discord.gg/GWVwSBBm48)

## 사용 언어

서버: go - Gin

스캐너: go

프런트엔드: vue - vue-element-plus-admin

> 업스트림 README는 서버를 `python - FastApi`로 적고 있습니다. 이 저장소의 서버는 Go로 다시 작성된 쪽입니다.

## 링크

- 공식 사이트: [https://www.scope-sentry.top](https://www.scope-sentry.top/en/)
- GitHub: [https://github.com/Autumn-27/ScopeSentry](https://github.com/Autumn-27/ScopeSentry)
- 스캐너 소스: [https://github.com/Autumn-27/ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan)
- UI 소스: [https://github.com/Autumn-27/ScopeSentry-UI](https://github.com/Autumn-27/ScopeSentry-UI)
- 플러그인 마켓: [Plugin Market](https://plugin.scope-sentry.top/en)
- 플러그인 템플릿: [https://github.com/Autumn-27/ScopeSentry-Plugin-Template](https://github.com/Autumn-27/ScopeSentry-Plugin-Template)

## 설치

```
git clone https://github.com/Autumn-27/ScopeSentry.git
cd ScopeSentry
# .env 파일에서 MongoDB와 Redis 계정 비밀번호를 바꾸세요.
docker-compose -f single-host-deployment.yml up -d
```

실행하면 mongodb, redis, scope-sentry(서버), scopesentry-scan(스캐너) 네 개의 컨테이너가 뜹니다. 기본으로 스캔 노드 하나가 함께 올라옵니다.

초기 사용자 비밀번호와 플러그인 2차 확인 비밀번호는 아래 명령으로 확인합니다.

```
docker logs scope-sentry
```

**노드 추가(선택)**

```
git clone https://github.com/Autumn-27/ScopeSentry-Scan.git
cd ScopeSentry-Scan/build
# .env 파일에서 MongoDB와 Redis 접속 정보를 수정하세요.
# NodeName은 노드 이름이며 노드마다 달라야 합니다(비워 두면 무작위로 생기고, 웹 화면에서 바꿀 수 있습니다).
docker-compose -f scan-docker-compose.yml up -d
```

## 이 브랜치를 로컬에서 띄우기

`single-host-deployment.yml` 은 Docker Hub 의 `autumn27/scopesentry:latest` 를 받아 씁니다. 그 이미지에는 이 브랜치의 변경이 들어 있지 않으므로, 직접 빌드해야 합니다.

```
git clone -b claude/awesome-fermi-x8dhsq https://github.com/naaamgi/ScopeSentry.git
cd ScopeSentry
docker compose -f local-test.yml up --build
```

빌드가 끝나면 http://localhost:8082 로 접속합니다. 초기 사용자 비밀번호와 플러그인 키는 아래로 확인합니다.

```
docker logs scope-sentry-local
```

`dockerfile.local` 이 프런트엔드와 서버를 컨테이너 안에서 빌드하므로 로컬에 node 나 go 를 설치하지 않아도 됩니다. 처음 빌드는 의존성을 받아오느라 몇 분 걸립니다.

정리할 때는 이렇게 합니다. `data-local/` 에 MongoDB 데이터가 남으므로, 초기 설치 상태부터 다시 보려면 함께 지우세요.

```
docker compose -f local-test.yml down
rm -rf data-local
```

### 소스에서 직접 빌드할 때 주의할 점

`cmd/main/static` 에 커밋되어 있는 프런트엔드 빌드 산출물은 릴리스 CI 가 매번 새로 덮어쓰는 스냅샷이라 소스보다 **뒤처져 있습니다**. `go build` 만 하면 예전 화면이 바이너리에 들어가고, 한국어와 지역 프로필이 보이지 않습니다. 직접 빌드할 때는 CI 와 같은 순서를 따르세요.

```
cd frontend && npm ci && npx vite build --mode pro && cd ..
cp -a frontend/dist-pro/. cmd/main/static/
go build -o ScopeSentry ./cmd/main/main.go
```

`cmd/main` 패키지에는 `main` 함수가 두 개(`main.go`, `basic_usage.go`) 있어 `go build ./cmd/main` 은 실패합니다. goreleaser 와 같이 `main.go` 를 지정해야 합니다.

## 언어 설정

오른쪽 위 언어 아이콘에서 **한국어 / English / 简体中文** 중에 고릅니다. 고른 값은 브라우저에 저장되고, 화면 문자열과 Element Plus 컴포넌트, 그리고 서버가 돌려주는 응답 메시지까지 같이 바뀝니다. 서버에는 요청마다 `Accept-Language` 헤더로 전달됩니다.

언어를 바꾸면 화면이 다시 읽힙니다.

## 지역 프로필

ScopeSentry는 원래 중국 환경을 전제로 만들어져 있습니다. 자산에 ICP 등록번호 필드가 있고, 미니프로그램(微信 小程序)이 자산 타입으로 들어 있으며, 기본 민감정보 규칙에 중국 주민번호·휴대전화·은행카드 패턴이 포함되어 있습니다.

**설정 → 시스템 설정 → 지역 프로필**에서 어느 지역 자료를 쓸지 고릅니다. 언어 설정과는 별개입니다. 영어 화면에서 국내 규칙을 쓰거나, 한국어 화면에서 중국 규칙을 쓰는 조합도 됩니다.

| 프로필 | 등록번호 필드 라벨 | 미니프로그램 탭 | 기본 활성 민감정보 규칙 |
| --- | --- | --- | --- |
| `KR` (기본값) | 사업자등록번호 | 숨김 | 주민등록번호, 외국인등록번호, 휴대전화, 사업자등록번호, 운전면허번호, Toss Payments 시크릿 키, Kakao REST API 키 |
| `CN` | ICP | 보임 | 중국 주민번호, 중국 휴대전화 |
| `GLOBAL` | 등록번호 | 숨김 | 어느 지역 규칙도 켜지 않음 |

이 포크는 국내 환경을 기준으로 하므로 기본값이 `KR`입니다. 업스트림 그대로의 동작을 원하면 `CN`을 고르세요. 프로필을 바꾸면 그 지역 규칙이 켜지고 다른 지역 규칙은 꺼지며, 화면이 다시 읽힙니다.

**새로 설치하는 경우** 설치 시점부터 `KR` 프로필과 국내 규칙이 함께 들어갑니다.

**이미 쓰던 설치본을 업그레이드하는 경우** 저장된 프로필 값이 있으면 그 값을 그대로 씁니다. 값이 아예 없었던 설치본(이 기능 이전 버전)만 `KR`이 들어가고, 그때 규칙 상태도 `KR`에 맞춰집니다. 즉 중국 환경으로 쓰던 설치본은 업그레이드 후 중국 PII 규칙이 꺼지므로, 그대로 쓰려면 프로필을 `CN`으로 되돌리세요.

개별 규칙을 손으로 켜 둔 것은 프로필을 다시 저장해도 남습니다. 규칙 상태를 다시 맞추는 것은 프로필이 실제로 **바뀌는** 순간에만 일어납니다.

### 저장되는 값에 대한 주의

국내 민감정보 규칙은 기존 172개 규칙과 같은 방식으로 **매칭된 값을 평문으로 저장**합니다. `KR` 프로필로 스캔하면 MongoDB의 `SensitiveResult` 컬렉션과 CSV 내보내기 결과에 실제 주민등록번호가 남을 수 있습니다. 해당 컬렉션과 내보낸 파일의 접근 통제를 함께 챙기세요.

여권번호 규칙은 기본으로 꺼져 있습니다. 영문자 한 자 뒤에 숫자 여덟 자라는 형태가 다른 문자열과 너무 많이 겹치기 때문입니다. 사업자등록번호 규칙은 같은 이유로 하이픈이 들어간 형태(`123-45-67890`)만 찾습니다.

### 스캔 대상 문법

프로젝트와 작업의 대상 입력란에는 도메인·IP 외에 접두사를 붙인 대상도 넣을 수 있습니다.

| 접두사 | 뜻 | 예시 |
| --- | --- | --- |
| `CMP:` | 회사명 | `CMP:예시주식회사` |
| `ICP:` | 중국 ICP 등록번호 | `ICP:京ICP证1234号` |
| `BRN:` | 국내 사업자등록번호 | `BRN:123-45-67890` |
| `APP:` | 앱 이름 | `APP:예시앱` |
| `APP-ID:` | 앱 패키지명 | `APP-ID:com.example.app` |
| `CIDR:` | 네트워크 대역을 한 노드에서 스캔 | `CIDR:192.168.0.0/18` |

접두사는 지역 프로필과 무관하게 모두 받습니다. 대상 입력란의 안내문만 프로필에 맞는 것을 보여줍니다.

자산 검색에서는 `icp` 대신 `brn`을 써도 같은 필드를 찾습니다. 예: `brn="123-45-67890"`.

> **중요**: 서버는 이 접두사를 파싱해서 보관할 뿐이고, "등록번호 → 회사 → 보유 도메인" 같은 실제 수집은 스캐너 플러그인이 합니다. `BRN:` 대상을 넣어도 국내 등록 자료를 조회하는 플러그인이 없으면 자산은 늘지 않습니다. 중국 쪽 경로(ICP 등록 자료 조회)는 [ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan)의 플러그인이 담당합니다.

## 플러그인 흐름도

<img src="流程图.svg"/>

## 현재 기능

- 플러그인 시스템(확장으로 어떤 도구든 추가)
- 서브도메인 열거
- 서브도메인 테이크오버 탐지
- 포트 스캔
- 자산 식별
- 디렉터리 스캔
- 취약점 스캔
- 민감정보 노출 탐지
- URL 추출
- 크롤러
- 페이지 모니터링
- 사용자 지정 웹 핑거프린트
- POC 가져오기
- 자산 그룹
- 다중 노드 스캔
- Webhook

## 예정

- 취약 비밀번호 점검

## 설치 안내

자세한 설치 방법은 [공식 사이트](https://www.scope-sentry.top)를 보세요.

## 라이선스

이 프로젝트의 모든 브랜치는 AGPL-3.0을 따르며, 다음 추가 조건이 함께 적용됩니다.

1. 이 소프트웨어를 상업적으로 쓰려면 별도의 상업용 라이선스가 필요합니다.
2. 회사, 조직, 영리 단체는 이 소프트웨어를 사용·배포·수정하기 전에 상업용 라이선스를 받아야 합니다.
3. 개인과 비영리 단체는 AGPL-3.0 조건에 따라 자유롭게 쓸 수 있습니다.
4. 상업용 라이선스 문의는 rainy-autumn@outlook.com 으로 보내 주세요.
