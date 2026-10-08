export interface GuideSection { heading: string; items: string[]; example?: string }
export interface Guide { title: string; summary: string; sections: GuideSection[] }

export const guides: Record<string, Guide> = {
  '/dashboard': {
    title: '홈', summary: '마지막으로 수집된 자산과 작업 현황을 요약합니다.',
    sections: [{ heading: '처음 확인할 순서', items: ['노드 관리에서 스캐너가 실행 중인지 확인합니다.', '스캔 작업에서 완료·진행·실패 상태를 확인합니다.', '자산 정보에서 결과가 어느 탭에 기록됐는지 확인합니다.'] }]
  },
  '/asset-information': {
    title: '자산 정보', summary: '스캔 결과를 자산 종류별로 조회합니다.',
    sections: [
      { heading: '조회', items: ['상단 탭에서 자산·IP·서브도메인·URL·취약점 등 결과 종류를 고릅니다.', '검색과 좌측 필터로 범위를 좁히고 상세에서 응답 내용과 태그를 확인합니다.'] },
      { heading: '관계도', items: ['도메인 → 호스트 → IP → 포트·서비스 관계를 표시합니다. 같은 DNS 이름은 하나의 노드로 묶습니다.', 'top-down에서는 IP 아래 서비스가 한 묶음으로 표시됩니다. 묶음을 클릭하면 개별 포트·서비스를 펼치거나 접습니다. 자유 배치에서는 노드를 끌어 옮길 수 있습니다.', '한 번에 최대 500건의 자산을 불러옵니다. 검색은 불러온 자산에만 적용됩니다. 휠로 확대·축소하고 빈 곳을 끌어 이동합니다. 배치 선택은 이 브라우저에 저장됩니다.'] },
      { heading: '민감정보와 유출 계정', items: ['민감정보 탭은 스캔한 페이지에서 규칙과 일치한 비밀값·개인정보 후보를 보여 줍니다.', '외부 유출 사고에서 수집된 계정 목록을 조회하는 기능이나 데이터 연동은 현재 제공하지 않습니다.'] }
    ]
  },
  '/task-management/ScanTask': {
    title: '스캔 작업', summary: '대상, 템플릿, 실행 노드를 묶어 작업을 실행합니다.',
    sections: [
      { heading: '실행 순서', items: ['새 작업에서 이름과 대상(한 줄에 하나), 스캔 템플릿, 노드를 지정합니다.', '시작 후 작업 진행률에서 모듈별 상태를 확인하고, 완료 후 결과에서 수집 자산을 봅니다.', '보기는 작업 설정, 작업 메뉴는 재실행·시작·중지·삭제에 사용합니다.'] },
      { heading: '대량 대상', items: ['처음에는 작은 범위와 기본 템플릿으로 시험합니다.', '프로젝트로 동기화하면 선택한 작업의 결과를 기존 또는 새 프로젝트에 연결합니다.'] }
    ]
  },
  '/task-management/ScheduledTask': {
    title: '예약 작업', summary: '같은 범위를 주기적으로 다시 스캔합니다.',
    sections: [{ heading: '설정', items: ['대상·템플릿·노드를 정한 뒤 실행 주기를 지정합니다.', '저장 후 목록에서 실행 상태를 확인하고 필요하면 일시 중지합니다.', '실행 결과는 스캔 작업과 자산 정보에서 확인합니다.'] }]
  },
  '/task-management/ScanTemplate': {
    title: '스캔 템플릿', summary: '작업마다 사용할 모듈과 플러그인 구성을 저장합니다.',
    sections: [{ heading: '구성', items: ['기본 템플릿의 구성을 참고해 필요한 단계만 켭니다.', '포트 스캔·크롤러·취약점 검사는 실행 시간과 자원 사용량이 커질 수 있습니다.', '저장한 뒤 새 스캔 작업에서 해당 템플릿을 선택해야 적용됩니다.'] }]
  },
  '/plugin-management': {
    title: '플러그인 관리', summary: '스캔 플러그인과 서버 플러그인은 이 설치의 DB에서 관리합니다.',
    sections: [
      { heading: '직접 관리', items: ['새 플러그인에서 이름·모듈·버전·소개·소스 코드를 입력하거나 파일 가져오기를 사용합니다.', '스캔 플러그인을 저장·수정하면 전체 노드에 설치 메시지를 보냅니다. 노드 관리 → 플러그인에서 설치 상태를 확인합니다.', '실패하면 플러그인 로그와 노드 로그를 확인합니다. 서버 플러그인은 서버에서 등록됩니다.'] },
      { heading: '외부 마켓', items: ['마켓은 api.scope-sentry.top에서 목록·소스를 받아오는 선택 기능입니다. 원격 버전과 로컬 버전이 다르면 업데이트를 표시합니다.', '설치·업데이트 버튼을 누를 때 원격 소스를 가져와 로컬 DB에 저장합니다. 자동 업데이트는 하지 않습니다.', '원격 소개는 제공자가 작성한 언어로 보입니다. 직접 만든 플러그인의 소개는 한국어로 작성할 수 있습니다.'] }
    ]
  },
  '/node-management': {
    title: '노드 관리', summary: '노드는 스캔을 수행하는 별도 프로세스·컨테이너입니다.',
    sections: [
      { heading: '지표 해석', items: ['CPU·메모리는 스캐너가 운영체제에서 읽은 시스템 사용률(%)입니다. 메모리 35%는 35GB가 아닙니다.', 'Docker 환경에서는 컨테이너만의 사용량과 다를 수 있습니다. 컨테이너별 사용량은 docker stats로 따로 확인합니다.', '최대 작업 수는 Module Config의 maxGoroutineCount, 작업 수는 현재 실행 중, 완료 수는 스캐너가 보고한 처리 완료 건수입니다.'] },
      { heading: '작업 버튼', items: ['플러그인 창의 설치·점검 표시는 스캔 완료 여부가 아니라 해당 노드의 플러그인 상태입니다. 번호는 이름순 목록 위치이며 플러그인 고유 ID가 아닙니다.', '스캔이 끝나도 플러그인을 제거할 필요가 없습니다. 다음 작업에서 쓰지 않을 때만 제거하세요. 점검 실패가 보이면 노드·플러그인 로그를 확인한 뒤 다시 점검하거나 다시 설치합니다.', '설정에서는 노드 이름·모듈별 동시 처리 수·실행 상태를 편집합니다. 재시작 전 작업 상태를 확인합니다.'] }
    ]
  },
  '/project-management/project-detail': {
    title: '프로젝트 상세', summary: '프로젝트에 묶인 자산과 작업을 확인합니다.',
    sections: [{ heading: '활용', items: ['스캔 작업의 프로젝트로 동기화 기능으로 결과를 연결합니다.', '프로젝트 태그와 자산 범위를 기준으로 결과를 추적합니다.'] }]
  },
  '/project-management': {
    title: '프로젝트 관리', summary: '여러 작업의 결과를 대상이나 업무 단위로 묶습니다.',
    sections: [{ heading: '시작', items: ['프로젝트 이름과 태그를 정해 생성합니다.', '스캔 작업에서 결과를 프로젝트로 동기화합니다.', '상세 화면에서 연결된 자산과 작업을 확인합니다.'] }]
  },
  '/poc-management': {
    title: 'POC 관리', summary: '취약점 검증에 쓰는 POC 템플릿을 보관합니다.',
    sections: [
      { heading: '추가와 실행', items: ['새로 만들기에는 Python 코드가 아니라 Nuclei YAML 템플릿 전체를 입력합니다. 예제 넣기로 구조를 보고 id를 고유하게 바꾸세요.', 'info.name·info.severity·info.tags가 목록의 이름·위험 등급·태그로 저장됩니다. http 요청과 matchers는 검사 대상에 맞게 수정하세요.', '가져오기는 Nuclei YAML 파일을 포함한 ZIP을 받습니다. 저장만으로는 스캔하지 않습니다. 스캔 템플릿에서 취약점 검사 모듈을 켠 뒤 작업을 실행하고 자산 정보 → 취약점에서 결과를 확인합니다.'] },
      { heading: '정리', items: ['이름 검색과 위험 등급 필터로 범위를 좁힙니다.', '삭제할 항목의 체크박스를 확인합니다.'] }
    ]
  },
  '/fingerprint-management': {
    title: '핑거프린트 관리', summary: '웹·제품 식별에 쓰는 YAML 규칙을 관리합니다.',
    sections: [
      { heading: '직접 수정', items: ['이름으로 검색한 뒤 YAML과 사용 상태를 수정합니다. 새 규칙은 새로 만들기에서 추가합니다.', '저장·삭제하면 전체 노드에 핑거프린트 새로고침 메시지를 보냅니다.', '제품명은 결과의 식별자입니다. 기존 이름을 번역하면 과거 결과와 검색이 달라질 수 있습니다.'] },
      { heading: '업데이트 버튼', items: ['숫자는 외부 api.scope-sentry.top에서 로컬 버전 이후 등록된 규칙 수입니다.', '확인 후 원격 YAML을 받아 fingerprint_id 기준으로 추가·덮어쓰기하고 최신 시각을 버전으로 저장합니다.', '같은 ID의 원격 업데이트는 직접 수정한 내용을 덮어쓸 수 있습니다. 직접 만든 규칙은 별도 ID로 보관됩니다.', '중국어 제품명·탐지 문자열은 원격 데이터에 포함됩니다. 탐지 패턴을 번역하면 식별이 깨질 수 있습니다.'] }
    ]
  },
  '/sensitive-information-rules': {
    title: '민감정보 규칙', summary: '스캔 결과에서 비밀값·개인정보 후보를 찾는 규칙입니다.',
    sections: [{ heading: '운영', items: ['탐지 패턴과 사용 상태를 확인합니다.', '지역 프로필에 따라 기본 활성 규칙이 달라집니다. 프로필 변경 후 이 목록을 다시 검토합니다.', '스캔 템플릿에서 민감정보 검사를 켠 작업에 적용됩니다.'] }]
  },
  '/dictionary-management/port': {
    title: '포트 사전', summary: '포트 스캔 플러그인이 참조하는 포트 목록입니다.',
    sections: [{ heading: '활용', items: ['포트 묶음의 이름과 목록을 확인·수정합니다.', '플러그인 파라미터의 {port.이름} 참조가 실제 사전 이름과 일치해야 합니다.'] }]
  },
  '/dictionary-management': {
    title: '사전 관리', summary: '스캔 플러그인이 참조할 단어 목록입니다.',
    sections: [{ heading: '활용', items: ['용도별 사전을 추가하거나 파일을 가져옵니다.', '플러그인 파라미터의 {dict.종류.이름} 참조가 실제 사전과 일치해야 합니다.', '큰 사전은 스캔 시간과 요청량을 늘립니다. 작은 사전으로 먼저 시험합니다.'] }]
  },
  '/configuration/system': {
    title: '시스템 설정', summary: '저장한 설정은 DB에 기록되고 전체 노드에 갱신 메시지가 전달됩니다.',
    sections: [
      { heading: '기본 항목', items: ['시간대는 Asia/Seoul처럼 IANA 형식으로 입력합니다.', '지역 프로필 KR/CN/GLOBAL은 지역별 민감정보 기본 규칙과 자산 표시를 바꿉니다. 화면 언어와 별개입니다.'] },
      { heading: 'Module Config 읽기', items: ['YAML의 maxGoroutineCount는 동시에 처리할 대상 수의 상한입니다.', '각 모듈의 goroutineCount는 해당 단계의 동시 처리 수입니다. portScan은 포트 검사, webCrawler는 크롤링, vulnerabilityScan은 취약점 검사입니다.', '기존 YAML 전체를 지우지 말고 필요한 숫자만 바꿉니다. 공백 2칸 들여쓰기를 유지합니다.'] },
      { heading: '권장 시작값', items: ['2 vCPU·4GB 노드라면 maxGoroutineCount 2, portScan 1, webCrawler 1, vulnerabilityScan 1부터 시험합니다.', '4 vCPU·8GB 이상이라면 현재 기본값(각각 3, 2, 2, 2)을 먼저 사용합니다.', '작업 하나를 돌린 뒤 노드 관리의 CPU·메모리와 로그를 보고 한 항목씩 조정합니다.'], example: 'maxGoroutineCount: 2\nportScan:\n  goroutineCount: 1\nwebCrawler:\n  goroutineCount: 1\nvulnerabilityScan:\n  goroutineCount: 1' },
      { heading: '중복 제거 설정', items: ['선택한 결과 종류에서 정해진 필드가 같은 기록을 묶고 가장 최근 _id 한 건만 남깁니다. 나머지는 데이터베이스에서 삭제됩니다. 실행 전에 DB 백업을 확인하세요.', '웹 자산은 URL·상태 코드·본문 해시, 그 외 자산은 호스트·IP·프로토콜로 비교합니다. 서브도메인은 호스트·유형·IP 목록, URL 결과는 출력값, 취약점은 URL·취약점 ID·매칭값을 기준으로 비교합니다.', '현재 버전은 주기와 사용 여부를 저장하지만 자동 실행은 연결되지 않았습니다. 「지금 한 번 실행」을 켜고 저장해야 백그라운드에서 한 번 실행됩니다.', '민감정보와 페이지 모니터링 항목은 현재 서버의 중복 제거 처리 대상에 포함되지 않아 화면에서 비활성화했습니다.'] }
    ]
  },
  '/configuration/subfinder': {
    title: 'subfinder 설정', summary: 'subfinder는 수동 정보원에서 서브도메인을 수집합니다.',
    sections: [{ heading: 'API 키 입력', items: ['YAML은 정보원별 API 키 목록입니다. []는 키가 없다는 뜻입니다. 빈칸을 모두 채울 필요는 없습니다.', '계정이 있는 정보원만 키를 넣습니다. 키가 없어도 일부 소스는 작동하지만 수집 범위가 달라집니다.', '저장 후 작은 도메인으로 시험합니다. 키가 포함된 화면·파일은 공유하지 않습니다.'], example: 'github:\n  - YOUR_GITHUB_TOKEN\nsecuritytrails: []' }]
  },
  '/configuration/rad': {
    title: 'rad 설정', summary: '브라우저 기반 웹 크롤러의 탐색 범위와 자원 사용량을 조절합니다.',
    sections: [
      { heading: '처음 조절할 값', items: ['max_depth는 링크를 따라갈 단계, max_page_concurrent는 동시 방문 페이지 수입니다.', 'max_page_visit_per_site는 사이트별 최대 방문 페이지 수, max_page_visit는 전체 상한입니다.', '소형 노드는 깊이 2·동시 페이지 2·사이트당 100페이지부터 시작해 조절합니다.'], example: 'max_depth: 2\nmax_page_concurrent: 2\nmax_page_visit_per_site: 100\nenable_image: false' },
      { heading: '범위 제한', items: ['hostname_allowed/disallowed와 path_allowed/disallowed는 방문 범위를 제한합니다.', 'proxy가 필요 없으면 빈 문자열을 유지합니다. Docker에서 force_sandbox를 켜면 브라우저가 실행되지 않을 수 있습니다.', '기존 YAML에서 필요한 값만 바꾸고 저장합니다. 다음 크롤러 작업에서 결과를 확인합니다.'] }
    ]
  },
  '/configuration/apikey': {
    title: 'API Key 관리', summary: '외부 자동화가 ScopeSentry API를 호출할 때 사용합니다.',
    sections: [{ heading: '운영', items: ['연동별로 키를 분리해 발급하고 필요한 연동에만 전달합니다.', '사용하지 않는 키는 삭제합니다. 노출됐으면 폐기하고 새로 만듭니다.'] }]
  },
  '/about': {
    title: '정보', summary: '프로젝트 버전과 외부 연결 주소를 확인합니다.',
    sections: [{ heading: '참고', items: ['스캐너는 별도 컨테이너이므로 서버 UI 버전과 다를 수 있습니다.', '외부 플러그인 마켓은 선택 기능입니다. 직접 만든 플러그인은 플러그인 관리에서 추가할 수 있습니다.'] }]
  }
}
