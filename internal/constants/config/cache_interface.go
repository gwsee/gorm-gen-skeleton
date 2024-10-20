package config

import (
	"gorm-gen-skeleton/internal/constants/container"
)

// CacheInterface config cache
type CacheInterface interface {
	container.ContainerInterface
	FuzzyDelete(key string)
}
