package configschema

import (
	"context"

	"github.com/jtomic1/config-schema-service/internal/repository"
	pb "github.com/jtomic1/config-schema-service/proto"
)

type Server struct {
	pb.UnimplementedConfigSchemaServiceServer
}

type ConfigSchemaRequest interface {
	GetNamespace() string
	GetSchemaName() string
	GetVersion() string
}

func NewServer() *Server {
	return &Server{}
}

func getConfigSchemaKey(req ConfigSchemaRequest) string {
	return req.GetNamespace() + "-" + req.GetSchemaName() + "-" + req.GetVersion()
}

func (s *Server) SaveConfigSchema(ctx context.Context, in *pb.SaveConfigSchemaRequest) (*pb.SaveConfigSchemaResponse, error) {
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	var status int32 = 0
	message := ""
	if err != nil {
		status = 13
		message = "Error while instantiating database client!"
	}
	err = repoClient.SaveConfigSchema(getConfigSchemaKey(in.GetSchemaDetails()), in.GetUser(), in.GetSchema())
	if err != nil {
		status = 13
		message = "Error while saving schema!"
	}
	if message == "" {
		message = "Schema saved successfully!"
	}
	return &pb.SaveConfigSchemaResponse{
		Status:  status,
		Message: message,
	}, nil
}

func (s *Server) GetConfigSchema(ctx context.Context, in *pb.GetConfigSchemaRequest) (*pb.GetConfigSchemaResponse, error) {
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	var status int32 = 0
	message := ""
	if err != nil {
		status = 13
		message = "Error while instantiating database client!"
	}
	key := getConfigSchemaKey(in.GetSchemaDetails())
	schemaData, err := repoClient.GetConfigSchema(key)
	if err != nil {
		status = 13
		message = "Error while retrieving schema!"
	}
	if schemaData == nil {
		message = "No schema with key '" + key + "' found!"
	} else {
		message = "Schema retrieved successfully!"
	}
	return &pb.GetConfigSchemaResponse{
		Status:     status,
		Message:    message,
		SchemaData: schemaData,
	}, nil
}

func (s *Server) DeleteConfigSchema(ctx context.Context, in *pb.DeleteConfigSchemaRequest) (*pb.DeleteConfigSchemaResponse, error) {
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	var status int32 = 0
	message := ""
	if err != nil {
		status = 13
		message = "Error while instantiating database client!"
	}
	if err := repoClient.DeleteConfigSchema(getConfigSchemaKey(in.GetSchemaDetails())); err != nil {
		status = 13
		message = "Error while deleting schema!"
	}
	message = "Schema deleted successfully!"
	return &pb.DeleteConfigSchemaResponse{
		Status:  status,
		Message: message,
	}, nil
}
