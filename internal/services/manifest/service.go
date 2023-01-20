package manifest

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type Service struct {
	refReader    ReferenceReader
	gitReader    GitReader
	secretReader SecretReader
}

func New(refReader ReferenceReader, gitReader GitReader, secretReader SecretReader) *Service {
	return &Service{
		refReader:    refReader,
		gitReader:    gitReader,
		secretReader: secretReader,
	}
}

// GetManifest returns the manifest from the git repository
func (w *Service) GetManifest(ctx context.Context, id string) (entity.ManifestWork, error) {
	ref, err := w.refReader.GetReference(ctx, id)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to get manifest reference: %w", err)
	}

	manifest, err := w.gitReader.GetManifest(ctx, ref)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to get manifest: %w", err)
	}

	// for each secret in the manifest get the value from vault
	for i := 0; i < len(manifest.Secrets); i++ {
		secret := &manifest.Secrets[i]
		s, err := w.secretReader.GetSecret(ctx, secret.Path, secret.Key)
		if err != nil {
			return entity.ManifestWork{}, fmt.Errorf("unable to read secret %q from vault: %w", secret.Path, err)
		}
		secret.Value = s.Value
		secret.Hash = s.Hash
	}

	return manifest, nil
}

// GetManifests return all the manifest from the whole git repository.
func (w *Service) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error) {
	refs, err := w.refReader.GetRepositoryReferences(ctx, repo)
	if err != nil {
		return []entity.ManifestWork{}, err
	}

	// for each ref get the real manifest and add devices, sets and namespaces
	manifests := make([]entity.ManifestWork, 0, len(refs))
	for _, ref := range refs {
		r := ref
		manifest, err := w.gitReader.GetManifest(ctx, ref)
		if err != nil {
			return []entity.ManifestWork{}, fmt.Errorf("unable to get manifest %q from repo %q: %w", ref.Path, repo.Id, err)
		}
		manifest.Reference = r
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}
