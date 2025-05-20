import { Button } from '../../common/components/ui/button'

import { PlusCircle } from 'lucide-react'
import { GridCardsMcpServers } from '../../features/servers/components/grid-cards-mcp-server'

export function RouteServers() {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-semibold">Installed</h1>
        <Button className="gap-2">
          <PlusCircle className="h-4 w-4" />
          Add a tool
        </Button>
      </div>
      <GridCardsMcpServers />
    </div>
  )
}
