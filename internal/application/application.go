package application

import (
	"go.uber.org/zap"

	"github.com/Vikot10/viarticles/internal/dto"
	"github.com/Vikot10/viarticles/internal/storage"
)

type Application struct {
	logger *zap.Logger
	store  *storage.Storage
}

type VkProvider interface {
	GetFaves() ([]*dto.Fave, error)
	SynchronizeFaves() error
}
