package genesis

// TestnetOperatedAccounts are the accounts for the initial genesis nodes hosted by Posichain for testnet.
var TestnetOperatedAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0x014223E0Cff9E54D21bDC9F1c9A9B04748b22763", BLSPublicKey: "22dc9781ccc8eb929a50d2a297be74e52b8431ece50fa5769c7a2966db94908253c2b4068659b884e645296074e9a780"},
	{Index: "1", ShardID: 0, Address: "0x795e816cbA12B9729Bf908Fc38cB05fFe60328d0", BLSPublicKey: "857c28c2dff132884e2aabe4e56b9e9e112857dfa5fc2e0f06a156720b281d5a2a70157d1c4e4fcdd13e1dd3e8c0d298"},
	{Index: "2", ShardID: 0, Address: "0xC30A48A6347D4290170996e93A6538866aE1B9E2", BLSPublicKey: "e89959cb16a9dae9ca7ea8a357796d29339435153a145b529d98a68015b1121454813e5f1957338838fdcb2e3db5cd06"},
	{Index: "3", ShardID: 0, Address: "0x10A4464216Ac2350A86B6Aa8649F60af003AdcF5", BLSPublicKey: "827d3a1024e34acc9d1f31fdd6cf64ea52197b157377873ae486893554282cfdcf3b1e1520eb7cb6b640f5fdbbf88280"},
}

// TestnetFoundationalAccounts are the accounts for the initial foundational nodes for testnet.
var TestnetFoundationalAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0x3eb46931E23a6949F8aAEA1d7C36B3b4D5e54da2", BLSPublicKey: "b656eccda0ad6be57e8880221bab2179ac9bd00a142978cd7f1e2a0793afa8f95b4597d96a7f50d7d9a16a723637900c"},
}
