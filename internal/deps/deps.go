package deps

import (
	"subscription-server/internal/contracts"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
)

type Deps struct {
	Storage       storage.Storage
	Logger        logger.Logger
	AppleService  contracts.Service
	GoogleService contracts.Service
}
