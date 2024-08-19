package constant

import "golang.org/x/exp/constraints"

const (
	Arbitrum        = "arbitrum"
	ArbitrumSepolia = "arbitrum-sepolia"
	Astar           = "astar"
	AstarZkevm      = "astar-zkevm"
	Aurora          = "aurora"
	Avalanche       = "avalanche"
	Base            = "base"
	Blast           = "blast"
	Bnb             = "bnb"
	Bsc             = "bsc"
	Celo            = "celo"
	CrabDvm         = "crab-dvm"
	Cronos          = "cronos"
	DarwiniaDvm     = "darwinia-dvm"
	Ethereum        = "ethereum"
	Fantom          = "fantom"
	Gnosis          = "gnosis"
	Kava            = "kava"
	Linea           = "linea"
	Mantle          = "mantle"
	Merlin          = "merlin"
	Mode            = "mode"
	Moonbeam        = "moonbeam"
	OpBNB           = "opBNB"
	Optimism        = "op"
	Polygon         = "polygon"
	PolygonZkEvm    = "polygon-zkEvm"
	PulseChain      = "pulseChain"
	Rootstock       = "rootstock"
	Scroll          = "scroll"
	Sepolia         = "sepolia"
	Zircuit         = "zircuit"
	ZkLinkNova      = "zkLink-nova"
	Zksync          = "zksync"
)

const (
	ChainNameKeyInCtx    ContextKey = "chain_name"
	LogKeyInCtx          ContextKey = "log"
	NoticeFieldsKeyInCtx ContextKey = "notice_fields"
	FeeTestKeyInCtx      ContextKey = "fee_test"
)

var (
	chainName2Id = map[string]int{
		Aurora:          1313161554,
		Arbitrum:        42161,
		ArbitrumSepolia: 421614,
		Astar:           592,
		AstarZkevm:      3776,
		Avalanche:       43114,
		Base:            8453,
		Blast:           81457,
		Bsc:             56,
		Celo:            42220,
		CrabDvm:         44,
		Cronos:          25,
		DarwiniaDvm:     46,
		Ethereum:        1,
		Fantom:          250,
		Gnosis:          100,
		Kava:            2222,
		Linea:           59144,
		Mantle:          5000,
		Merlin:          4200,
		Mode:            34443,
		Moonbeam:        1284,
		OpBNB:           204,
		Optimism:        10,
		Polygon:         137,
		PolygonZkEvm:    1101,
		PulseChain:      369,
		Rootstock:       30,
		Scroll:          534352,
		Sepolia:         11155111,
		Zircuit:         48900,
		ZkLinkNova:      810180,
		Zksync:          324,
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
