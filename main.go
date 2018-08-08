package main

import (
	"auservices/api"
	"auservices/msghandler"
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
	path := "./config/keys.dev.json"
	if len(os.Args) < 2 {
		log.Println(fmt.Sprintf("no configuration file path specified in argument, using default path: %s", path))
	} else {
		path = os.Args[1]
	}

	cfg, err := utilities.LoadConfiguration(path)
	if err != nil {
		log.Panicln("no configuration file found with path specified:", path)
	}
	port := cfg.ApplicationPort
	ltsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	resources := []msghandler.MessageHandler{}
	//register subscribers
	subscribers := msghandler.GetEventSubscribers().Subscribers
	for _, s := range subscribers {
		defer s.Close()
		resources = append(resources, s)
	}
	//register publishers
	pub, err := msghandler.CreateMessagePublisher(
		msghandler.MessageHanderInfo{},
	)
	if err != nil {
		log.Fatalf("could not initial pubhisher")
	}
	defer pub.Close()
	resources = append(resources, pub)

	si := msghandler.Server{
		MsgPublisher: pub,
	}
	grpcServer := grpc.NewServer()
	api.RegisterPingServer(grpcServer, &si)
	api.RegisterCategoryServicesServer(grpcServer, &si)
	api.RegisterItemLineServiceServer(grpcServer, &si)

	log.Println("serving at", port)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func(rs []msghandler.MessageHandler) {
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
