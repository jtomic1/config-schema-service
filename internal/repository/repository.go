package repository

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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

func (repo *EtcdRepository) SaveConfigSchema(key string, schema string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err := repo.client.Put(ctx, key, schema)
	cancel()
	return err
}

func (repo *EtcdRepository) GetConfigSchema(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := repo.client.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
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
