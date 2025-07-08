package redis

import (
	"context"
	"strconv"
	"time"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"

	"github.com/go-redis/redis/v8"

	// in-app-cache
	"github.com/gin-contrib/cache/persistence"
)

func NewRedisCacheConnection(redisUrl string) (*redis.Client, context.Context) {
	// redis://username:password@host:port/db
	_, username, password, host, port, db := utils.UrlStringToOptions(redisUrl)
	// convert db to int
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		logger.CLogger.Errorln("Error parse redis db to int: db option is not a number ", err)
	}

	rContext := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: username,
		Password: password,
		DB:       dbInt,
	})
	// ping redis for check connection
	_, err = rdb.Ping(rContext).Result()
	if err != nil {
		logger.CLogger.Errorln("Error redis initial ping: ", err)
	}
	return rdb, rContext
}

func NewInAppCacheStore(defaultTTL time.Duration) *persistence.InMemoryStore {
	return persistence.NewInMemoryStore(defaultTTL)
}
