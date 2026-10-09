import { ss } from '@/utils/storage'

const LOCAL_NAME = 'accountManager'

export interface AccountInfo {
  username: string
  password?: string
  name?: string
  cToken: string
  loginStatus: boolean
}

export interface AccountManagerState {
  accounts: AccountInfo[]
}

export function getLocalState(): AccountManagerState {
  const localState: AccountManagerState | null = ss.get(LOCAL_NAME)
  return { accounts: [], ...localState }
}

export function setLocalState(state: AccountManagerState): void {
  ss.set(LOCAL_NAME, state)
}
