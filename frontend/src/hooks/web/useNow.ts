import { dateUtil } from '@/utils/dateUtil'
import { reactive, toRefs } from 'vue'
import { tryOnMounted, tryOnUnmounted } from '@vueuse/core'
import { useI18n } from '@/hooks/web/useI18n'

// dayjs 의 day() 반환값(0=일요일)과 같은 순서
const WEEKDAY_KEYS = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday']

export const useNow = (immediate = true) => {
  const { t } = useI18n()

  let timer: IntervalHandle

  const state = reactive({
    year: 0,
    month: 0,
    week: '',
    day: 0,
    hour: '',
    minute: '',
    second: 0,
    meridiem: ''
  })

  const update = () => {
    const now = dateUtil()

    const h = now.format('HH')
    const m = now.format('mm')
    const s = now.get('s')

    state.year = now.get('y')
    state.month = now.get('M') + 1
    state.week = t(`weekday.${WEEKDAY_KEYS[now.day()]}`)
    state.day = now.get('date')
    state.hour = h
    state.minute = m
    state.second = s

    state.meridiem = now.format('A')
  }

  function start() {
    update()
    clearInterval(timer)
    timer = setInterval(() => update(), 1000)
  }

  function stop() {
    clearInterval(timer)
  }

  tryOnMounted(() => {
    immediate && start()
  })

  tryOnUnmounted(() => {
    stop()
  })

  return {
    ...toRefs(state),
    start,
    stop
  }
}
