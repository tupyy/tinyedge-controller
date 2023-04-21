package postgres_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	pgRepo "github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"gorm.io/gorm"
)

var (
	configuration = `
kind: configuration
version: v1

name: conf

description: |
  blabla
`
	workload = `
kind: workload
version: v1

name: dasdaot

description: |
  blabla

selectors:
  namespaces:
    - test
    - ggg
  sets:
    - ttt
    - fff
  devices:
    - toto

secrets:
  - id: nginx-password
    path: nginx
    key: data

resources:
  - $ref: /dep/configmap.yaml
  - $ref: /dep/nginx.yaml
  - $ref: /dep/postgres.yaml
`
)

var _ = Describe("Device repository", Ordered, func() {
	var (
		pgClient  pg.Client
		rawClient pg.Client
		//repo       *pgRepo.ManifestRepository
		deviceRepo *pgRepo.DeviceRepo
		gormDB     *gorm.DB
		folderTmp  string
	)

	BeforeAll(func() {
		var err error
		pgClient, err = pg.New(pg.ClientParams{
			Host:     "localhost",
			Port:     5433,
			DBName:   "tinyedge",
			User:     "postgres",
			Password: "postgres",
		})
		Expect(err).To(BeNil())

		rawClient, err = pg.New(pg.ClientParams{
			Host:     "localhost",
			Port:     5433,
			DBName:   "tinyedge",
			User:     "postgres",
			Password: "postgres",
		})
		Expect(err).To(BeNil())

		// repo, err = pgRepo.NewManifestRepository(pgClient)
		// Expect(err).To(BeNil())

		deviceRepo, err = pgRepo.NewDeviceRepo(pgClient)
		Expect(err).To(BeNil())

		config := gorm.Config{
			SkipDefaultTransaction: true, // No need transaction for those use cases.
		}

		gormDB, err = rawClient.Open(config)
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		tmpDir, err := os.MkdirTemp("", "repo-*")
		Expect(err).To(BeNil())

		workload, err := writeManifest(tmpDir, workload)
		Expect(err).To(BeNil())

		conf, err := writeManifest(tmpDir, configuration)
		Expect(err).To(BeNil())

		folderTmp = tmpDir
		gormDB.Exec(fmt.Sprintf("INSERT INTO repo (id,url,local_path) VALUES('id','url','%s');", folderTmp))
		gormDB.Exec(fmt.Sprintf(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
			('workload', 'workload', 'workload', 'id', '%s'),
			('workload2', 'workload','workload2','id','%s'),
			('configuration', 'configuration', 'configuration', 'id', '%s');`, workload, workload, conf))
	})

	Context("namespace", func() {
		It("creates successfully a namespace", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))
		})

		It("cannot change the default namespace when only one exists", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.UpdateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).NotTo(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			var n = struct {
				IsDefault bool
			}{}
			gormDB.Raw("SELECT is_default from namespace where id = ?;", "test").Scan(&n)
			Expect(n.IsDefault).To(BeTrue())
		})

		It("updating a namespace successfully change the default one", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test1",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.UpdateNamespace(context.TODO(), entity.Namespace{
				Name:      "test1",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			var n = struct {
				IsDefault bool
			}{}
			gormDB.Raw("SELECT is_default from namespace where id = ?;", "test1").Scan(&n)
			Expect(n.IsDefault).To(BeTrue())
		})

		It("updating a namespace successfully", func() {
			gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
			('configuration1', 'configuration', 'configuration', 'id', 'test');`)

			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.UpdateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration1",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			var n = struct {
				ConfigurationManifestID string
			}{}
			gormDB.Raw("SELECT configuration_manifest_id from namespace where id = ?;", "test").Scan(&n)
			Expect(n.ConfigurationManifestID).To(Equal("configuration1"))
		})

		It("cannot remove last namespace", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			err = deviceRepo.DeleteNamespace(context.TODO(), "test")
			Expect(err).ToNot(BeNil())
		})

		It("change default namespace successfully", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "isdefault",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			n := models.Namespace{}
			gormDB.Raw("SELECT * from namespace where is_default = true;").Scan(&n)
			Expect(n.ID).To(Equal("isdefault"))
			Expect(n.IsDefault.Value()).To(BeTrue())
		})

		It("removing the default namespace will set the next one to default", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test1",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			err = deviceRepo.DeleteNamespace(context.TODO(), "test")
			Expect(err).To(BeNil())

			count = 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			n := models.Namespace{}
			gormDB.Raw("SELECT * from namespace where is_default = true;").Scan(&n)
			Expect(n.ID).To(Equal("test1"))
			Expect(n.IsDefault.Value()).To(BeTrue())

		})

		It("get default namespace successfully", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "isdefault",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			n, err := deviceRepo.GetDefaultNamespace(context.TODO())
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("isdefault"))
			Expect(n.IsDefault).To(BeTrue())
		})

		It("get namespaces successfully", func() {
			err := gormDB.Exec(`INSERT INTO namespace (id, is_default, configuration_manifest_id) VALUES
				('first',true,'configuration'),
				('second',false,'configuration');`).Error
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			n, err := deviceRepo.GetNamespaces(context.TODO())
			Expect(err).To(BeNil())
			Expect(n).To(HaveLen(2))
		})

		It("get namespace with device successfully", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "test",
				IsDefault: true,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "otherone",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			err = deviceRepo.CreateDevice(context.TODO(), entity.Device{
				ID:          "device",
				EnrolStatus: entity.EnroledStatus,
				Registred:   true,
				NamespaceID: "otherone",
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(2))

			n, err := deviceRepo.GetNamespace(context.TODO(), "otherone")
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("otherone"))
			Expect(n.IsDefault).To(BeFalse())
			Expect(len(n.Devices)).To(Equal(1))
		})

		It("get namespace with workload and configuration successfully", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "otherone",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			tx := gormDB.Exec(`INSERT INTO namespaces_manifests (namespace_id, manifest_id) VALUES
			('otherone', 'workload'),
			('otherone', 'workload2');`)
			Expect(tx.Error).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			n, err := deviceRepo.GetNamespace(context.TODO(), "otherone")
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("otherone"))
			Expect(n.IsDefault).To(BeTrue())
			Expect(n.Configuration.Id).To(Equal("configuration"))
			Expect(n.Configuration.GetKind().String()).To(Equal("configuration"))

			// workload
			Expect(len(n.Workloads)).To(Equal(2))
			Expect(n.Workloads[0].GetID()).To(Equal("workload"))
		})

		It("get namespace with devices", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "otherone",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			tx := gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
			('device1', true, true, 'otherone'),
			('device2', true, true, 'otherone');`)
			Expect(tx.Error).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			n, err := deviceRepo.GetNamespace(context.TODO(), "otherone")
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("otherone"))
			Expect(n.IsDefault).To(BeTrue())
			Expect(n.Configuration.Id).To(Equal("configuration"))
			Expect(n.Configuration.GetKind().String()).To(Equal("configuration"))

			// devices
			Expect(len(n.Devices)).To(Equal(2))
			Expect(n.Devices[0]).To(Equal("device1"))
			Expect(n.Devices[1]).To(Equal("device2"))
		})

		It("get namespace with one device", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "otherone",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			tx := gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
			('device2', true, true, 'otherone');`)
			Expect(tx.Error).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			n, err := deviceRepo.GetNamespace(context.TODO(), "otherone")
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("otherone"))
			Expect(n.IsDefault).To(BeTrue())
			Expect(n.Configuration.Id).To(Equal("configuration"))
			Expect(n.Configuration.GetKind().String()).To(Equal("configuration"))

			// devices
			Expect(len(n.Devices)).To(Equal(1))
			Expect(n.Devices[0]).To(Equal("device2"))
		})

		It("get namespace with one set", func() {
			err := deviceRepo.CreateNamespace(context.TODO(), entity.Namespace{
				Name:      "otherone",
				IsDefault: false,
				Configuration: entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			tx := gormDB.Exec(`INSERT INTO device_set (id, namespace_id) VALUES
			('set1', 'otherone');`)
			Expect(tx.Error).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from namespace;").Scan(&count)
			Expect(count).To(Equal(1))

			n, err := deviceRepo.GetNamespace(context.TODO(), "otherone")
			Expect(err).To(BeNil())
			Expect(n.Name).To(Equal("otherone"))
			Expect(n.IsDefault).To(BeTrue())
			Expect(n.Configuration.Id).To(Equal("configuration"))
			Expect(n.Configuration.GetKind().String()).To(Equal("configuration"))

			// devices
			Expect(len(n.Devices)).To(Equal(0))
			Expect(len(n.Sets)).To(Equal(1))
			Expect(n.Sets[0]).To(Equal("set1"))
		})
	})

	Context("set", func() {
		BeforeEach(func() {
			gormDB.Exec(`INSERT INTO namespace (id, is_default, configuration_manifest_id) VALUES
			('namespace1', true, 'configuration'),
			('namespace2', false, 'configuration');`)
		})

		It("create successfully a set", func() {
			err := deviceRepo.CreateSet(context.TODO(), entity.Set{
				Name:        "set1",
				NamespaceID: "namespace1",
				Configuration: &entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			count := 0
			gormDB.Raw("SELECT count(*) from device_set;").Scan(&count)
			Expect(count).To(Equal(1))
		})

		It("unable to create set without namespace", func() {
			err := deviceRepo.CreateSet(context.TODO(), entity.Set{
				Name: "set1",
				Configuration: &entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).ToNot(BeNil())
		})

		It("successfully retrieve a set with configuration", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set1', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			set, err := deviceRepo.GetSet(context.TODO(), "set1")
			Expect(err).To(BeNil())

			Expect(set.Configuration).ToNot(BeNil())
			Expect(set.Configuration.Id).To(Equal("configuration"))
			Expect(set.Configuration.GetKind().String()).To(Equal("configuration"))
		})

		It("successfully retrieve a set with devices", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set1', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			eerr = gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id, device_set_id) VALUES
			('device1', true, true, 'namespace1', 'set1'),
			('device2', true, true, 'namespace1', 'set1');`).Error
			Expect(eerr).To(BeNil())

			set, err := deviceRepo.GetSet(context.TODO(), "set1")
			Expect(err).To(BeNil())
			Expect(len(set.Devices)).To(Equal(2))
			Expect(set.Devices[0]).To(Equal("device1"))
			Expect(set.Devices[1]).To(Equal("device2"))
		})

		It("successfully retrieve a set with workloads", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set1', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			eerr = gormDB.Exec(`INSERT INTO sets_manifests (manifest_id, device_set_id) VALUES
			('workload', 'set1'),
			('workload2', 'set1');`).Error
			Expect(eerr).To(BeNil())

			set, err := deviceRepo.GetSet(context.TODO(), "set1")
			Expect(err).To(BeNil())
			Expect(len(set.Workloads)).To(Equal(2))
			Expect(set.Workloads[0].GetID()).To(Equal("workload"))
			Expect(set.Workloads[1].GetID()).To(Equal("workload2"))
		})

		It("successfully retrieve all sets", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set', 'namespace2','configuration'),
				('set1', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			sets, err := deviceRepo.GetSets(context.TODO())
			Expect(err).To(BeNil())
			Expect(len(sets)).To(Equal(2))
			Expect([]string{sets[0].Name, sets[1].Name}).Should(ContainElement("set1"))
			Expect([]string{sets[0].Name, sets[1].Name}).Should(ContainElement("set"))
		})

		It("successfully delete set", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			err := deviceRepo.DeleteSet(context.TODO(), "set")
			Expect(err).To(BeNil())

			count := 0
			eerr = gormDB.Raw("SELECT count(*) from device_set;").Scan(&count).Error
			Expect(eerr).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("successfully update set", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			err := deviceRepo.UpdateSet(context.TODO(), entity.Set{
				Name:        "set",
				NamespaceID: "namespace1",
			})
			Expect(err).To(BeNil())

			set, err := deviceRepo.GetSet(context.TODO(), "set")
			Expect(err).To(BeNil())
			Expect(set.Configuration).To(BeNil())
		})

		It("successfully update set with new namespace id", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			err := deviceRepo.UpdateSet(context.TODO(), entity.Set{
				Name:        "set",
				NamespaceID: "namespace2",
				Configuration: &entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).To(BeNil())

			set, err := deviceRepo.GetSet(context.TODO(), "set")
			Expect(err).To(BeNil())
			Expect(set.Configuration).NotTo(BeNil())
			Expect(set.NamespaceID).To(Equal("namespace2"))
		})

		It("unable to update set when namespace does not exists", func() {
			eerr := gormDB.Exec(`INSERT INTO device_set (id, namespace_id, configuration_manifest_id) VALUES
				('set', 'namespace1', 'configuration');`).Error
			Expect(eerr).To(BeNil())

			err := deviceRepo.UpdateSet(context.TODO(), entity.Set{
				Name:        "set",
				NamespaceID: "nonamespace",
				Configuration: &entity.Configuration{
					ObjectMeta: entity.ObjectMeta{
						Id: "configuration",
					},
				},
			})
			Expect(err).NotTo(BeNil())
		})

		Context("device", func() {
			It("successfully retrieve a device", func() {
				tx := gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
				('device', 'enroled', true, 'namespace1');`)
				Expect(tx.Error).To(BeNil())

				device, err := deviceRepo.GetDevice(context.TODO(), "device")
				Expect(err).To(BeNil())
				Expect(device.ID).To(Equal("device"))
				Expect(device.NamespaceID).To(Equal("namespace1"))
				Expect(device.EnrolStatus.String()).To(Equal("enroled"))
				Expect(device.SetID).To(BeNil())
			})

			It("successfully retrieve a device with workloads", func() {
				tx := gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
				('device', 'enroled', true, 'namespace1');`)
				Expect(tx.Error).To(BeNil())

				tx = gormDB.Exec(`INSERT INTO devices_manifests (device_id, manifest_id) VALUES
				('device', 'workload'),
				('device', 'workload2');`)
				Expect(tx.Error).To(BeNil())

				device, err := deviceRepo.GetDevice(context.TODO(), "device")
				Expect(err).To(BeNil())
				Expect(device.ID).To(Equal("device"))
				Expect(device.NamespaceID).To(Equal("namespace1"))
				Expect(device.SetID).To(BeNil())
				Expect(len(device.Workloads)).To(Equal(2))
				Expect([]string{device.Workloads[0].GetID(), device.Workloads[1].GetID()}).Should(ContainElement("workload"))
				Expect([]string{device.Workloads[0].GetID(), device.Workloads[1].GetID()}).Should(ContainElement("workload2"))
			})

			It("successfully retrieve all devices", func() {
				tx := gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
				('device1', 'enroled', true, 'namespace1'),
				('device', 'enroled', true, 'namespace1');`)
				Expect(tx.Error).To(BeNil())

				devices, err := deviceRepo.GetDevices(context.TODO())
				Expect(err).To(BeNil())
				Expect(len(devices)).To(Equal(2))
				Expect([]string{devices[0].ID, devices[1].ID}).Should(ContainElement("device"))
				Expect([]string{devices[0].ID, devices[1].ID}).Should(ContainElement("device1"))
			})

			It("successfully creates a device", func() {
				err := deviceRepo.CreateDevice(context.TODO(), entity.Device{
					ID:          "device",
					EnrolStatus: entity.PendingEnrolStatus,
					Registred:   false,
					NamespaceID: "namespace1",
					Configuration: &entity.Configuration{
						ObjectMeta: entity.ObjectMeta{Id: "configuration"},
					},
				})
				Expect(err).To(BeNil())

				device, err := deviceRepo.GetDevice(context.TODO(), "device")
				Expect(err).To(BeNil())
				Expect(device.EnrolStatus.String()).To(Equal("pending"))
				Expect(device.Registred).To(BeFalse())
				Expect(device.NamespaceID).To(Equal("namespace1"))
			})

			It("successfully delete a device", func() {
				err := deviceRepo.CreateDevice(context.TODO(), entity.Device{
					ID:          "device",
					EnrolStatus: entity.PendingEnrolStatus,
					Registred:   false,
					NamespaceID: "namespace1",
					Configuration: &entity.Configuration{
						ObjectMeta: entity.ObjectMeta{Id: "configuration"},
					},
				})
				Expect(err).To(BeNil())

				err = deviceRepo.DeleteDevice(context.TODO(), "device")
				Expect(err).To(BeNil())

				count := 1
				eerr := gormDB.Raw("SELECT count(*) from device;").Scan(&count).Error
				Expect(eerr).To(BeNil())
				Expect(count).To(Equal(0))
			})

			It("successfully update a device", func() {
				device := entity.Device{
					ID:          "device",
					EnrolStatus: entity.PendingEnrolStatus,
					Registred:   false,
					NamespaceID: "namespace1",
					Configuration: &entity.Configuration{
						ObjectMeta: entity.ObjectMeta{Id: "configuration"},
					},
				}

				err := deviceRepo.CreateDevice(context.TODO(), device)
				Expect(err).To(BeNil())

				device.EnrolStatus = entity.EnroledStatus
				err = deviceRepo.UpdateDevice(context.TODO(), device)
				Expect(err).To(BeNil())

				d, err := deviceRepo.GetDevice(context.TODO(), device.ID)
				Expect(err).To(BeNil())
				Expect(d.EnrolStatus.String()).To(Equal("enroled"))

			})
		})
	})

	AfterEach(func() {
		// clean the db
		gormDB.Exec("DELETE FROM manifest;")
		gormDB.Exec("DELETE FROM device;")
		gormDB.Exec("DELETE FROM namespace;")
		gormDB.Exec("DELETE FROM device_set;")
		gormDB.Exec("DELETE FROM repo;")
		os.RemoveAll(folderTmp)
	})

	AfterAll(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pgClient.Shutdown(ctx)
		rawClient.Shutdown(ctx)
	})

})

func writeManifest(folder, content string) (string, error) {
	data := bytes.NewBufferString(content).Bytes()
	// write workload
	f, err := os.CreateTemp(folder, "file-*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(data)
	return path.Base(f.Name()), err
}
