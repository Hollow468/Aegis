import api from './index'
import type { LoginRequest, LoginResponse } from '../types'

export function login(data: LoginRequest): Promise<LoginResponse> {
  return api.post('/auth/login', data)
}

export function refreshToken(): Promise<LoginResponse> {
  return api.post('/auth/refresh')
}
