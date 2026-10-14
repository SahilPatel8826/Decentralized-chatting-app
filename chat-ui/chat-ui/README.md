# Chat UI: a test client for your Go chat backend

React + Vite. Talks to **user_service (:8080)** and **chat_service (:8081)**.
Built to exercise every feature the backend has and to show you what it really sends.

## Run

1. Start both Go services (Postgres running, `.env` files in place).
2. `npm install`
3. `npm run dev` → open http://localhost:5173

Change the ports in `vite.config.js` if yours differ. The dev proxy exists
because the Go services send no CORS headers.

## Test with two users

Each browser **tab** keeps its own login (sessionStorage), so open two tabs:
register `alice` in one and `bob` in the other. The backend allows **one socket
per user**, so never log in as the same user in two tabs.

Open the **Debug** bar at the bottom. It lists every HTTP call and every
WebSocket frame, raw. Look there first when something looks wrong.

## What to check, feature by feature

| Feature | How to try it | What you should see |
|---|---|---|
| Register / login | Auth page | Lands in chat; Debug shows `access_token` + `refresh_token` |
| Token refresh | Stay logged in >15 min and send a message | One 401 in Debug, then `/auth/refresh`, then the retry succeeds |
| User search | New chat → type a name | Matching users from `GET /users/search?q=` |
| Private chat | Alice: New chat → pick Bob | Conversation appears for Alice; Bob after his list reloads |
| Send + live push | Alice sends | Bob's tab shows it instantly (`message` frame) |
| Delivered / read ticks | Bob has the chat closed, then opens it | Alice's tick: ✓ sent → ✓✓ delivered → blue ✓✓ read |
| Typing | Alice types | Bob sees "alice typing…" (`typing_start`/`typing_stop` frames) |
| Presence | Bob closes his tab | Alice's dot for Bob turns offline |
| Delete message | Hover your own bubble → × | Gone for you; for the other user only if the backend sends `message_deleted` |
| Pagination | Send >30 messages, reopen the chat | "Load next page" button, pages of 30 |
| Group | New group → name + members | Group appears with "group" tag |
| Add / remove / leave | Members panel | Admin can add/remove; a non-admin should get an error toast |
| Profile | Profile → save | `PATCH /users/me` then the new name shows |
| Logout | Log out | Session cleared (the server does not revoke the token) |

## Backend quirks the UI works around (`src/api.js`)

- IDs come back as `ID` or `id` depending on the endpoint. Adapters normalise them.
- Results are wrapped: `conversation`, `conversations`, `message`, `messages`.
- Message history order depends on the backend. The UI sorts by message ID.
- The socket token goes in `?token=` (browsers cannot set WebSocket headers). This needs the `?token=` change in chat_service's auth middleware from Day 1.
- Sending a message is REST (`POST /messages`); the socket only delivers events and carries typing.
- An open chat marks incoming messages read (read implies delivered). Closed chats send "delivered".

## Layout

```
src/
  api.js          fetch wrapper, 401 → refresh → retry, response adapters
  session.js      tokens in sessionStorage
  useSocket.js    WebSocket with reconnect
  debug.js        log feeding the Debug panel
  Chat.jsx        all chat state + WebSocket event handling
  components/     AuthPage, Sidebar, ChatView, UserPicker, DebugPanel
```
