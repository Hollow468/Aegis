import api from './index'
import type { Route, RouteListResponse } from '../types'

export function getRoutes(): Promise<RouteListResponse> {
  return api.get('/routes')
}

export function getRoute(path: string): Promise<Route> {
  return api.get(`/routes/${encodeURIComponent(path)}`)
}

export function createRoute(route: Route): Promise<Route> {
  return api.post('/routes', route)
}

export function updateRoute(path: string, route: Route): Promise<Route> {
  return api.put(`/routes/${encodeURIComponent(path)}`, route)
}

export function deleteRoute(path: string): Promise<void> {
  return api.delete(`/routes/${encodeURIComponent(path)}`)
}
