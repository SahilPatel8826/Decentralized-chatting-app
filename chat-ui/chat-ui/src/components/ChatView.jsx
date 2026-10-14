import { useEffect, useRef, useState } from 'react'
import UserPicker from './UserPicker.jsx'

const tick = (s) => (s === 'sent' ? '✓' : '✓✓')
const timeOf = (iso) => {
  const d = iso ? new Date(iso) : null
  return d && !Number.isNaN(d.getTime())
    ? d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    : ''
}

export default function ChatView({
  conv,
  me,
  title,
  nameOf,
  presence,
  state,
  typingMap,
  onSend,
  onDelete,
  onMarkRead,
  onLoadMore,
  onTyping,
  onAdd,
  onRemove,
  onLeave,
}) {
  const items = state?.items || []
  const [text, setText] = useState('')
  const [showGroup, setShowGroup] = useState(false)
  const [, tickRender] = useState(0)
  const bottomRef = useRef(null)
  const readSent = useRef(new Set())
  const typingTimer = useRef(null)
  const typingOn = useRef(false)

  // re-render every second so "typing…" expires
  useEffect(() => {
    const t = setInterval(() => tickRender((n) => n + 1), 1000)
    return () => {
      clearInterval(t)
      clearTimeout(typingTimer.current)
    }
  }, [])

  // while this chat is open, incoming messages count as read
  useEffect(() => {
    for (const m of items) {
      if (m.senderId !== me.id && m.status !== 'read' && !readSent.current.has(m.id)) {
        readSent.current.add(m.id)
        onMarkRead(m)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [items])

  const lastId = items.length ? items[items.length - 1].id : 0
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: 'end' })
  }, [lastId])

  const typingNames = Object.entries(typingMap)
    .filter(([id, ts]) => id !== me.id && Date.now() - ts < 5000)
    .map(([id]) => nameOf(id))

  const otherId = conv.type === 'private' ? conv.participants.find((p) => p.userId !== me.id)?.userId : null
  const isAdmin = conv.participants.some((p) => p.userId === me.id && p.role === 'admin')

  const onChange = (e) => {
    setText(e.target.value)
    if (!typingOn.current) {
      typingOn.current = true
      onTyping(true)
    }
    clearTimeout(typingTimer.current)
    typingTimer.current = setTimeout(() => {
      typingOn.current = false
      onTyping(false)
    }, 2000)
  }

  const submit = (e) => {
    e.preventDefault()
    const content = text.trim()
    if (!content) return
    clearTimeout(typingTimer.current)
    if (typingOn.current) {
      typingOn.current = false
      onTyping(false)
    }
    onSend(content)
    setText('')
  }

  return (
    <div className="chat">
      <header className="chat-head">
        <div>
          <strong>{title}</strong>
          <div className="sub">
            {otherId ? (
              <>
                <i className={`dot ${presence[otherId] || 'unknown'}`} /> {presence[otherId] || 'presence unknown'}
              </>
            ) : (
              `${conv.participants.length} members${isAdmin ? ' · you are admin' : ''}`
            )}
            {' · '}conversation #{conv.id}
          </div>
        </div>
        {conv.type === 'group' && (
          <button onClick={() => setShowGroup((v) => !v)}>{showGroup ? 'Hide members' : 'Members'}</button>
        )}
      </header>

      {showGroup && conv.type === 'group' && (
        <div className="panel group">
          <ul>
            {conv.participants.map((p) => (
              <li key={p.userId}>
                {nameOf(p.userId)} <small>#{p.userId} · {p.role}</small>
                {p.userId !== me.id && <button onClick={() => onRemove(p.userId)}>Remove</button>}
              </li>
            ))}
          </ul>
          <UserPicker
            excludeIds={conv.participants.map((p) => p.userId)}
            onPick={onAdd}
            placeholder="Add a member…"
          />
          <button className="danger" onClick={onLeave}>
            Leave group
          </button>
        </div>
      )}

      <div className="messages">
        {state && !state.done && (
          <button className="more" disabled={state.loading} onClick={onLoadMore}>
            {state.loading ? 'Loading…' : `Load next page (${state.page + 1})`}
          </button>
        )}
        {state?.loading && items.length === 0 && <p className="muted">Loading…</p>}
        {!state?.loading && items.length === 0 && <p className="muted">No messages yet. Say hi.</p>}
        {items.map((m) => {
          const own = m.senderId === me.id
          return (
            <div key={m.id} className={`bubble ${own ? 'own' : 'other'}`}>
              {!own && conv.type === 'group' && <div className="who">{nameOf(m.senderId)}</div>}
              <div className="text">{m.content}</div>
              <div className="meta">
                <span>#{m.id}</span>
                <span>{timeOf(m.createdAt)}</span>
                {own && <span className={`tick ${m.status}`}>{tick(m.status)}</span>}
                {own && (
                  <button className="del" title="Delete" onClick={() => onDelete(m)}>
                    ×
                  </button>
                )}
              </div>
            </div>
          )
        })}
        <div ref={bottomRef} />
      </div>

      <div className="typing">{typingNames.length > 0 && `${typingNames.join(', ')} typing…`}</div>

      <form className="composer" onSubmit={submit}>
        <input placeholder="Type a message" value={text} onChange={onChange} />
        <button className="primary" disabled={!text.trim()}>
          Send
        </button>
      </form>
    </div>
  )
}
