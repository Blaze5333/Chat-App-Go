"use client"

import { useState, useRef, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useRouter } from "next/navigation"
import { Camera, Upload, ArrowLeft, User, AlertCircle, CheckCircle, LogOut } from "lucide-react"
import { apiClient } from "@/lib/api"
import { Alert, AlertDescription } from "@/components/ui/alert"

interface UserData {
  id: string
  username: string
  email: string
  image?: string
}

export default function SettingsPage() {
  const [userData, setUserData] = useState<UserData | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string>("")
  const [isUploading, setIsUploading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const fileInputRef = useRef<HTMLInputElement>(null)
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userStr = localStorage.getItem('user')
    
    if (!token || !userStr) {
      router.push('/login')
      return
    }

    try {
      const user = JSON.parse(userStr)
      setUserData(user)
      if (user.image) {
        setPreviewUrl(user.image)
      }
    } catch (error) {
      console.error('Error parsing user data:', error)
      router.push('/login')
    }
  }, [router])

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      if (file.size > 5 * 1024 * 1024) {
        setError("File size must be less than 5MB")
        return
      }
      
      if (!file.type.startsWith('image/')) {
        setError("Please select an image file")
        return
      }

      setSelectedFile(file)
      const url = URL.createObjectURL(file)
      setPreviewUrl(url)
      setError("")
    }
  }

  const handleUpload = async () => {
    if (!selectedFile) return

    setIsUploading(true)
    setError("")
    setSuccess("")

    try {
      const response = await apiClient.uploadImage(selectedFile)
      setSuccess("Profile image updated successfully!")
      
      if (userData) {
        const updatedUser = { ...userData, image: response.image_url }
        setUserData(updatedUser)
        localStorage.setItem('user', JSON.stringify(updatedUser))
      }
      
      setSelectedFile(null)
    } catch (error) {
      setError(error instanceof Error ? error.message : "Upload failed. Please try again.")
    } finally {
      setIsUploading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
  }

  if (!userData) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  const defaultImage = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 p-4">
      <div className="max-w-2xl mx-auto">
        <div className="mb-6">
          <Button
            onClick={() => router.push('/home')}
            variant="ghost"
            className="mb-4 text-gray-600 hover:text-gray-900"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Chat
          </Button>
        </div>

        <Card className="shadow-xl border-0">
          <CardHeader className="text-center space-y-4">
            <CardTitle className="text-2xl font-bold text-gray-900">Settings</CardTitle>
            <CardDescription className="text-gray-600">
              Manage your profile and preferences
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {error && (
              <Alert className="border-red-200 bg-red-50">
                <AlertCircle className="h-4 w-4 text-red-600" />
                <AlertDescription className="text-red-800">{error}</AlertDescription>
              </Alert>
            )}
            
            {success && (
              <Alert className="border-green-200 bg-green-50">
                <CheckCircle className="h-4 w-4 text-green-600" />
                <AlertDescription className="text-green-800">{success}</AlertDescription>
              </Alert>
            )}

            <div className="text-center space-y-4">
              <div className="relative inline-block">
                <Avatar className="w-24 h-24 border-4 border-white shadow-lg">
                  <AvatarImage 
                    src={previewUrl || userData.image || defaultImage} 
                    alt={userData.username}
                  />
                  <AvatarFallback className="text-2xl bg-blue-100 text-blue-600">
                    <User className="w-8 h-8" />
                  </AvatarFallback>
                </Avatar>
                <Button
                  onClick={() => fileInputRef.current?.click()}
                  size="sm"
                  className="absolute -bottom-2 -right-2 w-8 h-8 rounded-full p-0 bg-blue-600 hover:bg-blue-700"
                >
                  <Camera className="w-4 h-4" />
                </Button>
              </div>
              
              <div>
                <h3 className="text-lg font-semibold text-gray-900">{userData.username}</h3>
                <p className="text-sm text-gray-600">{userData.email}</p>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="file-upload" className="text-sm font-medium text-gray-700">
                  Profile Picture
                </Label>
                <Input
                  ref={fileInputRef}
                  id="file-upload"
                  type="file"
                  accept="image/*"
                  onChange={handleFileSelect}
                  className="hidden"
                />
                <div className="mt-2 flex items-center space-x-3">
                  <Button
                    onClick={() => fileInputRef.current?.click()}
                    variant="outline"
                    className="flex-1"
                  >
                    <Upload className="w-4 h-4 mr-2" />
                    Choose Image
                  </Button>
                  {selectedFile && (
                    <Button
                      onClick={handleUpload}
                      disabled={isUploading}
                      className="flex-1 bg-blue-600 hover:bg-blue-700"
                    >
                      {isUploading ? "Uploading..." : "Save Image"}
                    </Button>
                  )}
                </div>
                <p className="mt-1 text-xs text-gray-500">
                  PNG, JPG up to 5MB
                </p>
              </div>
            </div>

            <div className="pt-6 border-t border-gray-200">
              <Button
                onClick={handleLogout}
                variant="outline"
                className="w-full border-red-200 text-red-600 hover:bg-red-50 hover:border-red-300"
              >
                <LogOut className="w-4 h-4 mr-2" />
                Sign Out
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
