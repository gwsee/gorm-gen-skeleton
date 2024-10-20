package config

import "gorm-gen-skeleton/internal/constants/container"

type DriverInterface interface {
	container.ContainerInterface
	Listen()
}
