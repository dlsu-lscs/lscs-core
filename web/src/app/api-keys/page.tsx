"use client"

import { useState } from "react"
import { useAuth } from "@/hooks/use-auth"
import { api, APIKey } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { Loader2, Key, Trash2, Plus, AlertCircle } from "lucide-react"
import { useRouter } from "next/navigation"
import { AuthenticatedLayout } from "@/components/authenticated-layout"

export default function APIKeysPage() {
  const { user } = useAuth()
  const router = useRouter()
  const queryClient = useQueryClient()
  const [keyToRevoke, setKeyToRevoke] = useState<APIKey | null>(null)

  const isRND = user?.committee_id === "RND"

  const { data: apiKeys, isLoading, error } = useQuery({
    queryKey: ["api-keys"],
    queryFn: api.listApiKeys,
    enabled: isRND,
  })

  const revokeMutation = useMutation({
    mutationFn: (id: number) => api.revokeApiKey(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["api-keys"] })
      setKeyToRevoke(null)
    },
  })

  if (!isRND) {
    return (
      <AuthenticatedLayout>
        <div className="container mx-auto p-6">
          <div className="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded flex items-start gap-2">
            <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0" />
            <div>
              <p className="font-semibold">Access Denied</p>
              <p className="text-sm">API Key management is only available to RND committee members.</p>
            </div>
          </div>
        </div>
      </AuthenticatedLayout>
    )
  }

  const handleRevoke = () => {
    if (keyToRevoke) {
      revokeMutation.mutate(keyToRevoke.api_key_id)
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    })
  }

  const getKeyTypeLabel = (key: APIKey) => {
    if (key.is_admin) {
      return <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">Admin</span>
    } else if (key.is_dev) {
      return <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">Development</span>
    } else {
      return <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">Production</span>
    }
  }

  return (
    <AuthenticatedLayout>
      <div className="container mx-auto p-6 max-w-4xl">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">API Keys</h1>
            <p className="text-muted-foreground mt-1">
              Manage your API keys for external projects
            </p>
          </div>
          <Button onClick={() => router.push("/request-key")}>
            <Plus className="mr-2 h-4 w-4" />
            Create API Key
          </Button>
        </div>

        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : error ? (
          <div className="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded flex items-start gap-2">
            <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0" />
            <div>
              <p className="font-semibold">Error</p>
              <p className="text-sm">Failed to load API keys. Please try again later.</p>
            </div>
          </div>
        ) : apiKeys && apiKeys.length > 0 ? (
          <div className="space-y-4">
            {apiKeys.map((key) => (
              <Card key={key.api_key_id}>
                <CardHeader className="pb-3">
                  <div className="flex justify-between items-start">
                    <div className="flex items-center gap-2">
                      <Key className="h-5 w-5 text-muted-foreground" />
                      <CardTitle className="text-lg">
                        {key.project || "Unnamed Project"}
                      </CardTitle>
                    </div>
                    {getKeyTypeLabel(key)}
                  </div>
                  <CardDescription>
                    Created on {formatDate(key.created_at)}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2 text-sm">
                    {key.allowed_origin && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Allowed Origin:</span>
                        <span className="font-mono">{key.allowed_origin}</span>
                      </div>
                    )}
                    {key.expires_at && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Expires:</span>
                        <span>{formatDate(key.expires_at)}</span>
                      </div>
                    )}
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Key ID:</span>
                      <span className="font-mono text-xs">{key.api_key_id}</span>
                    </div>
                  </div>
                  <div className="mt-4 flex justify-end">
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => setKeyToRevoke(key)}
                    >
                      <Trash2 className="mr-2 h-4 w-4" />
                      Revoke
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <CardContent className="flex flex-col items-center justify-center h-64">
              <Key className="h-12 w-12 text-muted-foreground mb-4" />
              <p className="text-muted-foreground text-center mb-4">
                You haven&apos;t created any API keys yet.
              </p>
              <Button onClick={() => router.push("/request-key")}>
                <Plus className="mr-2 h-4 w-4" />
                Create Your First API Key
              </Button>
            </CardContent>
          </Card>
        )}

        {keyToRevoke && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <Card className="w-full max-w-md">
              <CardHeader>
                <CardTitle>Revoke API Key</CardTitle>
                <CardDescription>
                  Are you sure you want to revoke the API key for &quot;{keyToRevoke?.project || "Unnamed Project"}? 
                  This action cannot be undone. Applications using this key will stop working immediately.
                </CardDescription>
              </CardHeader>
              <CardContent className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setKeyToRevoke(null)}>
                  Cancel
                </Button>
                <Button
                  variant="destructive"
                  onClick={handleRevoke}
                  disabled={revokeMutation.isPending}
                >
                  {revokeMutation.isPending ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Revoking...
                    </>
                  ) : (
                    "Revoke Key"
                  )}
                </Button>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </AuthenticatedLayout>
  )
}
