<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, NPopconfirm, NSelect, NTag, useMessage } from 'naive-ui'
import { login } from '@/api'
import { getSiteSetting } from '@/api/openness'
import { useAccountStore, useAppStore, useAuthStore } from '@/store'
import { SvgIcon } from '@/components/common'
import { router } from '@/router'
import { t } from '@/locales'
import { languageOptions } from '@/utils/defaultData'
import type { Language } from '@/store/modules/app/helper'
import type { AccountInfo } from '@/store/modules/account/helper'
import { applySiteSetting } from '@/utils/cmn'

const authStore = useAuthStore()
const appStore = useAppStore()
const accountStore = useAccountStore()
const ms = useMessage()
const loading = ref(false)
const languageValue = ref<Language>(appStore.language)

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

const loginBackgroundUrl = ref('')

onMounted(() => {
  getSiteSetting<{ title: string, favicoUrl: string, loginBackgroundUrl: string }>().then(({ code, data }) => {
    if (code === 0) {
      applySiteSetting(data)
      loginBackgroundUrl.value = data?.loginBackgroundUrl || ''
    }
  }).catch(() => {})
})

const savedAccounts = computed(() => accountStore.accounts)
const isDark = computed(() => appStore.theme === 'dark')

const loginPost = async () => {
  if (!form.value.username || !form.value.password) {
    ms.warning(t('login.usernameOrPasswordRequired'))
    return
  }

  loading.value = true
  try {
    const res = await login<Login.LoginResponse>(form.value)
    if (res.code === 0) {
      authStore.setToken(res.data.token)
      authStore.setUserInfo(res.data)

      // 记住账号，方便下次一键登录
      accountStore.addAccount({
        username: form.value.username,
        password: form.value.password,
        name: res.data.name || form.value.username,
        cToken: res.data.token,
        loginStatus: true,
      })

      setTimeout(() => {
        ms.success(`Hi ${res.data.name},${t('login.welcomeMessage')}`)
        loading.value = false
        router.push({ path: '/' })
      }, 500)
    }
    else {
      loading.value = false
    }
  }
  catch (error) {
    loading.value = false
    console.log(error)
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}

// 一键登录已保存的账号
function handleQuickLogin(account: AccountInfo) {
  if (!account.password) {
    form.value.username = account.username
    form.value.password = ''
    return
  }
  form.value.username = account.username
  form.value.password = account.password
  loginPost()
}

// 添加账号（清空表单）
function handleAddAccount() {
  form.value = { username: '', password: '' }
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
}

function handleSwitchTheme() {
  appStore.setTheme(isDark.value ? 'light' : 'dark')
}
</script>

<template>
  <div
    class="login-container"
    :style="loginBackgroundUrl ? { backgroundImage: `url(${loginBackgroundUrl})`, backgroundSize: 'cover', backgroundPosition: 'center' } : {}"
  >
    <!-- 渐变光斑背景 -->
    <div class="blob blob-1" />
    <div class="blob blob-2" />
    <div class="blob blob-3" />
    <div class="blob blob-4" />

    <NCard class="login-card glass-card" style="border-radius: 20px;">
      <div class="mb-5 flex items-center justify-between">
        <div class="flex items-center">
          <NButton quaternary size="small" circle @click="handleSwitchTheme">
            {{ isDark ? '☀️' : '🌙' }}
          </NButton>
        </div>

        <div class="flex items-center">
          <div class="mr-2">
            <SvgIcon icon="ion-language" style="width: 20;height: 20;" />
          </div>
          <div class="min-w-[100px]">
            <NSelect v-model:value="languageValue" size="small" :options="languageOptions" @update-value="handleChangeLanuage" />
          </div>
        </div>
      </div>

      <div class="login-title">
        <NGradientText :size="30" type="success" class="!font-bold">
          {{ $t('common.appName') }}
        </NGradientText>
      </div>

      <!-- 已保存的账号 -->
      <div v-if="savedAccounts.length > 0" class="mb-4">
        <div class="text-xs text-slate-400 mb-[5px]">
          {{ $t('login.savedAccounts') }}
        </div>
        <div class="flex flex-wrap gap-[5px]">
          <div
            v-for="account, index in savedAccounts" :key="index"
            class="flex items-center rounded-lg bg-white/60 dark:bg-zinc-800/60 px-[8px] py-[4px] cursor-pointer hover:bg-white"
            :title="$t('login.quickLogin')"
            @click="handleQuickLogin(account)"
          >
            <span class="text-sm mr-[5px]">{{ account.name || account.username }}</span>
            <NTag v-if="account.loginStatus" size="tiny" type="success" :bordered="false">
              {{ $t('login.quickLogin') }}
            </NTag>
            <NPopconfirm
              :positive-text="$t('common.confirm')"
              :negative-text="$t('common.cancel')"
              @positive-click="accountStore.deleteAccount(account)"
            >
              <template #trigger>
                <span class="ml-[5px] text-slate-400" @click.stop>
                  <SvgIcon icon="line-md:close-small" style="width:16px;height:16px;" />
                </span>
              </template>
              {{ $t('login.removeAccountConfirm') }}
            </NPopconfirm>
          </div>
        </div>
      </div>

      <NForm :model="form" label-width="100px" @keydown.enter="handleSubmit">
        <NFormItem>
          <NInput v-model:value="form.username" :placeholder="$t('login.usernamePlaceholder')">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem>
          <NInput v-model:value="form.password" type="password" :placeholder="$t('login.passwordPlaceholder')">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block :loading="loading" @click="handleSubmit">
            {{ $t('login.loginButton') }}
          </NButton>
        </NFormItem>

        <div v-if="savedAccounts.length > 0" class="flex justify-center mb-2">
          <NButton quaternary type="info" @click="handleAddAccount">
            {{ $t('login.addAccountButton') }}
          </NButton>
        </div>

        <div class="flex justify-center text-slate-300">
          Powered By <a href="https://github.com/zhf883680/sun-panel" target="_blank" class="ml-[5px] text-slate-500">Sun-Panel</a>
        </div>
      </NForm>
    </NCard>
  </div>
</template>

<style>
.login-container {
    position: relative;
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    overflow: hidden;
    background: #eef2ff;
}

/* 夜间模式 */
.dark .login-container {
    background: #16161a;
}

/* 渐变光斑 */
.blob {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.55;
    pointer-events: none;
}

.blob-1 {
    width: 420px;
    height: 420px;
    background: #7c9cff;
    top: -80px;
    left: -60px;
    animation: blob-float-1 14s ease-in-out infinite;
}

.blob-2 {
    width: 380px;
    height: 380px;
    background: #6ee7b7;
    bottom: -100px;
    right: -60px;
    animation: blob-float-2 16s ease-in-out infinite;
}

.blob-3 {
    width: 300px;
    height: 300px;
    background: #f0abfc;
    top: 40%;
    right: 15%;
    animation: blob-float-3 18s ease-in-out infinite;
}

.blob-4 {
    width: 260px;
    height: 260px;
    background: #fcd34d;
    bottom: 10%;
    left: 12%;
    animation: blob-float-1 20s ease-in-out infinite;
}

@keyframes blob-float-1 {
    0%, 100% { transform: translate(0, 0) scale(1); }
    50% { transform: translate(40px, 30px) scale(1.1); }
}

@keyframes blob-float-2 {
    0%, 100% { transform: translate(0, 0) scale(1); }
    50% { transform: translate(-50px, -20px) scale(1.08); }
}

@keyframes blob-float-3 {
    0%, 100% { transform: translate(0, 0) scale(1); }
    50% { transform: translate(20px, -40px) scale(0.95); }
}

/* 毛玻璃卡片 */
.glass-card {
    position: relative;
    z-index: 1;
    background: rgba(255, 255, 255, 0.72) !important;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.6);
    box-shadow: 0 20px 50px -20px rgba(30, 41, 59, 0.45);
}

.dark .glass-card {
    background: rgba(30, 30, 34, 0.72) !important;
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 20px 50px -20px rgba(0, 0, 0, 0.8);
}

@media (min-width: 600px) {
    .login-card {
        width: auto;
        margin: 0px 10px;
    }
    .login-button {
        width: 100%;
    }
}

.login-card {
    margin: 20px;
    min-width: 400px;
}

.login-title {
    text-align: center;
    margin: 20px;
}
</style>
