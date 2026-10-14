// Tiny event log shown in the Debug panel (HTTP calls + WebSocket frames).
let entries = []
let counter = 0
const subs = new Set()

export function log(entry) {
  entries = [{ id: ++counter, t: new Date(), ...entry }, ...entries].slice(0, 300)
  subs.forEach((f) => f(entries))
}
export function subscribe(f) {
  subs.add(f)
  f(entries)
  return () => subs.delete(f)
}
export function clearLog() {
  entries = []
  subs.forEach((f) => f(entries))
}
