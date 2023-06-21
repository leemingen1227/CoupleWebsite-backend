package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
    config := util.Config{
        TokenSymmetricKey:   util.RandomString(32),
        AccessTokenDuration: time.Minute,
    }
	asynqRedisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	redisOption := redis.Options{
		Addr: config.RedisAddress,
	}

	redisClient := redis.NewClient(&redisOption)

	taskDistributor := worker.NewRedisTaskDistributor(asynqRedisOpt)
    server, err := NewServer(config, store, taskDistributor, redisClient)
    require.NoError(t, err)

    return server
}

func TestMain(m *testing.M) {
    gin.SetMode(gin.TestMode)
    os.Exit(m.Run())
}