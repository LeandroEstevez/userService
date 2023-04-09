package main

import (
	"log"
	"userMicroService/api"
	"userMicroService/events"
	"userMicroService/util"
)

func main() {
	events.SetUp()
	// events.SetupProducer()
	util.SetUpConnAndStore()

	server, err := api.NewServer(util.Conf, util.Store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(util.Conf.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
