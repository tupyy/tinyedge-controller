package manifest_test

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/errors"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
)

func allManifestsFilter(m entity.Manifest) bool {
	return true
}

var _ = Describe("manifests", func() {
	var (
		gitReader            *manifest.GitReaderMock
		deviceReaderWriter   *manifest.DeviceReaderMock
		manifestReaderWriter *manifest.ManifestReaderWriterMock
		service              *manifest.Service
		db                   *db
	)

	Describe("CRUD relations and manifests", Ordered, func() {
		db = NewDB()

		BeforeEach(func() {
			deviceReaderWriter = &manifest.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID: id,
					}, nil
				},
				GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{
						Name: id,
					}, nil
				},
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{
						Name: id,
					}, nil
				},
			}
			manifestReaderWriter = &manifest.ManifestReaderWriterMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					return db.GetManifests(), nil
				},
				GetManifestFunc: func(ctx context.Context, id string) (entity.Manifest, error) {
					m, ok := db.GetManifest(id)
					if !ok {
						return m, fmt.Errorf("not found")
					}
					return m, nil
				},
				InsertManifestFunc: func(ctx context.Context, manifest entity.Manifest) error {
					db.InsertManifest(manifest)
					return nil
				},
				CreateRelationFunc: func(ctx context.Context, relation entity.Relation) error {
					db.InsertRelation(relation)
					return nil
				},
				UpdateManifestFunc: func(ctx context.Context, manifest entity.Manifest) error {
					db.InsertManifest(manifest)
					return nil
				},
				DeleteManifestFunc: func(ctx context.Context, id string) error {
					db.DeleteManifest(id)
					return nil
				},
				DeleteRelationFunc: func(ctx context.Context, relation entity.Relation) error {
					db.DeleteRelation(relation)
					return nil
				},
			}
		})

		Context("namespace", func() {
			It("successfully creates a relation for a namespace", func() {
				gitReader = &manifest.GitReaderMock{
					GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
						workload := entity.Workload{
							ObjectMeta: entity.ObjectMeta{
								Id:   "test",
								Hash: "hash",
							},
							Selectors: entity.Selectors{
								{
									Type:  entity.NamespaceSelector,
									Value: "namespace",
								},
							},
						}
						manifest := []entity.Manifest{}
						manifest = append(manifest, workload)
						return manifest, nil
					},
				}
				service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
				err := service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount := db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(1), "expect 1 relation")
				r, ok := db.GetRelation("testnamespace")
				Expect(ok).To(BeTrue())
				Expect(r.ResourceID).To(Equal("namespace"))
				Expect(r.ManifestID).To(Equal("test"))
			})

			It("unable to create relation when namespace is missing", func() {
				gitReader = &manifest.GitReaderMock{
					GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
						workload := entity.Workload{
							ObjectMeta: entity.ObjectMeta{
								Id:   "test",
								Hash: "hash",
							},
							Selectors: entity.Selectors{
								{
									Type:  entity.NamespaceSelector,
									Value: "namespace",
								},
							},
						}
						manifest := []entity.Manifest{}
						manifest = append(manifest, workload)
						return manifest, nil
					},
				}
				deviceReaderWriter.GetNamespaceFunc = func(ctx context.Context, id string) (entity.Namespace, error) {
					return entity.Namespace{}, errors.NewResourceNotFoundError("namespace", id)
				}

				service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
				err := service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount := db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(0), "expect 1 relation")
			})

			It("successfully delete a relation for a namespace", func() {
				gitReader = &manifest.GitReaderMock{
					GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
						workload := entity.Workload{
							ObjectMeta: entity.ObjectMeta{
								Id:   "test",
								Hash: "hash",
							},
							Selectors: entity.Selectors{
								{
									Type:  entity.NamespaceSelector,
									Value: "namespace",
								},
							},
						}
						manifest := []entity.Manifest{}
						manifest = append(manifest, workload)
						return manifest, nil
					},
				}
				service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
				err := service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount := db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(1), "expect 1 relation")
				r, ok := db.GetRelation("testnamespace")
				Expect(ok).To(BeTrue())
				Expect(r.ResourceID).To(Equal("namespace"))
				Expect(r.ManifestID).To(Equal("test"))

				gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash1",
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				}
				m, _ := db.GetManifest("test")
				w := m.(entity.Workload)
				w.Namespaces = append(w.Namespaces, "namespace")
				db.InsertManifest(w)

				err = service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount = db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(0), "expect 0 relation")
			})

			It("successfully updates a relation for a namespace", func() {
				gitReader = &manifest.GitReaderMock{
					GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
						workload := entity.Workload{
							ObjectMeta: entity.ObjectMeta{
								Id:   "test",
								Hash: "hash",
							},
							Selectors: entity.Selectors{
								{
									Type:  entity.NamespaceSelector,
									Value: "namespace",
								},
							},
						}
						manifest := []entity.Manifest{}
						manifest = append(manifest, workload)
						return manifest, nil
					},
				}
				service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
				err := service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount := db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(1), "expect 1 relation")
				r, ok := db.GetRelation("testnamespace")
				Expect(ok).To(BeTrue())
				Expect(r.ResourceID).To(Equal("namespace"))
				Expect(r.ManifestID).To(Equal("test"))

				gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash1",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.NamespaceSelector,
								Value: "namespace1",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				}
				m, _ := db.GetManifest("test")
				w := m.(entity.Workload)
				w.Namespaces = append(w.Namespaces, "namespace")
				db.InsertManifest(w)

				err = service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount = db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(1), "expect 1 relation")
				r, ok = db.GetRelation("testnamespace1")
				Expect(ok).To(BeTrue())
				Expect(r.ResourceID).To(Equal("namespace1"))
				Expect(r.ManifestID).To(Equal("test"))
			})

			It("successfully delete relation and manifest when manifest removed from git", func() {
				gitReader = &manifest.GitReaderMock{
					GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
						workload := entity.Workload{
							ObjectMeta: entity.ObjectMeta{
								Id:   "test",
								Hash: "hash",
							},
							Selectors: entity.Selectors{
								{
									Type:  entity.NamespaceSelector,
									Value: "namespace",
								},
							},
						}
						manifest := []entity.Manifest{}
						manifest = append(manifest, workload)
						return manifest, nil
					},
				}
				service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
				err := service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount := db.Count()
				Expect(mCount).To(Equal(1), "expect 1 manifest")
				Expect(rCount).To(Equal(1), "expect 1 relation")
				r, ok := db.GetRelation("testnamespace")
				Expect(ok).To(BeTrue())
				Expect(r.ResourceID).To(Equal("namespace"))
				Expect(r.ManifestID).To(Equal("test"))

				gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					manifest := []entity.Manifest{}
					return manifest, nil
				}
				m, _ := db.GetManifest("test")
				w := m.(entity.Workload)
				w.Namespaces = append(w.Namespaces, "namespace")
				db.InsertManifest(w)

				err = service.UpdateManifests(context.TODO(), entity.Repository{})
				Expect(err).To(BeNil())

				// expect one manifest and one relation
				mCount, rCount = db.Count()
				Expect(mCount).To(Equal(0), "expect 0 manifest")
				Expect(rCount).To(Equal(0), "expect 0 relation")
			})
		})
	})

	Context("sets", func() {
		It("successfully creates a relation for a set", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok := db.GetRelation("testset")
			Expect(ok).To(BeTrue())
			Expect(r.ResourceID).To(Equal("set"))
			Expect(r.ManifestID).To(Equal("test"))
		})

		It("unable to create a relation when set is missing", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			d := &manifest.DeviceReaderMock{
				GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
					return entity.Set{}, errors.NewResourceNotFoundError("set", id)
				},
			}

			service = manifest.New(d, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(0), "expect 0 relation")
		})

		It("successfully creates a relation for a set and a namespace", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
							{
								Type:  entity.NamespaceSelector,
								Value: "namespace",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(2), "expect 2 relations")
		})

		It("successfully delete a relation for a set", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok := db.GetRelation("testset")
			Expect(ok).To(BeTrue())
			Expect(r.ResourceID).To(Equal("set"))
			Expect(r.ManifestID).To(Equal("test"))

			gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
				workload := entity.Workload{
					ObjectMeta: entity.ObjectMeta{
						Id:   "test",
						Hash: "hash1",
					},
				}
				manifest := []entity.Manifest{}
				manifest = append(manifest, workload)
				return manifest, nil
			}
			m, _ := db.GetManifest("test")
			w := m.(entity.Workload)
			w.Sets = append(w.Sets, "set")
			db.InsertManifest(w)

			err = service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount = db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(0), "expect 0 relation")
		})

		It("successfully updates a relation for a set", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok := db.GetRelation("testset")
			Expect(ok).To(BeTrue())
			Expect(r.ResourceID).To(Equal("set"))
			Expect(r.ManifestID).To(Equal("test"))

			gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
				workload := entity.Workload{
					ObjectMeta: entity.ObjectMeta{
						Id:   "test",
						Hash: "hash1",
					},
					Selectors: entity.Selectors{
						{
							Type:  entity.SetSelector,
							Value: "set1",
						},
					},
				}
				manifest := []entity.Manifest{}
				manifest = append(manifest, workload)
				return manifest, nil
			}
			m, _ := db.GetManifest("test")
			w := m.(entity.Workload)
			w.Sets = append(w.Sets, "set")
			db.InsertManifest(w)

			err = service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount = db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok = db.GetRelation("testset1")
			Expect(ok).To(BeTrue())
		})
	})

	Context("devices", func() {
		It("successfully creates a relation for a device", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
		})

		It("unable to create relation when device is missing", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			d := &manifest.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errors.NewResourceNotFoundError("device", id)
				},
			}
			service = manifest.New(d, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(0), "expect 0 relation")
		})

		It("successfully creates a relation for 1 set,1 namespace and 1 device", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
							{
								Type:  entity.NamespaceSelector,
								Value: "namespace",
							},
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(3), "expect 3 relations")
		})

		It("successfully deletes a relation for the set only", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.SetSelector,
								Value: "set",
							},
							{
								Type:  entity.NamespaceSelector,
								Value: "namespace",
							},
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(3), "expect 3 relations")

			gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
				workload := entity.Workload{
					ObjectMeta: entity.ObjectMeta{
						Id:   "test",
						Hash: "hash1",
					},
					Selectors: entity.Selectors{
						{
							Type:  entity.NamespaceSelector,
							Value: "namespace",
						},
						{
							Type:  entity.DeviceSelector,
							Value: "device",
						},
					},
				}
				manifest := []entity.Manifest{}
				manifest = append(manifest, workload)
				return manifest, nil
			}
			m, _ := db.GetManifest("test")
			w := m.(entity.Workload)
			w.Devices = append(w.Devices, "device")
			w.Namespaces = append(w.Namespaces, "namespace")
			db.InsertManifest(w)

			err = service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount = db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(2), "expect 2 relation")
		})

		It("successfully delete a relation for a device", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok := db.GetRelation("testdevice")
			Expect(ok).To(BeTrue())
			Expect(r.ResourceID).To(Equal("device"))
			Expect(r.ManifestID).To(Equal("test"))

			gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
				workload := entity.Workload{
					ObjectMeta: entity.ObjectMeta{
						Id:   "test",
						Hash: "hash1",
					},
				}
				manifest := []entity.Manifest{}
				manifest = append(manifest, workload)
				return manifest, nil
			}
			m, _ := db.GetManifest("test")
			w := m.(entity.Workload)
			w.Devices = append(w.Devices, "device")
			db.InsertManifest(w)

			err = service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			mCount, rCount = db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(0), "expect 0 relation")
		})

		It("successfully updates a relation for a device", func() {
			gitReader = &manifest.GitReaderMock{
				GetManifestsFunc: func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
					workload := entity.Workload{
						ObjectMeta: entity.ObjectMeta{
							Id:   "test",
							Hash: "hash",
						},
						Selectors: entity.Selectors{
							{
								Type:  entity.DeviceSelector,
								Value: "device",
							},
						},
					}
					manifest := []entity.Manifest{}
					manifest = append(manifest, workload)
					return manifest, nil
				},
			}
			service = manifest.New(deviceReaderWriter, manifestReaderWriter, gitReader)
			err := service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount := db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")

			gitReader.GetManifestsFunc = func(ctx context.Context, repo entity.Repository, filterFn func(m entity.Manifest) bool) ([]entity.Manifest, error) {
				workload := entity.Workload{
					ObjectMeta: entity.ObjectMeta{
						Id:   "test",
						Hash: "hash1",
					},
					Selectors: entity.Selectors{
						{
							Type:  entity.DeviceSelector,
							Value: "device1",
						},
					},
				}
				manifest := []entity.Manifest{}
				manifest = append(manifest, workload)
				return manifest, nil
			}
			m, _ := db.GetManifest("test")
			w := m.(entity.Workload)
			w.Devices = append(w.Devices, "device")
			db.InsertManifest(w)

			err = service.UpdateManifests(context.TODO(), entity.Repository{})
			Expect(err).To(BeNil())

			// expect one manifest and one relation
			mCount, rCount = db.Count()
			Expect(mCount).To(Equal(1), "expect 1 manifest")
			Expect(rCount).To(Equal(1), "expect 1 relation")
			r, ok := db.GetRelation("testdevice1")
			Expect(ok).To(BeTrue())
			Expect(r.ResourceID).To(Equal("device1"))
		})
	})

	AfterEach(func() {
		db.Clear()
	})
})

type db struct {
	Manifests map[string]entity.Manifest
	Relations map[string]entity.Relation
}

func NewDB() *db {
	return &db{
		Manifests: make(map[string]entity.Manifest),
		Relations: make(map[string]entity.Relation),
	}
}

func (d *db) Clear() {
	d.Manifests = make(map[string]entity.Manifest)
	d.Relations = make(map[string]entity.Relation)
}

func (d *db) Count() (int, int) {
	return len(d.Manifests), len(d.Relations)
}

func (d *db) InsertManifest(m entity.Manifest) {
	d.Manifests[m.GetID()] = m
}

func (d *db) InsertRelation(r entity.Relation) {
	d.Relations[fmt.Sprintf("%s%s", r.ManifestID, r.ResourceID)] = r
}

func (d *db) GetManifest(id string) (entity.Manifest, bool) {
	m, ok := d.Manifests[id]
	if !ok {
		return entity.Workload{}, false
	}
	return m, ok
}

func (d *db) GetManifests() []entity.Manifest {
	m := make([]entity.Manifest, 0, len(d.Manifests))
	for _, mm := range d.Manifests {
		m = append(m, mm)
	}
	return m
}

func (d *db) GetRelation(id string) (entity.Relation, bool) {
	r, ok := d.Relations[id]
	if !ok {
		return entity.Relation{}, false
	}
	return r, ok
}

func (d *db) GetRelations() []entity.Relation {
	r := make([]entity.Relation, 0, len(d.Relations))
	for _, rr := range d.Relations {
		r = append(r, rr)
	}
	return r
}

func (d *db) DeleteManifest(id string) {
	delete(d.Manifests, id)
	for id := range d.Relations {
		if strings.Contains(id, id) {
			delete(d.Relations, id)
		}
	}
}

func (d *db) DeleteRelation(r entity.Relation) {
	delete(d.Relations, fmt.Sprintf("%s%s", r.ManifestID, r.ResourceID))
}
