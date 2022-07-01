package genesis

// DevnetOperatedAccounts are the accounts for the initial genesis nodes hosted by Posichain for devnet.
var DevnetOperatedAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0x7F868BD9C67e3C83C6FEAB22ba029f69ed3E15d4", BLSPublicKey: "f47b6d2b91eb37a5c0b35520803b86901a086367506e423f1651464191cc95b67b95470b8bc41ea4ad57ead0739fc180"},
	{Index: "1", ShardID: 1, Address: "0xDC9ed4a0FBaF250E080fa02B850D66763AdF7e37", BLSPublicKey: "8fe989abe5334b53f85419d3e41225d615e60a5b536d280415e125df4f7f4136040d434c20b67a2848c2c24e5b06f38f"},
	{Index: "2", ShardID: 0, Address: "0xA616522dF16bc7aBb38C6CB99f2dA2919C48665C", BLSPublicKey: "25b597c52cfc071fd00b626844cb45c6e8a0252bd73841be21dbcb9623e136a3644d4a218dc5a6f76f7052e3569f0308"},
	{Index: "3", ShardID: 1, Address: "0x8DCAD03eadf7039d44BE8b669D70Ce6c16Ef6f66", BLSPublicKey: "77c865f73d76fc4554b63f0a0d0e34a6a1272d02578044a6e14f2e55bbc6dde1ef7ac0046f81c21752f91bb823e8e819"},
}

// DevnetFoundationalAccounts are the accounts for the initial foundational nodes for devnet.
var DevnetFoundationalAccounts = []DeployAccount{
	{Index: "0", ShardID: 0, Address: "0xC90c19Bb0498070136aB0dE7866592CDeE375Ceb", BLSPublicKey: "8274b0ad76ad77a1ebf19365cfe64564781ad988d276c1bd056f4b836e90dbd9bc30303432f2109721ede934ffbafb14"},
	{Index: "1", ShardID: 1, Address: "0xFcF6664d04b339650B5acC7fb140758788899f75", BLSPublicKey: "8cb7325a2f2d9532a657a3404de9c5c5eaa2c45e1b1e4e72a2fe2051d48ec3525d7146cbf79b3f72f9faccffcc0da105"},
	{Index: "2", ShardID: 0, Address: "0x4655b4ce5842c20127c5d0EA36FA1d3cb7786d13", BLSPublicKey: "683dee10231b0b40323373d903b799c2eaadd7b2bba5ce09b3bcddf8153bdd93edc40b8b75dbed762bfae4f16fdd4c8c"},
	{Index: "3", ShardID: 1, Address: "0xb2eADEe811387055aFe38A61BeFe5a76ad612785", BLSPublicKey: "5f42f1ac23adf59dc6d535d0f422967cf79f5b95f2489de7a127c20e44b9d659e60eed2b19a7bf6fcb229b4f1018a60f"},
}
