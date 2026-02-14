import { isNullOrWhitespace } from '@/utils'
import { getToken } from '@/utils/auth/token'

const WHITE_LIST = ['/login', '/404']

export function createAuthGuard(router: any) {
  router.beforeEach((to: any) => {
    const token = getToken()

    // User NOT logged in
    if (isNullOrWhitespace(token)) {
      // Allow white-listed routes
      if (WHITE_LIST.includes(to.path)) {
        return true
      }

      // Redirect to login
      return {
        path: '/login',
        query: { redirect: to.fullPath },
      }
    }

    // User logged in but trying to access login page
    if (to.path === '/login') {
      return { path: '/' }
    }

    return true
  })
}
