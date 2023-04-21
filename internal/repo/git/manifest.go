package git

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	reader "github.com/tupyy/tinyedge-controller/internal/repo/manifest"
	"go.uber.org/zap"
)

const (
	pattern = "\\.manifest\\.y[a]?ml$"
)

var (
	manifestPattern *regexp.Regexp
)

func init() {
	manifestPattern = regexp.MustCompile(pattern)
}

func getManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error) {
	files, err := findManifestFiles(ctx, repo.LocalPath, manifestPattern)
	if err != nil {
		return nil, fmt.Errorf("unable to search for manifest files in repo %q: %w", repo.LocalPath, err)
	}
	manifests := make([]entity.Manifest, 0, len(files))

	for _, file := range files {
		manifest, err := getManifest(ctx, repo, file)
		if err != nil {
			zap.S().Errorf("unable to parse manifest file %q in repo %q: %w", file, repo.LocalPath, err)
			continue
		}
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func getManifest(ctx context.Context, repo entity.Repository, filepath string) (entity.Manifest, error) {
	file := path.Join(repo.LocalPath, filepath)

	_, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("unable to find file %q in repo %q", filepath, repo.LocalPath)
	}

	return parseManifest(ctx, file, func(m entity.Manifest) entity.Manifest {
		if m.GetKind() == entity.WorkloadManifestKind {
			w, _ := m.(entity.Workload)
			w.Id = hash(filepath)[:12]
			w.Repository = repo
			return w
		}
		return m
	})
}

func parseManifest(ctx context.Context, filepath string, transformFn func(entity.Manifest) entity.Manifest) (entity.Manifest, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read manifest file %q: %w", filepath, err)
	}

	return reader.ReadManifest(bytes.NewBuffer(content), transformFn)
}

// findManifestFiles returns the list of all manifestworks files found in the repo
func findManifestFiles(ctx context.Context, path string, manifestPattern *regexp.Regexp) ([]string, error) {
	searchFn := func(ctx context.Context, root string, output chan string, errCh chan error, doneCh chan struct{}, pattern *regexp.Regexp) {
		err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}
			if !info.IsDir() && pattern.MatchString(info.Name()) {
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

	go searchFn(searchCtx, path, result, errCh, doneCh, manifestPattern)
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

func hash(base string, s ...string) string {
	hash := sha256.New()
	hash.Write([]byte(base))
	for _, ss := range s {
		hash.Write([]byte(ss))
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
