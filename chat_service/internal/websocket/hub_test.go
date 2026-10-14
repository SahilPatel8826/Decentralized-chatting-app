package websocket

import (
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"
)

func newClient(id string) *Client {
	return &Client{UserID: id, Send: make(chan []byte, 8)}
}

// recv returns (message, true) or (nil, false) if nothing arrives,
// or the channel is closed.
func recv(c *Client) ([]byte, bool) {
	select {
	case m, ok := <-c.Send:
		return m, ok
	case <-time.After(100 * time.Millisecond):
		return nil, false
	}
}

func TestRegisterAndIsOnline(t *testing.T) {
	h := NewHub()
	c := newClient("1")

	if h.IsOnline("1") {
		t.Fatal("user should be offline before register")
	}
	h.Register(c)
	if !h.IsOnline("1") {
		t.Fatal("user should be online after register")
	}
}

func TestRegisterReplacesOldConnection(t *testing.T) {
	h := NewHub()
	oldC, newC := newClient("1"), newClient("1")

	h.Register(oldC)
	h.Register(newC)

	if _, ok := <-oldC.Send; ok {
		t.Fatal("old connection's Send channel should be closed")
	}
	if !h.IsOnline("1") {
		t.Fatal("user should still be online via the new connection")
	}
}

func TestUnregisterIgnoresStaleClient(t *testing.T) {
	h := NewHub()
	oldC, newC := newClient("1"), newClient("1")

	h.Register(oldC)
	h.Register(newC)
	h.Unregister(oldC) // stale: must not remove the new connection

	if !h.IsOnline("1") {
		t.Fatal("stale Unregister removed the active connection")
	}
	select {
	case _, ok := <-newC.Send:
		if !ok {
			t.Fatal("new connection's Send channel was closed")
		}
	default:
	}
}

func TestSendToUser(t *testing.T) {
	h := NewHub()
	c := newClient("1")
	h.Register(c)

	h.SendToUser("1", []byte("hi"))
	if m, ok := recv(c); !ok || string(m) != "hi" {
		t.Fatalf("expected hi, got %q ok=%v", m, ok)
	}

	// offline user: must not panic
	h.SendToUser("999", []byte("nobody"))
}

func TestSendToUserDropsFullClient(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "1", Send: make(chan []byte, 1)}
	h.Register(c)

	h.SendToUser("1", []byte("a")) // fills buffer
	h.SendToUser("1", []byte("b")) // buffer full -> dropped

	if h.IsOnline("1") {
		t.Fatal("slow client should have been dropped")
	}
}

func TestTypingOnlyToOtherParticipants(t *testing.T) {
	h := NewHub()
	h.Participants = func(conversationID uint) []string {
		return []string{"1", "2"}
	}
	c1, c2, c3 := newClient("1"), newClient("2"), newClient("3")
	h.Register(c1)
	h.Register(c2)
	h.Register(c3)

	h.BroadcastTyping(10, "1", true)

	m, ok := recv(c2)
	if !ok {
		t.Fatal("participant 2 should receive typing event")
	}
	var p TypingPayload
	if err := json.Unmarshal(m, &p); err != nil {
		t.Fatal(err)
	}
	if p.Type != "typing_start" || p.ConversationID != 10 || p.UserID != "1" {
		t.Fatalf("unexpected payload: %+v", p)
	}
	if _, ok := recv(c1); ok {
		t.Fatal("sender should not receive own typing event")
	}
	if _, ok := recv(c3); ok {
		t.Fatal("non-participant must not receive typing event")
	}
}

func TestTypingIgnoredFromNonParticipant(t *testing.T) {
	h := NewHub()
	h.Participants = func(uint) []string { return []string{"1", "2"} }
	c1, c2, c3 := newClient("1"), newClient("2"), newClient("3")
	h.Register(c1)
	h.Register(c2)
	h.Register(c3)

	h.BroadcastTyping(10, "3", true)

	for _, c := range []*Client{c1, c2, c3} {
		if _, ok := recv(c); ok {
			t.Fatalf("user %s got an event from a non-participant", c.UserID)
		}
	}
}

func TestPresenceOnlyToContacts(t *testing.T) {
	h := NewHub()
	h.Contacts = func(userID string) []string { return []string{"2"} }
	c1, c2, c3 := newClient("1"), newClient("2"), newClient("3")
	h.Register(c1)
	h.Register(c2)
	h.Register(c3)

	h.BroadcastPresence("1", "online")

	m, ok := recv(c2)
	if !ok {
		t.Fatal("contact should receive presence")
	}
	var p PresencePayload
	if err := json.Unmarshal(m, &p); err != nil {
		t.Fatal(err)
	}
	if p.Type != "presence" || p.UserID != "1" || p.Status != "online" {
		t.Fatalf("unexpected payload: %+v", p)
	}
	if _, ok := recv(c3); ok {
		t.Fatal("non-contact must not receive presence")
	}
}

func TestNilLookupsDoNotPanic(t *testing.T) {
	h := NewHub()
	h.BroadcastTyping(1, "1", true)
	h.BroadcastPresence("1", "online")
}

// Run with: go test -race ./...
func TestConcurrentAccess(t *testing.T) {
	h := NewHub()
	h.Contacts = func(string) []string { return []string{"1", "2", "3"} }
	h.Participants = func(uint) []string { return []string{"1", "2", "3"} }

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		id := strconv.Itoa(i % 5)
		go func() {
			defer wg.Done()
			c := &Client{UserID: id, Send: make(chan []byte, 256)}
			h.Register(c)
			h.Unregister(c)
		}()
		go func() {
			defer wg.Done()
			h.SendToUser(id, []byte("x"))
			h.IsOnline(id)
			h.OnlineUsers()
		}()
		go func() {
			defer wg.Done()
			h.BroadcastPresence(id, "online")
			h.BroadcastTyping(1, "1", true)
		}()
	}
	wg.Wait()
}
