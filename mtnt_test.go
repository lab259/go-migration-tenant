package mtnt_test

import (
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/rlog"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestMigrationTenant(t *testing.T) {
	rlog.SetOutput(ginkgo.GinkgoWriter)
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "Migration Tenant Test Suite")
}
