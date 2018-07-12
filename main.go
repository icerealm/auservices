package main

import (
	"auservices/api"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	port := 7777
	ltsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pub, err := api.CreateMessagePublisher(
		api.MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial pubhisher")
		os.Exit(1)
	}
	defer pub.Close()

	sub, err := api.CreateMessageSubscriber(
		api.MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial subscriber")
		os.Exit(1)
	}
	defer sub.Close()
	sub.SubscribeEvent("test-chan", nil) //test

	resources := []api.MessageHandler{pub, sub}

	si := api.Server{
		MsgPublisher: pub,
	}
	grpcServer := grpc.NewServer()
	api.RegisterPingServer(grpcServer, &si)
	api.RegisterCategoryServicesServer(grpcServer, &si)

	log.Println("serving at", port)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func(rs []api.MessageHandler) {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			log.Println("SIGTERM,SIGINT interupt - cleanup...")
			for _, r := range rs {
				r.Close()
			}
		} else {
			log.Println(s)
		}
		os.Exit(1)
	}(resources)

	if err := grpcServer.Serve(ltsnr); err != nil {
		log.Fatalf("Failed to server: %s", err)
	}
}
