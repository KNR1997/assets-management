import type { ProfileUpdateInput, UserMeResponse } from '@/types'
import { API_ENDPOINTS } from './api-endpoints'
import { HttpClient } from './http-client'

export const userClient = {
  update: (variables: ProfileUpdateInput) => {
    return HttpClient.patch(API_ENDPOINTS.PROFILE, variables)
  },
  me: () => {
    return HttpClient.get<UserMeResponse>(API_ENDPOINTS.ME)
  },
}
