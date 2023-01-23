package database

import (
	"context"
	"github.com/go-redis/redis/v9"
	"net"
	"strconv"
)

type RedisConn struct {
	Client   *redis.Client
	address  string
	password string
	db       int
}

type RedisOptions struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisConn(opts RedisOptions) *RedisConn {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: opts.Password, // no password set
		DB:       opts.DB,       // use default DB
	})
	return &RedisConn{
		Client:   client,
		address:  address,
		password: opts.Password,
		db:       opts.DB,
	}
}

func (r *RedisConn) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
