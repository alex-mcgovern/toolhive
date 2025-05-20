import { useQueryServers } from '../hooks/use-query-servers'
import { CardMcpServer } from './card-mcp-server'

export function GridCardsMcpServers() {
  const { data, isLoading, error, isError } = useQueryServers()

  if (isLoading === true) {
    return <div className="text-center py-4">Loading servers...</div>
  }

  if (isError === true || !data) {
    return (
      <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-md">
        <p className="font-medium">Error</p>
        <p>{error?.message ?? 'An error occurred'}</p>
        <p>The API server might not be running, or network connectivity issues occurred.</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      {data.map((server) => (
        <CardMcpServer key={server.ID} {...server} />
      ))}
    </div>
  )
}
