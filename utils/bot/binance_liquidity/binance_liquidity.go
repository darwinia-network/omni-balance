package binance_liquidity

import (
	"context"
	"encoding/json"
	"omni-balance/utils/bot"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"strings"

	log "omni-balance/utils/logging"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	binance_connector "github.com/binance/binance-connector-go"

	"github.com/shopspring/decimal"
)

func init() {
	bot.Register(BinanceLiquidity{}.Name(), BinanceLiquidity{})
}

type BinanceLiquidity struct {
	Bot bot.Bot
}

type Binance struct {
	Key    string `json:"key" yaml:"key" help:"API key"`
	Secret string `json:"secret" yaml:"secret" help:"API secret"`
}

func (g BinanceLiquidity) Name() string {
	return "binance"
}

func (b BinanceLiquidity) GetClient(config configs.Config) (*binance_connector.Client, error) {
	var c Binance
	if err := config.GetProvidersConfig(b.Name(), configs.CEX, &c); err != nil {
		return nil, err
	}
	client := binance_connector.NewClient(c.Key, c.Secret)
	return client, nil
}

func (b BinanceLiquidity) Balance(ctx context.Context, args bot.Params) (decimal.Decimal, error) {
	client, err := b.GetClient(args.Conf)
	if err != nil {
		return decimal.Zero, err
	}
	balances, err := client.NewUserAssetService().Do(ctx)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "get binance balance")
	}
	for _, v := range balances {
		if !strings.EqualFold(v.Asset, args.Info.TokenName) {
			continue
		}
		return decimal.NewFromString(v.Free)
	}
	return decimal.Zero, nil
}

func (b BinanceLiquidity) Check(ctx context.Context, args bot.Params) ([]bot.Task, bot.ProcessType, error) {
	botConfigRaw := args.Conf.GetBotConfigUnderWallet(args.Info.Wallet.GetAddress().Hex(), b.Name())
	var (
		botConfig = struct {
			// tokenName: chainName
			OnlyWayTokens map[string]string `json:"onlyOneWayTokens" yaml:"onlyOneWayTokens" help:"Only way tokens"`
			// sourceToken: targetToken
			ChangeTokens map[string]string `json:"changeTokens" yaml:"changeTokens" help:"Change tokens"`
		}{}
	)
	if botConfigRaw != nil {
		data, _ := json.Marshal(botConfigRaw)
		_ = json.Unmarshal(data, &botConfig)
	}

	if len(botConfig.OnlyWayTokens) != 0 && botConfig.OnlyWayTokens[args.Info.TokenName] != args.Info.Chain && botConfig.OnlyWayTokens[args.Info.TokenName] != "" {
		log.Infof("token %s only way to %s, not %s", args.Info.TokenName, botConfig.OnlyWayTokens[args.Info.TokenName], args.Info.Chain)
		return nil, bot.Queue, nil
	}
	var (
		balance          decimal.Decimal
		tokenConfig      = args.Conf.GetTokenInfoOnChain(args.Info.TokenName, args.Info.Chain)
		isTarget2Binance = strings.Contains(args.Info.Chain, "binance")
	)
	if tokenConfig.Name == "" {
		return nil, bot.Queue, errors.Errorf("%s tokne not found in config", args.Info.TokenName)
	}
	if isTarget2Binance {
		binanceBalance, err := b.Balance(ctx, args)
		if err != nil {
			return nil, bot.Queue, err
		}
		balance = binanceBalance
	} else {
		chainBalance, err := args.Info.Wallet.GetExternalBalance(ctx, common.HexToAddress(tokenConfig.ContractAddress), tokenConfig.Decimals, args.Client)
		if err != nil {
			return nil, bot.Queue, err
		}
		balance = chainBalance
	}

	threshold := args.Conf.GetTokenThreshold(args.Info.Wallet.GetAddress().Hex(), args.Info.TokenName, args.Info.Chain)
	if balance.GreaterThan(threshold) {
		log.Infof("token %s balance %s, threshold %s, skip rebalance.", args.Info.TokenName, balance.String(), threshold.String())
		return nil, bot.Queue, nil
	}

	var (
		config = args.Conf
		info   = args.Info
		tasks  []bot.Task
	)
	amount := config.GetTokenPurchaseAmount(info.Wallet.GetAddress().Hex(), info.TokenName, info.Chain)

	client, err := b.GetClient(args.Conf)
	if err != nil {
		return nil, bot.Queue, err
	}
	coins, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		return nil, bot.Queue, errors.Wrap(err, "get all coins info")
	}
	var coinInfo *binance_connector.CoinInfo
	for index, v := range coins {
		if v.Coin != constant.GetBinanceTokenName(args.Info.TokenName) {
			continue
		}
		if !v.DepositAllEnable {
			continue
		}
		if !v.WithdrawAllEnable {
			continue
		}
		coinInfo = coins[index]
	}
	if coinInfo == nil {
		return nil, bot.Queue, errors.Errorf("token %s not found in binance", args.Info.TokenName)
	}
	var (
		chainName      = args.Info.Chain
		minWithdrawFee decimal.Decimal
	)
	if !isTarget2Binance {
		for _, v := range coinInfo.NetworkList {
			if !v.WithdrawEnable || !v.DepositEnable {
				continue
			}
			newChainName := constant.GetBinanceChainName(v.Network)
			if !args.Conf.ChainExists(newChainName) {
				continue
			}
			if strings.EqualFold(newChainName, args.Info.Chain) {
				chainName = newChainName
				break
			}
			withdrawFee := decimal.RequireFromString(v.WithdrawFee)
			if withdrawFee.LessThanOrEqual(decimal.Zero) || minWithdrawFee.Equal(decimal.Zero) {
				minWithdrawFee = withdrawFee
				chainName = newChainName
				continue
			}
		}
	}

	if chainName == "" {
		return nil, bot.Queue, errors.Errorf("token %s not found in binance", args.Info.TokenName)
	}
	tokenInChainConfig := args.Conf.GetChainConfig(chainName)
	if tokenInChainConfig.Name == "" || len(tokenInChainConfig.RpcEndpoints) == 0 {
		return nil, bot.Queue, errors.Errorf("token %s chain %s config not found", args.Info.TokenName, chainName)
	}
	var tokenInName = args.Info.TokenName
	if len(botConfig.ChangeTokens) != 0 && botConfig.ChangeTokens[args.Info.TokenName] != "" {
		tokenInName = botConfig.ChangeTokens[args.Info.TokenName]
	}

	if balance.LessThanOrEqual(threshold) && isTarget2Binance {
		if balance.Add(amount).LessThanOrEqual(threshold) {
			newAmount := threshold.Add(threshold.Mul(decimal.RequireFromString("0.01")))
			log.Infof("The binance current balance is %s, amount in config is %s, balance(%s) + amount(%s) <= threshold(%s), so set amount to %s", balance, amount, balance, amount, threshold, newAmount)
			amount = newAmount
		}
		log.Infof("binance balance is less than threshold, start rebalance, current balance %s, threshold %s tokenInName %s", balance.String(), threshold.String(), tokenInName)
		tasks = append(tasks, bot.Task{
			TokenOutChainName: args.Info.Chain,
			Wallet:            info.Wallet.GetAddress().Hex(),
			TokenOutName:      info.TokenName,
			TokenInName:       tokenInName,
			Amount:            amount,
			ProviderType:      configs.CEX,
			ProviderName:      b.Name(),
			Remark:            "chain_to_binance",
			CurrentBalance:    balance,
		})
	}

	if balance.LessThanOrEqual(threshold) && !isTarget2Binance {
		if balance.Add(amount).LessThanOrEqual(threshold) {
			newAmount := threshold.Add(threshold.Mul(decimal.RequireFromString("0.01")))
			log.Infof("The chain current balance is %s, amount in config is %s, balance(%s) + amount(%s) <= threshold(%s), so set amount to %s", balance, amount, balance, amount, threshold, newAmount)
			amount = newAmount
		}
		log.Infof("%s chain balance is less than threshold, start rebalance, current balance %s, threshold %s", args.Info.Chain, balance.String(), threshold.String())
		tasks = append(tasks, bot.Task{
			TokenOutChainName: args.Info.Chain,
			TokenInChainName:  "binance",
			Wallet:            info.Wallet.GetAddress().Hex(),
			TokenOutName:      info.TokenName,
			TokenInName:       tokenInName,
			Amount:            amount,
			ProviderType:      configs.CEX,
			ProviderName:      b.Name(),
			Remark:            "binance_to_chain",
			CurrentBalance:    balance,
		})
	}
	return tasks, bot.Queue, nil
}
