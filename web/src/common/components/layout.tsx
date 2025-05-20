import type { ReactNode } from 'react'
import { TopNav } from './top-nav'

export function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-col h-screen">
      <TopNav />
      <main className="flex-1 overflow-auto p-4">{children}</main>
    </div>
  )
}
