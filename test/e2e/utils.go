package e2e

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cucumber/godog/gherkin"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func parseDocString(manifest *gherkin.DocString) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	if manifest.ContentType == "yaml" {
		err := yaml.Unmarshal([]byte(manifest.Content), u)
		if err != nil {
			return nil, err
		}
	} else if manifest.ContentType == "json" {
		err := json.Unmarshal([]byte(manifest.Content), u)
		if err != nil {
			return nil, err
		}
	} else if manifest.ContentType == "" {
		return nil, errors.New("No content-type specified. Supported types are json, yaml.")
	} else {
		return nil, errors.New(fmt.Sprintf("Unrecognised content-type %s. Supported types are json, yaml.", manifest.ContentType))
	}
	return u, nil
}

func createUniqueNamespace(ctx *framework.TestCtx) error {
	var err error
	nsName, err = uuid()
	if err != nil {
		return err
	}
	err = f.Client.Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName,
		},
		Spec: corev1.NamespaceSpec{},
	}, &framework.CleanupOptions{TestContext: ctx})
	if err != nil {
		return err
	}
	return nil
}

func uuid() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u), err
}
