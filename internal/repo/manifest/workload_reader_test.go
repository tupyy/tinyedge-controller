package manifest

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

var (
	manifest = `
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

func TestManifestReader(t *testing.T) {
	RegisterTestingT(t)

	content := bytes.NewBufferString(manifest).Bytes()
	manifest, err := parseWorkloadManifest(content)
	Expect(err).To(BeNil())
	Expect(manifest).ToNot(BeNil())
	w, ok := manifest.(entity.Workload)
	Expect(ok).To(BeTrue())
	Expect(w).NotTo(BeNil())

	Expect(len(w.Resources)).To(Equal(3))
	Expect(w.Resources[0]).To(Equal("folder/dep/configmap.yaml"))
	Expect(len(w.Selectors)).To(Equal(5))
	Expect(w.GetKind()).To(Equal(entity.WorkloadManifestKind))
}
