package git

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"

	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
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

type manifestReader struct {
	repo   entity.Repository
	reader func(r io.Reader) entity.Workload
}

func (m *manifestReader) GetManifest(ctx context.Context, filepath string) (entity.Manifest, error) {
	file := path.Join(m.repo.LocalPath, filepath)
	_, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("unable to find file %q in repo %q", filepath, m.repo.LocalPath)
	}
	return m.parseManifest(ctx, filepath, m.repo.LocalPath)
}

func (m *manifestReader) GetManifests(ctx context.Context) ([]entity.Manifest, error) {
	files, err := m.findManifestFiles(ctx, m.repo.LocalPath, manifestPattern)
	if err != nil {
		return nil, fmt.Errorf("unable to search for manifest files in repo %q: %w", m.repo.LocalPath, err)
	}
	manifests := make([]entity.Manifest, 0, len(files))

	for _, file := range files {
		manifest, err := m.parseManifest(ctx, file, m.repo.LocalPath)
		if err != nil {
			zap.S().Errorf("unable to parse manifest file %q in repo %q: %w", file, m.repo.LocalPath, err)
			continue
		}
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func (m *manifestReader) parseManifest(ctx context.Context, filename, basePath string) (entity.Manifest, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read manifest file %q: %w", filename, err)
	}

	kind, err := m.getKind(content)
	if err != nil {
		return nil, fmt.Errorf("unable to read manifest %q kind: %w", filename, err)
	}

	switch kind {
	case "workload":
		manifest, err := parseWorkloadManifest(content, filename, basePath)
		if err != nil {
			return nil, fmt.Errorf("unable to parse workload manifest %q: %w", filename, err)
		}
		return manifest, nil
	default:
		return nil, fmt.Errorf("unsupported kind: configuration")
	}
}

// findManifestFiles returns the list of all manifestworks files found in the repo
func (m *manifestReader) findManifestFiles(ctx context.Context, path string, manifestPattern *regexp.Regexp) ([]string, error) {
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

func (w *manifestReader) getKind(content []byte) (string, error) {
	type anonymousStruct struct {
		Kind string `yaml:"kind"`
	}
	var a anonymousStruct
	if err := goyaml.Unmarshal(content, &a); err != nil {
		return "", fmt.Errorf("unknown struct: %s", err)
	}
	return a.Kind, nil
}
