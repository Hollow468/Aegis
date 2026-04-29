import { ref } from 'vue'

export function useConfirm() {
  const visible = ref(false)
  const title = ref('')
  const message = ref('')
  let resolveFn: ((value: boolean) => void) | null = null

  function confirm(t: string, msg: string): Promise<boolean> {
    title.value = t
    message.value = msg
    visible.value = true
    return new Promise((resolve) => {
      resolveFn = resolve
    })
  }

  function handleConfirm() {
    visible.value = false
    resolveFn?.(true)
    resolveFn = null
  }

  function handleCancel() {
    visible.value = false
    resolveFn?.(false)
    resolveFn = null
  }

  return { visible, title, message, confirm, handleConfirm, handleCancel }
}
