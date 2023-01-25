package edge_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEdge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Edge Suite")
}
