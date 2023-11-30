package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ih "github.com/priestess-dev/infra/http"
	"github.com/priestess-dev/widgets/recent-activities/cache"
	"github.com/priestess-dev/widgets/recent-activities/client"
	"github.com/priestess-dev/widgets/recent-activities/config"
	"github.com/priestess-dev/widgets/recent-activities/server"
	"github.com/priestess-dev/widgets/recent-activities/service"
	"net/http"
	"os"
)

func main() {
	// load config from file
	wd, err := os.Getwd()
	println("cwd: ", wd)
	appConfig := config.NewConfig("cmd/config.dev.yaml")
	fmt.Printf("[launcher] config: %+v, %v\n", appConfig.Redis.Addr, os.Getenv("REDIS_ADDR"))
	githubClient := client.NewClient()
	httpServer := server.NewServer(appConfig)
	redisCache, err := cache.NewRedis(appConfig)
	if err != nil {
		fmt.Printf("[launcher] load without redis: %s\n", err.Error())
	}

	githubService := service.NewService(appConfig, githubClient, redisCache)

	httpServer.AddRoutes(
		ih.EndpointConfig{
			Path:   "/activities",
			Method: "GET",
			Handler: func(ctx *gin.Context) {
				resp, err := githubService.ListEvents()
				if err != nil {
					ctx.Error(err)
					return
				}
				ctx.JSON(http.StatusOK, resp)
			},
		},
	)

	err = httpServer.Start()
	if err != nil {
		return
	}
}
