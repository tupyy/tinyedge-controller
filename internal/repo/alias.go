package repo

import (
	"github.com/tupyy/tinyedge-controller/internal/repo/cache"
	"github.com/tupyy/tinyedge-controller/internal/repo/git"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"github.com/tupyy/tinyedge-controller/internal/repo/vault/certificate"
	"github.com/tupyy/tinyedge-controller/internal/repo/vault/secret"
)

type (
	Device      postgres.DeviceRepo
	Manifest    postgres.ManifestRepository
	Repository  postgres.Repository
	Git         git.GitRepo
	Certificate certificate.CertficateRepo
	Secret      secret.Repository
	MemCache    cache.MemCacheRepo
)

var (
	NewDeviceRepo  = postgres.NewDeviceRepo
	NewManifest    = postgres.NewManifestRepository
	NewRepository  = postgres.NewRepository
	NewGit         = git.New
	NewCertificate = certificate.New
	NewSecret      = secret.New
	NewMemCache    = cache.NewCacheRepo
)
