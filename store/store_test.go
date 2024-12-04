package store

import (
	"context"
	"testing"

	"github.com/mizmorr/songslib/pkg/logger"
)

func TestInit(t *testing.T) {
	log := logger.Get()
	ctx := context.WithValue(context.Background(), "logger", log)

	_, err := New(ctx)
	if err != nil {
		t.Error(err)
	}
}
