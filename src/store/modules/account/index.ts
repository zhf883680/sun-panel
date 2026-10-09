import { defineStore } from 'pinia'
import type { AccountInfo } from './helper'
import { getLocalState, setLocalState } from './helper'
import { getInfo as getUserInfoByToken } from '@/api/system/user'
import { logoutByToken } from '@/api'

export const useAccountStore = defineStore('account-manager-store', {
  state: () => getLocalState(),

  actions: {
    // 添加账号（同账号会先移除旧的记录）
    addAccount(account: AccountInfo) {
      const exist = this.getAccountInfo(account.username)
      if (exist) {
        if (exist.loginStatus)
          this.logoutAndDeleteAccount(exist)
        else
          this.deleteAccount(exist)
      }
      this.accounts.push(account)
      this.recordState()
      return true
    },

    // 是否已保存该账号
    isExistAccount(username: string) {
      return this.accounts.some(account => account.username === username)
    },

    deleteAccount(account: AccountInfo) {
      for (let i = 0; i < this.accounts.length; i++) {
        const item = this.accounts[i]
        if (item.username === account.username || item.cToken === account.cToken) {
          this.accounts.splice(i, 1)
          break
        }
      }
      this.recordState()
    },

    getAccountInfo(username: string): AccountInfo | null {
      const index = this.accounts.findIndex(account => account.username === username)
      return index !== -1 ? this.accounts[index] : null
    },

    // 检查已保存账号的登录状态
    async loginCheck() {
      for (let i = 0; i < this.accounts.length; i++) {
        const account = this.accounts[i]
        if (account.loginStatus) {
          try {
            const { data } = await getUserInfoByToken<{ name: string }>(account.cToken)
            this.accounts[i].name = data.name || ''
          }
          catch {
            this.accounts[i].loginStatus = false
          }
        }
      }
      this.recordState()
    },

    // 退出该账号并删除本地记录
    async logoutAndDeleteAccount(account: AccountInfo) {
      const cToken = account.cToken
      this.deleteAccount(account)
      return logoutByToken(cToken)
    },

    clearAccount() {
      this.accounts = []
      this.recordState()
    },

    recordState() {
      setLocalState(this.$state)
    },
  },
})
