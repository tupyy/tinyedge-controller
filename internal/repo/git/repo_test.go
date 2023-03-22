package git_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	gitRepo "github.com/tupyy/tinyedge-controller/internal/repo/git"
)

var (
	manifest1 = `
version: v1
kind: workload

name: manifest1

description: |
  blabla

selectors:

secrets:
  - id: nginx-password
    path: nginx
    key: data

resources:
  - $ref: /dep/configmap.yaml
  - $ref: /dep/nginx.yaml
  - $ref: /dep/postgres.yaml
`

	manifest2 = `
version: v1
kind: workload

name: manifest2

description: |
  blabla

selectors:

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

var _ = Describe("Git repository", func() {
	Context("repo operations", func() {
		var (
			commit      string
			tmpDir      string
			cloneDir    string
			filename    string
			workingTree *git.Worktree
		)
		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "git-*")
			Expect(err).To(BeNil())

			fs := osfs.New(tmpDir)

			r, err := git.Init(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), fs)
			Expect(err).To(BeNil())

			filename = filepath.Join(tmpDir, "example-git-file")
			err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)

			w, err := r.Worktree()
			Expect(err).To(BeNil())
			workingTree = w

			_, err = w.Add("example-git-file")
			Expect(err).To(BeNil())

			c, err := w.Commit("first commit", &git.CommitOptions{
				Author: &object.Signature{
					Name:  "John Doe",
					Email: "j@doe.org",
					When:  time.Now(),
				},
			})
			commit = c.String()
			cloneDir, err = os.MkdirTemp("", "git-clone-*")
			Expect(err).To(BeNil())
		})

		It("clone successfully the repo", func() {
			repo := entity.Repository{
				Id:       "test",
				Url:      tmpDir,
				AuthType: entity.NoRepositoryAuthType,
			}

			r := gitRepo.New(cloneDir)
			clone, err := r.Clone(context.TODO(), repo)
			Expect(err).To(BeNil())
			Expect(clone.Branch).To(Equal("master"))
			fmt.Println(commit)
		})

		It("successfully get the head sha", func() {
			repo := entity.Repository{
				Id:       "test",
				Url:      tmpDir,
				AuthType: entity.NoRepositoryAuthType,
			}

			r := gitRepo.New(cloneDir)
			clone, err := r.Clone(context.TODO(), repo)
			Expect(err).To(BeNil())
			Expect(clone.Branch).To(Equal("master"))

			headSha, err := r.GetHeadSha(context.Background(), clone)
			Expect(headSha).To(Equal(commit))
		})

		It("successfully pull the repo", func() {
			repo := entity.Repository{
				Id:       "test",
				Url:      tmpDir,
				AuthType: entity.NoRepositoryAuthType,
			}

			r := gitRepo.New(cloneDir)
			clone, err := r.Clone(context.TODO(), repo)
			Expect(err).To(BeNil())
			Expect(clone.Branch).To(Equal("master"))

			headSha, err := r.GetHeadSha(context.Background(), clone)
			Expect(headSha).To(Equal(commit))

			// commit something else
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
			Expect(err).To(BeNil())
			_, err = f.Write([]byte("test"))
			Expect(err).To(BeNil())
			f.Close()

			newCommit, err := workingTree.Commit("second commit", &git.CommitOptions{
				Author: &object.Signature{
					Name:  "John Doe",
					Email: "j@doe.org",
					When:  time.Now(),
				},
			})

			// pull
			err = r.Pull(context.TODO(), clone)
			Expect(err).To(BeNil())
			headSha, err = r.GetHeadSha(context.Background(), clone)
			Expect(headSha).To(Equal(newCommit.String()))

		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
			os.RemoveAll(cloneDir)
		})
	})

	Context("find operations", func() {
		var (
			tmpDir   string
			cloneDir string
		)
		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "git-*")
			Expect(err).To(BeNil())

			fs := osfs.New(tmpDir)

			r, err := git.Init(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), fs)
			Expect(err).To(BeNil())

			os.Mkdir(path.Join(tmpDir, "folder1"), 0755)
			os.Mkdir(path.Join(tmpDir, "folder2"), 0755)

			file1 := filepath.Join(tmpDir, "folder1", "test.manifest.yaml")
			err = ioutil.WriteFile(file1, []byte(manifest1), 0644)

			file2 := filepath.Join(tmpDir, "folder2", "test1.manifest.yml")
			err = ioutil.WriteFile(file2, []byte(manifest2), 0644)

			w, err := r.Worktree()
			Expect(err).To(BeNil())

			_, err = w.Add("folder1/test.manifest.yaml")
			Expect(err).To(BeNil())

			_, err = w.Add("folder2/test1.manifest.yml")
			Expect(err).To(BeNil())

			_, err = w.Commit("first commit", &git.CommitOptions{
				Author: &object.Signature{
					Name:  "John Doe",
					Email: "j@doe.org",
					When:  time.Now(),
				},
			})

			cloneDir, err = os.MkdirTemp("", "git-clone-*")
			Expect(err).To(BeNil())
		})

		It("successfully finds all the manifest files in the repo", func() {
			repo := entity.Repository{
				Id:       "test",
				Url:      tmpDir,
				AuthType: entity.NoRepositoryAuthType,
			}

			r := gitRepo.New(cloneDir)
			clone, err := r.Clone(context.TODO(), repo)
			Expect(err).To(BeNil())
			Expect(clone.Branch).To(Equal("master"))

			manifests, err := r.GetManifests(context.TODO(), clone)
			Expect(len(manifests)).To(Equal(2))

			Expect(manifests[0].GetName()).To(Equal("manifest1"))
			Expect(manifests[1].GetName()).To(Equal("manifest2"))
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
			os.RemoveAll(cloneDir)
		})
	})
})
