<script setup lang="ts">
import { LoginForm } from './components'
import SetupForm from './components/SetupForm.vue'
import { setupStatusApi } from '@/api/login'
import { ThemeSwitch } from '@/components/ThemeSwitch'
import { LocaleDropdown } from '@/components/LocaleDropdown'
import { useI18n } from '@/hooks/web/useI18n'
import { underlineToHump } from '@/utils'
import { useAppStore } from '@/store/modules/app'
import { useDesign } from '@/hooks/web/useDesign'
import { onMounted, ref } from 'vue'
import { ElScrollbar } from 'element-plus'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('login')

const appStore = useAppStore()

const { t } = useI18n()

const isLogin = ref(true)
const needsSetup = ref(false)

onMounted(async () => {
  try {
    const result = await setupStatusApi()
    needsSetup.value = result.data.required
  } catch {
    // The login form remains available when an older server has no setup API.
  }
})

const toRegister = () => {
  isLogin.value = false
}

const toLogin = () => {
  isLogin.value = true
}
</script>

<template>
  <div
    :class="prefixCls"
    class="h-[100%] relative lt-xl:bg-[var(--login-bg-color)] lt-sm:px-10px lt-xl:px-10px lt-md:px-10px"
  >
    <ElScrollbar class="h-full">
      <div class="relative flex mx-auto min-h-100vh">
        <div
          :class="`${prefixCls}__left flex-1 relative p-30px lt-xl:hidden`"
        >
          <div class="flex items-center relative text-[var(--text-primary)]">
            <img src="@/assets/imgs/logo.png" alt="" class="w-48px h-48px mr-10px" />
            <span class="text-20px font-bold">{{ underlineToHump(appStore.getTitle) }}</span>
          </div>
          <div class="flex justify-center items-center h-[calc(100%-60px)]">
            <TransitionGroup
              appear
              tag="div"
              enter-active-class="animate__animated animate__bounceInLeft"
            >
              <div class="login-wordmark" key="1">Scope Sentry</div>
            </TransitionGroup>
          </div>
        </div>
        <div class="flex-1 p-30px lt-sm:p-10px bg-[var(--bg-page)] relative">
          <div
            class="flex justify-between items-center text-[var(--text-primary)] at-2xl:justify-end at-xl:justify-end"
          >
            <div class="flex items-center at-2xl:hidden at-xl:hidden">
              <img src="@/assets/imgs/logo.png" alt="" class="w-48px h-48px mr-10px" />
              <span class="text-20px font-bold">{{ underlineToHump(appStore.getTitle) }}</span>
            </div>

            <div class="flex justify-end items-center space-x-10px">
              <ThemeSwitch />
              <LocaleDropdown />
            </div>
          </div>
          <Transition appear enter-active-class="animate__animated animate__bounceInRight">
            <div
              class="h-full flex items-center m-auto w-[100%] at-2xl:max-w-500px at-xl:max-w-500px at-md:max-w-500px at-lg:max-w-500px"
            >
              <SetupForm v-if="needsSetup" class="m-auto" @complete="needsSetup = false" />
              <LoginForm
                v-else-if="isLogin"
                class="p-20px h-auto m-auto rounded-[var(--radius-lg)] bg-[var(--bg-card)]"
                @to-register="toRegister"
              />
            </div>
          </Transition>
        </div>
      </div>
    </ElScrollbar>
  </div>
</template>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-login';

.@{prefix-cls} {
  overflow: auto;

  &__left {
    background: var(--accent-bg);
    border-right: 1px solid var(--border);
  }
}
.login-wordmark { color: var(--accent); font-size: clamp(36px, 4vw, 72px); font-weight: 700; letter-spacing: -.05em; }
</style>
