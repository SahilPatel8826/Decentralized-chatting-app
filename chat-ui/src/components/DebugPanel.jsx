import { useEffect, useState } from 'react'
import { subscribe, clearLog } from '../debug'

const summary = (e) =>
  e.kind === 'http'
    ? `${e.method} ${e.url} → ${e.status}`
    : `WS ${e.dir}${e.data && e.data.type ? ' · ' + e.data.type : ''}`

export default function DebugPanel() {
  const [open, setOpen] = useState(false)
  const [entries, setEntries] = useState([])
  const [filter, setFilter] = useState('all')

  useEffect(() => subscribe(setEntries), [])

  const shown = entries.filter((e) => filter === 'all' || e.kind === filter)

  return (
    <div className={`debug ${open ? 'open' : ''}`}>
      <div className="debug-bar">
        <button onClick={() => setOpen(!open)}>Debug ({entries.length})</button>
        {open && (
          <>
            {['all', 'http', 'ws'].map((f) => (
              <button key={f} className={filter === f ? 'on' : ''} onClick={() => setFilter(f)}>
                {f}
              </button>
            ))}
            <button onClick={clearLog}>clear</button>
          </>
        )}
      </div>
      {open && (
        <div className="debug-list">
          {shown.map((e) => (
            <details key={e.id} className={e.kind === 'http' && e.status >= 400 ? 'bad' : ''}>
              <summary>
                <small>{e.t.toLocaleTimeString()}</small> {summary(e)}
              </summary>
              <pre>
                {JSON.stringify(e.kind === 'http' ? { request: e.body, response: e.data } : e.data ?? null, null, 2)}
              </pre>
            </details>
          ))}
        </div>
      )}
    </div>
  )
}
