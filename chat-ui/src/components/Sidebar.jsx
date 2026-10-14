import { useState } from 'react'
import UserPicker from './UserPicker.jsx'

export default function Sidebar({
  me,
  convs,
  activeId,
  titleOf,
  presence,
  wsStatus,
  onOpen,
  onNewPrivate,
  onNewGroup,
  onSaveProfile,
  onLogout,
}) {
  const [panel, setPanel] = useState(null) // 'private' | 'group' | 'profile'
  const [groupName, setGroupName] = useState('')
  const [members, setMembers] = useState([])
  const [username, setUsername] = useState('')
  const [bio, setBio] = useState('')

  const toggle = (name) => {
    setPanel(panel === name ? null : name)
    if (name === 'profile' && me) {
      setUsername(me.username)
      setBio(me.bio)
    }
  }

  const otherUserId = (c) => c.participants.find((p) => p.userId !== me?.id)?.userId

  return (
    <aside className="sidebar">
      <div className="me">
        <div>
          <strong>{me?.username || '…'}</strong> <small>#{me?.id}</small>
          <div className={`ws ${wsStatus}`}>socket: {wsStatus}</div>
        </div>
        <button onClick={onLogout}>Log out</button>
      </div>

      <div className="actions">
        <button className={panel === 'private' ? 'on' : ''} onClick={() => toggle('private')}>
          New chat
        </button>
        <button className={panel === 'group' ? 'on' : ''} onClick={() => toggle('group')}>
          New group
        </button>
        <button className={panel === 'profile' ? 'on' : ''} onClick={() => toggle('profile')}>
          Profile
        </button>
      </div>

      {panel === 'private' && (
        <div className="panel">
          <UserPicker
            excludeIds={[me?.id]}
            onPick={(u) => {
              onNewPrivate(u)
              setPanel(null)
            }}
          />
        </div>
      )}

      {panel === 'group' && (
        <div className="panel">
          <input placeholder="group name" value={groupName} onChange={(e) => setGroupName(e.target.value)} />
          <UserPicker
            excludeIds={[me?.id, ...members.map((m) => m.id)]}
            onPick={(u) => setMembers((p) => [...p, u])}
            placeholder="Add members…"
          />
          <div className="chips">
            {members.map((m) => (
              <span key={m.id} className="chip">
                {m.username || m.id}
                <button onClick={() => setMembers((p) => p.filter((x) => x.id !== m.id))}>×</button>
              </span>
            ))}
          </div>
          <button
            className="primary"
            disabled={!groupName.trim() || members.length === 0}
            onClick={() => {
              onNewGroup(groupName.trim(), members.map((m) => m.id))
              setGroupName('')
              setMembers([])
              setPanel(null)
            }}
          >
            Create group
          </button>
        </div>
      )}

      {panel === 'profile' && (
        <div className="panel">
          <input placeholder="username" value={username} onChange={(e) => setUsername(e.target.value)} />
          <input placeholder="bio" value={bio} onChange={(e) => setBio(e.target.value)} />
          <button
            className="primary"
            onClick={() => {
              onSaveProfile({ username, bio, avatar_url: '' })
              setPanel(null)
            }}
          >
            Save
          </button>
        </div>
      )}

      <ul className="convs">
        {convs.length === 0 && <li className="muted">No conversations yet.</li>}
        {convs.map((c) => {
          const other = c.type === 'private' ? otherUserId(c) : null
          return (
            <li key={c.id}>
              <button className={c.id === activeId ? 'active' : ''} onClick={() => onOpen(c.id)}>
                <span className="title">
                  {other && <i className={`dot ${presence[other] || 'unknown'}`} />}
                  {c.type === 'group' && <span className="tag">group</span>}
                  {titleOf(c)}
                </span>
                {c.unread > 0 && <span className="badge">{c.unread}</span>}
              </button>
            </li>
          )
        })}
      </ul>
    </aside>
  )
}
