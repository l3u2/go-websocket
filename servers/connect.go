package servers

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go-websocket/define/retcode"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type Controller struct {
}

type renderData struct {
	ClientId string `json:"clientId"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade error: %v", err)
		http.NotFound(w, r)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//解析参数
	clientId := r.FormValue("clientId")
	if len(clientId) == 0 {
		_ = Render(conn, "", "", retcode.SYSTEM_ID_ERROR, "客户端ID不能为空", []string{})
		_ = conn.Close()
		return
	}

	clientSocket := NewClient(clientId, conn)

	//读取客户端消息
	clientSocket.Read()

	//返回客户端信息
	if err = ConnRender(conn, renderData{ClientId: clientId}); err != nil {
		_ = conn.Close()
		return
	}

	// 用户连接事件
	Manager.Connect <- clientSocket
}
