package handler

import (
	"encoding/json"
	"net/http"
	"omni-balance/internal/daemons/market"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils/bot"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/provider"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func APIKey(conf configs.Config, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != conf.ApiKey || conf.ApiKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

type GateLiquidityParams struct {
	TokenName string          `json:"token_name"`
	FromChain string          `json:"from_chain"`
	ToChain   string          `json:"to_chain"`
	Amount    decimal.Decimal `json:"amount"`
	Address   string          `json:"address"`
}

func GateLiquidity(conf configs.Config) func(w http.ResponseWriter, r *http.Request) {
	return APIKey(conf, func(w http.ResponseWriter, r *http.Request) {
		args := GateLiquidityParams{}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if args.TokenName == "" || args.FromChain == "" || args.ToChain == "" || args.Amount.LessThanOrEqual(decimal.Zero) || args.Address == "" {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}
		if constant.GetChainId(args.FromChain) == 0 || constant.GetChainId(args.ToChain) == 0 {
			http.Error(w, "invalid chain", http.StatusBadRequest)
			return
		}
		if conf.GetTokenInfoOnChain(args.TokenName, args.FromChain).Name == "" {
			http.Error(w, "token not support,make sure you have set the token address in the config file", http.StatusBadRequest)
			return
		}
		if conf.GetTokenInfoOnChain(args.TokenName, args.ToChain).Name == "" {
			http.Error(w, "token not support, make sure you have set the token address in the config file", http.StatusBadRequest)
			return
		}
		wallet := conf.GetWalletConfig(common.HexToAddress(args.Address).Hex())
		if wallet.Address == "" {
			http.Error(w, "Address not support, make sure you have set the address in the config file", http.StatusBadRequest)
			return
		}

		o := &models.Order{
			Wallet:          wallet.Address,
			TokenInName:     args.TokenName,
			TokenOutName:    args.TokenName,
			SourceChainName: args.FromChain,
			TargetChainName: args.ToChain,
			Amount:          args.Amount,
			Status:          provider.TxStatusPending,
			ProviderType:    configs.CEX,
			ProviderName:    "gate.io",
			ProcessType:     string(bot.Queue),
			Remark:          "gate_liquidity",
			TaskId:          uuid.NewV4().String(),
		}
		if err := db.DB().Create(o).Error; err != nil {
			o.GetLogs().Errorf("create order error: %s", err.Error())
			http.Error(w, "create order error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		market.PushTask(market.Task{
			Id:          o.TaskId,
			ProcessType: bot.Queue,
		})
		logrus.Infof("create %s from %s to %s success", args.TokenName, args.FromChain, args.ToChain)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})
}
