"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { Shield, AlertCircle, CheckCircle, ArrowLeft } from "lucide-react"
import { apiClient } from "@/lib/api"
import { Alert, AlertDescription } from "@/components/ui/alert"

export default function VerifyOTPPage() {
  const [otp, setOtp] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const router = useRouter()
  const searchParams = useSearchParams()
  const email = searchParams.get("email")

  useEffect(() => {
    if (!email) {
      router.push("/register")
    }
  }, [email, router])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email) return

    setIsLoading(true)
    setError("")
    setSuccess("")

    try {
      await apiClient.verifyOTP({ email, otp })
      setSuccess("Email verified successfully! Redirecting to login...")
      setTimeout(() => {
        router.push("/login")
      }, 2000)
    } catch (error) {
      setError(error instanceof Error ? error.message : "Verification failed. Please try again.")
    } finally {
      setIsLoading(false)
    }
  }

  const handleOtpChange = (value: string) => {
    if (value.length <= 6) {
      setOtp(value)
    }
  }

  if (!email) {
    return null
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100 p-4">
      <Card className="w-full max-w-md shadow-xl border-0">
        <CardHeader className="text-center space-y-4">
          <div className="mx-auto w-12 h-12 bg-green-600 rounded-full flex items-center justify-center">
            <Shield className="w-6 h-6 text-white" />
          </div>
          <CardTitle className="text-2xl font-bold text-gray-900">Verify Your Email</CardTitle>
          <CardDescription className="text-gray-600">
            We sent a 6-digit code to <br />
            <span className="font-medium text-gray-800">{email}</span>
          </CardDescription>
        </CardHeader>
        <CardContent>
          {error && (
            <Alert className="mb-4 border-red-200 bg-red-50">
              <AlertCircle className="h-4 w-4 text-red-600" />
              <AlertDescription className="text-red-800">{error}</AlertDescription>
            </Alert>
          )}
          
          {success && (
            <Alert className="mb-4 border-green-200 bg-green-50">
              <CheckCircle className="h-4 w-4 text-green-600" />
              <AlertDescription className="text-green-800">{success}</AlertDescription>
            </Alert>
          )}
          
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="otp" className="text-sm font-medium text-gray-700 text-center block">
                Enter 6-digit code
              </Label>
              <Input
                id="otp"
                type="text"
                placeholder="000000"
                value={otp}
                onChange={(e) => handleOtpChange(e.target.value.replace(/\D/g, ""))}
                maxLength={6}
                required
                className="h-12 text-center text-lg font-mono tracking-widest border-gray-200 focus:border-green-500 focus:ring-green-500"
              />
              <p className="text-xs text-gray-500 text-center">
                Code expires in 10 minutes
              </p>
            </div>

            <Button
              type="submit"
              disabled={isLoading || otp.length !== 6}
              className="w-full h-11 bg-green-600 hover:bg-green-700 text-white font-medium"
            >
              {isLoading ? "Verifying..." : "Verify Email"}
            </Button>
          </form>

          <div className="mt-6 text-center space-y-3">
            <p className="text-sm text-gray-600">
              Didn't receive the code?{" "}
              <Button variant="link" className="p-0 h-auto font-medium text-green-600">
                Resend code
              </Button>
            </p>
            
            <Link 
              href="/register" 
              className="inline-flex items-center text-sm text-gray-600 hover:text-green-600 transition-colors"
            >
              <ArrowLeft className="w-4 h-4 mr-1" />
              Back to registration
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
