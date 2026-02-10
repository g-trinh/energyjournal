/* ═══════════════════════════════════════════════════════════════════════════
   Auth API Client — typed service layer for /users* endpoints
   ═══════════════════════════════════════════════════════════════════════════ */

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api'

// ── Request types ──────────────────────────────────────────────────────────

export interface CreateUserRequest {
  email: string
  password: string
  confirmPassword: string
  timezone?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RefreshRequest {
  refreshToken: string
}

// ── Response types ─────────────────────────────────────────────────────────

export interface AuthTokensResponse {
  idToken: string
  refreshToken: string
  expiresIn: string
}

export interface CreateUserAcceptedResponse {
  message: string
  status: string
}

export interface ActivationResponse {
  message: string
}

// ── Normalized result ──────────────────────────────────────────────────────

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: string }

// ── Service methods ────────────────────────────────────────────────────────

async function request<T>(
  path: string,
  options: RequestInit,
): Promise<ApiResult<T>> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    })

    if (res.ok) {
      const data: T = await res.json()
      return { ok: true, data }
    }

    // Return generic error — never expose backend detail
    return { ok: false, error: 'request_failed' }
  } catch {
    return { ok: false, error: 'network_error' }
  }
}

export function createUser(
  body: CreateUserRequest,
): Promise<ApiResult<CreateUserAcceptedResponse>> {
  return request<CreateUserAcceptedResponse>('/users', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function login(
  body: LoginRequest,
): Promise<ApiResult<AuthTokensResponse>> {
  return request<AuthTokensResponse>('/users/login', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function activateUser(
  token: string,
): Promise<ApiResult<ActivationResponse>> {
  return request<ActivationResponse>(
    `/users/activate?token=${encodeURIComponent(token)}`,
    { method: 'POST' },
  )
}

export function refreshTokens(
  body: RefreshRequest,
): Promise<ApiResult<AuthTokensResponse>> {
  return request<AuthTokensResponse>('/users/refresh', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
