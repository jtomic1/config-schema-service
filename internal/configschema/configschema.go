package configschema

import (
	"context"

	"github.com/jtomic1/config-schema-service/internal/repository"
	pb "github.com/jtomic1/config-schema-service/proto"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
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

func (s *Server) ValidateConfiguration(ctx context.Context, in *pb.ValidateConfigurationRequest) (*pb.ValidateConfigurationResponse, error) {
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
	validationResult, err := validateConfiguration(in.GetConfiguration(), schemaData.GetSchema())
	if err != nil {
		message = "Error while validating configuration!"
	}
	if validationResult.Valid() && message == "" {
		message = "The configuration is valid!"
	} else {
		message = validationResult.Errors()[0].String()
	}

	return &pb.ValidateConfigurationResponse{
		Status:  status,
		Message: message,
		IsValid: validationResult.Valid(),
	}, nil
}

func validateConfiguration(configuration string, schema string) (*gojsonschema.Result, error) {
	configurationJson, err := yaml.YAMLToJSON([]byte(configuration))
	if err != nil {
		return nil, err
	}
	schemaJson, err := yaml.YAMLToJSON([]byte(schema))
	if err != nil {
		return nil, err
	}
	configLoader := gojsonschema.NewStringLoader(string(configurationJson))
	schemaLoader := gojsonschema.NewStringLoader(string(schemaJson))
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
