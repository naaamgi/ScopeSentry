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

이 저장소를 직접 빌드한 버전은 처음 웹 화면에 접속하면 관리자 계정을 직접 만듭니다. 이후 그 계정으로 로그인해 플러그인도 관리합니다. 기존 설치의 계정은 그대로 유지됩니다. 위 `single-host-deployment.yml`의 업스트림 이미지는 이 변경이 적용되지 않으므로 아래 로컬 빌드 방법을 사용하세요.

새 설치에서는 첫 계정 생성이 완료될 때까지 웹 화면을 신뢰할 수 있는 네트워크에만 노출하세요.

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

빌드가 끝나면 http://localhost:8082 에서 `admin` / `admin`으로 로그인합니다. 이 계정은 `local-test.yml`에서만 생성하며 기존 계정이나 이미 있는 `admin` 계정의 비밀번호는 바꾸지 않습니다. 포트 8082는 로컬 컴퓨터에만 열립니다. 플러그인에는 별도 키가 필요하지 않습니다.

`dockerfile.local` 이 프런트엔드와 서버를 컨테이너 안에서 빌드하므로 로컬에 node 나 go 를 설치하지 않아도 됩니다. 처음 빌드는 의존성을 받아오느라 몇 분 걸립니다.

컨테이너는 네 개가 뜹니다.

| 컨테이너 | 역할 |
| --- | --- |
| `scopesentry-mongodb-local` | 자산·작업·설정 저장소 |
| `scopesentry-redis-local` | 노드와 서버 사이의 작업 큐 |
| `scope-sentry-local` | 서버와 웹 화면 (이 브랜치를 빌드한 것) |
| `scopesentry-scan-local` | 스캔 노드. 업스트림 이미지를 그대로 씁니다 |

스캔 노드가 붙었는지는 **노드 관리** 화면에서 `node-local` 이 실행 중으로 보이는지로 확인합니다. 보이지 않으면 `docker logs scopesentry-scan-local` 을 보세요.

정리할 때는 이렇게 합니다. `data-local/` 에 MongoDB 데이터가 남으므로, 초기 설치 상태부터 다시 보려면 함께 지우세요.

```
docker compose -f local-test.yml down
rm -rf data-local
```

### 테스트 도메인으로 스캔해 보기

스캔은 **작업(Task)** 단위로 돌아갑니다. 작업에는 **대상**과 **스캔 템플릿**(어떤 모듈을 어떤 순서로 돌릴지)이 필요합니다. 설치 직후 기본 템플릿이 하나 들어 있으므로 그대로 쓰면 됩니다.

스캔 대상은 본인이 권한을 가진 도메인이어야 합니다. 남의 도메인을 스캔하면 안 됩니다. 연습용으로는 `scanme.nmap.org`(Nmap 프로젝트가 스캔을 허용해 둔 호스트)나 본인이 띄운 테스트 서버를 쓰세요.

1. **노드 확인** — 노드 관리에서 `node-local` 이 실행 중인지 봅니다. 노드가 없으면 작업이 생성만 되고 아무 일도 일어나지 않습니다.
2. **작업 관리 → 스캔 작업 → 새 작업**
3. 입력
   - **작업 이름**: 아무 이름 (예: `첫 테스트`)
   - **대상**: 한 줄에 하나. 예: `scanme.nmap.org`
   - **노드 선택**: `node-local`
   - **스캔 템플릿**: 기본으로 들어 있는 템플릿
4. **제출** → 목록에서 진행률이 올라가는지 봅니다.
5. 결과는 **자산 정보** 화면의 각 탭(자산 / 서브도메인 / URL / 취약점 …)에 쌓입니다.

처음 돌릴 때 알아둘 점:

- **첫 스캔은 느립니다.** 스캔 노드가 subfinder, httpx, nuclei 같은 도구와 POC 템플릿을 내려받습니다. 진행률이 한동안 0 이어도 정상입니다. `docker logs -f scopesentry-scan-local` 로 지켜보세요.
- **외부 네트워크가 필요합니다.** 서브도메인 수집과 POC 다운로드가 인터넷을 씁니다. 사내망이라면 막힐 수 있습니다.
- **루트 도메인을 넣으면 전체 파이프라인을 한 번에 돌리지 마세요.** 분산 작업이 대상 하나를 단위로 나뉘기 때문에, 루트 도메인 하나를 받은 노드가 거기서 찾은 서브도메인까지 혼자 처리합니다. 2단계로 나누는 방법은 `SKILL.md` 의 "루트 도메인 전체 정보 수집" 항목에 있습니다.
- 노드가 하나뿐이니 동시 실행 수를 올려도 한계가 있습니다. **설정 → 시스템 설정**의 Module Config 에서 조절합니다.

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

## 화면에 남아 있는 중국어

화면 문자열과 서버 응답 메시지는 한국어로 번역되어 있지만, **원본 데이터**에는 중국어가 남아 있습니다. 제품 식별자와 탐지 패턴은 원본을 유지해야 합니다. 별도의 한국어 표시명을 추가하는 방식은 가능합니다.

### 웹 핑거프린트

핑거프린트 관리 화면 목록과 자산의 "애플리케이션·구성요소" 컬럼에 중국어 제품명이 보일 수 있습니다. 건수는 외부 규칙 업데이트에 따라 달라집니다.

번역하지 않는 이유:

- **이름이 식별자입니다.** 핑거프린트 이름은 자산의 `technologies` 필드에 그대로 저장되고, 검색에서 `app="..."` 로 찾습니다. 이름을 바꾸면 기존 스캔 결과와 매칭이 끊어집니다.
- **규칙 안의 중국어는 탐지 패턴입니다.** 중국 제품의 페이지 본문에 있는 중국어 문자열을 찾는 규칙이라, 번역하면 그 제품을 못 찾습니다.
- **원격 ID의 수정은 덮어써질 수 있습니다.** 핑거프린트 업데이트는 `api.scope-sentry.top`에서 변경된 YAML을 가져와 같은 `fingerprint_id`에 덮어씁니다.

중국 제품을 쓰지 않는다고 규칙을 지울 이유도 없습니다. 탐지 범위만 줄어듭니다.

### 플러그인 마켓

플러그인 마켓 화면은 `api.scope-sentry.top` 의 목록을 읽고, 사용자가 설치·업데이트를 누를 때 원격 소스를 가져옵니다. 중국어 이름과 소개는 **원격 제공자가 작성한 메타데이터**입니다. 현재 화면은 이를 자동 번역하지 않습니다. 마켓을 열지 않아도 로컬에 저장된 플러그인과 스캔은 사용할 수 있습니다.

#### 독립적으로 추가·수정하는 방법

1. **플러그인 관리 → 스캐너 플러그인(또는 서버 플러그인) → 새 플러그인**에서 이름, 모듈, 버전, 한국어 소개와 도움말을 입력하고 소스 코드 탭에 코드를 넣습니다. 기존 플러그인은 `수정`에서 버전과 코드를 바꿉니다. 파일을 받았다면 `가져오기`를 사용합니다.
2. 스캐너 플러그인을 저장·수정하면 서버가 전체 노드에 `install_plugin` 메시지를 보냅니다. **노드 관리 → 플러그인**에서 설치·검사 상태를 확인하고, 실패하면 해당 플러그인과 노드 로그를 확인합니다. 서버 플러그인은 서버 프로세스에 등록됩니다.
3. 새 플러그인 작성에는 저장소의 `plugin-template/` 와 [업스트림 플러그인 템플릿](https://github.com/Autumn-27/ScopeSentry-Plugin-Template)을 참고합니다. 스캐너 플러그인은 별도 [ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan)에서 실행되므로, 스캐너의 모듈 계약과 의존 도구도 확인해야 합니다.
4. 마켓에 등록하지 않은 자체 플러그인은 이 설치의 DB에서 관리합니다. **마켓 자동 업데이트 대상이 아닙니다.** 버전과 코드를 직접 수정하거나 자신의 배포 파일을 가져와 갱신합니다.

마켓의 업데이트 표시는 원격과 로컬의 **버전 문자열이 다른지** 비교한 결과입니다. 실제 업데이트는 버튼을 눌러야 수행됩니다. 마켓 자체를 교체하려면 `frontend/src/api/plugins/index.ts`의 목록·내보내기 API 주소와 `frontend/src/views/Plugins/components/PluginMarket.vue`의 상세 링크를 자신의 서비스로 바꿔야 합니다. 대량 배포·서명 검증·호환성 검사는 별도 설계가 필요합니다. 외부 서비스 없이 운영하려면 이 마켓 연결을 사용하지 않고 자체 플러그인 등록·가져오기 기능을 쓰면 됩니다.

### 핑거프린트 업데이트의 실제 동작

`업데이트` 숫자는 로컬 `FingerVersion` 이후의 원격 변경 건수입니다. 버튼을 누르면 `api.scope-sentry.top`에서 YAML을 받아 `fingerprint_id`로 **upsert**(기존 ID는 덮어쓰기, 없으면 추가)하고 노드에 새로고침을 알립니다. 성공한 항목 중 가장 최신 시각을 `FingerVersion`으로 저장합니다. 따라서 원격 ID가 붙은 규칙을 직접 수정한 경우 다음 업데이트로 수정 내용이 덮어써질 수 있습니다. 자체 규칙은 `새로 만들기`에서 별도 ID로 추가해 관리하세요.

### 스캔 노드가 남기는 로그

노드 로그와 플러그인 로그에 중국어가 섞입니다. 스캐너는 [ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan) 이라는 별도 저장소이고, 이 포크는 서버와 화면만 바꿨습니다.

### 소스 코드 주석

Go 와 Vue 소스의 주석 약 2,800줄은 중국어 그대로입니다. 동작에 영향이 없고, 업스트림 변경을 가져올 때 충돌 면적만 키우기 때문에 손대지 않았습니다.

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

프로젝트와 작업의 대상 입력란에는 도메인·IP 외에 접두사를 붙인 대상도 넣을 수 있습니다. 접두사 목록은 스캐너([ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan))의 `targetparser` 가 아는 것과 같아야 합니다.

| 접두사 | 뜻 | 예시 | 기본 설치에서 동작하는가 |
| --- | --- | --- | --- |
| `CIDR:` | 네트워크 대역을 한 노드에서 스캔 | `CIDR:192.168.0.0/18` | 동작함 |
| `APP:` | 앱 이름 | `APP:예시앱` | 받는 플러그인 없음 |
| `APP-ID:` | 앱 패키지명 | `APP-ID:com.example.app` | 받는 플러그인 없음 |
| `CMP:` | 회사명 | `CMP:예시주식회사` | 받는 플러그인 없음 |
| `ICP:` | 중국 ICP 등록번호 | `ICP:京ICP证1234号` | 받는 플러그인 없음 |

**`CMP:` 와 `ICP:` 는 기본 설치에서 아무 일도 하지 않습니다.** 스캐너의 `targetparser` 가 이 대상을 `types.Company` / `types.ICP` 로 바꿔 흘려보내는데, `assetmapping` 모듈의 `switch` 에 해당 `case` 가 빈 몸통으로 들어 있고 기본 플러그인(`httpx`)은 이 타입을 받지 않습니다. "등록번호 → 회사 → 보유 도메인" 경로를 쓰려면 플러그인 마켓에서 추가 플러그인(ENScan 계열)을 설치해야 합니다.

**국내 사업자등록번호용 접두사는 없습니다.** `BRN:` 을 추가했다가 되돌렸습니다. 스캐너의 `targetparser` 는 모르는 접두사를 `host:port` 처리 분기로 떨어뜨려 콜론으로 쪼개기 때문에, `BRN:123-45-67890` 은 호스트 `BRN` + 포트 `123-45-67890` 짜리 쓰레기 자산이 됩니다. 접두사를 늘리려면 스캐너 쪽 지원이 먼저 필요합니다.

자산 검색에서는 `icp` 대신 `brn` 을 써도 같은 필드를 찾습니다. 예: `brn="123-45-67890"`. 이쪽은 조회 시점의 별칭일 뿐이라 안전합니다.

### 업스트림과의 관계

이 저장소의 `frontend/` 는 [ScopeSentry-UI](https://github.com/Autumn-27/ScopeSentry-UI) 를 **통째로 복사해 둔 것**입니다. 두 쪽을 비교해 보면 이 브랜치의 변경을 빼고 실제로 다른 파일은 세 개(`src/router/index.ts`, `src/views/Asset/components/assetInfo.vue`, `vite.config.ts`)뿐입니다.

그래서 화면 쪽 변경은 **업스트림이 UI 저장소를 다시 복사해 오면 머지되지 않고 덮여 사라집니다.** `ko.ts` 를 포함해 이 브랜치가 건드린 `frontend/` 파일 36개가 그렇습니다. 업스트림을 따라갈 생각이라면:

- 화면 번역은 `ScopeSentry-UI` 쪽에 올리는 것이 맞습니다.
- 여기서 유지한다면 업스트림을 가져올 때마다 `frontend/` 가 덮였는지 확인하고 다시 얹어야 합니다. `git log --oneline -- frontend/` 로 복사 시점을 찾을 수 있습니다.

서버 쪽(`internal/`) 변경은 이 저장소가 원본이라 평범하게 머지됩니다. 스캐너는 [ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan) 으로 분리되어 있고 이 브랜치는 건드리지 않았습니다.

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
