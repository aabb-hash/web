package ws

import (
	"maps"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	username string
	Send     chan []byte
}

var connections map[uint8]*Client = make(map[uint8]*Client)

func Serve(userID uint8, messageType byte, data []byte) {
	switch messageType {
	case 1:
		Move(userID, data)
	case 4:
		HandleGameRequest(userID, data)
	}
}

func (client *Client) write() {
	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}
}

func CloseConnection(userID uint8) {
	send := connections[userID].Send
	conn := connections[userID].Conn

	delete(connections, userID)

	close(send)
	conn.Close()

	delete(positions, userID)
	SendUserBroadcast(2, nil, userID)
}

func NewConnection(conn *websocket.Conn, userID uint8, username string) {
	client := &Client{
		Conn:     conn,
		username: username,
		Send:     make(chan []byte, 256),
	}

	go client.write()
	connections[userID] = client
	Create(userID, username)
}

func compileUserMessage(messageType byte, data []byte, senderID uint8) []byte {
	message := []byte{messageType}
	if data != nil {
		message = append(message, data...)
	}
	message = append(message, byte(senderID))

	return message
}

func SendUserMessage(userID uint8, messageType byte, data []byte, senderID uint8) {
	client := connections[userID]
	sClient := connections[senderID]

	if client != nil && sClient != nil {
		client.Send <- compileUserMessage(messageType, data, senderID)
	}
}

func SendUserBroadcast(messageType byte, data []byte, senderID uint8) {
	message := compileUserMessage(messageType, data, senderID)

	for userID, client := range maps.Clone(connections) {
		if userID != senderID {
			select {
			case client.Send <- message:
			default:
			}
		}
	}
}

func SendMessage(userID uint8, messageType byte, data []byte) {
	client := connections[userID]

	if client != nil {
		client.Send <- append([]byte{messageType}, data...)
	}
}

func GetUsernameByID(id uint8) string {
	client := connections[id]
	if client != nil {
		return client.username
	}

	return ""
}
