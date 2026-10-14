import { useState } from 'react'
import * as api from '../api'

export default function AuthPage({ onAuthed }) {
  const [mode, setMode] = useState('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  const submit = async (e) => {
    e.preventDefault()
    setError('')
    setBusy(true)
    try {
      if (mode === 'register') {
        await api.register(username, email, password)
      }
      await api.login(email, password)
      onAuthed()
    } catch (err) {
      setError(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="auth">
      <form className="card" onSubmit={submit}>
        <h1>Chat test client</h1>
        <div className="tabs">
          <button type="button" className={mode === 'login' ? 'on' : ''} onClick={() => setMode('login')}>
            Log in
          </button>
          <button type="button" className={mode === 'register' ? 'on' : ''} onClick={() => setMode('register')}>
            Register
          </button>
        </div>
        {mode === 'register' && (
          <input placeholder="username" value={username} onChange={(e) => setUsername(e.target.value)} required />
        )}
        <input type="email" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        <input
          type="password"
          placeholder="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        {error && <p className="error">{error}</p>}
        <button className="primary" disabled={busy}>
          {busy ? '…' : mode === 'login' ? 'Log in' : 'Register & log in'}
        </button>
        <p className="hint">Each browser tab keeps its own login, so open two tabs for two users.</p>
      </form>
    </div>
  )
}
