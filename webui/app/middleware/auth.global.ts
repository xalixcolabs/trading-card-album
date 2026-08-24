export default defineNuxtRouteMiddleware(async (to) => {
    if (import.meta.server) return
    if (to.path === '/login') return
    const isLogin = useCookie("jwt")
    if (!isLogin.value) return navigateTo('/login')
})
