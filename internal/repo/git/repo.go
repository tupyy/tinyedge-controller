package git

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	manifestv1 "github.com/tupyy/tinyedge-controller/internal/repo/models/manifest/v1"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	k8syaml "sigs.k8s.io/yaml"
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
	if r.LocalPath == "" {
		return g.clone(ctx, r)
	}

	_, err := g.openRepository(ctx, r)
	if err != nil {
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
	err = w.Pull(&git.PullOptions{
		RemoteName:      "origin",
		SingleBranch:    true,
		InsecureSkipTLS: true,
	})
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

// GetManifestWorks returns all the manifest work found in the repo.
// Only the valid manifest works are returned
func (g *GitRepo) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error) {
	manifestFiles, err := g.findManifestFiles(ctx, repo.LocalPath)
	if err != nil {
		return nil, err
	}

	manifests := make([]entity.ManifestWork, 0, len(manifestFiles))
	for _, m := range manifestFiles {
		manifest, err := g.getManifest(ctx, m, repo.LocalPath)
		if err != nil {
			zap.S().Errorf("unable to get manifest file", "error", err, "filepath", m, "repo_id", repo.Id)
			continue
		}
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func (g *GitRepo) GetManifest(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error) {
	return g.getManifest(ctx, ref.Path, ref.Repo.LocalPath)
}

func (g *GitRepo) getManifest(ctx context.Context, filename, basePath string) (entity.ManifestWork, error) {
	computeHash := func(b []byte) string {
		hash := sha256.New()
		hash.Write(b)
		return fmt.Sprintf("%x", hash.Sum(nil))
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to read manifest file %q: %w", filename, err)
	}
	manifest, err := g.parse(content, filename, basePath)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to parse manifest work %q: %w", filename, err)
	}

	manifest.Hash = computeHash(content)
	manifest.Id = computeHash(bytes.NewBufferString(filename).Bytes())
	manifest.Path = filename

	return manifest, nil

}

func (g *GitRepo) clone(ctx context.Context, r entity.Repository) (entity.Repository, error) {
	zap.S().Infof("clone repo %q to local storage %q", r.Url, g.localStorage)
	clone, err := git.PlainClone(path.Join(g.localStorage, r.Id), false, &git.CloneOptions{
		URL:               r.Url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	mainBranch, err := g.getMainBranch(clone)
	if err != nil {
		return entity.Repository{}, err
	}
	newRepo := entity.Repository{
		Id:        r.Id,
		Url:       r.Url,
		Branch:    mainBranch,
		LocalPath: path.Join(g.localStorage, r.Id),
	}
	headSha, err := g.GetHeadSha(ctx, newRepo)
	if err != nil {
		return newRepo, err
	}
	newRepo.TargetHeadSha = headSha
	zap.S().Debugw("successfully cloned repo", "remote", r.Url, "local", newRepo.LocalPath)
	return newRepo, nil
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

// parse parses the manifest file and verify that all the resources defined are valid k8s manifestv1.
// Returns false if one resource is not a valid ConfigMap or Pod.
func (g *GitRepo) parse(content []byte, filename, basePath string) (entity.ManifestWork, error) {
	var manifest manifestv1.Manifest

	if err := goyaml.Unmarshal(content, &manifest); err != nil {
		return entity.ManifestWork{}, err
	}

	e := entity.ManifestWork{
		Version:     manifest.Version,
		Description: manifest.Description,
		Name:        manifest.Name,
		Selectors:   make([]entity.Selector, 0),
		Secrets:     make([]entity.Secret, 0, len(manifest.Secrets)),
		Valid:       true,
	}

	for i := 0; true; i++ {
		keepGoing := false
		if i < len(manifest.Selector.Namespaces) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.NamespaceSelector,
				Value: manifest.Selector.Namespaces[i],
			})
			keepGoing = true
		}
		if i < len(manifest.Selector.Sets) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.SetSelector,
				Value: manifest.Selector.Sets[i],
			})
			keepGoing = true
		}
		if i < len(manifest.Selector.Devices) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.DeviceSelector,
				Value: manifest.Selector.Devices[i],
			})
			keepGoing = true
		}
		if !keepGoing {
			break
		}
	}

	for _, s := range manifest.Secrets {
		e.Secrets = append(e.Secrets, entity.Secret{
			Path: s.Path,
			Id:   s.Name,
			Key:  s.Key,
		})
	}

	for _, r := range manifest.Resources {
		configmaps, pods, err := g.parseResource(r.Ref, basePath)
		if err != nil {
			e.Valid = false
			return e, err
		}
		e.ConfigMaps = append(e.ConfigMaps, configmaps...)
		e.Pods = append(e.Pods, pods...)
	}

	return e, nil
}

// findManifestFiles returns the list of all manifestworks files found in the repo
func (g *GitRepo) findManifestFiles(ctx context.Context, path string) ([]string, error) {
	searchFn := func(ctx context.Context, root string, output chan string, errCh chan error, doneCh chan struct{}, filename string) {
		err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}
			if !info.IsDir() && info.Name() == filename {
				output <- path
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			return nil
		})
		if err != nil {
			errCh <- err
		}
		doneCh <- struct{}{}
	}

	result := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	searchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go searchFn(searchCtx, path, result, errCh, doneCh, entity.ManifestWorkFilename)
	manifestWorks := make([]string, 0)

	keep := true
	for keep {
		select {
		case manifestWorkFile := <-result:
			manifestWorks = append(manifestWorks, manifestWorkFile)
		case err := <-errCh:
			zap.S().Errorf("error during manifest work file search in %q: %q", path, err)
		case <-doneCh:
			keep = false
		}
	}

	return manifestWorks, nil
}

func (g *GitRepo) parseResource(filename string, basePath string) ([]v1.ConfigMap, []v1.Pod, error) {
	filepath := path.Join(basePath, filename)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to read resource file %q: %w", filepath, err)
	}

	parts, err := g.splitYAML(content)
	if err != nil {
		return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to decode resource file %q: %w", filepath, err)
	}

	configMaps := make([]v1.ConfigMap, 0, len(parts))
	pods := make([]v1.Pod, 0, len(parts))

	allowedKinds := "ConfigMap|Pods"
	for _, part := range parts {
		kind, err := g.getKind(part)
		if err != nil {
			return configMaps, pods, fmt.Errorf("unable to get \"kind\" from yaml with error %q", err)
		}
		if kind == "" || !strings.Contains(allowedKinds, kind) {
			zap.S().Warnf("kind %q not allowed in manifest work and it will be ignored", kind)
			continue
		}
		switch kind {
		case "ConfigMap":
			var c v1.ConfigMap
			err := k8syaml.Unmarshal(part, &c)
			if err != nil {
				return configMaps, pods, fmt.Errorf("unable to unmarshal part %q: %v", string(part), err)
			}
			configMaps = append(configMaps, c)
		case "Pod":
			var p v1.Pod
			err := k8syaml.Unmarshal(part, &p)
			if err != nil {
				return configMaps, pods, fmt.Errorf("unable to unmarshal part %q: %v", string(part), err)
			}
			pods = append(pods, p)
		}
	}

	return configMaps, pods, nil
}

func (g *GitRepo) splitYAML(resources []byte) ([][]byte, error) {
	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func (g *GitRepo) getKind(content []byte) (string, error) {
	type anonymousStruct struct {
		Kind string `yaml:"kind"`
	}
	var a anonymousStruct
	if err := goyaml.Unmarshal(content, &a); err != nil {
		return "", fmt.Errorf("unknown struct: %s", err)
	}
	return a.Kind, nil
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
