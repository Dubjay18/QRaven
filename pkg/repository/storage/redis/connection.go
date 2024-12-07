package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/url"
	"qraven/internal/config"
	"qraven/pkg/repository/storage"
	"qraven/utils"
	"strconv"
)

var (
	Ctx     = context.Background()
	KeyName = "EmailQueue"
)

func ConnectToRedis(logger *utils.Logger, configDatabases config.Redis) *redis.Client {
	dbsCV := configDatabases
	utils.LogAndPrint(logger, "connecting to redis server")
	connectedServer := connectToDb(dbsCV.REDIS_HOST, dbsCV.REDIS_PORT, dbsCV.REDIS_DB, dbsCV.REDIS_PASSWORD, logger)

	utils.LogAndPrint(logger, "connected to redis server")

	storage.DB.Redis = connectedServer

	return connectedServer
}

func connectToDb(host, port, db string, password string, logger *utils.Logger) *redis.Client {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utils.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utils.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
	dbInst, err := strconv.Atoi(db)
	if err != nil {
		utils.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbInst,
	})

	if err := redisClient.Ping(Ctx).Err(); err != nil {
		utils.LogAndPrint(logger, fmt.Sprintln(addr))
		utils.LogAndPrint(logger, fmt.Sprintln("Redis db error: ", err))
	}

	pong, _ := redisClient.Ping(Ctx).Result()
	utils.LogAndPrint(logger, fmt.Sprintln("Redis says: ", pong))

	return redisClient
}
