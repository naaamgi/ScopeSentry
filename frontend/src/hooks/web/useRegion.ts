import { useRegionStore } from '@/store/modules/region'
import { useI18n } from '@/hooks/web/useI18n'

// 지역 프로필에 따라 달라지는 화면 조각을 한곳에 모아 둔다.
//
// 값은 setup 시점에 한 번 읽는다. 지역 프로필은 설정 화면에서 저장하고 나면
// 화면을 다시 불러야 반영되는데, 자산 표의 컬럼 정의(crudSchemas)가 setup 에서
// 한 번만 만들어지기 때문이다. 언어 전환도 같은 이유로 페이지를 새로 읽는다.
export const useRegion = () => {
  const regionStore = useRegionStore()
  const { t } = useI18n()

  const region = regionStore.getRegion

  return {
    region,

    // 미니프로그램 자산 탭을 보여줄지 여부.
    showMiniProgram: regionStore.showMiniProgram,

    // 자산의 등록번호 필드 라벨. 저장 필드는 어느 지역에서나 icp 하나다.
    registrationLabel: t(`asset.registration.${region}`),

    // 검색 문법 표에 쓸 키워드. 서버는 icp 와 brn 을 모두 받는다.
    registrationKeyword: region === 'KR' ? 'brn' : 'icp',

    // 검색 문법 표의 설명문.
    registrationSearchHelp: region === 'KR' ? t('searchHelp.brn') : t('searchHelp.icp'),

    // 스캔 대상 입력란의 안내문. 지역 전용 등록번호 줄을 덧붙인다.
    targetHint: t('task.msgTarget') + t(`task.msgTargetRegistration.${region}`)
  }
}
