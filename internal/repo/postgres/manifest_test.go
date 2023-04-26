package postgres_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	pgRepo "github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"gorm.io/gorm"
)

var _ = Describe("Manifest repository", Ordered, func() {
	var (
		pgClient  pg.Client
		rawClient pg.Client
		repo      *pgRepo.ManifestRepository
		gormDB    *gorm.DB
		folderTmp string
		workload  string
	)

	BeforeAll(func() {
		var err error
		port, _ := strconv.Atoi(getEnvVar("POSTGRES_PORT", "5433"))
		pgClient, err = pg.New(pg.ClientParams{
			Host:     getEnvVar("POSTGRES_HOST", "localhost"),
			Port:     uint(port),
			DBName:   getEnvVar("POSTGRES_DB", "tinyedge"),
			User:     getEnvVar("POSTGRES_USER", "postgres"),
			Password: getEnvVar("POSTGRES_PWD", "postgres"),
		})
		Expect(err).To(BeNil())

		rawClient, err = pg.New(pg.ClientParams{
			Host:     getEnvVar("POSTGRES_HOST", "localhost"),
			Port:     uint(port),
			DBName:   getEnvVar("POSTGRES_DB", "tinyedge"),
			User:     getEnvVar("POSTGRES_USER", "postgres"),
			Password: getEnvVar("POSTGRES_PWD", "postgres"),
		})
		Expect(err).To(BeNil())

		repo, err = pgRepo.NewManifestRepository(pgClient)
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

		workload, err = writeManifest(tmpDir, workload)
		Expect(err).To(BeNil())

		conf, err := writeManifest(tmpDir, configuration)
		Expect(err).To(BeNil())

		folderTmp = tmpDir
		gormDB.Exec(fmt.Sprintf("INSERT INTO repo (id,url,local_path) VALUES('id','url','%s');", folderTmp))
		gormDB.Exec(fmt.Sprintf(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
			('configuration', 'configuration', 'configuration', 'id', '%s');`, conf))
		gormDB.Exec(`INSERT INTO namespace (id,is_default, configuration_manifest_id) VALUES 
			('namespace1', false, 'configuration'),
			('namespace', true, 'configuration');`)
		gormDB.Exec(`INSERT into device_set (id, namespace_id, configuration_manifest_id) VALUES 
			('set', 'namespace', 'configuration'),
			('set1', 'namespace1', 'configuration');`)
		gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id, device_set_id) VALUES
			('device1', true, true, 'namespace', 'set'),
			('device2', true, true, 'namespace1', 'set1');`)
	})

	Context("crud manifests", func() {
		It("successfully insert a manifest", func() {
			manifest := entity.Workload{
				ObjectMeta: entity.ObjectMeta{
					Id: "workload",
				},
				TypeMeta: entity.TypeMeta{
					Kind: entity.WorkloadManifestKind,
				},
				Repository: entity.Repository{
					Id: "id",
				},
				Path: workload,
			}
			err := repo.InsertManifest(context.TODO(), manifest)
			Expect(err).To(BeNil())

			count := 0
			rerr := gormDB.Raw("SELECT count(*) from manifest;").Scan(&count).Error
			Expect(rerr).To(BeNil())
			Expect(count).To(Equal(2))

			m := models.Manifest{}
			rerr = gormDB.Raw("SELECT * from manifest where id = 'workload'").Scan(&m).Error
			Expect(rerr).To(BeNil())
			Expect(m.Path).To(Equal(workload))
			Expect(m.RepoID).To(Equal("id"))
			Expect(m.RefType).To(Equal("workload"))
		})

		It("successfully update a manifest", func() {
			manifest := entity.Workload{
				ObjectMeta: entity.ObjectMeta{
					Id: "workload",
				},
				TypeMeta: entity.TypeMeta{
					Kind: entity.WorkloadManifestKind,
				},
				Repository: entity.Repository{
					Id: "id",
				},
				Path: workload,
			}
			err := repo.InsertManifest(context.TODO(), manifest)
			Expect(err).To(BeNil())

			manifest.Path = "test"
			err = repo.UpdateManifest(context.TODO(), manifest)
			Expect(err).To(BeNil())

			count := 0
			rerr := gormDB.Raw("SELECT count(*) from manifest;").Scan(&count).Error
			Expect(rerr).To(BeNil())
			Expect(count).To(Equal(2))

			m := models.Manifest{}
			rerr = gormDB.Raw("SELECT * from manifest where id = 'workload'").Scan(&m).Error
			Expect(rerr).To(BeNil())
			Expect(m.Path).To(Equal("test"))
			Expect(m.RepoID).To(Equal("id"))
			Expect(m.RefType).To(Equal("workload"))
		})

		It("successfully delete a manifest", func() {
			manifest := entity.Workload{
				ObjectMeta: entity.ObjectMeta{
					Id: "workload",
				},
				TypeMeta: entity.TypeMeta{
					Kind: entity.WorkloadManifestKind,
				},
				Repository: entity.Repository{
					Id: "id",
				},
				Path: workload,
			}
			err := repo.InsertManifest(context.TODO(), manifest)
			Expect(err).To(BeNil())

			count := 0
			rerr := gormDB.Raw("SELECT count(*) from manifest;").Scan(&count).Error
			Expect(rerr).To(BeNil())
			Expect(count).To(Equal(2))

			err = repo.DeleteManifest(context.TODO(), manifest.Id)
			Expect(err).To(BeNil())

			count = 3
			rerr = gormDB.Raw("SELECT count(*) from manifest;").Scan(&count).Error
			Expect(rerr).To(BeNil())
			Expect(count).To(Equal(1))
		})

		Context("relations", func() {
			It("creates successfully relation between namespace and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				rerr := repo.CreateRelation(context.TODO(), entity.NewNamespaceRelation("namespace", "workload"))
				Expect(rerr).To(BeNil())

				m := models.NamespacesManifests{}
				err = gormDB.Raw("select * from namespaces_manifests where manifest_id = 'workload';").Scan(&m).Error
				Expect(err).To(BeNil())
				Expect(m.NamespaceID).To(Equal("namespace"))
			})

			It("creates successfully relation between set and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				rerr := repo.CreateRelation(context.TODO(), entity.NewSetRelation("set", "workload"))
				Expect(rerr).To(BeNil())

				m := models.SetsManifests{}
				err = gormDB.Raw("select * from sets_manifests where manifest_id = 'workload';").Scan(&m).Error
				Expect(err).To(BeNil())
				Expect(m.DeviceSetID).To(Equal("set"))
			})

			It("creates successfully relation between device and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				err = gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
				('device', 'enroled', true, 'namespace1');`).Error
				Expect(err).To(BeNil())

				rerr := repo.CreateRelation(context.TODO(), entity.NewDeviceRelation("device", "workload"))
				Expect(rerr).To(BeNil())

				m := models.DevicesManifests{}
				err = gormDB.Raw("select * from devices_manifests where manifest_id = 'workload';").Scan(&m).Error
				Expect(err).To(BeNil())
				Expect(m.DeviceID).To(Equal("device"))
			})

			It("delete successfully relation between namespace and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				err = gormDB.Exec(`INSERT INTO namespaces_manifests (namespace_id, manifest_id) VALUES
					('namespace', 'workload');`).Error
				Expect(err).To(BeNil())

				rerr := repo.DeleteRelation(context.TODO(), entity.NewNamespaceRelation("namespace", "workload"))
				Expect(rerr).To(BeNil())

				count := 0
				err = gormDB.Raw("select * from namespaces_manifests where manifest_id = 'workload';").Scan(&count).Error
				Expect(count).To(BeZero())
			})

			It("delete successfully relation between set and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				err = gormDB.Exec(`INSERT INTO sets_manifests (device_set_id, manifest_id) VALUES
					('set', 'workload');`).Error
				Expect(err).To(BeNil())

				rerr := repo.DeleteRelation(context.TODO(), entity.NewSetRelation("set", "workload"))
				Expect(rerr).To(BeNil())

				count := 0
				err = gormDB.Raw("select * from sets_manifests where manifest_id = 'workload';").Scan(&count).Error
				Expect(count).To(BeZero())
			})

			It("delete successfully relation between device and manifest", func() {
				err := gormDB.Exec(`INSERT INTO manifest (id, ref_type, name, repo_id, path) VALUES
					('workload', 'workload', 'workload', 'id', 'test');`).Error
				Expect(err).To(BeNil())

				err = gormDB.Exec(`INSERT INTO device (id, enroled, registered, namespace_id) VALUES
				('device', 'enroled', true, 'namespace1');`).Error
				Expect(err).To(BeNil())

				err = gormDB.Exec(`INSERT INTO devices_manifests (device_id, manifest_id) VALUES
					('device', 'workload');`).Error
				Expect(err).To(BeNil())

				rerr := repo.DeleteRelation(context.TODO(), entity.NewDeviceRelation("device", "workload"))
				Expect(rerr).To(BeNil())

				count := 0
				err = gormDB.Raw("select * from device_manifests where manifest_id = 'workload';").Scan(&count).Error
				Expect(count).To(BeZero())
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