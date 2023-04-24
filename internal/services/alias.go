package services

import (
	"github.com/tupyy/tinyedge-controller/internal/services/auth"
	"github.com/tupyy/tinyedge-controller/internal/services/certificate"
	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
	"github.com/tupyy/tinyedge-controller/internal/services/device"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"github.com/tupyy/tinyedge-controller/internal/services/repository"
)

type (
	Manifest      = manifest.Service
	Device        = device.Service
	Configuration = configuration.Service
	Repository    = repository.Service
	Edge          = edge.Service
	Auth          = auth.Service
	Certificate   = certificate.Service
)

var (
	NewManifest      = manifest.New
	NewDevice        = device.New
	NewConfiguration = configuration.New
	NewRepository    = repository.NewRepositoryService
	NewEdge          = edge.New
	NewAuth          = auth.New
	NewCertificate   = certificate.New
)
