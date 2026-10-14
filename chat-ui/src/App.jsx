import { useEffect, useState } from 'react'
import { getTokens, onSessionEnd } from './session'
import AuthPage from './components/AuthPage.jsx'
import Chat from './Chat.jsx'
import DebugPanel from './components/DebugPanel.jsx'

export default function App() {
  const [authed, setAuthed] = useState(!!getTokens())
  useEffect(() => onSessionEnd(() => setAuthed(false)), [])

  return (
    <>
      {authed ? <Chat /> : <AuthPage onAuthed={() => setAuthed(true)} />}
      <DebugPanel />
    </>
  )
}
