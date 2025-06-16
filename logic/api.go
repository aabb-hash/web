package logic

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aabb-hash/web/db"
	"github.com/aabb-hash/web/net"
	"github.com/aabb-hash/web/util"
)

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	path, _ := strings.CutPrefix(r.URL.Path, "/api/")

	switch path {
	case "settings":
		if r.Method != http.MethodPut || !util.CheckLoginSession(w, r) {
			return
		}

		username, _ := r.Cookie("username")

		newUsername := r.FormValue("username")
		newPassword := r.FormValue("password")

		if newPassword != "" {
			db.NewPassword(username.Value, util.HashPassword(newPassword))
		}

		fmt.Println("getting data to save")
		file, header, err := r.FormFile("newAvatar")
		if err == nil {
			data, err := io.ReadAll(file)
			if err != nil {
				panic(err)
			}

			fmt.Println("saving img")
			if strings.HasSuffix(header.Filename, ".jpg") {
				data = append(data, 1)
			} else if strings.HasSuffix(header.Filename, "png") {
				data = append(data, 0)
			} else {
				w.WriteHeader(400)
				return
			}

			db.NewAvatar(username.Value, data)
		}

		if newUsername != "" {
			db.NewUsername(username.Value, newUsername)
		}
	case "avatar":
		if r.Method != http.MethodGet || !util.CheckLoginSession(w, r) {
			return
		}

		id := r.FormValue("id")
		var avatar []byte

		if id == "" {
			usernameCookie, _ := r.Cookie("username")
			avatar = db.GetAvatar(usernameCookie.Value)
		} else {
			numericID, err := strconv.ParseUint(id, 10, 8)
			if err != nil {
				avatar = nil
			} else {
				avatar = db.GetAvatar(net.GetUsernameByID(uint8(numericID)))
			}
		}

		if avatar == nil {
			width, height := 150, 150
			img := image.NewRGBA(image.Rect(0, 0, width, height))

			red := color.RGBA{255, 0, 0, 255}
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					img.Set(x, y, red)
				}
			}

			w.Header().Set("Content-Type", "image/png")
			png.Encode(w, img)
			return
		}

		var mimeType string
		if avatar[len(avatar)-1] == 1 {
			mimeType = "image/jpeg"
		} else {
			mimeType = "image/png"
		}

		w.Header().Set("Content-Type", mimeType)

		data := avatar[:len(avatar)-1]
		w.Write(data)
		return
	case "stats-to-opponent":
		if r.Method != http.MethodGet || !util.CheckLoginSession(w, r) {
			return
		}

		username, _ := r.Cookie("username")
		opponent := r.FormValue("opponent")

		if opponent == "" {
			w.WriteHeader(400)
			return
		}

		win, draw, lose := db.GetStatsToOpponent(username.Value, opponent)

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%d;%d;%d", win, draw, lose)
		return
	}

	w.WriteHeader(204)
}
