package testhelper

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"net/url"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/database"
	"github.com/caarlos0/env/v6"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type TestDatabase struct {
	db           *sqlx.DB
	DefaultURL   string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"` //nolint:lll // default connection string
	newURL       string
	databaseName string
	logger       *zerolog.Logger
}

type TestMigrationConfig struct {
	MigrationDIR   string
	SeederDIR      string
	SkipMigrations bool
}

func NewTestDatabase(logger *zerolog.Logger) *TestDatabase {
	databaseName := fmt.Sprintf("test-%s", uuid.NewV4().String())
	testdb := &TestDatabase{
		db:           nil,
		databaseName: databaseName,
		logger:       logger,
	}

	_ = env.Parse(testdb)

	testdb.buildURLConnection()

	return testdb
}

func (testdb *TestDatabase) Create(testMigrationConfig ...TestMigrationConfig) (string, error) {
	var migrationConfig TestMigrationConfig

	if len(testMigrationConfig) > 0 {
		migrationConfig = testMigrationConfig[0]
	}

	if err := testdb.createNewDatabase(); err != nil {
		return "", err
	}

	migrationDB := database.New(testdb.newURL, 1, testdb.logger).Connect()

	defer func() {
		migrationDB.Close()
		testdb.db.Close()
	}()

	if migrationConfig.SkipMigrations {
		return testdb.newURL, nil
	}

	migrationDir := "scripts/db/migrations"
	seederDir := "scripts/db/migrations/test_seeders"

	if migrationConfig.MigrationDIR != "" {
		migrationDir = migrationConfig.MigrationDIR
	}

	if migrationConfig.SeederDIR != "" {
		seederDir = migrationConfig.SeederDIR
	}

	err := testdb.execMigrations(migrationDB.DB, "migrations", migrationDir)

	if err != nil {
		return "", err
	}

	err = testdb.execMigrations(migrationDB.DB, "seeds_migrations", seederDir)

	if err != nil {
		return "", err
	}

	return testdb.newURL, nil
}

func (testdb *TestDatabase) Drop() (err error) {
	testdb.db = database.New(testdb.DefaultURL, 5, testdb.logger).Connect()
	defer func() {
		if err := testdb.db.Close(); err != nil {
			return
		}
	}()

	_, err = testdb.db.Exec(fmt.Sprintf("DROP DATABASE \"%s\"", testdb.databaseName))

	return
}

func (testdb *TestDatabase) createNewDatabase() (err error) {
	testdb.db = database.New(testdb.DefaultURL, 5, testdb.logger).Connect()

	_, err = testdb.db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", testdb.databaseName))

	return
}

func (testdb *TestDatabase) buildURLConnection() {
	url, err := url.Parse(testdb.DefaultURL)
	if err != nil {
		testdb.logger.Error().Err(err).Msg("failed to parse url connection")
		panic(err)
	}
	url.Path = testdb.databaseName
	testdb.newURL = url.String()
}

func (testdb *TestDatabase) execMigrations(database *sql.DB, table string, directory string) error {
	rootDirectory, err := getRootDirectory()

	if err != nil {
		return err
	}

	seeders := &migrate.FileMigrationSource{
		Dir: rootDirectory + directory,
	}

	migrate.SetTable(table)
	migrations, err := migrate.Exec(database, "postgres", seeders, migrate.Up)

	if err != nil {
		return err
	}

	testdb.logger.Info().
		Str("directory", directory).
		Int("migrations", migrations).
		Msg("Migrations applied")

	return nil
}

func getRootDirectory() (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dirs := strings.Split(dir, "/")
	rootDir := ""

	for _, dir := range dirs {
		rootDir += dir + "/"

		info, _ := os.Stat(rootDir + "go.mod")

		if info != nil {
			if info.Name() == "go.mod" {
				break
			}
		}
	}

	return rootDir, nil
}
