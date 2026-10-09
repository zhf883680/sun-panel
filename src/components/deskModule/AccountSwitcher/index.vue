<script setup lang="ts">
import { computed, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NModal, NPopover, NTag, useDialog, useMessage } from 'naive-ui'
import { SvgIcon } from '@/components/common'
import { login, logout } from '@/api'
import { useAccountStore, useAuthStore } from '@/store'
import { t } from '@/locales'
import type { AccountInfo } from '@/store/modules/account/helper'

// 非 PRO 授权时的账号数量限制
const MAX_ACCOUNT = 5

const authStore = useAuthStore()
const accountStore = useAccountStore()
const ms = useMessage()
const dialog = useDialog()

const show = ref(false)
const addAccountShow = ref(false)
const addLoading = ref(false)
const form = ref({ username: '', password: '' })

const accounts = computed(() => accountStore.accounts)

function isCurrent(account: AccountInfo) {
  return account.cToken === authStore.token
}

// 切换账号
async function handleSwitch(account: AccountInfo) {
  if (isCurrent(account))
    return

  if (!account.loginStatus) {
    ms.warning(t('deskModule.accountSwitcher.failed'))
    return
  }

  authStore.setToken(account.cToken)
  show.value = false
  window.location.reload()
}

// 单账号退出
function handleLogout(account: AccountInfo) {
  dialog.warning({
    title: t('common.warning'),
    content: t('deskModule.accountSwitcher.logoutConfirmText', { username: account.name || account.username }),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await accountStore.logoutAndDeleteAccount(account)
      // 退出的是当前账号，需要跳回登录页
      if (isCurrent(account)) {
        authStore.removeToken()
        window.location.href = '/'
      }
    },
  })
}

// 退出全部账号
function handleLogoutAll() {
  dialog.warning({
    title: t('common.warning'),
    content: t('deskModule.accountSwitcher.logoutAllConfirmText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await logout()
      }
      catch {
        // 忽略退出失败
      }
      accountStore.clearAccount()
      authStore.removeToken()
      window.location.href = '/'
    },
  })
}

function handleOpenAddAccount() {
  if (accounts.value.length >= MAX_ACCOUNT) {
    ms.warning(t('deskModule.accountSwitcher.addAccountQuantityLimit', { number: MAX_ACCOUNT }))
    return
  }
  form.value = { username: '', password: '' }
  addAccountShow.value = true
}

// 添加账号（登录成功后保存）
async function handleAddAccount() {
  if (!form.value.username || !form.value.password) {
    ms.warning(t('login.usernameOrPasswordRequired'))
    return
  }

  if (accountStore.isExistAccount(form.value.username)) {
    ms.warning(t('deskModule.accountSwitcher.accountExist'))
    return
  }

  addLoading.value = true
  try {
    const res = await login<Login.LoginResponse & { name: string }>(form.value)
    if (res.code === 0) {
      accountStore.addAccount({
        username: form.value.username,
        password: form.value.password,
        name: res.data.name || form.value.username,
        cToken: res.data.token,
        loginStatus: true,
      })
      addAccountShow.value = false
      ms.success(t('common.success'))
    }
  }
  catch {
    // 错误提示由请求层处理
  }
  addLoading.value = false
}
</script>

<template>
  <NPopover v-model:show="show" trigger="click" placement="bottom-end" style="max-width: 280px;" :show-arrow="false">
    <template #trigger>
      <NButton color="#2a2a2a6b" :title="t('deskModule.accountSwitcher.accountSwitch')">
        <template #icon>
          <SvgIcon class="text-white font-xl" icon="lucide:users" />
        </template>
      </NButton>
    </template>

    <div class="w-[240px]">
      <div class="text-slate-500 text-xs mb-[5px]">
        {{ t('deskModule.accountSwitcher.loggedInAccount') }}
      </div>

      <div v-if="accounts.length === 0" class="text-center text-slate-400 text-xs py-[10px]">
        {{ t('deskModule.accountSwitcher.empty') }}
      </div>

      <div
        v-for="account, index in accounts" :key="index"
        class="flex items-center justify-between rounded-lg p-[5px] cursor-pointer hover:bg-slate-100 dark:hover:bg-zinc-700"
        @click="handleSwitch(account)"
      >
        <div class="flex items-center overflow-hidden">
          <div class="mr-[5px] truncate" :class="account.loginStatus ? '' : 'text-slate-400 line-through'">
            {{ account.name || account.username }}
          </div>
          <NTag v-if="isCurrent(account)" size="tiny" type="success" :bordered="false">
            {{ t('deskModule.accountSwitcher.current') }}
          </NTag>
        </div>
        <NButton quaternary size="tiny" @click.stop="handleLogout(account)">
          <template #icon>
            <SvgIcon icon="tabler-logout" />
          </template>
        </NButton>
      </div>

      <div class="flex justify-between mt-[10px]">
        <NButton quaternary size="tiny" type="primary" @click="handleOpenAddAccount">
          {{ t('deskModule.accountSwitcher.addAccount') }}
        </NButton>
        <NButton v-if="accounts.length > 0" quaternary size="tiny" type="error" @click="handleLogoutAll">
          {{ t('deskModule.accountSwitcher.logoutAll') }}
        </NButton>
      </div>
    </div>
  </NPopover>

  <NModal
    v-model:show="addAccountShow" preset="card" size="small"
    style="max-width: 400px;border-radius: 1rem;" :title="t('deskModule.accountSwitcher.addAccount')"
  >
    <NForm @keydown.enter="handleAddAccount">
      <NFormItem :label="t('common.username')">
        <NInput v-model:value="form.username" type="text" :placeholder="t('login.usernamePlaceholder')" />
      </NFormItem>
      <NFormItem :label="t('common.password')">
        <NInput v-model:value="form.password" type="password" :placeholder="t('login.passwordPlaceholder')" />
      </NFormItem>
    </NForm>
    <template #footer>
      <NButton type="primary" :loading="addLoading" style="float: right;" @click="handleAddAccount">
        {{ t('common.confirm') }}
      </NButton>
    </template>
  </NModal>
</template>
