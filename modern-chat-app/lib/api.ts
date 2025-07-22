const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface VerifyOTPRequest {
  email: string
  otp: string
}

export interface LoginResponse {
  message: string
  user: {
    id: string
    username: string
    email: string
    token: string
    image?: string
  }
}

export interface RegisterResponse {
  message: string
  user: {
    id: string
    username: string
    email: string
  }
}

export interface VerifyOTPResponse {
  message: string
}

export interface UploadImageResponse {
  message: string
  image_url: string
}

export interface SearchUserResponse {
  email: string
  username: string
  id: string
  image?: string
}

export interface ConversationResponse {
  message: string
  data: Conversation[]
}

export interface Participant {
  id: string
  username: string
  email: string
  image?: string
}

export interface Conversation {
  _id: string
  participants: Participant[]
  last_message?: Message | null
  created_at: string
  updated_at: string
  room_id: string
}

export interface Message {
  _id: string
  room_id: string
  username: string
  content: string
  user_id: string
  created_at: string
}

export interface ErrorResponse {
  error: string
  message: string
}

class ApiClient {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    // Add auth token if available
    const token = localStorage.getItem('token')
    if (token) {
      config.headers = {
        ...config.headers,
        'Authorization': `Bearer ${token}`,
      }
    }

    const response = await fetch(url, config)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error', message: 'An error occurred' }))
      throw new Error(errorData.message || errorData.error || 'Request failed')
    }

    return response.json()
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return this.request<LoginResponse>('/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async register(data: RegisterRequest): Promise<RegisterResponse> {
    return this.request<RegisterResponse>('/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async verifyOTP(data: VerifyOTPRequest): Promise<VerifyOTPResponse> {
    return this.request<VerifyOTPResponse>('/verify_otp', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async uploadImage(file: File): Promise<UploadImageResponse> {
    const url = `${this.baseUrl}/upload_image`
    const formData = new FormData()
    formData.append('file', file)
    
    const token = localStorage.getItem('token')
    const config: RequestInit = {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: formData,
    }

    const response = await fetch(url, config)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error', message: 'An error occurred' }))
      throw new Error(errorData.message || errorData.error || 'Upload failed')
    }

    return response.json()
  }

  getGoogleAuthUrl(): string {
    return `${this.baseUrl}/auth/google`
  }

  async searchUser(email: string): Promise<SearchUserResponse> {
    return this.request<SearchUserResponse>(`/users/search?email=${encodeURIComponent(email)}`, {
      method: 'GET',
    })
  }

  async getConversations(): Promise<ConversationResponse> {
    return this.request<ConversationResponse>('/conversation', {
      method: 'GET',
    })
  }

  async createRoom(userId: string): Promise<any> {
    return this.request<any>(`/create_room/${userId}`, {
      method: 'POST',
    })
  }

  async getRoomMessages(roomId: string): Promise<any> {
    return this.request<any>(`/get_room_messages/${roomId}`, {
      method: 'GET',
    })
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
