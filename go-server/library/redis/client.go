package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type ClientConf struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

type Client struct {
	*redis.Client
}

func NewClient(cf *ClientConf) (cli *Client, err error) {
	ocli := redis.NewClient(&redis.Options{
		Addr:     cf.Addr,
		Password: cf.Password,
		DB:       cf.DB,
		PoolSize: cf.PoolSize,
	})

	if _, err = ocli.Ping().Result(); err != nil {
		err = fmt.Errorf("client ping: %w", err)
		return
	}

	cli = &Client{
		Client: ocli,
	}

	return
}

func (cli *Client) Close() (err error) {
	if err = cli.Client.Close(); err != nil {
		err = fmt.Errorf("client close: %w", err)
		return
	}

	return
}
