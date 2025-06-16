package logic

import (
	"net/http"

	"github.com/aabb-hash/web/util"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	util.GetHeader("home", "home", w)
	util.GetHtmlContent("home", nil, w)

	util.GetHtmlContent("footer", nil, w)
}
