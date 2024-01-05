# Quasar - Remote Configuration Schema Service

This remote service offers a suite of functionalities designed to simplify YAML schema management and validation processes. The service specializes in CRUD operations for YAML schemas, schema validation against specified structures, and accessing version history for a particular schema. 

The following remote procedures are available:

 - **ConfigSchemaService/SaveConfigSchema**
 - **ConfigSchemaService/GetConfigSchema**
 - **ConfigSchemaService/DeleteConfigSchema**
 - **ConfigSchemaService/ValidateConfiguration**
 - **ConfigSchemaService/GetConfigSchemaVersions**

## Installation Guide

Prerequisites:
 - Go 1.20: Ensure you have Go version 1.20 installed on your machine. You can download and install it from the [official Go website](https://go.dev/).
 - etcd: Install and run etcd, the distributed key-value store. [Installation guide](https://etcd.io/docs/v3.5/install/)

Installation:
 1. Clone this repository.
 2. Navigate to the following location
    ```
    cd <repository_directory>\cmd\server
    ```
 3. Start etcd using the appropriate commands for your operating system. For example:
    ```
    etcd
    ```
 4. Run the server by double-clicking server.exe or running the command
    ```
    ./server.exe
    ```
Once the Go application is running, the service will be available.

*Notes*

 - Ensure etcd is running and accessible at localhost:2379 before starting the Go application.
 - The default port for the server is 50051


## ConfigSchemaService/SaveConfigSchema
This procedure is used to create a new schema. 
### Request
**SaveConfigSchema** accepts a message of type **SaveConfigSchemaRequest**, which consists of the following fields, all of which are <u>required</u>.
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| user    | [User](#user)  | User which has requested to save the schema |
| schema_details    | [ConfigSchemaDetails](#config-schema-details)  | Details regarding the schema (namespace, name of the schema, and schema version) |
|schema | string | YAML string representing the schema. Must be convertible into a valid JSON Schema format.|
### Response
**SaveConfigSchema** returns a message of type **SaveConfigSchemaResponse**, which consists of the following fields
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| status    | int32  | [gRPC Status Code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) |
| message   | string  | Response details |

### Example Usage
#### Example 1 - Valid Request
The following example demonstrates a successful request with no errors. 

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  },
  "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n      - city\n      - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n      - name\n      - age\n    type: object\nrequired:\n  - address\n  - person\ntype: object"
}
```
Response:
```json
{
	"status": 0,
	"message": "Schema saved successfully!"
}
```
The value provided in the "schema" field is stored in JSON format under the key **my_namespace/person_address_schema/v1.0.0**

#### Example 2 - Invalid Version
Consider the previous example. Since a schema under the key **my_namespace/person_address_schema/v1.0.0** has already been stored, trying to store another schema with the version that doesn't succeed **v1.0.0** would result in an error. In this example, version **v0.1.0** is used. Due to this restriction, it would also be impossible to send multiple requests where schema details are the same.

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v0.1.0"
  },
  "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n      - city\n      - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n      - name\n      - age\n    type: object\nrequired:\n  - address\n  - person\ntype: object"
}
```
Response:
```json
{
	"status": 3,
	"message": "Provided version is not latest! Please provide a version that succeeds 'v1.0.0'!"
}
```
#### Example 3 - Invalid Schema Value
In the following example, the value for "schema" is invalid (**Lorem ipsum** - an ordinary string). The server will reject this request because this string cannot be converted into a valid JSON Schema.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v2.0.0"
  },
  "schema": "Lorem ipsum"
}
```
Response:
```json
{
	"status": 3,
	"message": "schema is invalid"
}
```

#### Example 4 - Missing Mandatory Fields
In the following example, the **user** field is omitted from the request. Since it is a mandatory field, the server will reject this request.

Request:
```json 
{
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v2.0.0"
  },
  "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n      - city\n      - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n      - name\n      - age\n    type: object\nrequired:\n  - address\n  - person\ntype: object"
}
```
Response:
```json
{
	"status": 3,
	"message": "User cannot be empty!"
}
```
#### Example 5 - Invalid Character In Schema Details
Currently, only one character is considered illegal when providing schema details - the forward slash. Since it is used as a separator when generating keys for the database, the user is prohibited from including it inside the "schema_details". In this example, it is included inside the namespace.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my/namespace",
    "schema_name": "person_address_schema",
    "version": "v2.0.0"
  },
  "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n      - city\n      - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n      - name\n      - age\n    type: object\nrequired:\n  - address\n  - person\ntype: object"
}
```
Response:
```json
{
	"status": 3,
	"message": "Schema details must not contain '/'!"
}
```
#### Example 6 - Invalid SemVer In Schema Details
The value provided for the schema version must adhere to certain rules. Check out [ConfigSchemaDetails](#config-schema-details) section for more information about the restrictions imposed on this field. The following example provides an ordinary integer in the version field, which is invalid since it doesn't have the "v" prefix.
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "3"
  },
  "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n      - city\n      - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n      - name\n      - age\n    type: object\nrequired:\n  - address\n  - person\ntype: object"
}
```
Response:
```json
{
	"status": 3,
	"message": "Schema version must be a valid SemVer string with 'v' prefix!"
}
```
## ConfigSchemaService/GetConfigSchema
This procedure is used to retrieve a schema.
### Request
**GetConfigSchema** accepts a message of type **GetConfigSchemaRequest**, which consists of the following fields, all of which are <u>required</u>.
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| user    | [User](#user)  | User which has requested to get the schema |
| schema_details    | [ConfigSchemaDetails](#config-schema-details)  | Details regarding the schema (namespace, name of the schema, and schema version) |
### Response
**GetConfigSchema** returns a message of type **GetConfigSchemaResponse**, which consists of the following fields
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| status    | int32  | [gRPC Status Code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) |
| message   | string  | Response details |
|schema_data|[ConfigSchemaData](#config-schema-data)|Contains the schema value, as well as the creation time and the author

### Example Usage
#### Example 1 - Valid Request
The following example demonstrates a successful request with no errors. 

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 0,
  "message": "Schema retrieved successfully!",
  "schema_data": {
    "user": {
      "username": "johndoe",
      "email": "johndoe@example.com"
    },
    "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n    - city\n    - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n    - name\n    - age\n    type: object\nrequired:\n- address\n- person\ntype: object\n",
    "creation_time": {
      "seconds": "1704365029",
      "nanos": 54870600
    }
  }
}
```
#### Example 2 - Non-existent Key
In the following example, the query contains a schema name which does not exist.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "car_schema",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 0,
  "message": "No schema with key 'my_namespace/car_schema/v1.0.0' found!",
  "schema_data": null
}
```

#### Example 3 - Missing Mandatory Fields
In the following example, the **schema_name** field is omitted from the request. Since it is a mandatory field, the server will reject this request.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 3,
  "message": "Schema name cannot be empty!",
  "schema_data": null
}
```
## ConfigSchemaService/DeleteConfigSchema
This procedure is used to physically delete a schema.
### Request
**DeleteConfigSchema** accepts a message of type **DeleteConfigSchemaRequest**, which consists of the following fields, all of which are <u>required</u>.
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| user    | [User](#user)  | User which has requested to delete the schema |
| schema_details    | [ConfigSchemaDetails](#config-schema-details)  | Details regarding the schema (namespace, name of the schema, and schema version) |
### Response
**DeleteConfigSchema** returns a message of type **DeleteConfigSchemaResponse**, which consists of the following fields
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| status    | int32  | [gRPC Status Code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) |
| message   | string  | Response details |

### Example Usage
#### Example 1 - Valid Request
The following example demonstrates a successful request with no errors. 

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 0,
  "message": "Schema deleted successfully!"
}
```
#### Example 2 - Non-existent Key
In the following example, the query contains a schema name which does not exist.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "car_schema",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 0,
  "message": "No schema with key 'my_namespace/car_schema/v1.0.0' found!"
}
```

#### Example 3 - Missing Mandatory Fields
In the following example, the **username** field is omitted from the request. Since it is a mandatory field, the server will reject this request.

Request:
```json 
{
  "user": {
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  }
}
```
Response:
```json
{
  "status": 3,
  "message": "User's username cannot be empty!"
}
```
## ConfigSchemaService/ValidateConfiguration
This procedure is used to validate a provided YAML  string (configuration) against a previously created schema.
### Request
**ValidateConfiguration** accepts a message of type **ValidateConfigurationRequest**, which consists of the following fields, all of which are <u>required</u>.
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| user    | [User](#user)  | User which has requested to validate the configuration |
| schema_details    | [ConfigSchemaDetails](#config-schema-details)  | Details regarding the schema (namespace, name of the schema, and schema version) |
|configuration|string|A YAML string which represents the configuration that should be validated
### Response
**ValidateConfiguration** returns a message of type **ValidateConfigurationResponse**, which consists of the following fields
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| status    | int32  | [gRPC Status Code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) |
| message   | string  | Response details |
|is_valid | boolean | Validation result (true if the configuration is valid, false otherwise)

### Example Usage
#### Example 1 - Valid Request, Valid Configuration
The following example demonstrates a successful request with no errors for a valid configuration. 

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  },
  "configuration": "person:\n  name: John Doe\n  age: 23\naddress:\n  city: New York\n  country: USA"
}
```
Response:
```json
{
  "status": 0,
  "message": "The configuration is valid!",
  "is_valid": true
}
```
#### Example 2 - Valid Request, Invalid Configuration
The following example demonstrates a successful request with no errors for an invalid configuration. In the following example, the value for "age" is provided as a string, whereas the schema requires the "age" to be an integer. 

Request: 
```json
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema",
    "version": "v1.0.0"
  },
  "configuration": "person:\n  name: John Doe\n  age: twenty\naddress:\n  city: New York\n  country: USA"
}
```
Response:
```json
{
  "status": 0,
  "message": "person.age: Invalid type. Expected: integer, given: string",
  "is_valid": false
}
```

#### Example 3 - Non-existent Key
In the following example, the query contains a schema name which does not exist.

Request:
```json 
{
  "user": {
    "username": "johndoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "car_schema",
    "version": "v1.0.0"
  },
  "configuration": "person:\n  name: John Doe\n  age: 20\naddress:\n  city: New York\n  country: USA"
}
```
Response:
```json
{
  "status": 3,
  "message": "No schema with key 'my_namespace/car_schema/v1.0.0' found!",
  "is_valid": false
}
```
<br>
Omitting a required field is handled in the same manner as in previous endpoints. Naturally, the "is_valid" field in this case is always going to be false.

## ConfigSchemaService/GetConfigSchemaVersions
This procedure is used to retrieve all schemas under the given namespace and schema name. Schema array in the response is sorted in ascending order with respect to schemas' semantic version.
### Request
**GetConfigSchemaVersions** accepts a message of type **ConfigSchemaVersionsRequest**, which consists of the following fields, all of which are <u>required</u>.
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| user    | [User](#user)  | User which has requested to validate the configuration |
| schema_details    | [ConfigSchemaDetails](#config-schema-details)  | Details regarding the schema (namespace and schema name). Note that in this case, the "version" field is NOT required |
### Response
**GetConfigSchemaVersions** returns a message of type **ConfigSchemaVersionsResponse**, which consists of the following fields
|parameter| type  |                    description              |
|---------|-------|---------------------------------------------|
| status    | int32  | [gRPC Status Code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) |
| message   | string  | Response details |
|schema_versions | Array of [ConfigSchema](#config-schema) objects| Sorted array of [ConfigSchema](#config-schema), which includes schema details and schema value, as well as the author and creation time for each version
### Example Usage
#### Example 1 - Valid Request
The following example demonstrates a successful request with no errors. 

Request: 
```json
{
  "user": {
    "username": "johdnoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "person_address_schema"
  }
}
```
Response:
```json
{
  "schema_versions": [
    {
      "schema_details": {
        "namespace": "my_namespace",
        "schema_name": "person_address_schema",
        "version": "v1.0.0"
      },
      "schema_data": {
        "user": {
          "username": "johndoe",
          "email": "johndoe@example.com"
        },
        "schema": ,
        "creation_time": {
          "seconds": "1704447813",
          "nanos": 832872900
        }
      }
    },
    {
      "schema_details": {
        "namespace": "my_namespace",
        "schema_name": "person_address_schema",
        "version": "v2.0.0"
      },
      "schema_data": {
        "user": {
          "username": "johndoe",
          "email": "johndoe@example.com"
        },
        "schema": "properties:\n  address:\n    properties:\n      town:\n        type: string\n      country:\n        type: string\n    required:\n    - city\n    - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n    - name\n    - age\n    type: object\nrequired:\n- address\n- person\ntype: object\n",
        "creation_time": {
          "seconds": "1704449682",
          "nanos": 683805700
        }
      }
    },
    {
      "schema_details": {
        "namespace": "my_namespace",
        "schema_name": "person_address_schema",
        "version": "v3.0.0"
      },
      "schema_data": {
        "user": {
          "username": "johndoe",
          "email": "johndoe@example.com"
        },
        "schema": "properties:\n  address:\n    properties:\n      city:\n        type: string\n      country:\n        type: string\n    required:\n    - city\n    - country\n    type: object\n  person:\n    properties:\n      age:\n        type: integer\n      name:\n        type: string\n    required:\n    - name\n    - age\n    type: object\nrequired:\n- address\n- person\ntype: object\n",
        "creation_time": {
          "seconds": "1704449731",
          "nanos": 538727700
        }
      }
    }
  ],
  "status": 0,
  "message": "Schema versions retrieved successfully!"
}
```
#### Example 2 - Non-existent Prefix
The following example demonstrates a successful request with no schemas found for the given prefix. 

Request: 
```json
{
  "user": {
    "username": "johdnoe",
    "email": "johndoe@example.com"
  },
  "schema_details": {
    "namespace": "my_namespace",
    "schema_name": "car_schema"
  }
}
```
Response:
```json
{
	"schema_versions": [],
	"status": 0,
	"message": "No schema with prefix 'my_namespace/car_schema' found!"
}
```
<br>
Omitting a required field is handled in the same manner as in previous endpoints. Naturally, the "schema_versions" field in this case is always going to be an empty array.

## Custom Types
This section further describes custom types and messages which are defined in the service.
### <a name="user"></a> User
|property| type  |restrictions|             description              |
|---------|-------|-------|--------------------------------------|
| username    | string |Cannot be empty|User's username |
| email   | string |Cannot be empty| User's email |
---
### <a name="config-schema-details"></a> ConfigSchemaDetails
|property| type  |restrictions|          description              |
|---------|-------|---------|------------------------------------|
| namespace    | string | Cannot be empty<br>Cannot contain "/"| Namespace which the schema belongs to|
| schema_name   | string | Cannot be empty<br>Cannot contain "/" | Schema name |
|version|string|Cannot be empty*<br>Cannot contain "/"<br>Must be a valid SemVer string with "v" prefix [(more info about accepted version inputs)](https://pkg.go.dev/golang.org/x/mod/semver#pkg-overview)|Schema version|

**Note: Version CAN be omitted when sending a request to **ConfigSchemaService/GetConfigSchemaVersions** endpoint*

---
### <a name="config-schema-data"></a> ConfigSchemaData
|property| type  |   restrictions  |               description              |
|---------|-------|-------|-------------------------------------|
| user    | [User](#user) |Cannot be empty | User which has created the schema|
| schema   | string  |Must be a non-empty YAML string which can be converted to a valid JSON Schema| Schema value in YAML format |
|creation_time|[timestamppb.Timestamp](https://pkg.go.dev/google.golang.org/protobuf/types/known/timestamppb#Timestamp)| Cannot be empty|Time at which the schema was created|
---
### <a name="config-schema"></a> ConfigSchema
|property| type  |   restrictions  |               description              |
|---------|-------|-------|-------------------------------------|
| schema_details    | [ConfigSchemaDetails](#config-schema-details) |Cannot be empty | Schema details|
| schema_data| [ConfigSchemaData](#config-schema-data)  |Cannot be empty| Schema data |
