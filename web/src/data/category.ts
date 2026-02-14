import { computed } from 'vue'
import { useMessage } from 'naive-ui'
import { categoryClient } from './client/category'
import { API_ENDPOINTS } from './client/api-endpoints'
import type { Category, CategoryQueryOptions } from '@/types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

export const useCategoriesQuery = (options: Partial<CategoryQueryOptions>) => {
  const { data, error, isPending } = useQuery<Category[], Error>({
    queryKey: [API_ENDPOINTS.CATEGORIES, options],
    queryFn: () => categoryClient.all(options as CategoryQueryOptions),
  })
  // @ts-ignore
  const categories = computed<Category[]>(() => data.value ?? []) // todo -> fix
  return {
    categories,
    error,
    loading: isPending,
  }
}

export const useCreateCategoryMutation = () => {
  const queryClient = useQueryClient()
  const message = useMessage()

  return useMutation({
    mutationFn: categoryClient.create,

    onSuccess: () => {
      message.success('Created successfully')

      queryClient.invalidateQueries({
        queryKey: [API_ENDPOINTS.CATEGORIES],
      })
    },
    onError: (error: Error) => {
      console.error('Create category failed:', error)
    },
  })
}

export const useUpdateCategoryMutation = () => {
  const queryClient = useQueryClient()
  const message = useMessage()

  return useMutation({
    mutationFn: categoryClient.patch,
    onSuccess: () => {
      message.success('Updated successfully')

      queryClient.invalidateQueries({
        queryKey: [API_ENDPOINTS.CATEGORIES],
      })
    },
  })
}

export const useDeleteCategoryMutation = () => {
  const queryClient = useQueryClient()
  const message = useMessage()

  return useMutation({
    mutationFn: categoryClient.delete,
    onSuccess: () => {
      message.success('Deleted successfully')

      queryClient.invalidateQueries({
        queryKey: [API_ENDPOINTS.CATEGORIES],
      })
    },
  })
}
