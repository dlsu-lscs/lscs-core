"use client"

import { useEffect } from "react"
import { useAuth } from "@/hooks/use-auth"
import { Loader2 } from "lucide-react"

export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading) {
      if (isAuthenticated) {
        window.location.href = "/dashboard"
      } else {
        window.location.href = "/login"
      }
    }
  }, [isAuthenticated, isLoading])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <Loader2 className="h-8 w-8 animate-spin text-primary" />
    </div>
  )
}
