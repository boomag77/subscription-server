package deps

import (
	"subscription-server/internal/logger"
	"subscription-server/internal/service"
	"subscription-server/internal/storage"
)

type Deps struct {
	Storage       storage.Storage
	Logger        logger.Logger
	AppleService  service.Service
	GoogleService service.Service
}
