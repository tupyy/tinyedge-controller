package postgres_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

func getEnvVar(name string, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}
