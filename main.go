package main

import (
	"database/sql"
	"os"

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

	store := db.NewStore(conn)

	//connect to redis
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	runGinServer(config, store, taskDistributor)
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

func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := api.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Msg("cannot create server:")
    }
    
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Msg("cannot start server:")
    }
}