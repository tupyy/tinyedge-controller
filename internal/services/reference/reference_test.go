package reference_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/reference"
)

var _ = Describe("References", func() {
	var (
		gitReader             *reference.GitReaderMock
		deviceReaderWriter    *reference.DeviceReaderMock
		referenceReaderWriter *reference.ReferenceReaderWriterMock
		service               *reference.Service
	)

	Describe("Update", func() {
		It("updates the reference with success", func() {
			gitReader = &reference.GitReaderMock{
				GetReferencesFunc: func(ctx context.Context, repo entity.Repository) ([]entity.Reference, error) {
					return []entity.Reference{}, nil
				},
			}
			deviceReaderWriter = &reference.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, nil
				},
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, nil
				},
			}
			referenceReaderWriter = &reference.ReferenceReaderWriterMock{
				GetReferencesFunc: func(ctx context.Context, repo entity.Repository) ([]entity.Reference, error) {
					return []entity.Reference{}, nil
				},
			}
			service = reference.New(deviceReaderWriter, referenceReaderWriter, gitReader)
			err := service.UpdateReferences(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())
		})
	})
})
