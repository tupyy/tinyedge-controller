package vault

import (
	"context"
	"fmt"
	"log"
	"sync"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"go.uber.org/zap"
)

// renewResult is a bitmask which could contain one or more of the values below
type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiringAuthToken // will be revoked soon
)

type VaultParameters struct {
	// connection parameters
	Address          string
	ApproleRoleID    string
	ApproleSecretID  string
	SecretsMountPath string
	PKIMountPath     string
	PKIRoleID        string
}

type Vault struct {
	Client     *vault.Client
	parameters VaultParameters
	wg         sync.WaitGroup
}

// NewVaultAppRoleClient logs in to Vault using the AppRole authentication
// method, returning an authenticated client and the auth token itself, which
// can be periodically renewed.
func NewVaultAppRoleClient(ctx context.Context, parameters VaultParameters) (*Vault, error) {
	zap.S().Infof("connecting to vault @ %q", parameters.Address)

	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = parameters.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	parameters.PKIMountPath = "pki_int"
	parameters.PKIRoleID = "tinyedge-role"
	vault := &Vault{
		Client:     client,
		parameters: parameters,
	}

	token, err := vault.login(ctx)
	if err != nil {
		return nil, fmt.Errorf("vault login error: %w", err)
	}

	vault.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Errorw("unable to renew lease", "error", err)
			}
		}()
		vault.periodicallyRenewLeases(ctx, token)
		vault.wg.Done()
	}()

	zap.S().Infof("connected to vault @ %q", parameters.Address)

	return vault, nil
}

// A combination of a RoleID and a SecretID is required to log into Vault
// with AppRole authentication method. The SecretID is a value that needs
// to be protected, so instead of the app having knowledge of the SecretID
// directly, we have a trusted orchestrator (simulated with a script here)
// give the app access to a short-lived response-wrapping token.
//
// ref: https://www.vaultproject.io/docs/concepts/response-wrapping
// ref: https://learn.hashicorp.com/tutorials/vault/secure-introduction?in=vault/app-integration#trusted-orchestrator
// ref: https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices
func (v *Vault) login(ctx context.Context) (*vault.Secret, error) {
	zap.S().Debugf("logging in to vault with approle auth; role id: %s", v.parameters.ApproleRoleID)

	approleSecretID := &approle.SecretID{
		FromString: v.parameters.ApproleSecretID,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		v.parameters.ApproleRoleID,
		approleSecretID,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize approle authentication method: %w", err)
	}

	authInfo, err := v.Client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using approle auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no approle info was returned after login")
	}

	zap.S().Debug("logging in to vault with approle auth: success!")

	return authInfo, nil
}

// Once you've set the token for your Vault client, you will need to
// periodically renew it. Likewise, the database credentials lease will expire
// at some point and also needs to be renewed periodically.
//
// A function like this one should be run as a goroutine to avoid blocking.
// Production applications may also need to be more tolerant of failures and
// retry on errors rather than exiting.
//
// Additionally, enterprise Vault users should be aware that due to eventual
// consistency, the API may return unexpected errors when running Vault with
// performance standbys or performance replication, despite the client having
// a freshly renewed token. See the link below for several ways to mitigate
// this which are outside the scope of this code sample.
//
// ref: https://www.vaultproject.io/docs/enterprise/consistency#vault-1-7-mitigations
func (v *Vault) periodicallyRenewLeases(ctx context.Context, authToken *vault.Secret) {
	zap.S().Debug("renew / recreate secrets loop: begin")
	defer zap.S().Debug("renew / recreate secrets loop: end")

	currentAuthToken := authToken

	for {
		renewed, err := v.renewLeases(ctx, currentAuthToken)
		if err != nil {
			log.Fatalf("renew error: %v", err) // simplified error handling
		}

		if renewed&exitRequested != 0 {
			return
		}

		if renewed&expiringAuthToken != 0 {
			zap.S().Info("auth token: can no longer be renewed; will log in again")

			authToken, err := v.login(ctx)
			if err != nil {
				panic(err)
			}

			currentAuthToken = authToken
		}
	}
}

// renewLeases is a blocking helper function that uses LifetimeWatcher
// instances to periodically renew the given secrets when they are close to
// their 'token_ttl' expiration times until one of the secrets is close to its
// 'token_max_ttl' lease expiration time.
func (v *Vault) renewLeases(ctx context.Context, authToken *vault.Secret) (renewResult, error) {
	zap.S().Debug("renew cycle: begin")
	defer zap.S().Debug("renew cycle: end")

	// auth token
	authTokenWatcher, err := v.Client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: authToken,
	})
	if err != nil {
		return renewError, fmt.Errorf("unable to initialize auth token lifetime watcher: %w", err)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	// monitor events from both watchers
	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

		// DoneCh will return if renewal fails, or if the remaining lease
		// duration is under a built-in threshold and either renewing is not
		// extending it or renewing is disabled.  In both cases, the caller
		// should attempt a re-read of the secret. Clients should check the
		// return value of the channel to see if renewal was successful.
		case err := <-authTokenWatcher.DoneCh():
			// Leases created by a token get revoked when the token is revoked.
			return expiringAuthToken, err

		// RenewCh is a channel that receives a message when a successful
		// renewal takes place and includes metadata about the renewal.
		case info := <-authTokenWatcher.RenewCh():
			zap.S().Debugw("auth token: successfully renewed", "remaining duration", info.Secret.Auth.LeaseDuration)
		}
	}
}
