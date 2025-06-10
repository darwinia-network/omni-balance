package constant

var (
	binanceChainName = map[string]string{
		"ARBITRUM": Arbitrum,
		"ETH":      Ethereum,
	}
	binanceTokenName2ChainName = map[string]string{}
)

func GetBinanceChainName(chainName string) string {
	if name, ok := binanceChainName[chainName]; ok {
		return name
	}
	return chainName
}

func GetChainName2BinanceTokenName(binanceChainName string) string {
	var chainName2BinanceTokenName = map[string]string{}
	for k, v := range binanceTokenName2ChainName {
		chainName2BinanceTokenName[v] = k
	}
	if name, ok := chainName2BinanceTokenName[binanceChainName]; ok {
		return name
	}
	return binanceChainName
}

func GetBinanceTokenName(chainTokenName string) string {
	var chainTokenName2BinanceTokenName = map[string]string{}
	for k, v := range binanceTokenName2ChainName {
		chainTokenName2BinanceTokenName[v] = k
	}
	if name, ok := chainTokenName2BinanceTokenName[chainTokenName]; ok {
		return name
	}
	return chainTokenName
}
