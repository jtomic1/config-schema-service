package validators

import (
	"errors"

	pb "github.com/jtomic1/config-schema-service/proto"
)

func IsUserValid(user *pb.User) (bool, error) {
	if user == nil {
		return false, errors.New("User cannot be empty!")
	} else if user.Email == "" {
		return false, errors.New("User's email cannot be empty!")
	} else if user.Username == "" {
		return false, errors.New("User's username cannot be empty!")
	}
	return true, nil
}

func IsSchemaValid(schema string) (bool, error) {
	if schema == "" {
		return false, errors.New("Schema cannot be empty!")
	}
	return true, nil
}

func IsConfigurationValid(configuration string) (bool, error) {
	if configuration == "" {
		return false, errors.New("Configuration cannot be empty!")
	}
	return true, nil
}

func AreSchemaDetailsValid(schemaDetails *pb.ConfigSchemaDetails) (bool, error) {
	if schemaDetails == nil {
		return false, errors.New("Schema details cannot be empty!")
	} else if schemaDetails.Namespace == "" {
		return false, errors.New("Schema namespace cannot be empty!")
	} else if schemaDetails.SchemaName == "" {
		return false, errors.New("Schema name cannot be empty!")
	} else if schemaDetails.Version == "" {
		return false, errors.New("Schema version cannot be empty!")
	}
	return true, nil
}

func IsSaveSchemaRequestValid(saveRequest *pb.SaveConfigSchemaRequest) (bool, error) {
	userValid, userErr := IsUserValid(saveRequest.GetUser())
	if userErr != nil {
		return false, userErr
	}
	schemaDetailsValid, schemaDetailsErr := AreSchemaDetailsValid(saveRequest.GetSchemaDetails())
	if schemaDetailsErr != nil {
		return false, schemaDetailsErr
	}
	schemaValid, schemaErr := IsSchemaValid(saveRequest.GetSchema())
	if schemaErr != nil {
		return false, schemaErr
	}
	requestValid := userValid && schemaDetailsValid && schemaValid
	return requestValid, nil
}

func IsGetSchemaRequestValid(getRequest *pb.GetConfigSchemaRequest) (bool, error) {
	userValid, userErr := IsUserValid(getRequest.GetUser())
	if userErr != nil {
		return false, userErr
	}
	schemaDetailsValid, schemaDetailsErr := AreSchemaDetailsValid(getRequest.GetSchemaDetails())
	if schemaDetailsErr != nil {
		return false, schemaDetailsErr
	}

	requestValid := userValid && schemaDetailsValid
	return requestValid, nil
}

func IsDeleteSchemaRequestValid(deleteRequest *pb.DeleteConfigSchemaRequest) (bool, error) {
	userValid, userErr := IsUserValid(deleteRequest.GetUser())
	if userErr != nil {
		return false, userErr
	}
	schemaDetailsValid, schemaDetailsErr := AreSchemaDetailsValid(deleteRequest.GetSchemaDetails())
	if schemaDetailsErr != nil {
		return false, schemaDetailsErr
	}

	requestValid := userValid && schemaDetailsValid
	return requestValid, nil
}

func IsValidateConfigurationRequestValid(validateRequest *pb.ValidateConfigurationRequest) (bool, error) {
	userValid, userErr := IsUserValid(validateRequest.GetUser())
	if userErr != nil {
		return false, userErr
	}
	schemaDetailsValid, schemaDetailsErr := AreSchemaDetailsValid(validateRequest.GetSchemaDetails())
	if schemaDetailsErr != nil {
		return false, schemaDetailsErr
	}
	configurationValid, configurationErr := IsConfigurationValid(validateRequest.GetConfiguration())
	if configurationErr != nil {
		return false, configurationErr
	}
	requestValid := userValid && schemaDetailsValid && configurationValid
	return requestValid, nil
}

func IsGetConfigSchemaVersionsValid(versionsRequest *pb.ConfigSchemaVersionsRequest) (bool, error) {
	userValid, userErr := IsUserValid(versionsRequest.GetUser())
	if userErr != nil {
		return false, userErr
	}
	schemaDetailsValid, schemaDetailsErr := AreSchemaDetailsValid(versionsRequest.GetSchemaDetails())
	if schemaDetailsErr != nil {
		return false, schemaDetailsErr
	}

	requestValid := userValid && schemaDetailsValid
	return requestValid, nil
}
