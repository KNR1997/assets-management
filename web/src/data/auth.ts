import { fetchMe } from './user'
import { useUserStore } from '@/store'
import { authClient } from './client/auth'
import { setToken } from '@/utils/auth/token'
import { useMutation } from '@tanstack/vue-query'
import type { LoginInput, AuthResponse } from '@/types'

export function useLogin() {
  return useMutation<AuthResponse, Error, LoginInput>({
    mutationFn: authClient.login,
    onSuccess: async (response) => {
      setToken(response.data)
      // fetch user detail (username, email)
      const me = await fetchMe()
      // store details global state
      const userStore = useUserStore()
      userStore.setUserInfo(me.data)
    },
  })
}
