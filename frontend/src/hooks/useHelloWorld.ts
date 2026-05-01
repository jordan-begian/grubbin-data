import { useQuery } from '@tanstack/react-query'
import { fetchGreeting } from '../services/api'

export function useHelloWorldGreeting(name?: string) {
  return useQuery({
    queryKey: ['greeting', name],
    queryFn: () => fetchGreeting(name),
    refetchInterval: 30000,
  })
}
