package genesis

import "testing"

func TestDevnetOperatedAccounts(t *testing.T) {
	testDeployAccounts(t, DevnetOperatedAccountsV0)
	testDeployAccounts(t, DevnetOperatedAccountsV1)
}

func TestDevnetFoundationalAccounts(t *testing.T) {
	testDeployAccounts(t, DevnetFoundationalAccounts)
}
