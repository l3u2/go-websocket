package routers

import (
	"go-websocket/api/send2client"
	"go-websocket/servers"
	"net/http"
)

func Init() {
	sendToClientHandler := &send2client.Controller{}
	http.HandleFunc("/api/send_to_client", sendToClientHandler.Run)
	servers.StartWebSocket()
	go servers.WriteMessage()
}
