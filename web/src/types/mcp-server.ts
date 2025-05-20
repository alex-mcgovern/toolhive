// Server types definition
export interface ServerLabels {
  [key: string]: string
}

export interface MCPServer {
  ID: string
  Name: string
  Image: string
  Status: string
  State: string
  Created: string
  Labels: ServerLabels
  Ports: unknown[]
}

export interface MCPServersResponse {
  servers: MCPServer[]
}
