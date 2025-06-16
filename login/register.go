package login

import (
	"fmt"
	"net/http"

	"github.com/aabb-hash/web/db"
	"github.com/aabb-hash/web/util"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if util.CheckLoginSession(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.FormValue("username")
		if db.UserExist(username) {
			fmt.Fprint(w, "a user with this username already exists!")
		} else {
			password := r.FormValue("password")
			hash := util.HashPassword(password)
			session := util.GenerateSession()

			db.SaveUser(username, hash, session)
			util.SetLoginCookies(w, session, username)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	util.GetHeader("register", "register", w)
	util.GetHtmlContent("register", nil, w)

	util.GetHtmlContent("footer", nil, w)
}
