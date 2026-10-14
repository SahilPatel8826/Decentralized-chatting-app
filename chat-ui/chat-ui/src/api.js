import { getTokens, setTokens, clearTokens, tokenExp } from './session'
import { log } from './debug'

const USER = '/user-api'
const CHAT = '/chat-api'

export class ApiError extends Error {
  constructor(status, message, data) {
    super(message)
    this.status = status
    this.data = data
  }
}

/* ---------- adapters: the backend mixes ID/id and wraps results ---------- */
const pick = (o, ...keys) => {
  for (const k of keys) if (o && o[k] !== undefined && o[k] !== null) return o[k]
  return undefined
}
export const normUser = (u) => ({
  id: String(pick(u, 'ID', 'id', 'Id')),
  username: pick(u, 'username', 'Username') ?? '',
  email: pick(u, 'email', 'Email') ?? '',
  bio: pick(u, 'bio', 'Bio') ?? '',
})
export const normMessage = (m) => ({
  id: Number(pick(m, 'ID', 'id')),
  conversationId: Number(pick(m, 'conversation_id', 'ConversationID')),
  senderId: String(pick(m, 'sender_id', 'SenderID')),
  content: pick(m, 'content', 'Content') ?? '',
  status: pick(m, 'status', 'Status') ?? 'sent',
  createdAt: pick(m, 'CreatedAt', 'created_at'),
})
export const normConv = (c) => ({
  id: Number(pick(c, 'ID', 'id')),
  type: pick(c, 'Type', 'type') ?? 'private',
  name: pick(c, 'Name', 'name') ?? '',
  participants: (pick(c, 'Participants', 'participants') || []).map((p) => ({
    userId: String(pick(p, 'UserID', 'user_id')),
    role: pick(p, 'Role', 'role') ?? 'member',
  })),
  activity: Date.parse(pick(c, 'UpdatedAt', 'updated_at') || '') || 0,
})

/* ---------- transport ---------- */
async function raw(base, path, { method = 'GET', body, token } = {}) {
  const headers = {}
  if (body) headers['Content-Type'] = 'application/json'
  if (token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(base + path, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })
  const text = await res.text()
  let data = null
  try {
    data = text ? JSON.parse(text) : null
  } catch {
    data = text
  }
  log({ kind: 'http', method, url: base + path, status: res.status, body, data })
  return { res, data }
}

let refreshing = null
export function refresh() {
  if (!refreshing) {
    refreshing = (async () => {
      const t = getTokens()
      if (!t?.refresh) throw new ApiError(401, 'no refresh token')
      const { res, data } = await raw(USER, '/auth/refresh', {
        method: 'POST',
        body: { refresh_token: t.refresh },
      })
      if (!res.ok) {
        clearTokens()
        throw new ApiError(401, data?.error || 'session expired')
      }
      setTokens({ access: data.access_token, refresh: data.refresh_token || t.refresh })
    })().finally(() => {
      refreshing = null
    })
  }
  return refreshing
}

// Refresh early if the access token is about to expire (used before WS reconnect).
export async function ensureFresh() {
  const t = getTokens()
  if (t && tokenExp(t.access) - Date.now() < 60_000) await refresh()
}

async function request(base, path, opts = {}) {
  const t = getTokens()
  let { res, data } = await raw(base, path, { ...opts, token: t?.access })
  if (res.status === 401 && t?.refresh) {
    await refresh()
    ;({ res, data } = await raw(base, path, { ...opts, token: getTokens()?.access }))
  }
  if (!res.ok) {
    throw new ApiError(res.status, data?.error || res.statusText || 'request failed', data)
  }
  return data
}

/* ---------- user service ---------- */
export const register = (username, email, password) =>
  request(USER, '/auth/register', { method: 'POST', body: { username, email, password } })

export async function login(email, password) {
  const data = await request(USER, '/auth/login', { method: 'POST', body: { email, password } })
  setTokens({ access: data.access_token, refresh: data.refresh_token })
}

export async function logout() {
  try {
    await request(USER, '/auth/logout', { method: 'POST' })
  } catch {
    /* logout is best-effort; the server does not revoke tokens yet */
  }
  clearTokens()
}

export const me = async () => normUser(await request(USER, '/users/me'))
export const getUser = async (id) => normUser(await request(USER, `/users/${id}`))
export const updateMe = (patch) => request(USER, '/users/me', { method: 'PATCH', body: patch })
export async function searchUsers(q) {
  const data = await request(USER, `/users/search?q=${encodeURIComponent(q)}`)
  const list = Array.isArray(data) ? data : data?.users || []
  return list.map(normUser)
}

/* ---------- chat service ---------- */
export async function listConversations() {
  const data = await request(CHAT, '/conversations')
  return (data?.conversations || []).map(normConv)
}
export const getConversation = async (id) =>
  normConv((await request(CHAT, `/conversations/${id}`)).conversation)
export const createPrivate = async (participantId) =>
  normConv(
    (await request(CHAT, '/conversations/private', {
      method: 'POST',
      body: { participant_id: String(participantId) },
    })).conversation,
  )
export const createGroup = async (name, participantIds) =>
  normConv(
    (await request(CHAT, '/conversations/group', {
      method: 'POST',
      body: { name, participants: participantIds.map(String) },
    })).conversation,
  )
export const addParticipant = (convId, userId) =>
  request(CHAT, `/conversations/${convId}/participants`, {
    method: 'POST',
    body: { user_id: String(userId) },
  })
export const removeParticipant = (convId, userId) =>
  request(CHAT, `/conversations/${convId}/participants/${userId}`, { method: 'DELETE' })
export const leaveGroup = (convId) =>
  request(CHAT, `/conversations/${convId}/leave`, { method: 'DELETE' })

export async function getMessages(convId, page, limit) {
  const data = await request(CHAT, `/conversations/${convId}/messages?page=${page}&limit=${limit}`)
  return (data?.messages || []).map(normMessage)
}
export const sendMessage = async (conversationId, content) =>
  normMessage(
    (await request(CHAT, '/messages', {
      method: 'POST',
      body: { conversation_id: conversationId, content },
    })).message,
  )
export const markRead = (id) => request(CHAT, `/messages/${id}/read`, { method: 'PATCH' })
export const markDelivered = (id) => request(CHAT, `/messages/${id}/delivered`, { method: 'PATCH' })
export const deleteMessage = (id) => request(CHAT, `/messages/${id}`, { method: 'DELETE' })
