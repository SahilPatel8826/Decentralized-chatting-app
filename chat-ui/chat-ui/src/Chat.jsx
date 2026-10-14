import { useCallback, useEffect, useRef, useState } from 'react'
import * as api from './api'
import { useSocket } from './useSocket'
import Sidebar from './components/Sidebar.jsx'
import ChatView from './components/ChatView.jsx'

const LIMIT = 30
const byActivity = (a, b) => b.activity - a.activity

export default function Chat() {
  const [me, setMe] = useState(null)
  const [convs, setConvs] = useState([]) // normConv + unread
  const [activeId, setActiveId] = useState(null)
  const [msgs, setMsgs] = useState({}) // convId -> { items, page, done, loading }
  const [users, setUsers] = useState({}) // userId -> user
  const [presence, setPresence] = useState({}) // userId -> online | offline
  const [typing, setTyping] = useState({}) // convId -> { userId: timestamp }
  const [toast, setToast] = useState('')

  const meRef = useRef(null)
  const activeRef = useRef(null)
  const knownUsers = useRef(new Set())
  const loadConvsRef = useRef(() => {})
  meRef.current = me
  activeRef.current = activeId

  const notify = useCallback((m) => {
    setToast(String(m))
    setTimeout(() => setToast(''), 4000)
  }, [])

  /* ---------- users ---------- */
  const ensureUsers = useCallback(async (ids) => {
    const need = [...new Set(ids)].filter((id) => id && !knownUsers.current.has(id))
    need.forEach((id) => knownUsers.current.add(id))
    await Promise.all(
      need.map(async (id) => {
        try {
          const u = await api.getUser(id)
          setUsers((p) => ({ ...p, [id]: u }))
        } catch {
          knownUsers.current.delete(id)
        }
      }),
    )
  }, [])
  const nameOf = (id) => users[id]?.username || `user ${id}`

  /* ---------- conversations ---------- */
  const loadConvs = useCallback(async () => {
    const list = await api.listConversations()
    const full = await Promise.all(
      list.map((c) => (c.participants.length ? c : api.getConversation(c.id).catch(() => c))),
    )
    setConvs((prev) =>
      full
        .map((c) => ({ ...c, unread: prev.find((p) => p.id === c.id)?.unread || 0 }))
        .sort(byActivity),
    )
    ensureUsers(full.flatMap((c) => c.participants.map((p) => p.userId)))
  }, [ensureUsers])
  loadConvsRef.current = () => loadConvs().catch((e) => notify(e.message))

  const refreshConv = useCallback(
    async (id) => {
      const c = await api.getConversation(id)
      setConvs((p) => p.map((x) => (x.id === id ? { ...c, unread: x.unread, activity: x.activity } : x)))
      ensureUsers(c.participants.map((p) => p.userId))
    },
    [ensureUsers],
  )

  const titleOf = (c) => {
    if (c.type === 'group') return c.name || `Group #${c.id}`
    const other = c.participants.find((p) => p.userId !== me?.id)
    return other ? nameOf(other.userId) : `Chat #${c.id}`
  }

  /* ---------- messages ---------- */
  const patchMsg = useCallback((convId, id, fn) => {
    setMsgs((p) => {
      const cur = p[convId]
      if (!cur) return p
      return { ...p, [convId]: { ...cur, items: cur.items.map((m) => (m.id === id ? fn(m) : m)) } }
    })
  }, [])

  const loadMessages = useCallback(
    async (convId, page = 1) => {
      setMsgs((p) => ({ ...p, [convId]: { items: [], page: 0, done: false, ...p[convId], loading: true } }))
      try {
        const list = await api.getMessages(convId, page, LIMIT)
        setMsgs((p) => {
          const cur = p[convId] || { items: [] }
          const map = new Map(cur.items.map((m) => [m.id, m]))
          list.forEach((m) => map.set(m.id, { ...map.get(m.id), ...m }))
          return {
            ...p,
            [convId]: {
              items: [...map.values()].sort((a, b) => a.id - b.id),
              page,
              done: list.length < LIMIT,
              loading: false,
            },
          }
        })
      } catch (e) {
        notify(e.message)
        setMsgs((p) => ({ ...p, [convId]: { ...(p[convId] || { items: [] }), loading: false } }))
      }
    },
    [notify],
  )

  const openConv = (id) => {
    setActiveId(id)
    setConvs((p) => p.map((c) => (c.id === id ? { ...c, unread: 0 } : c)))
    if (!msgs[id]) loadMessages(id, 1)
  }

  const bump = (convId, fromOther) =>
    setConvs((p) =>
      p
        .map((c) =>
          c.id === convId
            ? {
                ...c,
                activity: Date.now(),
                unread: fromOther && activeRef.current !== convId ? c.unread + 1 : c.unread,
              }
            : c,
        )
        .sort(byActivity),
    )

  const send = async (convId, content) => {
    try {
      const m = await api.sendMessage(convId, content)
      setMsgs((p) => {
        const cur = p[convId] || { items: [], page: 1, done: false }
        if (cur.items.some((x) => x.id === m.id)) return p
        return { ...p, [convId]: { ...cur, items: [...cur.items, m] } }
      })
      bump(convId, false)
    } catch (e) {
      notify(e.message)
    }
  }

  const removeMessage = async (m) => {
    try {
      await api.deleteMessage(m.id)
      setMsgs((p) => {
        const cur = p[m.conversationId]
        if (!cur) return p
        return { ...p, [m.conversationId]: { ...cur, items: cur.items.filter((x) => x.id !== m.id) } }
      })
    } catch (e) {
      notify(e.message)
    }
  }

  const markRead = async (m) => {
    try {
      await api.markRead(m.id)
      patchMsg(m.conversationId, m.id, (x) => ({ ...x, status: 'read' }))
    } catch {
      /* visible in the Debug panel */
    }
  }

  /* ---------- websocket events ---------- */
  const onEvent = useCallback(
    (ev) => {
      if (!ev || typeof ev !== 'object') return
      const myId = meRef.current?.id
      switch (ev.type) {
        case 'message': {
          const m = {
            id: ev.message_id,
            conversationId: ev.conversation_id,
            senderId: String(ev.sender_id),
            content: ev.content,
            status: ev.status || 'sent',
            createdAt: new Date().toISOString(),
          }
          setMsgs((p) => {
            const cur = p[m.conversationId]
            if (!cur || cur.items.some((x) => x.id === m.id)) return p
            return { ...p, [m.conversationId]: { ...cur, items: [...cur.items, m] } }
          })
          setConvs((p) => {
            if (!p.some((c) => c.id === m.conversationId)) loadConvsRef.current()
            return p
          })
          bump(m.conversationId, true)
          setTyping((p) => {
            const cur = { ...(p[m.conversationId] || {}) }
            delete cur[m.senderId]
            return { ...p, [m.conversationId]: cur }
          })
          // Open chat: the read-receipt effect handles it (read implies delivered).
          if (m.senderId !== myId && activeRef.current !== m.conversationId) {
            api.markDelivered(m.id).catch(() => {})
          }
          break
        }
        case 'message_read':
          patchMsg(ev.conversation_id, ev.message_id, (m) => ({ ...m, status: 'read' }))
          break
        case 'message_delivered':
          patchMsg(ev.conversation_id, ev.message_id, (m) =>
            m.status === 'read' ? m : { ...m, status: 'delivered' },
          )
          break
        case 'message_deleted':
          setMsgs((p) => {
            const cur = p[ev.conversation_id]
            if (!cur) return p
            return {
              ...p,
              [ev.conversation_id]: { ...cur, items: cur.items.filter((x) => x.id !== ev.message_id) },
            }
          })
          break
        case 'presence':
          setPresence((p) => ({ ...p, [String(ev.user_id)]: ev.status }))
          break
        case 'typing_start':
        case 'typing_stop':
          setTyping((p) => {
            const cur = { ...(p[ev.conversation_id] || {}) }
            if (ev.type === 'typing_start') cur[String(ev.user_id)] = Date.now()
            else delete cur[String(ev.user_id)]
            return { ...p, [ev.conversation_id]: cur }
          })
          break
        default:
          break
      }
    },
    [patchMsg],
  )

  const { status: wsStatus, send: wsSend } = useSocket(!!me, onEvent)

  /* ---------- bootstrap ---------- */
  useEffect(() => {
    ;(async () => {
      try {
        const u = await api.me()
        knownUsers.current.add(u.id)
        setMe(u)
        setUsers((p) => ({ ...p, [u.id]: u }))
        await loadConvs()
      } catch (e) {
        notify(e.message)
      }
    })()
  }, [loadConvs, notify])

  /* ---------- actions passed to the UI ---------- */
  const startPrivate = async (user) => {
    try {
      const c = await api.createPrivate(user.id)
      setUsers((p) => ({ ...p, [user.id]: user }))
      knownUsers.current.add(user.id)
      await loadConvs()
      openConv(c.id)
    } catch (e) {
      notify(e.message)
    }
  }
  const startGroup = async (name, ids) => {
    try {
      const c = await api.createGroup(name, ids)
      await loadConvs()
      openConv(c.id)
    } catch (e) {
      notify(e.message)
    }
  }
  const groupAction = async (fn, convId) => {
    try {
      await fn()
      await refreshConv(convId)
    } catch (e) {
      notify(e.message)
    }
  }
  const leave = async (convId) => {
    try {
      await api.leaveGroup(convId)
      setConvs((p) => p.filter((c) => c.id !== convId))
      setActiveId(null)
    } catch (e) {
      notify(e.message)
    }
  }
  const saveProfile = async (patch) => {
    try {
      await api.updateMe(patch)
      const u = await api.me()
      setMe(u)
      setUsers((p) => ({ ...p, [u.id]: u }))
    } catch (e) {
      notify(e.message)
    }
  }

  const active = convs.find((c) => c.id === activeId) || null

  return (
    <div className="app">
      <Sidebar
        me={me}
        convs={convs}
        activeId={activeId}
        titleOf={titleOf}
        presence={presence}
        wsStatus={wsStatus}
        onOpen={openConv}
        onNewPrivate={startPrivate}
        onNewGroup={startGroup}
        onSaveProfile={saveProfile}
        onLogout={api.logout}
      />
      <main className="main">
        {active ? (
          <ChatView
            key={active.id}
            conv={active}
            me={me}
            title={titleOf(active)}
            nameOf={nameOf}
            presence={presence}
            state={msgs[active.id]}
            typingMap={typing[active.id] || {}}
            onSend={(text) => send(active.id, text)}
            onDelete={removeMessage}
            onMarkRead={markRead}
            onLoadMore={() => loadMessages(active.id, (msgs[active.id]?.page || 0) + 1)}
            onTyping={(isTyping) =>
              wsSend({ type: isTyping ? 'typing_start' : 'typing_stop', conversation_id: active.id })
            }
            onAdd={(u) => groupAction(() => api.addParticipant(active.id, u.id), active.id).then(() => ensureUsers([u.id]))}
            onRemove={(userId) => groupAction(() => api.removeParticipant(active.id, userId), active.id)}
            onLeave={() => leave(active.id)}
          />
        ) : (
          <div className="empty">Pick a conversation or start a new one.</div>
        )}
      </main>
      {toast && <div className="toast">{toast}</div>}
    </div>
  )
}
