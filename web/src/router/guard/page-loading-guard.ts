export function createPageLoadingGuard(router: any) {
  router.beforeEach(() => {
    window.$loadingBar?.start()
  })

  router.afterEach(() => {
    setTimeout(() => {
      window.$loadingBar?.finish()
    }, 200)
  })

  router.onError(() => {
    window.$loadingBar?.error()
  })
}
