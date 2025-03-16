package repositories

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type RedisDriver struct {
	conf   *config.RedisConfig
	secret *config.RedisSecret
	client *redis.Client
}

func NewRedisDriver(
	lc fx.Lifecycle,
	config *config.Config,
	secrets *config.Secrets,
) *RedisDriver {
	driver := &RedisDriver{
		conf:   config.Redis,
		secret: secrets.Redis,
	}

	lc.Append(fx.Hook{
		OnStart: driver.Start,
		OnStop:  driver.Stop,
	})

	return driver
}

func (d *RedisDriver) Start(ctx context.Context) error {
	d.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", d.conf.Host, d.conf.Port),
		Username: d.secret.User,
		Password: d.secret.Password,
	})

	_, err := d.client.Do(ctx, "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		return fmt.Errorf("change config error: %w", err)
	}

	return nil
}

func (d *RedisDriver) Stop(ctx context.Context) error {
	return d.client.Close()
}

func (d *RedisDriver) Unwrap() *redis.Client {
	return d.client
}
