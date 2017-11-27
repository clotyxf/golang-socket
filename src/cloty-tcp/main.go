package main

import (
	"config"
	"global"
	"http/controller"
	"http/service"
	"log"
	"net/http"
	"os"
	"route"
	"strings"
)

func main() {
	//读取配置
	configFile := "config/env.yml"
	config.YmlFileRead(configFile)
	var localhost string

	global.App.Host = config.YamlConfig.Get("listen.host").String()
	global.App.Port = config.YamlConfig.Get("listen.port").String()
	localhost = global.App.Host + ":" + global.App.Port

	if len(os.Args) == 2 && os.Args[1] == "socket" {
		global.App.SocketHost = config.YamlConfig.Get("socket.host").String()
		global.App.SocketPort = config.YamlConfig.Get("socket.port").String()
		log.Printf("%v", "start socket connect")
		localhost = global.App.SocketHost + ":" + global.App.SocketPort

		hub := service.NewHub()
		go hub.Run()
		controller.RegisterSockets(hub)
	} else {
		global.App.InitPath()
		//router mux
		controller.RegisterRoutes()

		fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(global.App.ProjectRoot+"/static")))
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			reqURI := r.RequestURI
			if strings.HasSuffix(reqURI, "/") { //以/结尾的URL，直接返回404
				http.NotFound(w, r)
			} else {
				fileHandler.ServeHTTP(w, r)
			}
		})
	}
	//config read

	log.Printf("%v", localhost)

	s := "loadinig..."
	log.Printf("%v", s)

	log.Fatal(http.ListenAndServe(localhost, route.DefaultMux))
}
