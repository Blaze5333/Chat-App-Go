"use client"

import { useState, useEffect, useRef } from "react"
import { useParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ArrowLeft, Send, MoreVertical } from "lucide-react"
import { apiClient, Message } from "@/lib/api"

interface ChatMessage {
  _id: string
  room_id: string
  username: string
  content: string
  user_id: string
  created_at: string
  isOwn: boolean
}

export default function ChatPage() {
  const params = useParams()
  const router = useRouter()
  const roomId = params.roomId as string
  
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [newMessage, setNewMessage] = useState("")
  const [isConnected, setIsConnected] = useState(false)
  const [roomInfo, setRoomInfo] = useState<{name: string, email: string} | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  
  const wsRef = useRef<WebSocket | null>(null)

  const messagesEndRef = useRef<HTMLDivElement>(null)
  const currentUser = useRef<any>(null)

  useEffect(() => {
    // Check authentication
    const token = localStorage.getItem('token')
    const userStr = localStorage.getItem('user')
    
    if (!token || !userStr) {
      router.push('/login')
      return
    }
    
    currentUser.current = JSON.parse(userStr)
    
    // Get room info from session storage (passed from home page)
    const storedRoomInfo = sessionStorage.getItem(`room_${roomId}`)
    if (storedRoomInfo) {
      setRoomInfo(JSON.parse(storedRoomInfo))
    }
    
    // Load initial messages and connect to websocket
    initializeChat()
    
    // Cleanup on unmount
    return () => {
      if (wsRef.current) {
        console.log('Closing WebSocket connection')
        wsRef.current.close()
      }
    }
  }, [roomId, router])

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const initializeChat = async () => {
    try {
      setIsLoading(true)
      console.log('Initializing chat for room:', roomId)
      
      // First, get room messages
      await loadRoomMessages()
      
      // Then connect to websocket
      connectWebSocket()
      
    } catch (error) {
      console.error('Error initializing chat:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const loadRoomMessages = async () => {
    try {
      console.log('Loading room messages for room:', roomId)
      const response = await apiClient.getRoomMessages(roomId)
      console.log('Room messages response:', response)
      
      if (response.data && Array.isArray(response.data)) {
        const chatMessages: ChatMessage[] = response.data.map((msg: any) => ({
          _id: msg._id,
          room_id: msg.room_id,
          content: msg.content,
          username: msg.username,
          user_id: msg.user_id,
          created_at: msg.created_at,
          isOwn: msg.user_id === currentUser.current?.id
        }))
        console.log('Transformed chat messages:', chatMessages)
        setMessages(chatMessages)
      } else {
        console.log('No messages found in response')
        setMessages([])
      }
      
      // Extract room info from the response if available
      if (response.room_info) {
        setRoomInfo(response.room_info)
      }
    } catch (error) {
      console.error('Error loading room messages:', error)
      setMessages([])
    }
  }


  const connectWebSocket = () => {
    const userId = currentUser.current?.id
    const username = currentUser.current?.username || ""
    const wsUrl = `${process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080'}/join_room/${roomId}?user_id=${userId}&username=${username}`
   
    console.log('Connecting to WebSocket:', wsUrl)
    
    try {
      wsRef.current = new WebSocket(wsUrl)
    
      wsRef.current.onopen = () => {
        console.log('WebSocket connected successfully')
        setIsConnected(true)
      }
      
      wsRef.current.onmessage = (event) => {
        console.log('Raw WebSocket message received:', event.data)
        try {
          const message = JSON.parse(event.data)
          console.log('Parsed message:', message)
          const chatMessage: ChatMessage = {
            _id: message._id,
            room_id: message.room_id,
            content: message.content,
            username: message.username,
            user_id: message.user_id,
            created_at: message.created_at,
            isOwn: message.user_id === currentUser.current?.id
          }
          
          console.log('Adding message to chat:', chatMessage)
          setMessages(prev => {
            console.log('Previous messages:', prev)
            const newMessages = [...prev, chatMessage]
            console.log('New messages array:', newMessages)
            return newMessages
          })
        } catch (error) {
          console.error('Error parsing message:', error)
        }
      }
      
      wsRef.current.onclose = (event) => {
        console.log('WebSocket disconnected. Code:', event.code, 'Reason:', event.reason)
        setIsConnected(false)
        
        // Attempt to reconnect after 3 seconds
        setTimeout(() => {
          if (wsRef.current?.readyState === WebSocket.CLOSED) {
            console.log('Attempting to reconnect WebSocket...')
            connectWebSocket()
          }
        }, 3000)
      }
      
      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error)
        setIsConnected(false)
      }
    } catch (error) {
      console.error('Error creating WebSocket connection:', error)
    }
  }

  const sendMessage = () => {
    if (!newMessage.trim() || !wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.log('Cannot send message. Message empty:', !newMessage.trim(), 'WebSocket ready:', wsRef.current?.readyState === WebSocket.OPEN)
      return
    }
    
   
    
    try {
      wsRef.current.send(newMessage.trim())
      setNewMessage("")
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  const formatTime = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
  }

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">Loading chat...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.back()}
              className="p-2"
            >
              <ArrowLeft className="w-5 h-5" />
            </Button>
            
            <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
              <span className="text-blue-600 font-medium text-sm">
                {getInitials(roomInfo?.name || 'Unknown User')}
              </span>
            </div>
            
            <div>
              <h1 className="text-lg font-semibold text-gray-900">
                {roomInfo?.name || 'Chat Room'}
              </h1>
              <div className="flex items-center space-x-2">
                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-sm text-gray-500">
                  {isConnected ? 'Online' : 'Disconnected'}
                </span>
              </div>
            </div>
          </div>
          
          <Button variant="ghost" size="sm">
            <MoreVertical className="w-5 h-5" />
          </Button>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <p className="text-gray-500">No messages yet. Start the conversation!</p>
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message._id}
              className={`flex ${message.isOwn ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[70%] rounded-lg px-4 py-2 ${
                  message.isOwn
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-900 border border-gray-200'
                }`}
              >
                {!message.isOwn && (
                  <p className="text-sm font-medium text-gray-600 mb-1">
                    {message.username}
                  </p>
                )}
                <p className="text-sm">{message.content}</p>
                <p
                  className={`text-xs mt-1 ${
                    message.isOwn ? 'text-blue-100' : 'text-gray-500'
                  }`}
                >
                  {formatTime(message.created_at)}
                </p>
              </div>
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Message Input */}
      <div className="bg-white border-t border-gray-200 p-4">
        <div className="flex items-center space-x-3">
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            className="flex-1 border-gray-200 focus:border-blue-500 focus:ring-blue-500"
            disabled={!isConnected}
          />
          <Button
            onClick={sendMessage}
            disabled={!newMessage.trim() || !isConnected}
            className="bg-blue-600 hover:bg-blue-700"
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
        {!isConnected && (
          <p className="text-sm text-red-500 mt-2">
            Disconnected. Attempting to reconnect...
          </p>
        )}
      </div>
    </div>
  )
}
