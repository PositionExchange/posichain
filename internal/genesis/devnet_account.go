package genesis

// DevnetOperatedAccounts are the accounts for the initial genesis nodes hosted by Posichain for devnet.
var DevnetOperatedAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0x58AcF11BF43F2dDeE961C079ade3BAC9da83Dd0f", BLSPublicKey: "bad4b445ae3dbb130fc13f42d7f1065e7208fd05444a68a615768f4bedf6beb4ae2f0c4768810fb556cf2ecff49e330a"},
	{Index: "1", ShardID: 1, Address: "0x72948DE451561eb58e36d24E45FF64e0D97Ae7B0", BLSPublicKey: "31ad8c6c7effad3af47ba8f31b18ae1b0168864e88b47b100cb69e1eae040976b356b1ed5f13771ef87c13a84198798d"},
	{Index: "2", ShardID: 0, Address: "0xE247fd7041a714B047C40F0b06d1c83B9458Be11", BLSPublicKey: "de3bfaa629c0d7d1a167b930f040b4ca35a3ba72efe0dc8330ecf3bfceb1815d989871a13e41c04008715f97d62bb106"},
	{Index: "3", ShardID: 1, Address: "0xa3A6238a7139eFd6a94E678f4bb0706Bd513163E", BLSPublicKey: "6e2a288108bfb8c4933c053b8f8a5e399301d9bdd455e2cca3916eaa56f5ec5401b2d8b0b9aa370a46720512bcb24e09"},
}

// DevnetFoundationalAccounts are the accounts for the initial foundational nodes for devnet.
var DevnetFoundationalAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0x1BAAC0973CE99D8048936D5b102Dd55B62047ea0", BLSPublicKey: "3c10d84c6b7f3aa0204ed00fdc9b9406f612ed46c8356054087b8bb3fda63cb6840f4bacaaca118f1166ddcbb3c54e06"},
	{Index: "1", ShardID: 1, Address: "0x22bd5bbf88f7735DEF9378FD50D952BACD8c2339", BLSPublicKey: "210db574c3bcaf52da54fe0a09678502773d837d22d0200631fe32e860e2f4c7faa416ef7b4469e8c57da0e4a0edaa8b"},
	{Index: "2", ShardID: 0, Address: "0x237bE5DB39613077e04402310fa6f844acc04F25", BLSPublicKey: "deb404802a80b002557b952a499ffc48b9cddaffca3522357aee2cedcbab3202dcd35109948a11b9296c592d6f641508"},
	{Index: "3", ShardID: 1, Address: "0x8f943B9AF1c2189b734c317d83d791A479d65fCc", BLSPublicKey: "2db4e71fe73efd575b4279c7e0a9cf6b02ec2c3f0cdd9d1d645333c1e0191d26f93fd8a0aa38e968cfa78ca377cd5d8d"},
}
