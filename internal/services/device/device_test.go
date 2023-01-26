package device_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/device"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

var _ = Describe("Device", func() {
	Describe("Delete namespace", func() {
		It("correctly delete the default namespace", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: true, DeviceIDs: []string{"toto"}}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					namespaces := []entity.Namespace{
						{Name: "default", IsDefault: true},
						{Name: "other-one", IsDefault: false},
					}
					return namespaces, nil
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return nil
				},
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{ID: "toto", NamespaceID: "default"}, nil
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			namespace, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).To(BeNil())
			Expect(namespace.Name).To(Equal("default"))

			// calls to GetNamespaces should be one
			Expect(len(deviceReaderWriter.GetNamespacesCalls())).To(Equal(1))

			// should call GetDeviceFunc
			Expect(len(deviceReaderWriter.GetDeviceCalls())).To(Equal(1))

			// should call UpdateDevice to set the namespace
			updateDeviceCalls := deviceReaderWriter.UpdateCalls()
			Expect(len(updateDeviceCalls)).To(Equal(1))
			Expect(updateDeviceCalls[0].Device.NamespaceID).To(Equal("other-one"))

			// should change the default namespace
			updateCalls := deviceReaderWriter.UpdateNamespaceCalls()
			Expect(len(updateCalls)).To(Equal(1))
			Expect(updateCalls[0].Namespace.Name).To(Equal("other-one"))
			Expect(updateCalls[0].Namespace.IsDefault).To(BeTrue())

			// should expect call to delete namespace
			Expect(len(deviceReaderWriter.DeleteNamespaceCalls())).To(Equal(1))
		})

		It("correctly delete namespace not the default one", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: false}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					namespaces := []entity.Namespace{
						{Name: "default", IsDefault: false},
						{Name: "other-one", IsDefault: true},
					}
					return namespaces, nil
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			namespace, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).To(BeNil())
			Expect(namespace.Name).To(Equal("default"))

			// calls to GetNamespaces should be one
			Expect(len(deviceReaderWriter.GetNamespacesCalls())).To(Equal(1))

			// should not call update
			updateCalls := deviceReaderWriter.UpdateNamespaceCalls()
			Expect(len(updateCalls)).To(Equal(0))

			// should expect call to delete namespace
			Expect(len(deviceReaderWriter.DeleteNamespaceCalls())).To(Equal(1))
		})

		It("cannot delete the last namespace", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: false}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					namespaces := []entity.Namespace{
						{Name: "default", IsDefault: false},
					}
					return namespaces, nil
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("cannot delete the last namespace"))
		})

		It("return error when GetNamespaces fails", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: false}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					return nil, errors.New("error")
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).ToNot(BeNil())
		})

		It("return error when GetDevice fails", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: false, DeviceIDs: []string{"toto"}}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					namespaces := []entity.Namespace{
						{Name: "default", IsDefault: false},
						{Name: "other-one", IsDefault: true},
					}
					return namespaces, nil
				},
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errors.New("error")
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).ToNot(BeNil())
		})

		It("return error when UpdateDevice fails", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{Name: "default", IsDefault: false, DeviceIDs: []string{"toto"}}, nil
				},
				GetNamespacesFunc: func(ctx context.Context) ([]entity.Namespace, error) {
					namespaces := []entity.Namespace{
						{Name: "default", IsDefault: false},
						{Name: "other-one", IsDefault: true},
					}
					return namespaces, nil
				},
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, nil
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("error")
				},
				UpdateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetDefaultNamespaceFunc: func(ctx context.Context) (entity.Namespace, error) {
					// should return the other-one
					return entity.Namespace{Name: "other-one", IsDefault: true}, nil
				},
				DeleteNamespaceFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteNamespace(context.TODO(), "default")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Create namaspace", func() {
		It("correctly creates a namespace", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, errService.NewResourceNotFoundError("namespace", id)
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateNamespace(context.TODO(), entity.Namespace{Name: "default"})
			Expect(err).To(BeNil())
		})
		It("cannot create a namespace when it already exists", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateNamespace(context.TODO(), entity.Namespace{Name: "default"})
			Expect(err).ToNot(BeNil())
			Expect(err).To(BeAssignableToTypeOf(errService.ResourceAlreadyExists{}))
		})
		It("returns error from GetNamespace", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, errors.New("unknown")
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateNamespace(context.TODO(), entity.Namespace{Name: "default"})
			Expect(err).ToNot(BeNil())
		})
		It("returns error when creating", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateNamespaceFunc: func(ctx context.Context, namespace entity.Namespace) error {
					return errors.New("error")
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateNamespace(context.TODO(), entity.Namespace{Name: "default"})
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Create set", func() {
		It("correctly creates a set", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateSetFunc: func(ctx context.Context, set entity.Set) error {
					return nil
				},
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, errService.NewResourceNotFoundError("set", id)
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateSet(context.TODO(), entity.Set{Name: "default", NamespaceID: "default"})
			Expect(err).To(BeNil())
			Expect(len(deviceReaderWriter.CreateSetCalls())).To(Equal(1))
		})
		It("cannot create a set when it already exists", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateSet(context.TODO(), entity.Set{Name: "default"})
			Expect(err).ToNot(BeNil())
			Expect(err).To(BeAssignableToTypeOf(errService.ResourceAlreadyExists{}))
		})
		It("cannot create a set when namespaceID is missing", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, errService.NewResourceNotFoundError("set", id)
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateSet(context.TODO(), entity.Set{Name: "default"})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("namespace is missing"))
		})
		It("cannot create set when set's namespace doesn't exist", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, errService.NewResourceNotFoundError("namespace", id)
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateSet(context.TODO(), entity.Set{Name: "default"})
			Expect(err).ToNot(BeNil())
		})
		It("unable to create set", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				CreateSetFunc: func(ctx context.Context, set entity.Set) error {
					return errors.New("error")
				},
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, errService.NewResourceNotFoundError("set", id)
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.CreateSet(context.TODO(), entity.Set{Name: "default", NamespaceID: "default"})
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete set", func() {
		It("delete set successfuly", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, nil
				},
				DeleteSetFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteSet(context.TODO(), "id")
			Expect(err).To(BeNil())
		})
		It("cannot delete a set which doesn't exist", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, errService.NewResourceNotFoundError("set", id)
				},
				DeleteSetFunc: func(ctx context.Context, id string) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteSet(context.TODO(), "id")
			Expect(err).ToNot(BeNil())
		})
		It("delete set returns error", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, nil
				},
				DeleteSetFunc: func(ctx context.Context, id string) error {
					return errors.New("error")
				},
			}
			service := device.New(deviceReaderWriter)
			_, err := service.DeleteSet(context.TODO(), "id")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("Update device", func() {
		It("update correctly", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return nil
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.UpdateDevice(context.TODO(), entity.Device{})
			Expect(err).To(BeNil())
		})
		It("return error", func() {
			deviceReaderWriter := &device.DeviceReaderWriterMock{
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("error")
				},
			}
			service := device.New(deviceReaderWriter)
			err := service.UpdateDevice(context.TODO(), entity.Device{})
			Expect(err).ToNot(BeNil())
		})
	})
})
