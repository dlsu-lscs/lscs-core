"use client"

import { useEffect } from "react"
import { useAuth } from "@/hooks/use-auth"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"

export default function LoginPage() {
  const { isAuthenticated, isLoading, login } = useAuth()

  useEffect(() => {
    if (isAuthenticated && !isLoading) {
      window.location.href = "/dashboard"
    }
  }, [isAuthenticated, isLoading])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/10 to-secondary/10 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Welcome to LSCS Core</CardTitle>
          <CardDescription>
            Sign in with your DLSU email to continue
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button
            onClick={() => login()}
            className="w-full"
            size="lg"
          >
            Sign in with Google
          </Button>
          <p className="text-center text-sm text-muted-foreground">
            Only DLSU email addresses are allowed
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
