package main

import (
	"auservices/api"
	"auservices/utilities"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	_, err := utilities.LoadConfiguration(os.Args[1])
	if err != nil {
		log.Panicln("no configuration file found with path specified:", os.Args[1])
	}
	port := 7777
	ltsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	resources := []api.MessageHandler{}
	subscribers := api.GetEventSubscribers().Subscribers
	for _, s := range subscribers {
		defer s.Close()
		resources = append(resources, s)
	}

	pub, err := api.CreateMessagePublisher(
		api.MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial pubhisher")
		os.Exit(1)
	}
	defer pub.Close()
	resources = append(resources, pub)

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
