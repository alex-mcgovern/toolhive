
import type { MCPServer } from '@/types/mcp-server'
import { useQuery, type UseQueryResult } from '@tanstack/react-query'

export function useQueryServers(): UseQueryResult<MCPServer[], Error> {
  return useQuery({
    queryFn: async () => {
      try {
        const response = await fetch('http://127.0.0.1:8080/api/v1beta/servers')
        
        if (!response.ok) {
          throw new Error(`Failed to fetch servers: ${response.statusText}`)
        }
        
        const data = await response.json()
        return data.servers
      } catch (error) {
        throw error instanceof Error ? error : new Error('Unknown error occurred while fetching servers')
      }
    },
    queryKey: ['servers']
  })
}
