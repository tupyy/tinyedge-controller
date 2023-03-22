package git

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

type GitRepo struct {
	// localStorage is the path to a temporary folder where all the clones are kept
	localStorage string
}

func New(localStorage string) *GitRepo {
	return &GitRepo{localStorage: localStorage}
}

// Open opens the git repo. If the repo does not exists in the local storage it will be cloned from remote.
// Returns a new entity with updated information if the repo was cloned.
func (g *GitRepo) Open(ctx context.Context, r entity.Repository) (entity.Repository, error) {
	_, err := g.openRepository(ctx, r)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return entity.Repository{}, errService.NewResourceNotFoundError("repository", r.Id)
		}
		return entity.Repository{}, err
	}

	zap.S().Debugw("successfully open repo", "local", r.LocalPath)

	return r, nil
}

// Pull pull from origin the repo.
func (g *GitRepo) Pull(ctx context.Context, r entity.Repository) error {
	repo, err := g.openRepository(ctx, r)
	if err != nil {
		return fmt.Errorf("unable to pull from %q: %w", r.Url, err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	pullOptions := &git.PullOptions{
		RemoteName:      "origin",
		RemoteURL:       r.Url,
		SingleBranch:    true,
		InsecureSkipTLS: true,
	}
	if r.AuthType != entity.NoRepositoryAuthType {
		authMethod, err := g.getCredentials(ctx, r.Credentials, r.CredentialsSecretPath)
		if err != nil {
			return err
		}
		pullOptions.Auth = authMethod
	}

	err = w.Pull(pullOptions)
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		return fmt.Errorf("unable to pull from origin of repo %q: %w", r.Url, err)
	}
	return nil
}

// GetHeadSha returns the head sha for the specified repo.
// It does not pull before returning the sha.
func (g *GitRepo) GetHeadSha(ctx context.Context, r entity.Repository) (string, error) {
	repo, err := g.openRepository(ctx, r)
	if err != nil {
		return "", fmt.Errorf("unable to open repository %q: %w", r.Url, err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return "", err
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(r.Branch),
	})
	if err != nil {
		return "", fmt.Errorf("unable to checkout branch %q from repo %q", r.Branch, r.Url)
	}
	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("unable to read head from repo %q: %w", r.Url, err)
	}
	return ref.Hash().String(), err
}

// GetManifest return the manifest referred by ref
func (g *GitRepo) GetManifest(ctx context.Context, repo entity.Repository, filepath string) (entity.Manifest, error) {
	reader := &manifestReader{repo: repo}
	return reader.GetManifest(ctx, filepath)
}

// GetManifests returns all the manifest of a repo.
func (g *GitRepo) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error) {
	reader := &manifestReader{repo: repo}
	return reader.GetManifests(ctx)
}

func (g *GitRepo) Clone(ctx context.Context, repo entity.Repository) (entity.Repository, error) {
	zap.S().Infof("clone repo %q to local storage %q", repo.Url, g.localStorage)

	cloneOptions := &git.CloneOptions{
		URL:               repo.Url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	if repo.AuthType != entity.NoRepositoryAuthType {
		authMethod, err := g.getCredentials(ctx, repo.Credentials, repo.CredentialsSecretPath)
		if err != nil {
			return repo, err
		}
		cloneOptions.Auth = authMethod
	}

	clone, err := git.PlainClone(path.Join(g.localStorage, repo.Id), false, cloneOptions)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return entity.Repository{}, err
		}
		r, err := git.PlainOpen(path.Join(g.localStorage, repo.Id))
		if err != nil {
			return entity.Repository{}, err
		}
		clone = r
	}

	mainBranch, err := g.getMainBranch(clone)
	if err != nil {
		return entity.Repository{}, err
	}
	repo.LocalPath = path.Join(g.localStorage, repo.Id)
	repo.Branch = mainBranch

	headSha, err := g.GetHeadSha(ctx, repo)
	if err != nil {
		return repo, err
	}
	repo.TargetHeadSha = headSha
	zap.S().Debugw("successfully cloned repo", "remote", repo.AuthType, "local", repo.LocalPath)
	return repo, nil
}

// openRepository opens a repo from local storage.
func (g *GitRepo) openRepository(ctx context.Context, r entity.Repository) (*git.Repository, error) {
	repo, err := git.PlainOpen(r.LocalPath)
	if err != nil {
		zap.S().Infof("unable to open repo %q", r.LocalPath)
		return nil, err
	}
	return repo, nil
}

func (g *GitRepo) getMainBranch(r *git.Repository) (string, error) {
	// check if main branch is "main" or "master"
	w, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("unable to open repository %w", err)
	}
	// try main
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
	})
	if err == nil {
		return "main", nil
	}
	err = w.Checkout(&git.CheckoutOptions{})
	if err == nil {
		return "master", nil
	}
	return "", fmt.Errorf("no branch named \"main\" or \"master\" in repo")
}

func (g *GitRepo) getCredentials(ctx context.Context, fn entity.CredentialsFunc, secretPath string) (transport.AuthMethod, error) {
	credetials, err := fn(ctx, secretPath)
	if err != nil {
		return nil, err
	}

	switch v := credetials.(type) {
	case entity.SSHRepositoryAuth:
		publicKeys, err := ssh.NewPublicKeys("git", v.PrivateKey, v.Password)
		if err != nil {
			return nil, fmt.Errorf("unable to create auth method: %w", err)
		}
		return publicKeys, nil
	case entity.TokenRepositoryAuth:
		return &http.TokenAuth{Token: v.Token}, nil
	case entity.BasicRepositoryAuth:
		return &http.BasicAuth{Username: v.Username, Password: v.Password}, nil
	}
	return nil, errors.New("unknown auth method")
}
