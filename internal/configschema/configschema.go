package configschema

import (
	"context"

	"github.com/jtomic1/config-schema-service/internal/repository"
	"github.com/jtomic1/config-schema-service/internal/validators"
	pb "github.com/jtomic1/config-schema-service/proto"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/mod/semver"
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
	return req.GetNamespace() + "/" + req.GetSchemaName() + "/" + req.GetVersion()
}

func getConfigSchemaPrefix(req ConfigSchemaRequest) string {
	return req.GetNamespace() + "/" + req.GetSchemaName()
}

func (s *Server) SaveConfigSchema(ctx context.Context, in *pb.SaveConfigSchemaRequest) (*pb.SaveConfigSchemaResponse, error) {
	_, err := validators.IsSaveSchemaRequestValid(in)
	if err != nil {
		return &pb.SaveConfigSchemaResponse{
			Status:  3,
			Message: err.Error(),
		}, nil
	}
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.SaveConfigSchemaResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
		}, nil
	}
	latestVersion, err := repoClient.GetLatestVersionByPrefix(getConfigSchemaPrefix(in.GetSchemaDetails()))
	if err != nil {
		return &pb.SaveConfigSchemaResponse{
			Status:  13,
			Message: err.Error(),
		}, nil
	}
	if semver.Compare(in.GetSchemaDetails().GetVersion(), latestVersion) != 1 {
		return &pb.SaveConfigSchemaResponse{
			Status:  3,
			Message: "Provided version is not latest! Please provide a version that succeeds '" + latestVersion + "'!",
		}, nil
	}
	err = repoClient.SaveConfigSchema(getConfigSchemaKey(in.GetSchemaDetails()), in.GetUser(), in.GetSchema())
	if err != nil {
		return &pb.SaveConfigSchemaResponse{
			Status:  13,
			Message: err.Error(),
		}, nil
	}
	return &pb.SaveConfigSchemaResponse{
		Status:  0,
		Message: "Schema saved successfully!",
	}, nil
}

func (s *Server) GetConfigSchema(ctx context.Context, in *pb.GetConfigSchemaRequest) (*pb.GetConfigSchemaResponse, error) {
	_, err := validators.IsGetSchemaRequestValid(in)
	if err != nil {
		return &pb.GetConfigSchemaResponse{
			Status:     3,
			Message:    err.Error(),
			SchemaData: nil,
		}, nil
	}
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.GetConfigSchemaResponse{
			Status:     13,
			Message:    "Error while instantiating database client!",
			SchemaData: nil,
		}, nil
	}
	key := getConfigSchemaKey(in.GetSchemaDetails())
	schemaData, err := repoClient.GetConfigSchema(key)
	if err != nil {
		return &pb.GetConfigSchemaResponse{
			Status:     13,
			Message:    "Error while retrieving schema!",
			SchemaData: nil,
		}, nil
	}
	var message string
	if schemaData == nil {
		message = "No schema with key '" + key + "' found!"
	} else {
		message = "Schema retrieved successfully!"
	}
	return &pb.GetConfigSchemaResponse{
		Status:     0,
		Message:    message,
		SchemaData: schemaData,
	}, nil
}

func (s *Server) DeleteConfigSchema(ctx context.Context, in *pb.DeleteConfigSchemaRequest) (*pb.DeleteConfigSchemaResponse, error) {
	_, err := validators.IsDeleteSchemaRequestValid(in)
	if err != nil {
		return &pb.DeleteConfigSchemaResponse{
			Status:  3,
			Message: err.Error(),
		}, nil
	}
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.DeleteConfigSchemaResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
		}, nil
	}
	if err := repoClient.DeleteConfigSchema(getConfigSchemaKey(in.GetSchemaDetails())); err != nil {
		return &pb.DeleteConfigSchemaResponse{
			Status:  3,
			Message: err.Error(),
		}, nil
	} else {
		return &pb.DeleteConfigSchemaResponse{
			Status:  0,
			Message: "Schema deleted successfully!",
		}, nil
	}
}

func (s *Server) ValidateConfiguration(ctx context.Context, in *pb.ValidateConfigurationRequest) (*pb.ValidateConfigurationResponse, error) {
	isValid, err := validators.IsValidateConfigurationRequestValid(in)
	if err != nil {
		return &pb.ValidateConfigurationResponse{
			Status:  3,
			Message: err.Error(),
			IsValid: isValid,
		}, nil
	}
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.ValidateConfigurationResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
			IsValid: false,
		}, nil
	}
	key := getConfigSchemaKey(in.GetSchemaDetails())
	schemaData, err := repoClient.GetConfigSchema(key)
	if err != nil {
		return &pb.ValidateConfigurationResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
			IsValid: false,
		}, nil
	} else if schemaData == nil {
		return &pb.ValidateConfigurationResponse{
			Status:  3,
			Message: "No schema with key '" + key + "' found!",
			IsValid: false,
		}, nil
	}
	validationResult, err := validateConfiguration(in.GetConfiguration(), schemaData.GetSchema())
	if err != nil {
		return &pb.ValidateConfigurationResponse{
			Status:  3,
			Message: "Error while validating schema!",
			IsValid: false,
		}, nil
	}
	var message string
	if validationResult.Valid() && message == "" {
		message = "The configuration is valid!"
	} else {
		message = validationResult.Errors()[0].String()
	}

	return &pb.ValidateConfigurationResponse{
		Status:  0,
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

func (s *Server) GetConfigSchemaVersions(ctx context.Context, in *pb.ConfigSchemaVersionsRequest) (*pb.ConfigSchemaVersionsResponse, error) {
	_, err := validators.IsGetConfigSchemaVersionsValid(in)
	if err != nil {
		return &pb.ConfigSchemaVersionsResponse{
			Status:  3,
			Message: err.Error(),
		}, nil
	}
	repoClient, err := repository.NewClient()
	defer repoClient.Close()
	if err != nil {
		return &pb.ConfigSchemaVersionsResponse{
			Status:  13,
			Message: "Error while instantiating database client!",
		}, nil
	}
	key := getConfigSchemaPrefix(in.GetSchemaDetails())
	schemaVersions, err := repoClient.GetSchemasByPrefix(key)
	if err != nil {
		return &pb.ConfigSchemaVersionsResponse{
			Status:  13,
			Message: "Error while retrieving schema!",
		}, nil
	}
	var message string
	if schemaVersions == nil {
		message = "No schema with prefix '" + key + "' found!"
	} else {
		message = "Schema versions retrieved successfully!"
	}
	return &pb.ConfigSchemaVersionsResponse{
		Status:         0,
		Message:        message,
		SchemaVersions: schemaVersions,
	}, nil
}
