package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
	"github.com/hibiken/asynq")

func newTestServer(t *testing.T, store db.Store) *Server {
    config := util.Config{
        TokenSymmetricKey:   util.RandomString(32),
        AccessTokenDuration: time.Minute,
    }
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
    server, err := NewServer(config, store, taskDistributor)
    require.NoError(t, err)

    return server
}

func TestMain(m *testing.M) {
    gin.SetMode(gin.TestMode)
    os.Exit(m.Run())
}