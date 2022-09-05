package nodeconfig

var (
	mainnetBootNodes = []string{
		"/dnsaddr/bootstrap.posichain.org",
	}

	testnetBootNodes = []string{
		"/dnsaddr/bootstrap.t.posichain.org",
	}

	devnetBootNodes = []string{
		"/dnsaddr/bootstrap.d.posichain.org",
	}

	dockernetBootNodes = []string{
		"/dnsaddr/bootstrap.k.posichain.org",
	}

	stressBootNodes = []string{
		"/dnsaddr/bootstrap.s.posichain.org",
	}
)

const (
	mainnetDNSZone   = "posichain.org"
	testnetDNSZone   = "t.posichain.org"
	devnetDNSZone    = "d.posichain.org"
	dockernetDNSZone = "k.posichain.org"
	stressnetDNSZone = "s.posichain.org"
)

const (
	// DefaultLocalListenIP is the IP used for local hosting
	DefaultLocalListenIP = "127.0.0.1"
	// DefaultPublicListenIP is the IP used for public hosting
	DefaultPublicListenIP = "0.0.0.0"
	// DefaultP2PPort is the key to be used for p2p communication
	DefaultP2PPort = 9000
	// DefaultLegacyDNSPort is the default legacy DNS port. The actual port used is DNSPort - 3000. This is a
	// very bad design. Refactored to DefaultDNSPort
	DefaultLegacyDNSPort = 9000
	// DefaultDNSPort is the default DNS port for both remote node and local server.
	DefaultDNSPort = 6000
	// DefaultRPCPort is the default rpc port. The actual port used is 9000+500
	DefaultRPCPort = 9500
	// DefaultAuthRPCPort is the default rpc auth port. The actual port used is 9000+501
	DefaultAuthRPCPort = 9501
	// DefaultRosettaPort is the default rosetta port. The actual port used is 9000+700
	DefaultRosettaPort = 9700
	// DefaultWSPort is the default port for web socket endpoint. The actual port used is
	DefaultWSPort = 9800
	// DefaultAuthWSPort is the default port for web socket auth endpoint. The actual port used is
	DefaultAuthWSPort = 9801
	// DefaultPrometheusPort is the default prometheus port. The actual port used is 9000+900
	DefaultPrometheusPort = 9900
	// DefaultP2PConcurrency is the default P2P concurrency, 0 means is set the default value of P2P Discovery, the actual value is 10
	DefaultP2PConcurrency = 0
	// DefaultMaxConnPerIP is the maximum number of connections to/from a remote IP
	DefaultMaxConnPerIP = 10
	// DefaultMaxPeers is the maximum number of remote peers, with 0 representing no limit
	DefaultMaxPeers = 0
)

const (
	// DefaultRateLimit for RPC, the number of requests per second
	DefaultRPCRateLimit = 1000
)

const (
	// rpcHTTPPortOffset is the port offset for RPC HTTP requests
	rpcHTTPPortOffset = 500

	// rpcHTTPAuthPortOffset is the port offset for RPC Auth HTTP requests
	rpcHTTPAuthPortOffset = 501

	// rpcHTTPPortOffset is the port offset for rosetta HTTP requests
	rosettaHTTPPortOffset = 700

	// rpcWSPortOffSet is the port offset for RPC websocket requests
	rpcWSPortOffSet = 800

	// rpcWSAuthPortOffSet is the port offset for RPC Auth websocket requests
	rpcWSAuthPortOffSet = 801

	// prometheusHTTPPortOffset is the port offset for prometheus HTTP requests
	prometheusHTTPPortOffset = 900
)

// GetDefaultBootNodes get the default bootnode with the given network type
func GetDefaultBootNodes(networkType NetworkType) []string {
	switch networkType {
	case Mainnet:
		return mainnetBootNodes
	case Testnet:
		return testnetBootNodes
	case Devnet:
		return devnetBootNodes
	case Stressnet:
		return stressBootNodes
	case Dockernet:
		return dockernetBootNodes
	case Localnet:
		return nil
	}
	return nil
}

// GetDefaultDNSZone get the default DNS zone with the given network type
func GetDefaultDNSZone(networkType NetworkType) string {
	switch networkType {
	case Mainnet:
		return mainnetDNSZone
	case Testnet:
		return testnetDNSZone
	case Devnet:
		return devnetDNSZone
	case Stressnet:
		return stressnetDNSZone
	case Dockernet:
		return dockernetDNSZone
	case Localnet:
		return ""
	}
	return ""
}

// GetDefaultDNSPort get the default DNS port for the given network type
func GetDefaultDNSPort(NetworkType) int {
	return DefaultDNSPort
}

// GetRPCHTTPPortFromBase return the rpc HTTP port from base port
func GetRPCHTTPPortFromBase(basePort int) int {
	return basePort + rpcHTTPPortOffset
}

// GetRPCAuthHTTPPortFromBase return the rpc HTTP port from base port
func GetRPCAuthHTTPPortFromBase(basePort int) int {
	return basePort + rpcHTTPAuthPortOffset
}

// GetRosettaHTTPPortFromBase return the rosetta HTTP port from base port
func GetRosettaHTTPPortFromBase(basePort int) int {
	return basePort + rosettaHTTPPortOffset
}

// GetWSPortFromBase return the Websocket port from the base port
func GetWSPortFromBase(basePort int) int {
	return basePort + rpcWSPortOffSet
}

// GetWSAuthPortFromBase return the Websocket port from the base auth port
func GetWSAuthPortFromBase(basePort int) int {
	return basePort + rpcWSAuthPortOffSet
}

// GetPrometheusHTTPPortFromBase return the prometheus HTTP port from base port
func GetPrometheusHTTPPortFromBase(basePort int) int {
	return basePort + prometheusHTTPPortOffset
}
