package genesis

import "testing"

func TestDevnetOperatedAccounts(t *testing.T) {
	testDeployAccounts(t, DevnetOperatedAccounts)
}

func TestDevnetFoundationalAccounts(t *testing.T) {
	testDeployAccounts(t, DevnetFoundationalAccounts)
}
