"use client"

import { createContext, useContext, useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import MessageNotification from "@/components/MessageNotification"

interface NotificationData {
  UserId: string
  Content: string
  Username: string
  Type: string
}

interface NotificationContextType {
  onlineRef: React.MutableRefObject<WebSocket | null>
  setupNotificationSocket: (userId: string) => void
  closeNotificationSocket: () => void
  isConnected: boolean
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined)

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const [notification, setNotification] = useState<NotificationData | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const onlineRef = useRef<WebSocket | null>(null)
  const currentUserIdRef = useRef<string | null>(null)
  const router = useRouter()

  const setupNotificationSocket = (userId: string) => {
    if (!userId) return

    // Avoid duplicate connections for the same user
    if (currentUserIdRef.current === userId && isConnected && onlineRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected for user:', userId)
      return
    }

    currentUserIdRef.current = userId

    // Close existing connection if any
    if (onlineRef.current && onlineRef.current.readyState !== WebSocket.CLOSED) {
      onlineRef.current.close(1000, 'Reconnecting')
    }

    const onlineUrl = `${process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080'}/ws/join_app?user_id=${userId}`
    
    try {
      onlineRef.current = new WebSocket(onlineUrl)
      
      onlineRef.current.onopen = () => {
        console.log('Notification WebSocket connected')
        setIsConnected(true)
      }
      
      onlineRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('Notification received:', data)
          
          if (data.type === "notification") {
            setNotification(data)
          }
          
          // Dispatch online status updates for home page
          window.dispatchEvent(new CustomEvent('onlineStatusUpdate', { detail: data }))
        } catch (error) {
          console.error('Error parsing notification:', error)
        }
      }
      
      onlineRef.current.onerror = (error) => {
        console.warn('Notification WebSocket error - this is normal during navigation:', error.type || 'connection error')
        setIsConnected(false)
      }
      
      onlineRef.current.onclose = (event) => {
        console.log('Notification WebSocket closed:', event.code, event.reason)
        setIsConnected(false)
        
        // Only attempt to reconnect if it wasn't a manual close and the user is still logged in
        if (event.code !== 1000 && event.code !== 1001) {
          const token = localStorage.getItem('token')
          if (token && userId && currentUserIdRef.current === userId) {
            console.log('Attempting to reconnect in 3 seconds...')
            setTimeout(() => {
              const currentToken = localStorage.getItem('token')
              if (currentToken && currentUserIdRef.current === userId) {
                setupNotificationSocket(userId)
              }
            }, 3000)
          }
        }
      }
    } catch (error) {
      console.error('Error setting up notification socket:', error)
      setIsConnected(false)
    }
  }

  const closeNotificationSocket = () => {
    setIsConnected(false)
    currentUserIdRef.current = null
    if (onlineRef.current && onlineRef.current.readyState === WebSocket.OPEN) {
      onlineRef.current.close(1000, 'Manual close') // 1000 = normal closure
      onlineRef.current = null
    }
  }

  const handleViewChat = (userId: string) => {
    const chats = JSON.parse(localStorage.getItem('chats') || '[]')
    const chat = chats.find((c: any) => c.userId === userId)
    
    if (chat) {
      sessionStorage.setItem(`room_${chat.roomId}`, JSON.stringify({
        name: chat.name,
        email: chat.email
      }))
      router.push(`/chat/${chat.roomId}`)
    }
  }

  const handleCloseNotification = () => {
    setNotification(null)
  }

  useEffect(() => {
    // Handle page unload/tab close
    const handleBeforeUnload = () => {
      closeNotificationSocket()
    }
    
    window.addEventListener('beforeunload', handleBeforeUnload)
    
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
      closeNotificationSocket()
    }
  }, [])

  return (
    <NotificationContext.Provider value={{ 
      onlineRef, 
      setupNotificationSocket, 
      closeNotificationSocket,
      isConnected
    }}>
      {children}
      <MessageNotification 
        notification={notification}
        onClose={handleCloseNotification}
        onViewChat={handleViewChat}
      />
    </NotificationContext.Provider>
  )
}

export function useNotification() {
  const context = useContext(NotificationContext)
  if (context === undefined) {
    throw new Error('useNotification must be used within a NotificationProvider')
  }
  return context
}
