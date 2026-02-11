const ID_TOKEN_KEY = 'idToken'
const REFRESH_TOKEN_KEY = 'refreshToken'

export function getIdToken(): string | null {
  const token = localStorage.getItem(ID_TOKEN_KEY)
  return token && token.trim().length > 0 ? token : null
}

export function isAuthenticated(): boolean {
  return getIdToken() !== null
}

export function clearSession(): void {
  localStorage.removeItem(ID_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

export function persistSession(idToken: string, refreshToken: string): void {
  localStorage.setItem(ID_TOKEN_KEY, idToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
}
