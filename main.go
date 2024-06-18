package main

import (
	"context"
	"database/sql"

	_ "github.com/joho/godotenv/autoload"
	"github.com/qaultsabit/wallet/src/app/usecases"
	userUC "github.com/qaultsabit/wallet/src/app/usecases/user"
	"github.com/qaultsabit/wallet/src/infra/config"
	ms_log "github.com/qaultsabit/wallet/src/infra/log"
	"github.com/qaultsabit/wallet/src/infra/persistence/postgres"
	userRepo "github.com/qaultsabit/wallet/src/infra/persistence/postgres/user"
	"github.com/qaultsabit/wallet/src/interface/rest"
	"github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	conf := config.Make()

	isProd := false
	if conf.App.Environment == "PRODUCTION" {
		isProd = true
	}

	m := make(map[string]interface{})
	m["env"] = conf.App.Environment
	m["service"] = conf.App.Name
	logger := ms_log.NewLogInstance(
		ms_log.LogName(conf.Log.Name),
		ms_log.IsProduction(isProd),
		ms_log.LogAdditionalFields(m),
	)

	db, err := postgres.New(conf.SqlDb, logger)
	if err != nil {
		panic(err)
	}
	defer func(l *logrus.Logger, sqlDB *sql.DB, dbName string) {
		err := sqlDB.Close()
		if err != nil {
			l.Errorf("error closing sql database %s: %s", dbName, err)
		} else {
			l.Printf("sql database %s successfuly closed.", dbName)
		}
	}(logger, db.Conn.DB, db.Conn.DriverName())

	userRepository := userRepo.NewUserRepository(db.Conn)

	httpServer, err := rest.New(
		conf.Http,
		isProd,
		logger,
		usecases.AllUseCases{
			UserUC: userUC.NewUserUseCase(userRepository),
		},
	)
	if err != nil {
		panic(err)
	}

	httpServer.Start(ctx)
}
