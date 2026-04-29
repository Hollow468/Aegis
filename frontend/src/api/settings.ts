import api from './index'

export interface SystemSettings {
  server: { port: number; host: string; readTimeout: string; writeTimeout: string }
  etcd: { endpoints: string; prefix: string }
  redis: { addr: string; db: number }
  jwt: { expire: string }
  log: { level: string; format: string }
}

export function getSettings(): Promise<SystemSettings> {
  return api.get('/settings')
}

export function updateSettings(settings: Partial<SystemSettings>): Promise<void> {
  return api.put('/settings', settings)
}
