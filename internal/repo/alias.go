package repo

import "github.com/tupyy/tinyedge-controller/internal/repo/postgres"

type DeviceRepo = postgres.DeviceRepo

var (
	NewDeviceRepo           = postgres.NewDeviceRepo
	NewDeviceRepoWithReader = postgres.NewDeviceRepoWithReader
)
