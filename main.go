package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aabb-hash/web/db"
	"github.com/aabb-hash/web/logic"
	"github.com/aabb-hash/web/login"
	"github.com/aabb-hash/web/net"
	"github.com/aabb-hash/web/play"
	"github.com/aabb-hash/web/util"
	"github.com/gorilla/websocket"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathPrefix := strings.Split(path, "/")[1]

	switch pathPrefix {
	case "login":
		login.HandleLogin(w, r)
		return
	case "register":
		login.HandleRegister(w, r)
		return
	case "api":
		logic.HandleAPI(w, r)
		return
	}

	if path == "/" {
		if util.CheckLoginSession(w, r) {
			logic.HandleHome(w, r)
			return
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		util.GetHeader("main", "main", w)
		util.GetHtmlContent("footer", nil, w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("saopdko")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if util.CheckLoginSession(w, r) {
		cookie, _ := r.Cookie("username")
		path := strings.TrimPrefix(r.URL.Path, "/ws")

		net.NewConnection(cookie.Value, conn, path)
	}
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	if util.CheckLoginSession(w, r) {
		cookie, _ := r.Cookie("username")
		gameID := strings.TrimPrefix(r.URL.Path, "/game/")

		play.LoadGame(gameID, cookie.Value, w, r)
	}
}

func main() {
	/*net.Init()
	util.InitHtml()
	db.Init()

	mux := http.NewServeMux()
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc("/", requestHandler)

	wrappedMux := net.CheckRequest(mux)
	http.ListenAndServe(":80", wrappedMux)
	return*/

	util.InitHtml()
	db.Init()
	net.Init()
	play.Init()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/ws/", wsHandler)
	http.HandleFunc("/game/", gameHandler)
	http.HandleFunc("/", requestHandler)

	http.ListenAndServe(":44444", nil)
}
