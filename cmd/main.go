package main

import (
	"fmt"
	"food/api"
	"food/config"
	"food/pkg/logger"
	"food/service"

	postgres "food/storage/postgres"
	"food/storage/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	var loggerLevel = new(string)

	*loggerLevel = logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	pgconn, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic("postgres no connection: " + err.Error())
	}
	defer pgconn.CloseDB()

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: RedisAddr,
	// })

	newRedis := redis.New(cfg)
	services := service.New(pgconn, log, newRedis)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	api.NewApi(r, &cfg, pgconn, log, services)

	fmt.Println("Listening server", cfg.PostgresHost+cfg.HTTPPort)
	err = r.Run(cfg.HTTPPort)
	// err = r.Run(cfg.ServerHost + cfg.HTTPPort)
	if err != nil {
		panic(err)
	}

	// // Check Redis connection
	// _, err := rdb.Ping(ctx).Result()
	// if err != nil {
	//  log.Fatalf("Could not connect to Redis: %v", err)
	// }
}
