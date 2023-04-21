package postgres_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	pgRepo "github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"gorm.io/gorm"
)

var _ = Describe("Repository", Ordered, func() {
	var (
		pgClient  pg.Client
		rawClient pg.Client
		repo      *pgRepo.Repository
		gormDB    *gorm.DB
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

		repo, err = pgRepo.NewRepository(pgClient)
		Expect(err).To(BeNil())

		config := gorm.Config{
			SkipDefaultTransaction: true, // No need transaction for those use cases.
		}

		gormDB, err = rawClient.Open(config)
		Expect(err).To(BeNil())
	})

	Context("crud", func() {
		It("creates successfully a repo", func() {
			err := repo.InsertRepository(context.TODO(), entity.Repository{
				Id:                    "repo",
				AuthType:              entity.SSHRepositoryAuthType,
				Url:                   "url",
				LocalPath:             "/test",
				CredentialsSecretPath: "/secret",
				CurrentHeadSha:        "current",
				TargetHeadSha:         "target",
				PullPeriod:            2 * time.Second,
			})
			Expect(err).To(BeNil())

			r := models.Repo{}
			rerr := gormDB.Raw("select * from repo where id = 'repo';").Scan(&r).Error
			Expect(rerr).To(BeNil())
			Expect(r.AuthType.Value()).To(Equal("ssh"))
			Expect(r.URL).To(Equal("url"))
			Expect(r.LocalPath.Value()).To(Equal("/test"))
			Expect(r.AuthSecretPath.Value()).To(Equal("/secret"))
			Expect(r.CurrentHeadSha.Value()).To(Equal("current"))
			Expect(r.TargetHeadSha.Value()).To(Equal("target"))
			Expect(r.PullPeriodSeconds.Value()).To(Equal(int64(2)))
		})

		It("updates successfully a repo", func() {
			initialRepo := entity.Repository{
				Id:                    "repo",
				AuthType:              entity.SSHRepositoryAuthType,
				Url:                   "url",
				LocalPath:             "/test",
				CredentialsSecretPath: "/secret",
				CurrentHeadSha:        "current",
				TargetHeadSha:         "target",
				PullPeriod:            2 * time.Second,
			}
			err := repo.InsertRepository(context.TODO(), initialRepo)
			Expect(err).To(BeNil())

			initialRepo.Url = "newurl"
			err = repo.UpdateRepository(context.TODO(), initialRepo)
			Expect(err).To(BeNil())

			r, err := repo.GetRepository(context.TODO(), "repo")
			Expect(err).To(BeNil())
			Expect(r).To(Equal(initialRepo))
		})

		It("retrieve successfully a repo", func() {
			initialRepo := entity.Repository{
				Id:                    "repo",
				AuthType:              entity.SSHRepositoryAuthType,
				Url:                   "url",
				LocalPath:             "/test",
				CredentialsSecretPath: "/secret",
				CurrentHeadSha:        "current",
				TargetHeadSha:         "target",
				PullPeriod:            2 * time.Second,
			}
			err := repo.InsertRepository(context.TODO(), initialRepo)
			Expect(err).To(BeNil())

			r, err := repo.GetRepository(context.TODO(), "repo")
			Expect(err).To(BeNil())
			Expect(r).To(Equal(initialRepo))
		})

		It("retrieve successfully all repo", func() {
			initialRepo := entity.Repository{
				Id:                    "repo",
				AuthType:              entity.SSHRepositoryAuthType,
				Url:                   "url",
				LocalPath:             "/test",
				CredentialsSecretPath: "/secret",
				CurrentHeadSha:        "current",
				TargetHeadSha:         "target",
				PullPeriod:            2 * time.Second,
			}
			err := repo.InsertRepository(context.TODO(), initialRepo)
			Expect(err).To(BeNil())

			repos, err := repo.GetRepositories(context.TODO())
			Expect(err).To(BeNil())
			Expect(len(repos)).To(Equal(1))
			Expect(repos[0]).To(Equal(initialRepo))
		})
	})

	AfterEach(func() {
		// clean the db
		gormDB.Exec("DELETE FROM repo;")
	})

})
