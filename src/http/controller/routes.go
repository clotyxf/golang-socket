package controller

import (
	"http/service"
)

func RegisterRoutes() {
	new(indexController).RegisterRoute()
	new(loginController).RegisterRoute()
	new(socketController).RegisterRoute()
}

func RegisterSockets(h *service.Hub) {
	hub = h
	new(socketController).RegisterRoute()
}
