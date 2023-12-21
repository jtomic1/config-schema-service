package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jtomic1/config-schema-service/internal/configschema"
	pb "github.com/jtomic1/config-schema-service/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type configSchemaServer struct {
	pb.UnimplementedConfigSchemaServiceServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	configSchemaServer := configschema.NewServer()

	pb.RegisterConfigSchemaServiceServer(grpcServer, configSchemaServer)
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
