package e2e

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
)

var godogOpts godog.Options

func TestMain(m *testing.M) {
	// Default godog Options
	godogOpts = godog.Options{
		Concurrency: 1,
		Format:      "pretty",
		Output:      colors.Colored(os.Stdout),
		Randomize:   time.Now().UTC().UnixNano(),
	}
	godog.BindFlags("godog.", flag.CommandLine, &godogOpts)
	framework.MainEntry(m)
}
