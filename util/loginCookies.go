package util

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/aabb-hash/web/db"
)

func SetLoginCookies(w http.ResponseWriter, session []byte, username string) {
	sesssionHex := hex.EncodeToString(session)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sesssionHex,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func CheckLoginSession(w http.ResponseWriter, r *http.Request) bool {
	if _, err := r.Cookie("session"); err != nil {
		return false
	} else {
		session, _ := r.Cookie("session")
		username, _ := r.Cookie("username")

		sessionBytes, _ := hex.DecodeString(session.Value)
		return db.VerifyUserSession(username.Value, sessionBytes)
	}
}
