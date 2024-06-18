package postgres

import (
	"fmt"

	"github.com/qaultsabit/wallet/src/infra/config"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	Conn *sqlx.DB
}

func New(conf config.SqlDbConf, logger *logrus.Logger) (PostgresDB, error) {
	db := PostgresDB{}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Host,
		conf.Username,
		conf.Password,
		conf.Name,
		conf.Port,
	)

	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect to the database!")
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	logger.Printf("Connected to database.")
	return db, nil
}
