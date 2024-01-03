package repository

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	pb "github.com/jtomic1/config-schema-service/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/mod/semver"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/yaml"
)

var (
	endpoint = "localhost:2379"
	timeout  = 5 * time.Second
)

type EtcdRepository struct {
	client *clientv3.Client
}

func NewClient() (*EtcdRepository, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: timeout,
	})
	return &EtcdRepository{
		client: cli,
	}, err
}

func (repo *EtcdRepository) Close() {
	repo.client.Close()
}

func (repo *EtcdRepository) SaveConfigSchema(key string, user *pb.User, schema string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := repo.client.Get(ctx, key)
	if err != nil {
		return err
	}
	if res.Count > 0 {
		return errors.New("Key '" + key + "' already exists!")
	}
	schemaJson, err := yaml.YAMLToJSON([]byte(schema))
	if err != nil {
		return err
	}
	schemaData := &pb.ConfigSchemaData{
		User:         user,
		Schema:       string(schemaJson),
		CreationTime: timestamppb.New(time.Now()),
	}
	serializedData, err := json.Marshal(schemaData)
	if err != nil {
		return err
	}
	_, err = repo.client.Put(ctx, key, string(serializedData))
	return err
}

func (repo *EtcdRepository) GetConfigSchema(key string) (*pb.ConfigSchemaData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := repo.client.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var schemaData pb.ConfigSchemaData
	if err := json.Unmarshal(resp.Kvs[0].Value, &schemaData); err != nil {
		return nil, err
	}
	schemaYaml, err := yaml.JSONToYAML([]byte(schemaData.GetSchema()))
	if err != nil {
		return nil, err
	}
	schemaData.Schema = string(schemaYaml)
	return &schemaData, nil
}

func (repo *EtcdRepository) DeleteConfigSchema(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	if res.Deleted > 0 {
		return nil
	}
	return errors.New("No schema with key '" + key + "' found!")
}

func (repo *EtcdRepository) GetSchemasByPrefix(prefix string) ([]*pb.ConfigSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := repo.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	} else if res.Count == 0 {
		return nil, nil
	}
	schemas := make([]*pb.ConfigSchema, res.Count)
	for i, schemaKv := range res.Kvs {
		schemaDetails := getSchemaDetailsFromKey(string(schemaKv.Key))
		var schemaData pb.ConfigSchemaData
		if err := json.Unmarshal(schemaKv.Value, &schemaData); err != nil {
			return nil, err
		}
		schemaYaml, err := yaml.JSONToYAML([]byte(schemaData.GetSchema()))
		if err != nil {
			return nil, err
		}
		schemaData.Schema = string(schemaYaml)
		schemas[i] = &pb.ConfigSchema{
			SchemaDetails: schemaDetails,
			SchemaData:    &schemaData,
		}
	}
	sort.Slice(schemas, func(i, j int) bool {
		return semver.Compare(schemas[i].GetSchemaDetails().GetVersion(), schemas[j].GetSchemaDetails().GetVersion()) == -1
	})
	return schemas, nil
}

func (repo *EtcdRepository) GetLatestVersionByPrefix(prefix string) (string, error) {
	schemas, err := repo.GetSchemasByPrefix(prefix)
	if err != nil {
		return "", err
	}
	return schemas[len(schemas)-1].GetSchemaDetails().GetVersion(), nil
}

func getSchemaDetailsFromKey(key string) *pb.ConfigSchemaDetails {
	tokens := strings.Split(key, "/")
	return &pb.ConfigSchemaDetails{
		Namespace:  tokens[0],
		SchemaName: tokens[1],
		Version:    tokens[2],
	}
}
