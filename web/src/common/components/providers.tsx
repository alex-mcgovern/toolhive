import type { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { TooltipProvider } from './ui/tooltip'

const queryClient = new QueryClient()

interface AppProvidersProps {
  children: ReactNode
}

export function Providers({ children }: AppProvidersProps) {
  return (
    <BrowserRouter>
      <TooltipProvider delayDuration={0}>
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      </TooltipProvider>
    </BrowserRouter>
  )
}
