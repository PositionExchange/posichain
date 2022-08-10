package genesis

import "testing"

func TestDockerNetOperatedAccounts(t *testing.T) {
	testDeployAccounts(t, DockernetOperatedAccounts)
}

func TestDockerNetFoundationalAccounts(t *testing.T) {
	testDeployAccounts(t, DockernetFoundationalAccounts)
}
