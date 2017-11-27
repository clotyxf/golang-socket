package controller

import (
	"config"
	"global"
	"http/service"
	"log"
	"net/http"
	"route"
	//"golang.org/x/net/websocket"
)

type socketController struct{}

var hub *service.Hub

func (this socketController) RegisterRoute() {
	route.HandleFunc("/socket", this.SocketClient)
}

func (socketController) SocketClient(w http.ResponseWriter, r *http.Request) {
	//读取配置
	configFile := "config/env.yml"
	config.YmlFileRead(configFile)
	global.App.SocketHost = config.YamlConfig.Get("socket.host").String()
	global.App.SocketPort = config.YamlConfig.Get("socket.port").String()
	log.Printf("%v", "start socket_client connect")
	log.Printf("socket_ip:%v", global.App.SocketHost)
	log.Printf("socket_port:%v", global.App.SocketPort)

	service.ServeWs(hub, w, r)
}
