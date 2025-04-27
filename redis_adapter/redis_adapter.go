package redis_adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	Ttl         time.Duration `yaml:"ttl"`
}

var Cfg = Config{
	Addr:        "localhost:6379",
	Password:    "",
	User:        "",
	DB:          0,
	MaxRetries:  5,
	DialTimeout: 10 * time.Second,
	Timeout:     5 * time.Second,
	Ttl:         30 * time.Second,
}

type LRUCache struct {
	client *redis.Client
	ttl    time.Duration
}

var WeatherCache *LRUCache = nil

func NewLRUCache(ctx context.Context) (*LRUCache, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         Cfg.Addr,
		Password:     Cfg.Password,
		DB:           Cfg.DB,
		Username:     Cfg.User,
		MaxRetries:   Cfg.MaxRetries,
		DialTimeout:  Cfg.DialTimeout,
		ReadTimeout:  Cfg.Timeout,
		WriteTimeout: Cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return &LRUCache{client: db, ttl: Cfg.Ttl}, nil
}

func (c *LRUCache) Set(ctx context.Context, key string, value []byte) error {
	return c.client.Set(ctx, key, value, c.ttl).Err()
}

func (c *LRUCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()

	return val, err
}
