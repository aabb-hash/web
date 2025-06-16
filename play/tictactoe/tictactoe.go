package tictactoe

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/aabb-hash/web/db"
	"github.com/aabb-hash/web/util"
)

type Game struct {
	player1 uint8
	player2 uint8
	move    int8
	board   []uint8
}

var gameMap map[uint8]*Game = make(map[uint8]*Game)
var idToUsername map[uint8]string = make(map[uint8]string)
var playAgainMap map[uint8]uint8 = make(map[uint8]uint8)

func LoadGame(gameID uint8, w http.ResponseWriter, r *http.Request) {
	util.GetHeader("tictactoe", "tictactoe", w)
	util.GetHtmlContent("tictactoe", nil, w)
	util.GetHtmlContent("footer", nil, w)
}

func LoadPlayer(gameID uint8, userID uint8, username string) {
	game, exists := gameMap[gameID]
	if exists {
		game.player2 = userID
		idToUsername[userID] = username

		if rand.Intn(2) == 0 {
			game.player2 = game.player1
			game.player1 = userID
		}

		game.move = 0

		SendMessage(game.player1, 0, append([]byte{game.player2}, []byte(idToUsername[game.player2])...))
		SendMessage(game.player2, 0, append([]byte{game.player1}, []byte(idToUsername[game.player1])...))

		SendMessage(game.player1, 8, []byte{1})
		SendMessage(game.player2, 8, []byte{0})
	} else {
		gameMap[gameID] = &Game{
			player1: userID,
			player2: 0,
			move:    -1,
			board:   make([]uint8, 9),
		}

		idToUsername[userID] = username
	}
}

func HandleLeave(userID uint8) {
	for gameID, game := range gameMap {
		if game.player1 == userID {
			if game.move == -1 {
				SendMessage(game.player2, 8, []byte{6})
			} else {
				SendMessage(game.player2, 8, []byte{2})
				SendMessage(game.player2, 8, []byte{6})
			}

			cleanUp(gameID, game)
		} else if game.player2 == userID {
			if game.move == -1 {
				SendMessage(game.player1, 8, []byte{6})
			} else {
				SendMessage(game.player1, 8, []byte{2})
				SendMessage(game.player1, 8, []byte{6})
			}

			cleanUp(gameID, game)
		}
	}

	delete(idToUsername, userID)
}

func cleanUp(gameID uint8, game *Game) {
	delete(gameMap, gameID)
	delete(playAgainMap, game.player1)
	delete(playAgainMap, game.player2)
}

func PlayAgain(userID uint8, gameID uint8) {
	game, gameExists := gameMap[gameID]
	if !gameExists || !hasPlayer(game, userID) || game.move != -1 {
		return
	}

	player, hasAgain := playAgainMap[gameID]
	if hasAgain {
		if userID == player {
			return
		}

		delete(playAgainMap, game.player1)
		delete(playAgainMap, game.player2)

		player = game.player2
		if rand.Intn(2) == 0 {
			game.player2 = game.player1
			game.player1 = player
		}

		game.board = make([]uint8, 9)
		game.move = 0

		SendMessage(game.player1, 8, []byte{1})
		SendMessage(game.player2, 8, []byte{0})
	} else {
		if game.player1 == userID {
			SendMessage(game.player2, 8, []byte{5})
		} else {
			SendMessage(game.player1, 8, []byte{5})
		}

		playAgainMap[gameID] = userID
	}
}

func Move(userID uint8, data []byte) {
	game, exists := gameMap[data[1]]
	if exists {
		if game.move == -1 {
			return
		}

		move := data[0]
		if game.board[move] == 0 {
			if game.player1 == userID && game.move%2 == 0 {
				game.board[move] = 1
				game.move++
				SendMessage(game.player2, 7, data)
			} else if game.player2 == userID {
				game.board[move] = 2
				game.move++
				SendMessage(game.player1, 7, data)
			}

			if game.move > 4 {
				signal := checkWinner(game)
				switch signal {
				case 1:
					game.move = -1
					SendMessage(game.player1, 8, []byte{2})
					SendMessage(game.player2, 8, []byte{3})
					go saveMatchAndStats(game.player1, game.player2, false)
				case 2:
					game.move = -1
					SendMessage(game.player2, 8, []byte{2})
					SendMessage(game.player1, 8, []byte{3})
					go saveMatchAndStats(game.player2, game.player1, false)
				case 3:
					game.move = -1
					SendMessage(game.player1, 8, []byte{4})
					SendMessage(game.player2, 8, []byte{4})
					go saveMatchAndStats(game.player1, game.player2, true)
				}
			}
		}
	}
}

func checkWinner(game *Game) uint8 {
	board := game.board

	for x := 0; x < 3; x++ {
		if board[x*3] != 0 && board[x*3] == board[x*3+1] && board[x*3] == board[x*3+2] {
			return board[x*3]
		}
	}

	for y := 0; y < 3; y++ {
		if board[y] != 0 && board[y] == board[y+3] && board[y] == board[y+6] {
			return board[y*3]
		}
	}

	if board[0] != 0 && board[4] == board[0] && board[8] == board[0] {
		return board[0]
	}

	if board[2] != 0 && board[4] == board[2] && board[6] == board[2] {
		return board[2]
	}

	if game.move == 9 {
		return 3
	}

	return 0
}

func hasPlayer(game *Game, userID uint8) bool {
	return game.player1 == userID || game.player2 == userID
}

func saveMatchAndStats(player1 uint8, player2 uint8, draw bool) {
	username1 := idToUsername[player1]
	username2 := idToUsername[player2]

	fmt.Print(username1, username2)

	db.SaveMatch(username1, username2, draw)

	if draw {
		db.IncreaseStatistic(username1, "draw")
		db.IncreaseStatistic(username2, "draw")
	} else {
		db.IncreaseStatistic(username1, "win")
		db.IncreaseStatistic(username2, "lose")
	}
}
