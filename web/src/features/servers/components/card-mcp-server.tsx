import { Card, CardContent, CardHeader, CardTitle } from "@/common/components/ui/card";
import { StatusIcon } from "@/common/components/ui/status-icon";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/common/components/ui/tooltip";
import type { MCPServer } from "@/types/mcp-server";


// TODO: We don't know all the possible members of `state` yet
function CardStatusIcon({ state }: { state: string }) {
  switch (state.toLowerCase()) {
    case 'running':
      return <StatusIcon variant="success" />
    case 'exited':
      return <StatusIcon variant="destructive" />
    case 'paused':
      return <StatusIcon variant="warning" />
    default:
      return <StatusIcon variant="default" />
  }
}

function CardStatusTooltip({ state, status }: { state: string; status: string }) {
  return (
    <Tooltip>
      <TooltipTrigger>
        <CardStatusIcon state={state} />
      </TooltipTrigger>
      <TooltipContent>
        <p>{status}</p>
      </TooltipContent>
    </Tooltip>
  )
}

export function CardMcpServer({ Name, State, Status, Image }: MCPServer) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <CardStatusTooltip state={State} status={Status} />
          {Name}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <span className="text-muted-foreground font-mono text-sm">{Image}</span>
      </CardContent>
    </Card>
  )
}
