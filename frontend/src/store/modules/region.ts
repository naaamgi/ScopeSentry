import { defineStore } from 'pinia'
import { store } from '../index'
import { getSystemConfigurationApi } from '@/api/Configuration'

// 지역 프로필. 어느 지역 전용 자산 필드와 규칙을 쓸지 고르는 설정이며,
// 화면 언어(locale)와는 별개다. 서버의 system 설정에 들어 있다.
export type RegionProfile = 'CN' | 'KR' | 'GLOBAL'

export const REGION_PROFILES: RegionProfile[] = ['CN', 'KR', 'GLOBAL']

// 업스트림이 전제하던 환경. 서버 쪽 기본값과 같아야 한다.
const DEFAULT_REGION: RegionProfile = 'CN'

const normalize = (value: unknown): RegionProfile => {
  const upper = String(value ?? '')
    .trim()
    .toUpperCase()
  return (REGION_PROFILES as string[]).includes(upper) ? (upper as RegionProfile) : DEFAULT_REGION
}

interface RegionState {
  region: RegionProfile
}

export const useRegionStore = defineStore('region', {
  state: (): RegionState => ({
    region: DEFAULT_REGION
  }),
  getters: {
    getRegion(): RegionProfile {
      return this.region
    },
    // 미니프로그램은 중국 위챗 생태계의 자산 타입이라 다른 지역에서는 숨긴다.
    showMiniProgram(): boolean {
      return this.region === 'CN'
    }
  },
  actions: {
    setRegion(value: unknown) {
      this.region = normalize(value)
    },
    // 서버 설정을 읽어 저장된 값을 갱신한다. 값이 영속화되어 있어 화면은 요청을
    // 기다리지 않고 바로 그릴 수 있고, 이 호출은 그 값을 최신으로 맞춘다.
    async fetchRegion() {
      try {
        const res = await getSystemConfigurationApi()
        if (res?.code === 200) {
          this.setRegion(res.data?.region)
        }
      } catch {
        // 설정을 읽지 못해도 화면은 떠야 한다. 저장된 값이나 기본값을 그대로 쓴다.
      }
    }
  },
  persist: true
})

export const useRegionStoreWithOut = () => {
  return useRegionStore(store)
}
