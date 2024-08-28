package vkservice

import (
	"go.uber.org/zap"
)

type VkService struct {
	logger      *zap.Logger
	accessToken string
}

func New(logger *zap.Logger, accessToken string) *VkService {
	return &VkService{
		logger:      logger,
		accessToken: accessToken,
	}
}
