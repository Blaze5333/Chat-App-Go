"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Plus, Search, MessageCircle, MoreVertical, LogOut, Settings, User } from "lucide-react"
import { apiClient } from "@/lib/api"
import AddUserModal from "@/components/AddUserModal"

interface Chat {
  id: string
  name: string
  email: string
  lastMessage: string
  timestamp: string
  unreadCount: number
  avatar?: string
  roomId: string
}

export default function HomePage() {
  const router = useRouter()
  const [chats, setChats] = useState<Chat[]>([])
  const [searchQuery, setSearchQuery] = useState("")
  const [currentUser, setCurrentUser] = useState<any>(null)
  const [isAddUserModalOpen, setIsAddUserModalOpen] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userStr = localStorage.getItem('user')
    
    if (!token || !userStr) {
      window.location.href = '/login'
      return
    }
    
    try {
      const user = JSON.parse(userStr)
      setCurrentUser(user)
    } catch (error) {
      console.error('Error parsing user data:', error)
      window.location.href = '/login'
      return
    }
    
    fetchChats()
  }, [])

  const fetchChats = async () => {
    try {
      // Check if user is authenticated
      const token = localStorage.getItem('token')
      const userStr = localStorage.getItem('user')
      
      if (!token || !userStr) {
        window.location.href = '/login'
        return
      }

      const currentUser = JSON.parse(userStr)
      
      // Fetch conversations from your backend
      const response = await apiClient.getConversations()
      
      // Check if response.data exists and is an array
      if (!response.data || !Array.isArray(response.data)) {
        console.log('No conversations data received:', response)
        setChats([])
        return
      }
      
      // Transform the conversations data to match our Chat interface
      const transformedChats: Chat[] = response.data.map((conv) => {
        // Find the other participant (not the current user)
        const otherParticipant = conv.participants.find(p => p.id !== currentUser.id)
        
        // Use participant data from the conversation
        const participantName = otherParticipant?.username || 'Unknown User'
        const participantEmail = otherParticipant?.email || ''
        
        return {
          id: conv._id,
          roomId: conv._id,
          name: participantName,
          email: participantEmail,
          avatar: otherParticipant?.image,
          lastMessage: conv.last_message?.content || 'No messages yet',
          timestamp: conv.last_message 
            ? new Date(conv.last_message.created_at).toLocaleTimeString([], { 
                hour: '2-digit', 
                minute: '2-digit' 
              })
            : new Date(conv.created_at).toLocaleTimeString([], { 
                hour: '2-digit', 
                minute: '2-digit' 
              }),
          unreadCount: 0, // You can implement unread count logic later
        }
      })
      
      setChats(transformedChats)
    } catch (error) {
      console.error("Failed to fetch chats:", error)
      
      // Fallback to mock data if API fails
      const mockChats: Chat[] = [
        {
          id: "1",
          roomId: "mock-room-1",
          name: "John Doe",
          email: "john@example.com",
          lastMessage: "Hey, how are you doing?",
          timestamp: "2:30 PM",
          unreadCount: 2,
        },
        {
          id: "2",
          roomId: "mock-room-2",
          name: "Jane Smith",
          email: "jane@example.com",
          lastMessage: "Thanks for the help!",
          timestamp: "1:15 PM",
          unreadCount: 0,
        },
        {
          id: "3",
          roomId: "mock-room-3",
          name: "Mike Johnson",
          email: "mike@example.com",
          lastMessage: "See you tomorrow",
          timestamp: "12:45 PM",
          unreadCount: 1,
        },
      ]
      setChats(mockChats)
    }
  }

  const filteredChats = chats.filter(
    (chat) =>
      chat.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      chat.email.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const handleLogout = () => {
    // Clear authentication data
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    // Redirect to login page
    window.location.href = '/login'
  }

  const handleAddUser = () => {
    setIsAddUserModalOpen(true)
  }

  const handleUserAdded = () => {
    fetchChats()
  }

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
  }

  const handleChatClick = (chat: Chat) => {
    // Store room info for the chat page
    sessionStorage.setItem(`room_${chat.roomId}`, JSON.stringify({
      name: chat.name,
      email: chat.email
    }))
    
    // Navigate to chat room
    router.push(`/chat/${chat.roomId}`)
  }

  return (
    <div className="h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
              <MessageCircle className="w-4 h-4 text-white" />
            </div>
            <h1 className="text-xl font-semibold text-gray-900">Chats</h1>
          </div>
          <div className="flex items-center space-x-2">
            {currentUser && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="p-1">
                    <Avatar className="w-8 h-8">
                      <AvatarImage 
                        src={currentUser.image || "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"} 
                        alt={currentUser.username}
                      />
                      <AvatarFallback className="text-xs bg-blue-100 text-blue-600">
                        <User className="w-4 h-4" />
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem onClick={() => router.push('/settings')}>
                    <Settings className="w-4 h-4 mr-2" />
                    Settings
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={handleLogout} className="text-red-600">
                    <LogOut className="w-4 h-4 mr-2" />
                    Sign Out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
      </div>

      {/* Search Bar */}
      <div className="bg-white px-4 py-3 border-b border-gray-200">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <Input
            placeholder="Search chats..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 border-gray-200 focus:border-blue-500 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* Chat List */}
      <div className="flex-1 overflow-y-auto">
        {filteredChats.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <MessageCircle className="w-12 h-12 mb-4 text-gray-300" />
            <p className="text-lg font-medium mb-2">No chats yet</p>
            <p className="text-sm text-center px-8">Start a conversation by adding someone to your chats</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-100">
            {filteredChats.map((chat) => (
              <div 
                key={chat.id} 
                className="bg-white hover:bg-gray-50 px-4 py-4 cursor-pointer transition-colors"
                onClick={() => handleChatClick(chat)}
              >
                <div className="flex items-center space-x-3">
                  <Avatar className="w-12 h-12">
                    <AvatarImage 
                      src={chat.avatar || "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"} 
                      alt={chat.name}
                    />
                    <AvatarFallback className="bg-blue-100 text-blue-600 font-medium text-sm">
                      {getInitials(chat.name)}
                    </AvatarFallback>
                  </Avatar>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between mb-1">
                      <h3 className="font-medium text-gray-900 truncate">{chat.name}</h3>
                      <span className="text-xs text-gray-500">{chat.timestamp}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <p className="text-sm text-gray-600 truncate">{chat.lastMessage}</p>
                      {chat.unreadCount > 0 && (
                        <span className="bg-blue-600 text-white text-xs rounded-full px-2 py-1 min-w-[20px] text-center">
                          {chat.unreadCount}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Add Button */}
      <div className="absolute bottom-6 right-6">
        <Button
          onClick={handleAddUser}
          className="w-14 h-14 rounded-full bg-blue-600 hover:bg-blue-700 shadow-lg"
        >
          <Plus className="w-6 h-6" />
        </Button>
      </div>

      <AddUserModal
        open={isAddUserModalOpen}
        onOpenChange={setIsAddUserModalOpen}
        onUserAdded={handleUserAdded}
        existingChats={chats}
      />
    </div>
  )
}
