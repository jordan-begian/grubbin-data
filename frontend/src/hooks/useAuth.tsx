import { createContext, useContext, useState, useCallback } from 'react'
import { useMutation } from '@tanstack/react-query'
import type { UserResponse, LoginRequest, RegisterUserRequest } from '@/types/auth'
import { loginUser, registerUser } from '@/services/auth'

interface AuthContextType {
  user: UserResponse | null
  isAuthenticated: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterUserRequest) => Promise<void>
  logout: () => void
  isLoading: boolean
  error: string | null
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

const STORAGE_KEY = 'grubbin_auth_user'

function getStoredUser(): UserResponse | null {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    return stored ? JSON.parse(stored) : null
  } catch {
    return null
  }
}

function setStoredUser(user: UserResponse | null): void {
  if (user) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(user))
  } else {
    localStorage.removeItem(STORAGE_KEY)
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserResponse | null>(getStoredUser)
  const [error, setError] = useState<string | null>(null)

  const loginMutation = useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      setUser(data)
      setStoredUser(data)
      setError(null)
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const registerMutation = useMutation({
    mutationFn: registerUser,
    onSuccess: (data) => {
      setUser(data)
      setStoredUser(data)
      setError(null)
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const login = useCallback(
    async (credentials: LoginRequest) => {
      setError(null)
      await loginMutation.mutateAsync(credentials)
    },
    [loginMutation]
  )

  const register = useCallback(
    async (data: RegisterUserRequest) => {
      setError(null)
      await registerMutation.mutateAsync(data)
    },
    [registerMutation]
  )

  const logout = useCallback(() => {
    setUser(null)
    setStoredUser(null)
    setError(null)
  }, [])

  const isLoading = loginMutation.isPending || registerMutation.isPending
  const isAuthenticated = user !== null

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        login,
        register,
        logout,
        isLoading,
        error,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
