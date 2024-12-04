package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mizmorr/auth/config"
	"github.com/mizmorr/auth/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var (
	pgInstance *DB
	once       sync.Once
)

func Dial(ctx context.Context) (*DB, error) {
	conf := config.Get()
	log := logger.GetLoggerFromContext(ctx)
	if conf.PgURL == "" {
		return nil, errors.New("PG_URL is not set")
	}

	once.Do(func() {
		var (
			db  *gorm.DB
			err error
		)

		for conf.PgConnAttempts > 0 {
			db, err = gorm.Open(postgres.Open(conf.PgURL), &gorm.Config{})
			if err == nil {
				log.Info().Msg("Connected to PostgreSQL")
				break
			}

			conf.PgConnAttempts--
			log.Debug().Msg(fmt.Sprintf("Postgres is trying to connect, attempts left: %d", conf.PgConnAttempts))

			time.Sleep(conf.PgTimeout)
		}
		if err != nil {
			panic("Connection to PostgreSQL failed")
		}

		log.Debug().Msg("Creating DB if not exists..")

		err = createDBIfNotExists(db, conf.DBName)
		if err != nil {
			panic(err)
		}

		url := getUrlToDB(conf.PgURL, conf.DBName)
		dbCreated, err := gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			log.Info().Msg("Connection to db failed")
			panic(err)
		}

		pgInstance = &DB{dbCreated}
	})
	return pgInstance, nil
}

func createDBIfNotExists(db *gorm.DB, dbName string) error {
	var exists bool
	err := db.Raw(fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func getUrlToDB(urlToPostgres, dbname string) string {
	return strings.ReplaceAll(urlToPostgres, "localhost/?", fmt.Sprintf("localhost/%s?", dbname))
}
