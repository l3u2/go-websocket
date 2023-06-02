package servers

import (
	"fmt"
	"github.com/gorilla/websocket"
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
			fmt.Println("接收到的消息体", string(receiveMsg))
			if err != nil {
				if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					Manager.DisConnect <- c
					return
				} else if messageType != websocket.PingMessage {
					return
				}
			}
		}
	}()
}
