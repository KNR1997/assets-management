<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
// assets
import bgImg from '@/assets/images/login_bg.webp'
// utils
import { setToken } from '@/utils/auth/token'
// icons
import { Icon } from '@iconify/vue'
// hooks
import { useLogin } from '@/data/auth'
// components
import AppPage from '@/components/page/AppPage.vue'

const router = useRouter()
const { query } = useRoute()
const { t } = useI18n({ useScope: 'global' })
const message = useMessage()

const { mutateAsync: login, isPending } = useLogin()

const loginInfo = ref({
  email: '',
  password: '',
})

async function handleLogin() {
  const { email, password } = loginInfo.value
  if (!email || !password) {
    message.warning(t('views.login.message_input_username_password'))
    return
  }
  try {
    const response = await login({ email, password })
    // save token
    setToken(response.data)
    // toast message
    message.success(t('views.login.message_login_success'))
    // redirect
    if (query.redirect) {
      const path = query.redirect
      Reflect.deleteProperty(query, 'redirect')
      // @ts-ignore
      router.push({ path, query })
    } else {
      router.push('/')
    }
  } catch (e: any) {
    console.error('login error', e.error)
  }
}
</script>

<template>
  <AppPage :show-footer="true" bg-cover :style="{ backgroundImage: `url(${bgImg})` }">
    <div
      style="transform: translateY(25px)"
      class="m-auto max-w-1500 min-w-345 f-c-c rounded-10 bg-white bg-opacity-60 p-15 card-shadow"
      dark:bg-dark
    >
      <div hidden w-380 px-20 py-35 md:block>
        <!-- <icon-custom-front-page pt-10 text-300 color-primary></icon-custom-front-page> -->
        <Icon text-100 icon="logos:vue" />
      </div>
      <div w-320 flex-col px-20 py-35>
        <!-- <h5 f-c-c text-24 font-normal color="#6a6a6a">
          <icon-custom-logo mr-10 text-50 color-primary />{{ $t('app_name') }}
        </h5> -->
        <div mt-30>
          <n-input
            autofocus
            class="h-50 items-center pl-10 text-16"
            v-model:value="loginInfo.email"
            type="email"
            placeholder="Email"
          />
        </div>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.password"
            class="h-50 items-center pl-10 text-16"
            type="password"
            show-password-on="mousedown"
            placeholder="123456"
            :maxlength="20"
            @keypress.enter="handleLogin"
          />
        </div>

        <div mt-20>
          <n-button
            h-50
            w-full
            rounded-5
            text-16
            type="primary"
            @click="handleLogin"
            :loading="isPending"
          >
            Login
          </n-button>
        </div>
      </div>
    </div>
  </AppPage>
</template>
