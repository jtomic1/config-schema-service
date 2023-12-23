package repository

import (
	"context"
	"encoding/json"
	"time"

	pb "github.com/jtomic1/config-schema-service/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	schemaData := &pb.ConfigSchemaData{
		User:         user,
		Schema:       schema,
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
	return &schemaData, nil
}

func (repo *EtcdRepository) DeleteConfigSchema(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	res, err := repo.client.Delete(ctx, key)
	cancel()
	if res.Deleted > 0 {
		return nil
	}
	return err
}
