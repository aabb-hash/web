package tictactoe

import (
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
	case 4:
		PlayAgain(userID, data[0])
	case 7:
		Move(userID, data)
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

	HandleLeave(userID)
}

func NewConnection(conn *websocket.Conn, userID uint8, username string, gameID uint8) {
	client := &Client{
		Conn:     conn,
		username: username,
		Send:     make(chan []byte, 256),
	}

	go client.write()
	connections[userID] = client
	LoadPlayer(gameID, userID, username)
}

func SendMessage(userID uint8, messageType byte, data []byte) {
	client := connections[userID]

	if client != nil {
		client.Send <- append([]byte{messageType}, data...)
	}
}
