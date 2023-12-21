package repository

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	pb "github.com/jtomic1/config-schema-service/proto"
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
	if err != nil {
		log.Fatal(err)
	}
	return &EtcdRepository{
		client: cli,
	}, nil
}

func (repo *EtcdRepository) Close() {
	repo.client.Close()
}

func getConfigSchemaKey(req *pb.ConfigSchemaSaveRequest) string {
	return req.GetNamespace() + "-" + req.GetAppName() + "-" + req.GetConfigurationName() + "-" + req.GetVersion() + "-" + req.GetArch()
}

func (repo *EtcdRepository) SaveConfigSchema(req *pb.ConfigSchemaSaveRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	key := getConfigSchemaKey(req)
	resp, err := repo.client.Put(ctx, key, req.GetJsonSchema())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	return nil
}
