import api from './index'

export interface RateLimitRule {
  route: string
  strategy: string
  rate: number
  burst: number
}

export function getRateLimits(): Promise<RateLimitRule[]> {
  return api.get('/rate-limits')
}

export function addRateLimit(rule: RateLimitRule): Promise<void> {
  return api.post('/rate-limits', rule)
}

export function deleteRateLimit(route: string): Promise<void> {
  return api.delete(`/rate-limits/${encodeURIComponent(route)}`)
}
