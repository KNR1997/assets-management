import { useMessage } from 'naive-ui'
import { userClient } from './client/user'
import { API_ENDPOINTS } from './client/api-endpoints'
import { useMutation, useQueryClient } from '@tanstack/vue-query'

export const useUpdateProfile = () => {
  const queryClient = useQueryClient()
  const message = useMessage()

  return useMutation({
    mutationFn: userClient.update,
    onSuccess: () => {
      message.success('Updated successfully')

      queryClient.invalidateQueries({
        queryKey: [API_ENDPOINTS.PROFILE],
      })
    },
  })
}

export const fetchMe = () => {
  return userClient.me()
}
