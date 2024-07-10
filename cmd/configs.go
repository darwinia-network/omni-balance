package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"omni-balance/internal/daemons"
	"omni-balance/internal/db"
	"omni-balance/internal/handler"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	yaml_ncoder "github.com/zwgblue/yaml-encoder"
)

var (
	placeholder            sync.Map
	setPlaceholderFinished = make(chan struct{}, 1)
)

func startHttpServer(_ context.Context, port string) error {
	http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.EqualFold(request.Method, http.MethodPost) {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var args = make(map[string]interface{})
		if err := json.NewDecoder(request.Body).Decode(&args); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		for k, v := range args {
			placeholder.Store(k, v)
		}

		setPlaceholderFinished <- struct{}{}
	}))
	http.Handle("/api/gate/liquidity", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !utils.IsFinishedInit() {
			http.Error(w, "Not finished init", http.StatusInternalServerError)
			return
		}
		handler.GateLiquidity(*config)(w, r)
	}))
	http.Handle("/remove_order", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.EqualFold(request.Method, http.MethodPost) {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var order = struct {
			Id int `json:"id" form:"id"`
		}{}
		if err := json.NewDecoder(request.Body).Decode(&order); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		err := db.DB().Model(&models.Order{}).Where("id = ?", order.Id).Limit(1).Delete(&models.Order{}).Error
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
	server := &http.Server{
		Addr:    port,
		Handler: http.DefaultServeMux,
	}
	utils.Go(func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("http server error: %s", err)
		}
	})
	logrus.Infof("http server started on %s", port)
	return nil
}

func waitForPlaceholder(_ context.Context, configPath string) (newConfigPath string, err error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	<-setPlaceholderFinished
	placeholder.Range(func(key, value interface{}) bool {
		data = bytes.ReplaceAll(data, []byte(key.(string)), []byte(cast.ToString(value)))
		return true
	})
	newConfigPath = filepath.Join(os.TempDir(), ".omni-balance.config.yaml")
	if err := os.WriteFile(newConfigPath, data, 0644); err != nil {
		return "", err
	}
	return newConfigPath, err
}

func initConfig(ctx context.Context, enablePlaceholder bool, configPath, serverPort string) (err error) {
	err = startHttpServer(ctx, serverPort)
	if err != nil {
		return err
	}
	if enablePlaceholder {
		ports := strings.Split(serverPort, ":")
		if len(ports) < 2 {
			ports = append([]string{}, "", "8080")
		}
		logrus.Infof("waiting for placeholder, you can use `curl -X POST -d '{\"<you_placeholder>\":\"0x1234567890\"}' http://127.0.0.1:%s` to set placeholder", ports[1])
		configPath, err = waitForPlaceholder(context.Background(), configPath)
		if err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(configPath)
		}()
	}
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		return err
	}
	config.Init()
	return config.Check()
}

func CreateExampleConfig(exampleConfigPath string) error {
	var tasks = make(map[string]time.Duration)
	for _, v := range daemons.GetTaskConfig() {
		tasks[v.Name] = v.DefaultInterval
	}
	exampleConfig := configs.Config{
		Debug: true,
		Chains: []configs.Chain{
			{
				Id:   1,
				Name: "etnereum",
				RpcEndpoints: []string{
					"https://api.tatum.io/v3/blockchain/node/ethereum-mainnet",
					"https://ethereum-rpc.publicnode.com",
				},
				Tokens: []configs.Token{
					{
						Name:            "ETH",
						ContractAddress: constant.ZeroAddress.Hex(),
						Decimals:        18,
					},
					{
						Name:            "RING",
						ContractAddress: "0x9469D013805bFfB7D3DEBe5E7839237e535ec483",
						Decimals:        18,
					},
					{
						Name:            "USDT",
						ContractAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Decimals:        6,
					},
				},
			},
			{
				Id:   42161,
				Name: "arbitrum",
				RpcEndpoints: []string{
					"https://1rpc.io/arb",
					"https://arbitrum.llamarpc.com",
				},
				Tokens: []configs.Token{
					{
						Name:            "ETH",
						ContractAddress: constant.ZeroAddress.Hex(),
						Decimals:        18,
					},
					{
						Name:            "USDT",
						ContractAddress: "0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9",
						Decimals:        6,
					},
					{
						Name:            "RING",
						ContractAddress: "0x9e523234D36973f9e38642886197D023C88e307e",
						Decimals:        18,
					},
				},
			},
		},
		SourceTokens: []configs.SourceToken{
			{
				Name:   "USDT",
				Chains: []string{"etnereum", "arbitrum"},
			},
			{
				Name:   "ETH",
				Chains: []string{"etnereum", "arbitrum"},
			},
		},
		Providers: []configs.Provider{
			{
				Type: configs.CEX,
				Name: "gate.io",
				Config: map[string]interface{}{
					"key":    "<gate_api_key>",
					"secret": "<gate_api_secret>",
				},
			},
			{
				Type: configs.DEX,
				Name: "uniswap",
			},
			{
				Type: configs.Bridge,
				Name: "helixbridge",
			},
			{
				Type: configs.Bridge,
				Name: "darwinia-bridge",
			},
		},
		Wallets: []configs.Wallet{
			{
				Address: "0x43Ef13E84D9992d1461a1f90CAc4653658CEA4FD",
				Tokens: []configs.WalletToken{
					{
						Name:      "ETH",
						Amount:    decimal.RequireFromString("1"),
						Threshold: decimal.RequireFromString("3"),
						Chains:    []string{"ethereum"},
					},
				},
				PrivateKey: "<wallet1_private_key>",
			},
			{
				Address: "0x178D8546C5f78e01133858958355B06EC3406A1A",
				Tokens: []configs.WalletToken{
					{
						Name:      "RING",
						Amount:    decimal.RequireFromString("100000"),
						Threshold: decimal.RequireFromString("2000000"),
						Chains:    []string{"ethereum", "arbitrum"},
					},
				},
				PrivateKey: "<wallet2_private_key>",
			},
		},
		TaskInterval: tasks,
		Db: configs.DbConfig{
			Type:       configs.SQLite,
			SQLite:     &configs.Sqlite{Path: "./omni-balance.db"},
			MySQL:      new(configs.MysqlConfig),
			PostgreSQL: new(configs.MysqlConfig),
		},
	}
	exampleConfigData, err := yaml_ncoder.NewEncoder(exampleConfig,
		yaml_ncoder.WithComments(yaml_ncoder.CommentsOnHead),
		yaml_ncoder.WithOmitEmpty(false),
	).Encode()
	if err != nil {
		return errors.Wrap(err, "encode example config")
	}
	if err := os.WriteFile(exampleConfigPath, exampleConfigData, 0644); err != nil {
		return errors.Wrap(err, "write example config")
	}
	logrus.Infof("Example config file created: %s. In the example configuration file, some values are enclosed in '<>'."+
		" You can add the -p parameter at runtime and follow the prompts to replace these values.", exampleConfigPath)
	return nil
}
