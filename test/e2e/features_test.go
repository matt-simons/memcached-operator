package e2e

import (
	"testing"
	"time"

	"github.com/cucumber/godog"
	openshiftapi "github.com/openshift/api"
	"github.com/operator-framework/operator-sdk-samples/go/memcached-operator/pkg/apis"
	cachev1alpha1 "github.com/operator-framework/operator-sdk-samples/go/memcached-operator/pkg/apis/cache/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

var ctx *framework.TestCtx
var f *framework.Framework

func TestLifecycleFeatures(t *testing.T) {
	f = framework.Global
	ctx = framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: time.Minute, RetryInterval: time.Second})
	if err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	e2eutil.WaitForOperatorDeployment(t, f.KubeClient, f.Namespace, "memcached-operator", 1, 15*time.Second, time.Minute)
	time.Sleep(25 * time.Second)

	// CRITICAL
	// This is required so that we can serialize our CRDs
	MemcachedList := &cachev1alpha1.MemcachedList{}
	err = framework.AddToFrameworkScheme(apis.AddToScheme, MemcachedList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	openshiftapi.Install(f.Scheme)

	godogOpts.Paths = []string{
		"features",
	}
	godog.RunWithOptions("features", func(s *godog.Suite) {
		Scenarios(t, s)
	}, godogOpts)

	// Show operator logs
	//cmd := exec.Command("kubectl", "logs", "-lname=memcached-operator", "--tail=-1", "-n", f.Namespace)
	//stdoutStderr, err := cmd.Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Operator Log Output:\n\n%s\n", stdoutStderr)
}
