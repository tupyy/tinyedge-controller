package entity

import (
	"context"
	"time"
)

type RepositoryAuthType int

const (
	SSHRepositoryAuthType RepositoryAuthType = iota
	BasicRepositoryAuthType
	TokenRepositoryAuthType
	NoRepositoryAuthType
)

type CredentialsFunc func(ctx context.Context, path string) (interface{}, error)

// Repository holds the information about the git repository where the ManifestWork are to be found.
type Repository struct {
	Id                    string
	AuthType              RepositoryAuthType
	Credentials           CredentialsFunc
	CredentialsSecretPath string
	Url                   string
	Branch                string
	LocalPath             string
	CurrentHeadSha        string
	TargetHeadSha         string
	PullPeriod            time.Duration
	Manifests             []string
}

type SSHRepositoryAuth struct {
	PrivateKey []byte
	Password   string
}

type TokenRepositoryAuth struct {
	Token string
}

type BasicRepositoryAuth struct {
	Username string
	Password string
}
