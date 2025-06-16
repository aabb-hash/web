package logic

import (
	"net/http"

	"github.com/aabb-hash/web/util"
)

func HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	util.GetHeader("leaderboard", "leaderboard", w)
	util.GetHtmlContent("leaderboard", nil, w)

	util.GetHtmlContent("footer", nil, w)
}
