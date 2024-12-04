package store

import (
	"context"
	"time"

	"github.com/mizmorr/auth/store/migrations"
	"github.com/mizmorr/auth/store/pg"
	"github.com/mizmorr/songslib/pkg/util"
	"github.com/pkg/errors"
)

type Store struct {
	Pg         *pg.DB
	Repository any
}

var store Store

func New(ctx context.Context) (*Store, error) {
	logger := util.GetLogger(ctx)

	logger.Debug().Msg("Initializing PostgreSQL store")
	pg, err := pg.Dial(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "pg.Store: failed to connect to the database")
	}

	logger.Debug().Msg("Running PostgreSQL migrations")
	if err := migrations.AutoMigrate(ctx, pg.DB); err != nil {
		return nil, errors.Wrap(err, "pg.Store: failed to run migrations")
	}

	if pg != nil {
		store.Pg = pg
		go store.keepAlive(ctx)
		var repo interface{}
		store.Repository = repo
	}
	logger.Info().Msg("PostgreSQL store initialized successfully")
	return &store, nil
}

const KeepALiveTimeout = 5

func (store *Store) keepAlive(ctx context.Context) {
	logger := util.GetLogger(ctx)
	for {
		time.Sleep(time.Second * KeepALiveTimeout)
		var (
			lost_connection bool
			err             error
		)

		if store.Pg == nil {
			lost_connection = true
		}
		if lost_connection {
			logger.Debug().Msg("[store.keepAlive] Lost connection, is trying to reconnect...")
			store.Pg, err = pg.Dial(ctx)
			if err != nil {
				logger.Err(err)
			} else {
				logger.Debug().Msg("[store.keepAlive] Connection established")
			}
		}

	}
}
