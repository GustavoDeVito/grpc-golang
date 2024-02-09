package main

import (
	"log"
	"net"

	user "github.com/GustavoDeVito/grpc-golang/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	server := grpc.NewServer()

	service := &UserServer{}
	err = service.InitializeDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	user.RegisterUserServiceServer(server, service)

	reflection.Register(server)

	err = server.Serve(list)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
