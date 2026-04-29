import { ref, onMounted, onUnmounted } from 'vue'

export function usePolling(fn: () => void | Promise<void>, intervalMs: number) {
  const active = ref(false)
  let timer: ReturnType<typeof setInterval> | null = null

  function start() {
    if (timer) return
    active.value = true
    fn()
    timer = setInterval(() => { fn() }, intervalMs)
  }

  function stop() {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    active.value = false
  }

  onMounted(start)
  onUnmounted(stop)

  return { active, start, stop }
}
