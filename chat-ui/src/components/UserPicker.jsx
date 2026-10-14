import { useEffect, useState } from 'react'
import * as api from '../api'

export default function UserPicker({ excludeIds = [], onPick, placeholder = 'Search users…' }) {
  const [q, setQ] = useState('')
  const [results, setResults] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    if (!q.trim()) {
      setResults([])
      return undefined
    }
    const t = setTimeout(async () => {
      try {
        const list = await api.searchUsers(q)
        setResults(list.filter((u) => !excludeIds.includes(u.id)))
        setError('')
      } catch (e) {
        setError(e.message)
      }
    }, 300)
    return () => clearTimeout(t)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [q])

  return (
    <div className="picker">
      <input placeholder={placeholder} value={q} onChange={(e) => setQ(e.target.value)} />
      {error && <p className="error">{error}</p>}
      <ul>
        {results.map((u) => (
          <li key={u.id}>
            <button
              type="button"
              onClick={() => {
                onPick(u)
                setQ('')
                setResults([])
              }}
            >
              {u.username || u.email} <small>#{u.id}</small>
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
