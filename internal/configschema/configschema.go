package configschema

import (
	"context"
	"log"

	"github.com/jtomic1/config-schema-service/internal/repository"
	pb "github.com/jtomic1/config-schema-service/proto"
)

type Server struct {
	pb.UnimplementedConfigSchemaServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SaveConfigSchema(ctx context.Context, in *pb.ConfigSchemaSaveRequest) (*pb.ConfigSchemaSaveResponse, error) {
	log.Println("Invoked Save Config Schema: ", in)
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.ConfigSchemaSaveResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
		}, err
	}
	repoClient.SaveConfigSchema(in)
	return &pb.ConfigSchemaSaveResponse{
		Status:  0,
		Message: "Configuration saved successfully!",
	}, nil
}
