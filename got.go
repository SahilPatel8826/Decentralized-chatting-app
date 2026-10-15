package more

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow only our local frontend during development.
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" ||
			origin == "http://localhost:5173"
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	// Step 1: Upgrade HTTP connection to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Step 2: Close the connection when the handler exits.
	defer conn.Close()

	log.Println("Client connected")

	// Step 3: Continuously read messages.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}

		// Step 4: Print the received message.
		log.Printf("Received: %s", message)

		// Step 5: Send the same message back.
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWS)

	log.Println("WebSocket server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
