import type { ReactNode } from 'react'
import { Outlet } from 'react-router-dom'
import Topbar from './Topbar'

interface LayoutProps {
  children?: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <>
      <Topbar />
      {children ?? <Outlet />}
    </>
  )
}
