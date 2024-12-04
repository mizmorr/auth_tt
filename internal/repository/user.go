package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mizmorr/auth/internal/domain"
	"github.com/mizmorr/auth/store/pg"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return uuid.Nil, err
	}

	return user.GUID, nil
}

func (r *UserRepository) UpdateActiveSession(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	err := r.db.Model(&user).Update("ActiveSessionID", user.ActiveSessionID).Error
	if err != nil {
		return uuid.Nil, err
	}

	return user.GUID, nil
}
