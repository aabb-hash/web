package play

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/aabb-hash/web/play/tictactoe"
)

type Game struct {
	gamemodeID uint8
	player1    string
	player2    string
}

var recycleIDs []uint8
var availableIDs []uint8
var syncIDs sync.Mutex

var gameMap map[uint8]*Game = make(map[uint8]*Game)

func CreateGame(gamemodeID uint8, player1 string, player2 string) (uint8, bool) {
	game := &Game{
		gamemodeID: gamemodeID,
		player1:    player1,
		player2:    player2,
	}

	gameID, max := getNewGameID()

	if !max {
		gameMap[gameID] = game
	}

	return gameID, max
}

func LoadGame(id string, player string, w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseUint(id, 10, 8)
	if err != nil {
		return
	}

	game := gameMap[uint8(gameID)]
	if game == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if game.player1 == player || game.player2 == player {
		switch game.gamemodeID {
		case 1:
			tictactoe.LoadGame(uint8(gameID), w, r)
		}
	}
}

func IsValidGameMode(gamemode uint8) bool {
	return gamemode == 1 || gamemode == 2
}

func GetGamemode(gameID uint8) uint8 {
	game, exists := gameMap[gameID]
	if exists {
		return game.gamemodeID
	}

	return 0
}

func Init() {
	for i := 0; i < 256; i++ {
		availableIDs = append(availableIDs, uint8(i))
	}
}

func getNewGameID() (uint8, bool) {
	var gameID uint8
	if len(availableIDs) > 0 {
		gameID = availableIDs[0]
		availableIDs = availableIDs[1:]
	} else {
		syncIDs.Lock()
		availableIDs = append(availableIDs, recycleIDs...)
		recycleIDs = recycleIDs[:0]
		syncIDs.Unlock()

		if len(availableIDs) == 0 {
			return 0, true
		}

		gameID = availableIDs[0]
		availableIDs = availableIDs[1:]
	}

	return gameID, false
}
