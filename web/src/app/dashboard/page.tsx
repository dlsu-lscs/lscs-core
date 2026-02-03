"use client"

import { useAuth } from "@/hooks/use-auth"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { User, Mail, Phone, MapPin, School } from "lucide-react"
import Link from "next/link"

export default function DashboardPage() {
  const { user } = useAuth()

  if (!user) {
    return null
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Welcome back, {user.nickname || user.full_name}
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Profile</CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-sm">
                <span className="font-medium">{user.full_name}</span>
                {user.nickname && (
                  <span className="text-muted-foreground">
                    ({user.nickname})
                  </span>
                )}
              </div>
              {user.position_id && (
                <p className="text-sm text-muted-foreground">
                  {user.position_name}
                </p>
              )}
              {user.committee_id && (
                <p className="text-sm text-muted-foreground">
                  {user.committee_name}
                </p>
              )}
            </div>
            <Link href="/profile">
              <Button variant="outline" className="w-full mt-4">
                View Profile
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Contact</CardTitle>
            <Mail className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p className="text-sm">{user.email}</p>
              {user.telegram && (
                <p className="text-sm text-muted-foreground">
                  @{user.telegram}
                </p>
              )}
              {user.discord && (
                <p className="text-sm text-muted-foreground">{user.discord}</p>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Academic</CardTitle>
            <School className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {user.college && (
                <p className="text-sm text-muted-foreground">{user.college}</p>
              )}
              {user.program && (
                <p className="text-sm text-muted-foreground">{user.program}</p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
