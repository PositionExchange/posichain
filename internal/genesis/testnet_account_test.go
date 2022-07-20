package genesis

import "testing"

func TestTestnetOperatedAccounts(t *testing.T) {
	testDeployAccounts(t, TestnetOperatedAccounts)
}

func TestTestnetFoundationalAccounts(t *testing.T) {
	testDeployAccounts(t, TestnetFoundationalAccounts)
}
