import { request } from './client'

export async function getSummary(type = 'month') {
  const data = await request(`/summary?type=${type}`)
  return data?.summaries ?? []
}