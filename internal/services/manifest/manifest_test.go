package manifest_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
)

var _ = Describe("Manifest", func() {
	It("retrieve manifest", func() {
		gitReader := &manifest.GitReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error) {
				return entity.ManifestWork{
					Id: ref.Id,
					Secrets: []entity.Secret{
						{
							Id:   "secret",
							Path: "path",
							Key:  "key",
						},
					},
				}, nil
			},
		}
		secretReader := &manifest.SecretReaderMock{
			GetSecretFunc: func(ctx context.Context, path, key string) (entity.Secret, error) {
				return entity.Secret{
					Id:    "secret",
					Path:  path,
					Key:   key,
					Value: "value",
				}, nil
			},
		}

		service := manifest.New(nil, gitReader, secretReader)
		manifest, err := service.GetManifest(context.TODO(), entity.ManifestReference{Id: "reference"})
		Expect(err).To(BeNil())

		// should call GetSecret
		secretCalls := secretReader.GetSecretCalls()
		Expect(len(secretCalls)).To(Equal(1))
		Expect(secretCalls[0].Key).To(Equal("key"))
		Expect(len(manifest.Secrets)).To(Equal(1))
		manifestSecret := manifest.Secrets[0]
		Expect(manifestSecret.Value).To(Equal("value"))
	})

	It("retrieve manifest returns error on missing secret", func() {
		gitReader := &manifest.GitReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error) {
				return entity.ManifestWork{
					Id: ref.Id,
					Secrets: []entity.Secret{
						{
							Id:   "secret",
							Path: "path",
							Key:  "key",
						},
					},
				}, nil
			},
		}
		secretReader := &manifest.SecretReaderMock{
			GetSecretFunc: func(ctx context.Context, path, key string) (entity.Secret, error) {
				return entity.Secret{}, errors.New("unable to read secret")
			},
		}

		service := manifest.New(nil, gitReader, secretReader)
		_, err := service.GetManifest(context.TODO(), entity.ManifestReference{Id: "reference"})
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("unable to read secret"))
	})

	It("retrieve manifest returns error on missing manifest", func() {
		gitReader := &manifest.GitReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error) {
				return entity.ManifestWork{}, errors.New("unable to get manifest")
			},
		}

		service := manifest.New(nil, gitReader, nil)
		_, err := service.GetManifest(context.TODO(), entity.ManifestReference{Id: "ref"})
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("unable to get manifest"))
		Expect(len(gitReader.GetManifestCalls())).To(Equal(1))
	})

	It("read all manifests from repo", func() {
		refReader := &manifest.ReferenceReaderMock{
			GetReferencesFunc: func(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error) {
				return []entity.ManifestReference{
					{
						Id: "ref1",
					},
					{
						Id: "ref2",
					},
				}, nil
			},
		}
		gitReader := &manifest.GitReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error) {
				return entity.ManifestWork{
					Id: ref.Id,
					Secrets: []entity.Secret{
						{
							Id:   ref.Id,
							Path: ref.Id,
							Key:  "key",
						},
					},
				}, nil
			},
		}
		secretReader := &manifest.SecretReaderMock{
			GetSecretFunc: func(ctx context.Context, path, key string) (entity.Secret, error) {
				return entity.Secret{
					Id:    "secret",
					Path:  path,
					Key:   key,
					Value: "value",
				}, nil
			},
		}

		service := manifest.New(refReader, gitReader, secretReader)
		manifest, err := service.GetManifests(context.TODO(), entity.Repository{Id: "some repo"})
		Expect(err).To(BeNil())
		Expect(len(manifest)).To(Equal(2))

		// should call GetSecret twice
		secretCalls := secretReader.GetSecretCalls()
		Expect(len(secretCalls)).To(Equal(2))
		Expect(secretCalls[0].Path).To(Equal("ref1"))

		Expect(len(manifest[0].Secrets)).To(Equal(1))
		manifestSecret := manifest[0].Secrets[0]
		Expect(manifestSecret.Value).To(Equal("value"))
	})
})
