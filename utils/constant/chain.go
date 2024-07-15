package constant

import "golang.org/x/exp/constraints"

const (
	Zksync          = "zksync"
	Sepolia         = "sepolia"
	Polygon         = "polygon"
	PolygonZkEvm    = "polygon-zkEvm"
	Bsc             = "bsc"
	Blast           = "blast"
	DarwiniaDvm     = "darwinia-dvm"
	Base            = "base"
	Arbitrum        = "arbitrum"
	Gnosis          = "gnosis"
	Scroll          = "scroll"
	Optimism        = "op"
	AstarZkevm      = "astar-zkevm"
	Mantle          = "mantle"
	Linea           = "linea"
	CrabDvm         = "crab-dvm"
	Ethereum        = "ethereum"
	Merlin          = "merlin"
	Avalanche       = "avalanche"
	Mode            = "mode"
	Cronos          = "cronos"
	PulseChain      = "pulseChain"
	Kava            = "kava"
	ZkLinkNova      = "zkLink-nova"
	Rootstock       = "rootstock"
	Astar           = "astar"
	OpBNB           = "opBNB"
	Bnb             = "bnb"
	Celo            = "celo"
	Moonbeam        = "moonbeam"
	ArbitrumSepolia = "arbitrum-sepolia"
)

const (
	ChainNameKeyInCtx    ContextKey = "chain_name"
	LogKeyInCtx          ContextKey = "log"
	NoticeFieldsKeyInCtx ContextKey = "notice_fields"
	FeeTestKeyInCtx      ContextKey = "fee_test"
)

var (
	chainName2Id = map[string]int{
		ArbitrumSepolia: 421614,
		Moonbeam:        1284,
		Sepolia:         11155111,
		Zksync:          324,
		Polygon:         137,
		PolygonZkEvm:    1101,
		Bsc:             56,
		Blast:           81457,
		DarwiniaDvm:     46,
		Base:            8453,
		Arbitrum:        42161,
		Gnosis:          100,
		Scroll:          534352,
		Optimism:        10,
		AstarZkevm:      3776,
		Mantle:          5000,
		Linea:           59144,
		CrabDvm:         44,
		Ethereum:        1,
		Merlin:          4200,
		Avalanche:       43114,
		Mode:            34443,
		Cronos:          25,
		PulseChain:      369,
		Kava:            2222,
		ZkLinkNova:      810180,
		Rootstock:       30,
		Astar:           592,
		OpBNB:           204,
		Bnb:             56,
		Celo:            42220,
	}
	chainId2Name = make(map[int]string)
)

func init() {
	if len(chainId2Name) != 0 {
		return
	}
	for k, v := range chainName2Id {
		chainId2Name[v] = k
	}
}

func GetChainName[t constraints.Integer](chainId t) string {
	return chainId2Name[int(chainId)]
}

func GetChainId(chainName string) int {
	return chainName2Id[chainName]
}
