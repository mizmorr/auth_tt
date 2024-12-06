package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mizmorr/auth_tt/internal/domain"
	"github.com/mizmorr/auth_tt/internal/repository"
	"github.com/mizmorr/auth_tt/pkg/logger"
	"github.com/mizmorr/auth_tt/store/migrations"
	"github.com/mizmorr/auth_tt/store/pg"
	"github.com/pkg/errors"
)

type SessionRepo interface {
	Create(ctx context.Context, session *domain.Session) (string, error)
	Delete(ctx context.Context, session *domain.Session) error
}

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) (uuid.UUID, error)
	UpdateActiveSession(ctx context.Context, user *domain.User) (uuid.UUID, error)
}

var (
	_ UserRepo    = (*repository.UserRepository)(nil)
	_ SessionRepo = (*repository.SessionRepository)(nil)
)

type Store struct {
	Pg                *pg.DB
	SessionRepository SessionRepo
	UserRepository    UserRepo
}

var store Store

func New(ctx context.Context) (*Store, error) {
	logger := logger.GetLoggerFromContext(ctx)

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
		store.UserRepository = repository.NewUserRepository(pg)
		store.SessionRepository = repository.NewSessionRepository(pg)

	}
	logger.Info().Msg("PostgreSQL store initialized successfully")
	return &store, nil
}

const KeepALiveTimeout = 5

func (store *Store) keepAlive(ctx context.Context) {
	logger := logger.GetLoggerFromContext(ctx)
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
