package controllers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	junitReporter := reporters.NewJUnitReporter("reports/controllers.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Controller Suite", []Reporter{junitReporter})
}
