package genesis

import "testing"

func TestMainnetFoundationalNodeAccounts(t *testing.T) {
	for name, accounts := range map[string][]DeployAccount{
		"V0": FoundationalNodeAccounts,
	} {
		t.Run(name, func(t *testing.T) { testDeployAccounts(t, accounts) })
	}
}
