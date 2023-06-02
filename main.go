package main

import (
	"fmt"
	"go-websocket/pkg/setting"
	"go-websocket/routers"
	"go-websocket/tools/log"
	"net/http"
)

func init() {
	setting.Setup()
	log.Setup()
}

func main() {
	//初始化路由
	routers.Init()
	fmt.Printf("服务器启动成功，端口号：%s\n", setting.CommonSetting.HttpPort)

	if err := http.ListenAndServe(":"+setting.CommonSetting.HttpPort, nil); err != nil {
		panic(err)
	}
}
