package rediscron

import (
	"Open_IM/cmd/userscoreupload/redlock"
	"github.com/go-redis/redis/v8"
)

type RedisMutexBuilder interface {
	NewMutex(pfx string) (*redlock.RedLock, error)
	GetClient() *redis.Client
}

type Config struct {
	Address string `json:"address"`
	Secret  string `json:"secret"`
	DB      int    `json:"db"`
}
type redisMutexBuilder struct {
	*redis.Client
}

func NewRedisMutexBuilder(conf Config) (*redisMutexBuilder, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Secret,
		DB:       conf.DB,
	})
	return &redisMutexBuilder{Client: c}, nil
}

func (c *redisMutexBuilder) NewMutex(pfx string) (*redlock.RedLock, error) {
	return redlock.New(c.Client, pfx, 5000, true)
}
func (c *redisMutexBuilder) GetClient() *redis.Client {
	return c.Client
}
