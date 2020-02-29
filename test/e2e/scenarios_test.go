package e2e

import (
	"testing"

	"github.com/cucumber/godog"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
)

var nsName string
var scenarioCtx *framework.TestCtx

func Scenarios(t *testing.T, s *godog.Suite) {
	RegisterSteps(s)

	s.BeforeScenario(func(interface{}) {
		scenarioCtx = framework.NewTestCtx(t)
		createUniqueNamespace(scenarioCtx)
	})
	s.AfterScenario(func(interface{}, error) {
		scenarioCtx.Cleanup()
	})
}
