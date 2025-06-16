package login

import (
	"net/http"

	"github.com/aabb-hash/web/db"
	"github.com/aabb-hash/web/util"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if util.CheckLoginSession(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		hash := util.HashPassword(password)
		session := util.GenerateSession()

		if db.VerifyUserPassword(username, hash) {
			util.SetLoginCookies(w, session, username)
			db.NewSession(username, session)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	util.GetHeader("login", "login", w)
	util.GetHtmlContent("login", nil, w)

	util.GetHtmlContent("footer", nil, w)
}
