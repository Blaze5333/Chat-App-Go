"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog"
import { Search, UserPlus, AlertCircle, CheckCircle, User } from "lucide-react"
import { apiClient } from "@/lib/api"
import { Alert, AlertDescription } from "@/components/ui/alert"

interface AddUserModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onUserAdded: () => void
  existingChats: Array<{ email: string }>
  onlineUsers?: string[]
}

interface SearchedUser {
  id: string
  username: string
  email: string
  image?: string
  isOnline?: boolean
}

export default function AddUserModal({ open, onOpenChange, onUserAdded, existingChats, onlineUsers = [] }: AddUserModalProps) {
  const [searchEmail, setSearchEmail] = useState("")
  const [searchedUser, setSearchedUser] = useState<SearchedUser | null>(null)
  const [isSearching, setIsSearching] = useState(false)
  const [isAdding, setIsAdding] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")

  const handleSearch = async () => {
    if (!searchEmail.trim()) {
      setError("Please enter an email address")
      return
    }

    setIsSearching(true)
    setError("")
    setSearchedUser(null)

    try {
      const user = await apiClient.searchUser(searchEmail.trim())
      const userWithStatus = {
        ...user,
        isOnline: onlineUsers.includes(user.username)
      }
      setSearchedUser(userWithStatus)
    } catch (error) {
      setError(error instanceof Error ? error.message : "User not found")
      setSearchedUser(null)
    } finally {
      setIsSearching(false)
    }
  }

  const handleAddUser = async () => {
    if (!searchedUser) return

    setIsAdding(true)
    setError("")
    setSuccess("")

    try {
      await apiClient.createRoom(searchedUser.id)
      setSuccess(`${searchedUser.username} added to your conversations!`)
      
      setTimeout(() => {
        onUserAdded()
        onOpenChange(false)
        resetModal()
      }, 1500)
    } catch (error) {
      setError(error instanceof Error ? error.message : "Failed to add user")
    } finally {
      setIsAdding(false)
    }
  }

  const resetModal = () => {
    setSearchEmail("")
    setSearchedUser(null)
    setError("")
    setSuccess("")
    setIsSearching(false)
    setIsAdding(false)
  }

  const handleOpenChange = (open: boolean) => {
    onOpenChange(open)
    if (!open) {
      resetModal()
    }
  }

  const isUserAlreadyAdded = searchedUser ? existingChats.some(chat => chat.email === searchedUser.email) : false

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <UserPlus className="w-5 h-5 text-blue-600" />
            <span>Add New Contact</span>
          </DialogTitle>
          <DialogDescription>
            Search for a user by email to start a conversation
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
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

          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <div className="flex space-x-2">
              <Input
                id="email"
                type="email"
                placeholder="Enter email address"
                value={searchEmail}
                onChange={(e) => setSearchEmail(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                className="flex-1"
              />
              <Button 
                onClick={handleSearch} 
                disabled={isSearching || !searchEmail.trim()}
                size="sm"
                className="px-3"
              >
                {isSearching ? (
                  <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
                ) : (
                  <Search className="w-4 h-4" />
                )}
              </Button>
            </div>
          </div>

          {searchedUser && (
            <div className="p-4 border border-gray-200 rounded-lg bg-gray-50">
              <div className="flex items-center space-x-3 mb-3">
                <div className="relative">
                  <Avatar className="w-12 h-12">
                    <AvatarImage 
                      src={searchedUser.image || "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"} 
                      alt={searchedUser.username}
                    />
                    <AvatarFallback className="bg-blue-100 text-blue-600 font-medium">
                      <User className="w-5 h-5" />
                    </AvatarFallback>
                  </Avatar>
                  {searchedUser.isOnline ? (
                    <div className="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 border-2 border-white rounded-full"></div>
                  ) : (
                    <div className="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-gray-400 border-2 border-white rounded-full"></div>
                  )}
                </div>
                <div className="flex-1">
                  <div className="flex items-center space-x-2 mb-1">
                    <h3 className="font-medium text-gray-900">{searchedUser.username}</h3>
                    {searchedUser.isOnline ? (
                      <span className="text-xs text-green-600 font-medium bg-green-50 px-2 py-0.5 rounded-full">
                        Online
                      </span>
                    ) : (
                      <span className="text-xs text-gray-500 font-medium bg-gray-50 px-2 py-0.5 rounded-full">
                        Offline
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-600">{searchedUser.email}</p>
                </div>
              </div>
              
              <Button
                onClick={handleAddUser}
                disabled={isAdding || isUserAlreadyAdded}
                className="w-full"
                variant={isUserAlreadyAdded ? "secondary" : "default"}
              >
                {isAdding ? (
                  <>
                    <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2" />
                    Adding...
                  </>
                ) : isUserAlreadyAdded ? (
                  "Already in your conversations"
                ) : (
                  <>
                    <UserPlus className="w-4 h-4 mr-2" />
                    Add to Conversations
                  </>
                )}
              </Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
