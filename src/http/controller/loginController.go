package controller

import (
	"global"
	"log"
	"net/http"
	"route"
)

type loginController struct{}

func (self loginController) RegisterRoute() {
	route.HandleFunc("/login", self.Login)
}

func (loginController) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", global.App.Host)
	global.App.LoginTcpServer(1)
	log.Printf("%v", global.App.Host)
	log.Printf("%v", "visit login success")
}

func connectTcpServer() {

}
