package repository

import (
	"context"

	"github.com/mizmorr/auth_tt/internal/domain"
	"github.com/mizmorr/auth_tt/store/pg"
	"github.com/pkg/errors"
)

type SessionRepository struct {
	db *pg.DB
}

func NewSessionRepository(db *pg.DB) *SessionRepository { return &SessionRepository{db: db} }

func (s *SessionRepository) Create(ctx context.Context, session *domain.Session) (string, error) {
	err := s.db.Create(session).Error
	if err != nil {
		return "", errors.Wrap(err, "failed to create session")
	}

	return session.ID, nil
}

func (s *SessionRepository) Delete(ctx context.Context, session *domain.Session) error {
	return s.db.Delete(session).Error
}
