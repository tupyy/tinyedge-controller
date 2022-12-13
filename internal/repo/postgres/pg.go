package device

import (
	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"gorm.io/gorm"
)

type DeviceRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}
