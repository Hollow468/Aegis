import api from './index'
import type { UpstreamStatus } from '../types'

export function getUpstreams(): Promise<UpstreamStatus[]> {
  return api.get('/upstreams')
}

export function registerUpstream(service: string, address: string, weight: number): Promise<void> {
  return api.post('/upstreams', { service, address, weight })
}

export function deregisterUpstream(service: string, address: string): Promise<void> {
  return api.delete(`/upstreams/${encodeURIComponent(service)}/${encodeURIComponent(address)}`)
}
