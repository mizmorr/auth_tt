package migrations

import (
	"context"

	"github.com/mizmorr/auth/internal/domain"
	"github.com/mizmorr/songslib/pkg/util"
	"gorm.io/gorm"
)

func AutoMigrate(ctx context.Context, db *gorm.DB) error {
	logger := util.GetLogger(ctx)

	logger.Debug().Msg("Running migrations..")

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&domain.Session{}); err != nil {
		return err
	}

	logger.Info().Msg("Migrations completed successfully.")
	return nil
}
