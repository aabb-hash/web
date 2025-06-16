package net

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aabb-hash/web/net/ws"
	"github.com/aabb-hash/web/play"
	"github.com/aabb-hash/web/play/tictactoe"
	"github.com/gorilla/websocket"
)

var recycleIDs []uint8
var availableIDs []uint8
var syncIDs sync.Mutex

var connections map[uint8]string = make(map[uint8]string)
var transfers sync.Map

func listen(channelID uint8, userID uint8, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("x")
			break
		}

		messageType := message[0]
		data := message[1:]

		if messageType == 6 {
			go transfer(connections[userID], userID)
			closeConnection(channelID, userID, true)
			return
		}

		switch channelID {
		case 0:
			ws.Serve(userID, messageType, data)
		case 1:
			tictactoe.Serve(userID, messageType, data)
		}
	}

	closeConnection(channelID, userID, false)
}

func closeConnection(channelID uint8, userID uint8, transfer bool) {
	switch channelID {
	case 0:
		ws.CloseConnection(userID)
	case 1:
		tictactoe.CloseConnection(userID)
	}

	if !transfer {
		syncIDs.Lock()
		recycleIDs = append(recycleIDs, userID)
		syncIDs.Unlock()
	}
}

func transfer(username string, userID uint8) {
	transfers.Store(username, userID)
	time.Sleep(time.Second * 6)

	if _, ok := transfers.Load(username); !ok {
		transfers.Delete(username)
	}
}

func NewConnection(username string, conn *websocket.Conn, path string) {
	var userID uint8

	if id, ok := transfers.Load(username); ok {
		transfers.Delete(username)
		userID = id.(uint8)
	} else if len(availableIDs) > 0 {
		userID = availableIDs[0]
		availableIDs = availableIDs[1:]
	} else {
		syncIDs.Lock()
		availableIDs = append(availableIDs, recycleIDs...)
		recycleIDs = recycleIDs[:0]
		syncIDs.Unlock()

		if len(availableIDs) == 0 {
			conn.Close()
			return
		}

		userID = availableIDs[0]
		availableIDs = availableIDs[1:]
	}

	if strings.HasPrefix(path, "/game/") {
		gameID, err := strconv.ParseUint(strings.TrimPrefix(path, "/game/"), 10, 8)
		if err != nil {
			conn.Close()
			return
		}

		channelID := play.GetGamemode(uint8(gameID))
		switch channelID {
		case 0:
			conn.Close()
			return
		case 1:
			tictactoe.NewConnection(conn, userID, username, uint8(gameID))
		}

		go listen(channelID, userID, conn)
	} else {
		ws.NewConnection(conn, userID, username)
		go listen(0, userID, conn)
	}

	connections[userID] = username
}

func Init() {
	for i := 0; i < 256; i++ {
		availableIDs = append(availableIDs, uint8(i))
	}
}

func GetUsernameByID(userID uint8) string {
	username, exists := connections[userID]
	if exists {
		return username
	}

	return ""
}
