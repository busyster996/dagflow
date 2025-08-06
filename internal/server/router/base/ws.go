package base

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upGrader.Upgrade(w, r, nil)
}

func CloseWs(ws *websocket.Conn, message string) {
	if ws == nil {
		return
	}
	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, message)
	_ = ws.WriteControl(websocket.CloseMessage, closeMessage, time.Now().Add(3*time.Second))
	time.Sleep(1 * time.Second)
	_ = ws.Close()
}
