package servers

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go-websocket/pkg/setting"
	"go-websocket/tools/util"
	"net/http"
	"time"
)

//channel通道
var ToClientChan chan clientInfo

//channel通道结构体
type clientInfo struct {
	ClientId  string
	Cmd       string
	MessageId string
	Code      int
	Msg       string
	Data      *string
}

type RetData struct {
	MessageId string      `json:"messageId"`
	Cmd       string      `json:"cmd"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

func StartWebSocket() {
	websocketHandler := &Controller{}
	http.HandleFunc("/ws", websocketHandler.Run)
	go Manager.Start()
}

func Render(conn *websocket.Conn, messageId string, cmd string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:      code,
		MessageId: messageId,
		Cmd:       cmd,
		Msg:       message,
		Data:      data,
	})
}

//监听并发送给客户端信息
func WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		log.WithFields(log.Fields{
			"host":      setting.GlobalSetting.LocalHost,
			"port":      setting.CommonSetting.HttpPort,
			"clientId":  clientInfo.ClientId,
			"messageId": clientInfo.MessageId,
			"cmd":       clientInfo.Cmd,
			"code":      clientInfo.Code,
			"msg":       clientInfo.Msg,
			"data":      clientInfo.Data,
		}).Info("发送到本机")
		if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageId, clientInfo.Cmd, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				log.WithFields(log.Fields{
					"host":     setting.GlobalSetting.LocalHost,
					"port":     setting.CommonSetting.HttpPort,
					"clientId": clientInfo.ClientId,
					"msg":      clientInfo.Msg,
				}).Error("客户端异常离线：" + err.Error())
			}
		}
	}
}

//发送信息到指定客户端
func SendMessage2Client(clientId string, cmd string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()

	//如果是单机服务，则只发送到本机
	SendMessage2LocalClient(messageId, clientId, cmd, code, msg, data)

	return
}

//通过本服务器发送信息
func SendMessage2LocalClient(messageId, clientId string, cmd string, code int, msg string, data *string) {
	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.HttpPort,
		"clientId": clientId,
	}).Info("发送到通道")
	ToClientChan <- clientInfo{ClientId: clientId, MessageId: messageId, Cmd: cmd, Code: code, Msg: msg, Data: data}
	return
}
