package servers

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var PreCheckTime uint64

type Client struct {
	ClientId      string          // 标识ID
	Socket        *websocket.Conn // 用户连接
	ConnectTime   uint64          // 首次连接时间
	LastCheckTime uint64          // 最后心跳时间
	IsDeleted     bool            // 是否删除或下线
}

type SendData struct {
	Code int
	Msg  string
	Data *interface{}
}

func NewClient(clientId string, socket *websocket.Conn) *Client {
	return &Client{
		ClientId:      clientId,
		Socket:        socket,
		ConnectTime:   uint64(time.Now().Unix()),
		LastCheckTime: uint64(time.Now().Unix()),
		IsDeleted:     false,
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
				c.LastCheckTime = uint64(time.Now().Unix())
				// 启动协程循环检查所有客户端最后心跳时间是否大于120s，如果大于120s,做主动下线处理
				ClientDisConnect()
				// 向客户端发送pong消息
				if err = PongRender(c.Socket); err != nil {
					_ = c.Socket.Close()
					return
				}
			}
		}
	}()
}

func ClientDisConnect() {
	go func() {
		currentTime := uint64(time.Now().Unix())
		if currentTime-PreCheckTime > 120 {
			for clientId, conn := range Manager.AllClient() {
				if currentTime-conn.LastCheckTime > 120 {
					Manager.DisConnect <- conn
					log.Infof("%s 客户端两分钟没有活跃关闭连接", clientId)
				}
			}

			PreCheckTime = currentTime
		}
	}()
}
