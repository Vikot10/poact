package vkservice

import (
	"go.uber.org/zap"

	"github.com/Vikot10/viarticles/internal/storage"
)

type VKService struct {
	logger      *zap.Logger
	storage     *storage.Storage
	accessToken string
}
