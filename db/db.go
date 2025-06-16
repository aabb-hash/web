package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init() {
	init, err := sql.Open("sqlite", "data.sqlite")
	if err != nil {
		panic(err)
	}
	db = init

	createTables()
}

func createTables() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (username VARCHAR(20) PRIMARY KEY, password_hash BLOB NOT NULL, session BLOB NOT NULL, last_login INTEGER, avatar BLOB)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS statistics (username VARCHAR(20) PRIMARY KEY, win INTEGER NOT NULL, lose INTEGER NOT NULL, draw INTEGER NOT NULL)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS matches (id INTEGER PRIMARY KEY AUTOINCREMENT, player1 VARCHAR(20), player2 VARCHAR(20), draw BOOLEAN)")
	if err != nil {
		panic(err)
	}
}

func UserExist(username string) bool {
	var resp string

	stmt, _ := db.Prepare("SELECT username FROM users WHERE username = ? LIMIT 1")
	err := stmt.QueryRow(username).Scan(&resp)

	if err != nil {
		return false
	}

	stmt.Close()
	return resp == username
}

func VerifyUserSession(username string, session []byte) bool {
	var user string
	var lastLogin int64

	stmt, _ := db.Prepare("SELECT username, last_login FROM users WHERE username = ? AND session = ? LIMIT 1")
	err := stmt.QueryRow(username, session).Scan(&user, &lastLogin)

	if err != nil {
		return false
	}

	stmt.Close()
	return user == username && time.Unix(lastLogin, 0).After(time.Now().Add(-1*time.Hour))
}

func VerifyUserPassword(username string, hash []byte) bool {
	var user string

	stmt, _ := db.Prepare("SELECT username FROM users WHERE username = ? AND password_hash = ? LIMIT 1")
	err := stmt.QueryRow(username, hash).Scan(&user)

	if err != nil {
		return false
	}

	stmt.Close()
	return user == username
}

func NewSession(username string, session []byte) {
	stmt, _ := db.Prepare("UPDATE users SET session = ?, last_login = ? WHERE username = ?")
	stmt.Exec(session, time.Now().Unix(), username)
	stmt.Close()
}
func SaveUser(username string, hash []byte, session []byte) {
	stmt, _ := db.Prepare("INSERT INTO users VALUES (?, ?, ?, ?, ?)")
	stmt.Exec(username, hash, session, time.Now().Unix(), nil)
	stmt.Close()

	stmt, _ = db.Prepare("INSERT INTO statistics VALUES (?, 0, 0, 0)")
	stmt.Exec(username)
	stmt.Close()
}

func NewAvatar(username string, avatar []byte) {
	stmt, _ := db.Prepare("UPDATE users SET avatar = ? WHERE username = ?")
	stmt.Exec(avatar, username)
	stmt.Close()
}
func NewUsername(username string, newUsername string) {
	stmt, _ := db.Prepare("UPDATE users SET username = ? WHERE username = ?")
	stmt.Exec(newUsername, username)
	stmt.Close()
}
func NewPassword(username string, newPasswordHash []byte) {
	stmt, _ := db.Prepare("UPDATE users SET password_hash = ? WHERE username = ?")
	stmt.Exec(newPasswordHash, username)
	stmt.Close()
}

func GetAvatar(username string) []byte {
	var avatar []byte

	stmt, _ := db.Prepare("SELECT avatar FROM users WHERE username = ? LIMIT 1")
	err := stmt.QueryRow(username).Scan(&avatar)

	if err != nil {
		return nil
	}

	stmt.Close()
	return avatar
}

func IncreaseStatistic(username string, statistic string) {
	stmt, _ := db.Prepare("UPDATE statistics SET " + statistic + " = " + statistic + " + 1 WHERE username = ?")
	stmt.Exec(username)
	stmt.Close()
}
func SaveMatch(player1 string, player2 string, draw bool) {
	stmt, _ := db.Prepare("INSERT INTO matches(player1, player2, draw) VALUES(?, ?, ?)")
	stmt.Exec(player1, player2, draw)
	stmt.Close()
}

func GetStatsToOpponent(player string, opponent string) (int, int, int) {
	var wins int
	var draw int
	var lose int

	stmt, _ := db.Prepare("SELECT Count(player1) FROM matches WHERE player1 = ? AND player2 = ? AND draw = false")
	stmt.QueryRow(player, opponent).Scan(&wins)
	stmt.Close()

	stmt, _ = db.Prepare("SELECT Count(player1) FROM matches WHERE draw = true AND ((player1 = ? AND player2 = ?) OR (player1 = ? AND player2 = ?))")
	stmt.QueryRow(player, opponent, opponent, player).Scan(&draw)
	stmt.Close()

	stmt, _ = db.Prepare("SELECT Count(player1) FROM matches WHERE player1 = ? AND player2 = ? AND draw = false")
	stmt.QueryRow(opponent, player).Scan(&lose)
	stmt.Close()

	return wins, draw, lose
}
