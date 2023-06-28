package servers

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

type Client struct {
	ClientId    string          // 标识ID
	Socket      *websocket.Conn // 用户连接
	ConnectTime uint64          // 首次连接时间
	IsDeleted   bool            // 是否删除或下线
}

type SendData struct {
	Code int
	Msg  string
	Data *interface{}
}

func NewClient(clientId string, socket *websocket.Conn) *Client {
	return &Client{
		ClientId:    clientId,
		Socket:      socket,
		ConnectTime: uint64(time.Now().Unix()),
		IsDeleted:   false,
	}
}

func (c *Client) Read() {
	go func() {
		for {
			// 在这里完成客户端和服务端的会话
			messageType, receiveMsg, err := c.Socket.ReadMessage()

			if err != nil {
				if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					Manager.DisConnect <- c
					return
				} else if messageType != websocket.PingMessage {
					return
				}
			}

			msg := strings.Replace(string(receiveMsg), "\n", "", -1)
			msg = strings.Replace(msg, "\"", "", -1)
			if msg == "ping" {
				// 向客户端发送pong消息
				if err = PongRender(c.Socket); err != nil {
					_ = c.Socket.Close()
					return
				}
			}
		}
	}()
}
