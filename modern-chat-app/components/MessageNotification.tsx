"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { X, Eye } from "lucide-react"

interface NotificationData {
  UserId: string
  Content: string
  Username: string
  Type: string
}

interface MessageNotificationProps {
  notification: NotificationData | null
  onClose: () => void
  onViewChat: (userId: string) => void
}

export default function MessageNotification({ notification, onClose, onViewChat }: MessageNotificationProps) {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    if (notification) {
      setIsVisible(true)
      console.log('Notification received from message file:', notification)
      // Auto-dismiss after 5 seconds
      const timer = setTimeout(() => {
        handleClose()
      }, 5000)
      
      return () => clearTimeout(timer)
    } else {
      setIsVisible(false)
    }
  }, [notification])

  const handleClose = () => {
    setIsVisible(false)
    setTimeout(() => {
      onClose()
    }, 300) // Wait for animation to complete
  }

  const handleViewChat = () => {
    if (notification) {
      onViewChat(notification.UserId)
      handleClose()
    }
  }

  const getInitials = (name: string) => {
    return name
     
  }

  const truncateContent = (content: string, maxLength: number = 50) => {
    if (content?.length <= maxLength) return content
    return content?.substring(0, maxLength) + "..."
  }

  if (!notification) return null

  return (
    <div 
      className={`fixed top-4 right-4 z-50 bg-white border border-gray-200 rounded-lg shadow-lg p-4 min-w-[320px] max-w-[400px] transition-all duration-300 ease-in-out ${
        isVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
      }`}
    >
      <div className="flex items-start space-x-3">
        <Avatar className="w-10 h-10 flex-shrink-0">
          <AvatarImage 
            src="https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png" 
            alt={notification.username}
          />
          <AvatarFallback className="bg-blue-100 text-blue-600 font-medium text-sm">
            {getInitials(notification.username)}
          </AvatarFallback>
        </Avatar>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-1">
            <h4 className="font-medium text-gray-900 text-sm truncate">
              {notification.username}
            </h4>
            <Button
              variant="ghost"
              size="sm"
              className="p-1 h-6 w-6 text-gray-400 hover:text-gray-600"
              onClick={handleClose}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
          
          <p className="text-sm text-gray-600 mb-3 break-words">
            {truncateContent(notification.content)}
          </p>
          
          <div className="flex space-x-2">
            <Button
              size="sm"
              className="flex-1 h-8 text-xs"
              onClick={handleViewChat}
            >
              <Eye className="w-3 h-3 mr-1" />
              View Chat
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-8 text-xs"
              onClick={handleClose}
            >
              Dismiss
            </Button>
          </div>
        </div>
      </div>
      
      {/* Progress bar for auto-dismiss */}
      <div 
        className="absolute bottom-0 left-0 h-1 bg-blue-600 rounded-b-lg w-full"
        style={{
          animation: 'shrink 5s linear forwards'
        }}
      />
      
      <style jsx>{`
        @keyframes shrink {
          from { width: 100%; }
          to { width: 0%; }
        }
      `}</style>
    </div>
  )
}
