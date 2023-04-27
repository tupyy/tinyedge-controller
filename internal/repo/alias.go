package repo

import (
	"github.com/tupyy/tinyedge-controller/internal/repo/cache"
	"github.com/tupyy/tinyedge-controller/internal/repo/git"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"github.com/tupyy/tinyedge-controller/internal/repo/vault"
)

type (
	Device      postgres.DeviceRepo
	Manifest    postgres.ManifestRepository
	Repository  postgres.Repository
	Git         git.GitRepo
	Certificate vault.CertficateRepo
	Secret      vault.SecretRepository
	MemCache    cache.MemCacheRepo
)

var (
	NewDeviceRepo  = postgres.NewDeviceRepo
	NewManifest    = postgres.NewManifestRepository
	NewRepository  = postgres.NewRepository
	NewGit         = git.New
	NewCertificate = vault.NewCertificateRepository
	NewSecret      = vault.NewSecretRepository
	NewMemCache    = cache.NewCacheRepo
)
