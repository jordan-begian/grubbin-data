import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type {
  CreateDeliveryRequest,
  UpdateDeliveryRequest,
} from '@/types/delivery'
import {
  createDelivery,
  getDeliveries,
  updateDeliveries,
  deleteDeliveries,
} from '@/services/deliveries'

const DELIVERIES_QUERY_KEY = 'deliveries'

export function useDeliveries(userId: string | undefined) {
  const queryClient = useQueryClient()

  const deliveriesQuery = useQuery({
    queryKey: [DELIVERIES_QUERY_KEY, userId],
    queryFn: () => {
      if (!userId) throw new Error('User ID required')
      return getDeliveries(userId)
    },
    enabled: !!userId,
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateDeliveryRequest) => {
      if (!userId) throw new Error('User ID required')
      return createDelivery(userId, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [DELIVERIES_QUERY_KEY, userId],
      })
      toast.success('Delivery created successfully!')
    },
    onError: (err: Error) => {
      toast.error(`Failed to create delivery: ${err.message}`)
    },
  })

  const updateMutation = useMutation({
    mutationFn: (updates: UpdateDeliveryRequest[]) => {
      if (!userId) throw new Error('User ID required')
      return updateDeliveries(userId, updates)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [DELIVERIES_QUERY_KEY, userId],
      })
      toast.success('Delivery updated successfully!')
    },
    onError: (err: Error) => {
      toast.error(`Failed to update delivery: ${err.message}`)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (deliveryIds: string[]) => {
      if (!userId) throw new Error('User ID required')
      return deleteDeliveries(userId, deliveryIds)
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [DELIVERIES_QUERY_KEY, userId],
      })
      const count = variables.length
      toast.success(
        `${count} delivery${count > 1 ? 'ies' : 'y'} deleted successfully!`
      )
    },
    onError: (err: Error) => {
      toast.error(`Failed to delete delivery: ${err.message}`)
    },
  })

  return {
    deliveries: deliveriesQuery.data?.deliveries ?? [],
    stats: deliveriesQuery.data?.stats,
    isLoading: deliveriesQuery.isLoading,
    createDelivery: createMutation.mutateAsync,
    updateDeliveries: updateMutation.mutateAsync,
    deleteDeliveries: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  }
}
