import { getApiV1AuthMe } from '~/services/auth/auth'

export default defineNuxtRouteMiddleware(async () => {
  const profile = useProfile()
  if (profile.value?.is_admin) return

  try {
    const { data } = await getApiV1AuthMe()
    if (data.is_admin) {
      useProfile(data as any)
      return
    }
  } catch {
    return navigateTo('/login')
  }

  return navigateTo('/')
})