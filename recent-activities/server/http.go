package server

import (
	ih "github.com/priestess-dev/infra/http"
	"github.com/priestess-dev/widgets/recent-activities/config"
)

type Server struct {
	ih.Server
}

func NewServer(config *config.Config) ih.Server {
	return ih.NewServer(config.Host, config.Port)
}
