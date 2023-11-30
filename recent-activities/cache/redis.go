package cache

import (
	"github.com/priestess-dev/infra/cache"
	ir "github.com/priestess-dev/infra/cache/redis"
	"github.com/priestess-dev/widgets/recent-activities/config"
)

type redis struct {
	cache.Cache
}

func NewRedis(config *config.Config) (cache.Cache, error) {
	c, err := ir.NewCache(ir.Config{
		Addr:         string(config.Redis.Addr),
		Password:     string(config.Redis.Password),
		DB:           config.Redis.DB,
		Prefix:       string(config.Redis.Prefix),
		KeySeparator: string(config.Redis.KeySeparator),
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
