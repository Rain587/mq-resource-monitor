package main

import (
	"go-mq/api"
)

func main() {
	httpServer := api.HttpServer{
		Port: 8080,
	}
	httpServer.Run()
}
