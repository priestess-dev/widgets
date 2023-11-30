package client

import (
	"github.com/priestess-dev/infra/thirdparty/github"
	"net/http"
)

type Client struct {
	github.Client
}

func NewClient() github.Client {
	return github.NewClient(http.DefaultClient, nil)
}
