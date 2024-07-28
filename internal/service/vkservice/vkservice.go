package vkservice

import (
	"go.uber.org/zap"

	"github.com/Vikot10/viarticles/internal/storage"
)

type VkService struct {
	logger      *zap.Logger
	accessToken string
}

func New(logger *zap.Logger, storage *storage.Storage, accessToken string) *VkService {
	return &VkService{
		logger:      logger,
		accessToken: accessToken,
	}
}
