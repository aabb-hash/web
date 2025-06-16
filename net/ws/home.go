package ws

import (
	"fmt"

	"github.com/aabb-hash/web/play"
)

type GameRequest struct {
	targetUserID uint8
	gamemodeID   uint8
}

var positions map[uint8][]byte = make(map[uint8][]byte)
var gameRequests map[uint8][]*GameRequest = make(map[uint8][]*GameRequest)

func Move(id uint8, data []byte) {
	positions[id] = data
	SendUserBroadcast(1, data, id)
}

func Create(id uint8, username string) {
	data := []byte{0x00, 0x32, 0x00, 0x32}

	positions[id] = data

	for userID, userData := range positions {
		if userID != id {
			fmt.Println(id, userID)
			SendUserMessage(id, 0, append(userData, []byte(GetUsernameByID(userID))...), userID)
		}
	}

	data = append(data, []byte(username)...)
	SendUserBroadcast(0, data, id)
}

func HandleGameRequest(userID uint8, data []byte) {
	if len(data) < 2 {
		return
	}

	game := data[0]
	targetID := data[1]

	_, exists := positions[targetID]
	if !exists {
		SendMessage(userID, 5, []byte{0}) // target is no longer available
		return
	}

	requests, hasAny := gameRequests[targetID]
	if hasAny {
		for _, request := range requests {
			if userID == request.targetUserID {

				gameID, max := play.CreateGame(game, GetUsernameByID(targetID), GetUsernameByID(userID))

				if max {
					SendMessage(targetID, 5, []byte{3}) // game limit reached
					SendMessage(userID, 5, []byte{3})
				}

				SendMessage(targetID, 5, []byte{2, targetID, gameID}) // request accepted
				SendMessage(userID, 5, []byte{2, targetID, gameID})
				return
			}
		}
	}

	if play.IsValidGameMode(game) {
		requests, hasAny := gameRequests[userID]
		newRequest := &GameRequest{
			targetUserID: targetID,
			gamemodeID:   game,
		}

		if hasAny {
			requests = append(requests, newRequest)
			gameRequests[userID] = requests
		} else {
			gameRequests[userID] = []*GameRequest{newRequest}
		}

		SendUserMessage(targetID, 4, []byte{game}, userID) // forward request
		SendMessage(userID, 5, []byte{1, targetID})        // send request status "pending"
	}
}
