package models_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)

	junitReporter := reporters.NewJUnitReporter("reports/model.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Model Suite", []Reporter{junitReporter})
}
