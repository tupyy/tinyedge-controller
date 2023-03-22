package configuration_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
)

var _ = Describe("ConfigurationResponse", func() {
	It("serve configuration for device successful", func() {
		deviceReader := &configuration.DeviceReaderMock{
			GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
				return entity.Device{
					ID: "toto",
					Configuration: &entity.Configuration{
						ID: "conf_for_toto",
					},
					ManifestIDS: []string{"manifest_for_toto"},
				}, nil
			},
		}
		manifestReader := &configuration.ManifestReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.Reference) (entity.WorkloadManifest, error) {
				return entity.WorkloadManifest{
					Id: "manifest_for_toto",
				}, nil
			},
		}
		referenceReader := &configuration.ReferenceReaderMock{
			GetReferenceFunc: func(ctx context.Context, id string) (entity.Reference, error) {
				return entity.Reference{}, nil
			},
		}

		confService := configuration.New(deviceReader, manifestReader, referenceReader, nil)
		confResponse, err := confService.GetDeviceConfiguration(context.TODO(), "toto")
		Expect(err).To(BeNil())
		Expect(confResponse.Configuration.ID).To(Equal("conf_for_toto"))
	})

	It("serve configuration for device successful #2", func() {
		deviceReader := &configuration.DeviceReaderMock{
			GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
				setID := "set"
				return entity.Device{ID: "toto", SetID: &setID}, nil
			},
			GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
				return entity.Set{
					Configuration: &entity.Configuration{
						ID: "conf_for_toto",
					},
					ManifestIDS: []string{"manifest_for_toto"},
				}, nil
			},
		}
		manifestReader := &configuration.ManifestReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.Reference) (entity.WorkloadManifest, error) {
				return entity.WorkloadManifest{
					Id: "manifest_for_toto",
				}, nil
			},
		}
		referenceReader := &configuration.ReferenceReaderMock{
			GetReferenceFunc: func(ctx context.Context, id string) (entity.Reference, error) {
				return entity.Reference{}, nil
			},
		}

		confService := configuration.New(deviceReader, manifestReader, referenceReader, nil)
		confResponse, err := confService.GetDeviceConfiguration(context.TODO(), "toto")
		Expect(err).To(BeNil())
		Expect(confResponse.Configuration.ID).To(Equal("conf_for_toto"))
	})

	It("serve configuration for device successful #3", func() {
		deviceReader := &configuration.DeviceReaderMock{
			GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
				return entity.Device{ID: "toto", NamespaceID: "default"}, nil
			},
			GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
				return entity.Namespace{
					Configuration: entity.Configuration{
						ID: "conf_for_toto",
					},
					ManifestIDS: []string{"manifest_for_toto"},
				}, nil
			},
		}
		manifestReader := &configuration.ManifestReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.Reference) (entity.WorkloadManifest, error) {
				return entity.WorkloadManifest{
					Id: "manifest_for_toto",
				}, nil
			},
		}
		referenceReader := &configuration.ReferenceReaderMock{
			GetReferenceFunc: func(ctx context.Context, id string) (entity.Reference, error) {
				return entity.Reference{}, nil
			},
		}

		confService := configuration.New(deviceReader, manifestReader, referenceReader, nil)
		confResponse, err := confService.GetDeviceConfiguration(context.TODO(), "toto")
		Expect(err).To(BeNil())
		Expect(confResponse.Configuration.ID).To(Equal("conf_for_toto"))
	})

	It("serve configuration for device successful #4", func() {
		deviceReader := &configuration.DeviceReaderMock{
			GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
				return entity.Device{
					ID:          "toto",
					NamespaceID: "default",
					ManifestIDS: []string{"manifest_for_toto"},
				}, nil
			},
			GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
				return entity.Namespace{
					Configuration: entity.Configuration{
						ID: "conf_for_toto",
					},
					ManifestIDS: []string{"manifest_from namespace"},
				}, nil
			},
		}
		manifestReader := &configuration.ManifestReaderMock{
			GetManifestFunc: func(ctx context.Context, ref entity.Reference) (entity.WorkloadManifest, error) {
				return entity.WorkloadManifest{
					Id: ref.Id,
				}, nil
			},
		}
		referenceReader := &configuration.ReferenceReaderMock{
			GetReferenceFunc: func(ctx context.Context, id string) (entity.Reference, error) {
				return entity.Reference{}, nil
			},
		}

		confService := configuration.New(deviceReader, manifestReader, referenceReader, nil)
		confResponse, err := confService.GetDeviceConfiguration(context.TODO(), "toto")
		Expect(err).To(BeNil())
		Expect(confResponse.Configuration.ID).To(Equal("conf_for_toto"))
	})

})
