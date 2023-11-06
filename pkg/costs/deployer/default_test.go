package deployer

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	k8sdeployer "github.com/seal-io/walrus/pkg/k8s/deploy"
)

func TestHelm(t *testing.T) {
	kubeconfigPath := os.Getenv("TEST_KUBECONFIG")
	if kubeconfigPath == "" {
		t.Skip("environment TEST_KUBECONFIG isn't provided")
		return
	}

	kubeConfigContentByte, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return
	}

	ctx := context.Background()
	deployer, err := k8sdeployer.New(string(kubeConfigContentByte))
	assert.Nil(t, err, "error create helm")

	yaml, err := opencost("test", "docker.io")
	assert.Nil(t, err, "error create opencost yaml")

	app, err := prometheus("docker.io")
	assert.Nil(t, err, "error create prometheus app")

	err = deployer.EnsureChart(app, true)
	assert.Nil(t, err, "error ensure prometheus chart")

	err = deployer.EnsureYaml(ctx, yaml)
	assert.Nil(t, err, "error ensure opencost yaml")
}
