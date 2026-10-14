import { useCallback, useEffect, useRef, useState } from 'react'
import { getTokens } from './session'
import { ensureFresh } from './api'
import { log } from './debug'

// One WebSocket per logged-in user. Token goes in the query string
// because browsers cannot set headers on a WebSocket handshake.
export function useSocket(enabled, onEvent) {
  const wsRef = useRef(null)
  const handlerRef = useRef(onEvent)
  handlerRef.current = onEvent
  const [status, setStatus] = useState('closed')

  useEffect(() => {
    if (!enabled) return undefined
    let stopped = false
    let timer

    function connect() {
      const t = getTokens()
      if (!t || stopped) return
      setStatus('connecting')
      const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
      const ws = new WebSocket(
        `${proto}://${window.location.host}/ws?token=${encodeURIComponent(t.access)}`,
      )
      wsRef.current = ws
      ws.onopen = () => {
        setStatus('open')
        log({ kind: 'ws', dir: 'open' })
      }
      ws.onmessage = (e) => {
        let data
        try {
          data = JSON.parse(e.data)
        } catch {
          data = e.data
        }
        log({ kind: 'ws', dir: 'in', data })
        handlerRef.current(data)
      }
      ws.onclose = () => {
        setStatus('closed')
        log({ kind: 'ws', dir: 'close' })
        if (stopped) return
        timer = setTimeout(async () => {
          try {
            await ensureFresh()
          } catch {
            return // session ended; App will show the login page
          }
          connect()
        }, 2000)
      }
    }

    connect()
    return () => {
      stopped = true
      clearTimeout(timer)
      wsRef.current?.close()
    }
  }, [enabled])

  const send = useCallback((obj) => {
    const ws = wsRef.current
    if (ws && ws.readyState === 1) {
      ws.send(JSON.stringify(obj))
      log({ kind: 'ws', dir: 'out', data: obj })
    }
  }, [])

  return { status, send }
}
