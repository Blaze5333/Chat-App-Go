"use client"

import { useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"

export default function GoogleCallbackPage() {
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    const token = searchParams.get("token")
    const user = searchParams.get("user")
    const error = searchParams.get("error")

    if (error) {
      console.error("Google auth error:", error)
      router.push("/login?error=google_auth_failed")
      return
    }

    if (token && user) {
      try {
        const userData = JSON.parse(decodeURIComponent(user))
        localStorage.setItem("token", token)
        localStorage.setItem("user", JSON.stringify(userData))
        router.push("/home")
      } catch (error) {
        console.error("Error parsing user data:", error)
        router.push("/login?error=invalid_user_data")
      }
    } else {
      router.push("/login?error=missing_auth_data")
    }
  }, [router, searchParams])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p className="text-gray-600">Completing your sign in...</p>
      </div>
    </div>
  )
}
