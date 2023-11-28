package database

import (
	_ "github.com/lib/pq"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type (
	Database struct {
		URL          string
		MaxOpenConns int
		logger       *zerolog.Logger
	}
)

func New(url string, maxopenconns int, logger *zerolog.Logger) *Database {
	return &Database{url, maxopenconns, logger}
}

func (dbConfig *Database) Connect() *sqlx.DB {
	return dbConfig.connect("postgres")
}

func (dbConfig *Database) SQLServerConnect() *sqlx.DB {
	return dbConfig.connect("sqlserver")
}

func (dbConfig *Database) connect(driver string) *sqlx.DB {
	database, _ := sqlx.Open(driver, dbConfig.URL)
	if err := database.Ping(); err != nil {
		dbConfig.logger.Error().Err(err).Msg("failed to ping database")
		panic(err)
	}

	database.SetMaxOpenConns(dbConfig.MaxOpenConns)

	dbConfig.logger.Info().Msg("database connected with success")

	return database
}
