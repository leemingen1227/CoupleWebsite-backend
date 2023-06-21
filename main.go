package main

import (
	"database/sql"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hibiken/asynq"
	"github.com/leemingen1227/couple-server/api"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/mail"
	"github.com/leemingen1227/couple-server/util"
	"github.com/leemingen1227/couple-server/worker"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func main(){
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config:" + err.Error())
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:" + err.Error())
	}

	//run the migration
	//runDBMigration(config.MigrationURL, config.DBSource)

	//connect to redis for cache
	redisOptions := redis.Options{
		Addr: config.RedisAddress,
	}
	ctx := context.Background()
	redisClient := redis.NewClient(&redisOptions)
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Msg("Redis connection initialization failed:" + err.Error())
	}

	store := db.NewStore(conn)

	//connect to redis for asynq
	asynqRedisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(asynqRedisOpt)
	go runTaskProcessor(config, asynqRedisOpt, store)
	runGinServer(config, store, taskDistributor, redisClient)
}

func runDBMigration(migrationURL string, dbSource string) {
    migration, err := migrate.New(migrationURL, dbSource)
    if err != nil {
        log.Fatal().Msg("cannot create new migrate instance:")
    }

    if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal().Msg("failed to run migrate up:")
    }

    log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("starting task processor")
	err := processor.Start()
	if err != nil {
		log.Fatal().Msg("cannot start task processor")
	}
}

func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, redisClient *redis.Client) {
    server, err := api.NewServer(config, store, taskDistributor, redisClient)
    if err != nil {
        log.Fatal().Msg("cannot create server:")
    }
    
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Msg("cannot start server:")
    }
}