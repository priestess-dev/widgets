package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/priestess-dev/infra/cache"
	"github.com/priestess-dev/infra/thirdparty/github"
	"github.com/priestess-dev/widgets/recent-activities/config"
	"time"
)

type Service interface {
	ListEvents() (github.GetUserEventsResponse, error)
}

type service struct {
	client github.Client
	config *config.Config
	cache  cache.Cache
}

func (s *service) ListEvents() (github.GetUserEventsResponse, error) {
	if s.cache != nil {
		exist, raw, err := s.cache.Load(context.Background(), "github.events")
		if err == nil && exist {
			var resp github.GetUserEventsResponse
			err = json.Unmarshal([]byte(raw), &resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
	resp, err := s.client.ListPublicEventsForUser(github.GetUserEventsRequest{
		Username: string(s.config.Github.Username),
	})

	if err != nil {
		return nil, err
	}
	// save to redis
	if s.cache != nil {
		raw, err := json.Marshal(resp)
		if err == nil {
			err = s.cache.StoreEX(context.Background(), "github.events", string(raw), time.Second*s.config.Redis.TTL)
			if err != nil {
				fmt.Printf("Error save to redis %s\n", err.Error())
			}
		}
	}
	return resp, nil
}

func NewService(config *config.Config, client github.Client, cache cache.Cache) Service {
	return &service{
		client: client,
		config: config,
		cache:  cache,
	}
}
