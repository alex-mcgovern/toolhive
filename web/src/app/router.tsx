import { Routes, Route, Navigate } from 'react-router-dom'
import { RouteServers } from '../app/routes/route-servers'

export function Router() {
  return (
    <Routes>
      <Route path="/" element={<RouteServers />} />
      <Route path="/clients" element={<p>Clients page</p>} />
      <Route path="/store" element={<p>Store page</p>} />
      <Route path="/logs" element={<p>Logs page</p>} />
      <Route path="/secrets" element={<p>Secrets page</p>} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
