import api from './index'
import type { CircuitBreakerStatus } from '../types'

export function getCircuitBreakers(): Promise<CircuitBreakerStatus[]> {
  return api.get('/circuit-breakers')
}

export function resetCircuitBreaker(route: string): Promise<void> {
  return api.put(`/circuit-breakers/${encodeURIComponent(route)}/reset`)
}
