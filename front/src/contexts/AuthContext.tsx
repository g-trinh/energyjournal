import { clearSession, getIdToken, isAuthenticated } from '@/lib/session'
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { useNavigate } from 'react-router-dom'

export type AuthStatus = 'loading' | 'authenticated' | 'anonymous'

interface AuthContextValue {
  status: AuthStatus
  email: string | null
  isAuthenticated: boolean
  signIn: (email: string) => void
  signOut: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

interface MeResponse {
  email?: unknown
}

export function AuthProvider({ children }: AuthProviderProps) {
  const navigate = useNavigate()
  const [status, setStatus] = useState<AuthStatus>('loading')
  const [email, setEmail] = useState<string | null>(null)

  useEffect(() => {
    const abortController = new AbortController()

    async function initAuth() {
      if (!isAuthenticated()) {
        setStatus('anonymous')
        setEmail(null)
        return
      }

      setStatus('loading')
      const token = getIdToken()

      if (!token) {
        setStatus('anonymous')
        setEmail(null)
        return
      }

      try {
        const response = await fetch('/api/users/me', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
          signal: abortController.signal,
        })

        if (response.status === 401) {
          clearSession()
          setStatus('anonymous')
          setEmail(null)
          return
        }

        if (!response.ok) {
          setStatus('anonymous')
          setEmail(null)
          return
        }

        const payload: MeResponse = await response.json()
        setEmail(typeof payload.email === 'string' ? payload.email : null)
        setStatus('authenticated')
      } catch {
        setStatus('anonymous')
        setEmail(null)
      }
    }

    void initAuth()

    return () => {
      abortController.abort()
    }
  }, [])

  const signIn = useCallback((email: string) => {
    setEmail(email)
    setStatus('authenticated')
  }, [])

  const signOut = useCallback(() => {
    clearSession()
    setStatus('anonymous')
    setEmail(null)
    navigate('/auth', { replace: true })
  }, [navigate])

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      email,
      isAuthenticated: status === 'authenticated',
      signIn,
      signOut,
    }),
    [status, email, signIn, signOut],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }

  return context
}
