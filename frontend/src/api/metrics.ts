import api from './index'
import type { MetricsSummary } from '../types'

export function getMetricsSummary(): Promise<MetricsSummary> {
  return api.get('/metrics/summary')
}

export function getMetricsRaw(): Promise<string> {
  return api.get('/metrics', { responseType: 'text' })
}
