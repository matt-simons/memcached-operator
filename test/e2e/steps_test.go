package e2e

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const cleanTimeout = time.Second * 30
const cleanRetryInterval = time.Second * 3

func RegisterSteps(s *godog.Suite) {
	s.Step(`^the operator reconciles$`, theOperatorReconciles)

	// Resource manipulation (creation, update, delete)
	s.Step(`^I create a Resource:$`, iCreateAResource)
	s.Step(`^I create an invalid Resource:$`, iCreateAnInvalidResource)
	s.Step(`^I update a Resource:$`, iUpdateAResource)
	s.Step(`^I delete a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)" in namespace "([^"]*)"$`,
		iDeleteANamespacedResource)
	s.Step(`^I delete a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)"$`,
		iDeleteAClusterResource)

	// Resource assertions
	s.Step(`^there should exist a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)" in namespace "([^"]*)"$`,
		thereShouldExistANamespacedResource)
	s.Step(`^there should exist a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)"$`,
		thereShouldExistAClusterResource)
	s.Step(`^there should not exist a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)" in namespace "([^"]*)"$`,
		thereShouldNotExistANamespacedResource)
	s.Step(`^there should not exist a "([^"]*)" "([^"]*)" "([^"]*)" called "([^"]*)"$`,
		thereShouldNotExistAClusterResource)

	// ConfigMap resource assertions
	s.Step(`^there should exist a data entry called "([^"]*)" in ConfigMap "([^"]*)" in namespace "([^"]*)"$`,
		thereShouldExistADataEntryCalledInConfigMapInNamespace)
	s.Step(`^there should not exist a data entry called "([^"]*)" in ConfigMap "([^"]*)" in namespace "([^"]*)"$`,
		thereShouldNotExistADataEntryCalledInConfigMapInNamespace)

	// DaemonSet resource assertions
	s.Step(`^the DaemonSet "([^"]*)" in namespace "([^"]*)" should have rolled out a new version$`, theDaemonSetShouldHaveRolledOutANewVersion)

	// Deployment resource assertions
	s.Step(`^there should exist (\d+) ready pods for Deployment called "([^"]*)" in namespace "([^"]*)"$`, thereShouldExistReadyPodsForDeploymentCalled)
}

func theOperatorReconciles() error {
	time.Sleep(15 * time.Second)
	return nil
}

func getNamespace(s string) string {
	return strings.ReplaceAll(s, "NAMESPACE", nsName)
}

/*******************************************************************************
 * Resource creation, update, delete, list
 ******************************************************************************/

func iCreateAResource(manifest *gherkin.DocString) error {
	u, err := parseDocString(manifest)
	if err != nil {
		return nil
	}
	u.SetNamespace(getNamespace(u.GetNamespace()))
	return f.Client.Create(context.TODO(), u, &framework.CleanupOptions{TestContext: scenarioCtx, Timeout: cleanTimeout, RetryInterval: cleanRetryInterval})
}

func iCreateAnInvalidResource(manifest *gherkin.DocString) error {
	u, err := parseDocString(manifest)
	if err != nil {
		return nil
	}
	u.SetNamespace(nsName)
	err = f.Client.Create(context.TODO(), u, &framework.CleanupOptions{TestContext: scenarioCtx, Timeout: cleanTimeout, RetryInterval: cleanRetryInterval})
	if err == nil {
		return errors.New(fmt.Sprintf("Invalid resource %s called %s was created", u.GetKind(), u.GetName()))
	}
	return nil
}

func iUpdateAResource(manifest *gherkin.DocString) error {
	u, err := parseDocString(manifest)
	if err != nil {
		return nil
	}
	u.SetNamespace(nsName)
	return f.Client.Update(context.TODO(), u)
}

func iDeleteAClusterResource(version, group, kind, name string) error {
	return iDeleteANamespacedResource(version, group, kind, name, "")
}

func iDeleteANamespacedResource(version, group, kind, name, namespace string) error {
	namespace = getNamespace(namespace)
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: version,
	})
	u.SetNamespace(namespace)
	u.SetName(name)
	return f.Client.Delete(context.TODO(), u)
}

func iDeleteAllNamespacedResources(version, group, kind string) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: version,
	})
	return f.Client.DeleteAllOf(context.TODO(), u)
}

/*******************************************************************************
 * General resource assertions
 ******************************************************************************/

func thereShouldExistAClusterResource(version, group, kind, name string) error {
	return thereShouldExistANamespacedResource(version, group, kind, name, "")
}

func thereShouldExistANamespacedResource(version, group, kind, name, namespace string) error {
	namespace = getNamespace(namespace)
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: version,
	})
	return f.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, u)
}

func thereShouldNotExistAClusterResource(version, group, kind, name string) error {
	err := thereShouldNotExistANamespacedResource(version, group, kind, name, "")
	if err != nil {
		return errors.New(fmt.Sprintf("Found a %s resource called %s.", kind, name))
	}
	return nil
}

func thereShouldNotExistANamespacedResource(version, group, kind, name, namespace string) error {
	namespace = getNamespace(namespace)
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: version,
	})
	err := f.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, u)
	if err == nil {
		return errors.New(fmt.Sprintf("Found a %s resource called %s in namespace %s.", kind, name, namespace))
	}
	return nil
}

/*******************************************************************************
 * ConfigMap resource assertions
 ******************************************************************************/

func thereShouldExistADataEntryCalledInConfigMapInNamespace(entry, name, namespace string) error {
	namespace = getNamespace(namespace)
	entry = getNamespace(entry)
	var cm corev1.ConfigMap
	err := f.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, &cm)
	if err != nil {
		return err
	}
	if cm.Data == nil {
		return errors.New(fmt.Sprintf("ConfigMap %s in namespace %s has no entries", name, namespace))
	}
	if _, ok := cm.Data[entry]; !ok {
		return errors.New(fmt.Sprintf("Entry %s does not exist in ConfigMap %s in namespace %s", entry, name, namespace))
	}
	return nil
}

func thereShouldNotExistADataEntryCalledInConfigMapInNamespace(entry, name, namespace string) error {
	err := thereShouldExistADataEntryCalledInConfigMapInNamespace(entry, name, namespace)
	if err == nil {
		return errors.New(fmt.Sprintf("Entry %s exists in ConfigMap %s in namespace %s", entry, name, namespace))
	}
	return nil
}

/*******************************************************************************
 * DaemonSet resource assertions
 ******************************************************************************/

func theDaemonSetShouldHaveRolledOutANewVersion(name, namespace string) error {
	var ds appsv1.DaemonSet
	err := f.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, &ds)
	if err != nil {
		return err
	}
	if ds.Status.ObservedGeneration > 1 {
		return nil
	}
	return errors.New("DaemonSet did not roll out a new version")
}

/*******************************************************************************
 * Deployment resource assertions
 ******************************************************************************/

func thereShouldExistReadyPodsForDeploymentCalled(replicas int, name, namespace string) error {
	var d appsv1.Deployment
	err := f.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, &d)
	if err != nil {
		return err
	}
	if d.Status.Replicas == int32(replicas) && d.Status.UnavailableReplicas == 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("Incorrect number of ready pods for Deployment. Expected %d replicas was %d with %d Unavailable.", replicas, d.Status.Replicas, d.Status.UnavailableReplicas))
}
