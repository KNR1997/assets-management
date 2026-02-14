import { createAuthGuard } from './auth-guard'
import { createPageLoadingGuard } from './page-loading-guard'
import { createPageTitleGuard } from './page-title-guard'

export function setupRouterGuard(router: any) {
  createPageLoadingGuard(router)
  createAuthGuard(router)
  createPageTitleGuard(router)
}
