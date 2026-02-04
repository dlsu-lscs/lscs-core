"use client"

import { useState, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { useAuth } from "@/hooks/use-auth"
import { api, RequestKeyRequest, RequestKeyResponse } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2, Copy, ArrowLeft, AlertCircle } from "lucide-react"
import { AuthenticatedLayout } from "@/components/authenticated-layout"

type KeyType = "dev" | "prod" | "admin"

export default function RequestKeyPage() {
  const { isAuthenticated, isLoading: authLoading, user } = useAuth()
  const router = useRouter()
  const searchParams = useSearchParams()
  const [loading, setLoading] = useState(true)
  const [authorized, setAuthorized] = useState(false)
  const [submitted, setSubmitted] = useState(false)
  const [keyResponse, setKeyResponse] = useState<RequestKeyResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const [formData, setFormData] = useState<RequestKeyRequest>({
    project: "",
    allowed_origin: "",
    is_dev: false,
    is_admin: false,
  })

  useEffect(() => {
    if (authLoading) return

    if (!isAuthenticated) {
      const redirect = searchParams.get("redirect") || "/request-key"
      window.location.href = `/auth/google/login?redirect=${encodeURIComponent(redirect)}`
      return
    }

    const isRND = user?.committee_id === "RND"
    const isAVPPlus = ["AVP", "VP", "EVP", "PRES"].includes(user?.position_id || "")
    const isAuthorized = isRND && isAVPPlus

    if (!isAuthorized) {
      setAuthorized(false)
      setLoading(false)
    } else {
      setAuthorized(true)
      setLoading(false)
    }
  }, [isAuthenticated, authLoading, user, searchParams])

  const handleKeyTypeChange = (type: KeyType) => {
    setFormData((prev) => ({
      ...prev,
      is_dev: type === "dev",
      is_admin: type === "admin",
      allowed_origin: type === "admin" ? "" : prev.allowed_origin,
    }))
    setError(null)
  }

  const validateForm = (): string | null => {
    if (formData.is_dev) {
      if (formData.allowed_origin && !formData.allowed_origin.startsWith("http://localhost")) {
        return "For dev keys, allowed_origin must start with http://localhost"
      }
    } else if (!formData.is_admin) {
      if (!formData.allowed_origin) {
        return "allowed_origin is required for production keys"
      }
      if (formData.allowed_origin.includes("localhost")) {
        return "localhost is not a valid origin for production keys"
      }
      try {
        new URL(formData.allowed_origin)
      } catch {
        return "Invalid URL for allowed_origin"
      }
    }
    return null
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    const validationError = validateForm()
    if (validationError) {
      setError(validationError)
      return
    }

    setLoading(true)
    try {
      const response = await api.requestKey(formData)
      setKeyResponse(response)
      setSubmitted(true)
    } catch (err: any) {
      setError(err.message || "Failed to create API key")
    } finally {
      setLoading(false)
    }
  }

  const handleCopyKey = () => {
    if (keyResponse) {
      navigator.clipboard.writeText(keyResponse.api_key)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  if (authLoading || (loading && !submitted)) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!authorized && !submitted) {
    return (
      <AuthenticatedLayout>
        <div className="container mx-auto p-6">
          <div className="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded flex items-start gap-2">
            <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0" />
            <div>
              <p className="font-semibold">Access Denied</p>
              <p className="text-sm">
                API key request is only available to RND committee members with AVP position or higher.
              </p>
            </div>
          </div>
        </div>
      </AuthenticatedLayout>
    )
  }

  if (submitted && keyResponse) {
    return (
      <AuthenticatedLayout>
        <div className="container mx-auto p-6 max-w-2xl">
          <Button
            variant="ghost"
            className="mb-6"
            onClick={() => router.push("/api-keys")}
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to API Keys
          </Button>

          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-2xl">API Key Generated</CardTitle>
              <CardDescription>
                Your new API key has been generated successfully
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="bg-muted p-4 rounded-lg">
                <p className="text-sm text-muted-foreground mb-2">API Key</p>
                <div className="flex items-center gap-2">
                  <code className="flex-1 font-mono text-sm break-all">
                    {keyResponse.api_key}
                  </code>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleCopyKey}
                  >
                    <Copy className="h-4 w-4" />
                    {copied ? "Copied!" : "Copy"}
                  </Button>
                </div>
              </div>

              {keyResponse.expires_at && (
                <p className="text-sm text-muted-foreground">
                  Expires: {new Date(keyResponse.expires_at).toLocaleDateString()}
                </p>
              )}

              <p className="text-sm text-destructive">
                Make sure to copy your API key now. You won&apos;t be able to see it again.
              </p>

              <Button
                className="w-full"
                onClick={() => router.push("/api-keys")}
              >
                View API Keys
              </Button>
            </CardContent>
          </Card>
        </div>
      </AuthenticatedLayout>
    )
  }

  return (
    <AuthenticatedLayout>
      <div className="container mx-auto p-6 max-w-2xl">
        <Button
          variant="ghost"
          className="mb-6"
          onClick={() => router.push("/api-keys")}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to API Keys
        </Button>

        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-2xl">Request API Key</CardTitle>
            <CardDescription>
              Create a new API key for your external project
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label className="text-sm font-medium">Project Name (Optional)</label>
                <Input
                  value={formData.project}
                  onChange={(e) =>
                    setFormData({ ...formData, project: e.target.value })
                  }
                  placeholder="My LSCS Project"
                  className="mt-1"
                />
              </div>

              <div>
                <label className="text-sm font-medium">Key Type</label>
                <div className="flex gap-4 mt-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="keyType"
                      checked={!formData.is_dev && !formData.is_admin}
                      onChange={() => handleKeyTypeChange("prod")}
                    />
                    <span className="text-sm">Production</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="keyType"
                      checked={formData.is_dev}
                      onChange={() => handleKeyTypeChange("dev")}
                    />
                    <span className="text-sm">Development</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="radio"
                      name="keyType"
                      checked={formData.is_admin}
                      onChange={() => handleKeyTypeChange("admin")}
                    />
                    <span className="text-sm">Admin</span>
                  </label>
                </div>
              </div>

              {!formData.is_admin && (
                <div>
                  <label className="text-sm font-medium">
                    {formData.is_dev ? "Allowed Origin (for dev, localhost only)" : "Allowed Origin (Required for production)"}
                  </label>
                  <Input
                    value={formData.allowed_origin}
                    onChange={(e) =>
                      setFormData({ ...formData, allowed_origin: e.target.value })
                    }
                    placeholder={formData.is_dev ? "http://localhost:3000" : "https://example.com"}
                    className="mt-1"
                  />
                  {formData.is_dev && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Must start with http://localhost
                    </p>
                  )}
                  {!formData.is_dev && !formData.is_admin && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Must be a valid HTTPS URL, no localhost
                    </p>
                  )}
                </div>
              )}

              {error && (
                <div className="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded flex items-start gap-2">
                  <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0" />
                  <p className="text-sm">{error}</p>
                </div>
              )}

              <Button
                type="submit"
                className="w-full"
                disabled={loading}
              >
                {loading ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Creating...
                  </>
                ) : (
                  "Generate API Key"
                )}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthenticatedLayout>
  )
}
