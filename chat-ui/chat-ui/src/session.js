// sessionStorage = one login per browser TAB, so two tabs can be
// two different users (Alice / Bob) at the same time.
const KEY = 'chat_tokens'
const listeners = new Set()

export const getTokens = () => {
  try {
    return JSON.parse(sessionStorage.getItem(KEY))
  } catch {
    return null
  }
}
export const setTokens = (t) => sessionStorage.setItem(KEY, JSON.stringify(t))
export const clearTokens = () => {
  sessionStorage.removeItem(KEY)
  listeners.forEach((f) => f())
}
export const onSessionEnd = (f) => {
  listeners.add(f)
  return () => listeners.delete(f)
}

// Access-token expiry in ms (0 if it cannot be read).
export function tokenExp(access) {
  try {
    const part = access.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(part)).exp * 1000
  } catch {
    return 0
  }
}
