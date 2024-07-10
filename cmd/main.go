package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"omni-balance/internal/actions"
	"omni-balance/internal/daemons"
	_ "omni-balance/internal/daemons/bot"
	_ "omni-balance/internal/daemons/cross_chain"
	_ "omni-balance/internal/daemons/market"
	_ "omni-balance/internal/daemons/token_price"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"omni-balance/utils/notice"
	"omni-balance/utils/provider"
	_ "omni-balance/utils/provider/bridge/darwinia"
	_ "omni-balance/utils/provider/bridge/helix_liquidity_claim"
	_ "omni-balance/utils/provider/bridge/li"
	_ "omni-balance/utils/provider/bridge/okx"
	_ "omni-balance/utils/provider/bridge/routernitro"
	_ "omni-balance/utils/provider/cex/gate"
	_ "omni-balance/utils/provider/dex/uniswap"
	_ "omni-balance/utils/wallets/safe"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	config        = new(configs.Config)
	defaultsUsage = flag.Usage
	ctx, cancel   = context.WithCancel(context.TODO())
	// fill it when building
	version       string
	commitMessage string
	commitTime    string
)

func Usage(_ *cli.Context) error {
	defaultsUsage()
	fmt.Printf("Supported providers:\n")
	for providerType, providerFns := range provider.ListProviders() {
		fmt.Printf(" %s:\n", providerType)
		for _, fn := range providerFns {
			providerObj, err := fn(*config, true)
			if err != nil {
				logrus.Panicf("init provider error: %v", err)
			}
			fmt.Printf("  %s:\n", providerObj.Name())
			for _, v := range providerObj.Help() {
				fmt.Printf("   %s\n", v)
			}
		}
	}
	return nil
}

func Action(cli *cli.Context) error {
	fmt.Println(cli.String("conf"))
	if err := initConfig(ctx, cli.Bool("placeholder"), cli.String("conf"), cli.String("port")); err != nil {
		return errors.Wrap(err, "init config")
	}
	logrus.Infof("version: %s, commit: %s, commitTime: %s", version, commitMessage, commitTime)
	if config.Debug {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:          true,
			ForceQuote:             true,
			DisableLevelTruncation: false,
			QuoteEmptyFields:       true,
		})
	}

	if !config.Debug {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	if err := notice.Init(notice.Type(config.Notice.Type), config.Notice.Config, config.Notice.Interval); err != nil {
		logrus.Warnf("init notice error: %v", err)
	}

	if err := db.InitDb(*config); err != nil {
		return errors.Wrap(err, "init db")
	}

	if err := db.DB().AutoMigrate(
		new(models.Order),
		new(models.OrderProcess),
		new(models.TokenPrice)); err != nil {
		return errors.Wrap(err, "auto migrate db")
	}

	if err := daemons.Run(ctx, *config); err != nil {
		return errors.Wrap(err, "run daemons")
	}

	utils.FinishInit()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	if os.Getenv("is_local") != "true" {
		time.Sleep(time.Second * 5)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "omni-balance"
	app.Action = Action
	app.Commands = []*cli.Command{
		{
			Name:        "gate_liquidity",
			Usage:       "Create an order for the liquidity of Gate.",
			Description: "Create an order for the liquidity of Gate.",
			Action:      actions.DoGateLiquidity,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Usage: "omni-balance server port",
					Value: "http://127.0.0.1:8080",
				},
				&cli.StringFlag{
					Name:     "apiKey",
					Required: true,
					Usage:    "omni-balance API key in config",
				},
				&cli.StringFlag{
					Name:     "tokenName",
					Required: true,
					Usage:    "a token name must be in config",
				},
				&cli.StringFlag{
					Name:     "fromChain",
					Required: true,
					Usage:    "deposit to the Gate exchange from that chain.",
				},
				&cli.StringFlag{
					Name:     "toChain",
					Required: true,
					Usage:    "Withdraw to that chain from the Gate exchange.",
				},
				&cli.StringFlag{
					Name:     "amount",
					Required: true,
					Usage:    "The number of tokens that need to be rebalanced.",
				},
				&cli.StringFlag{
					Name:     "address",
					Required: true,
					Usage:    "Use that address to rebalance, and note that this address must exist in the configuration file.",
				},
			},
		},
		{
			Name:  "del_order",
			Usage: "delete order by id",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "id",
					Usage: "order id",
				},
				&cli.StringFlag{
					Name:  "server",
					Usage: "server host, example: http://127.0.0.1:8080",
					Value: "http://127.0.0.1:8080",
				},
			},
			Action: func(c *cli.Context) error {
				u, err := url.Parse(c.String("server"))
				if err != nil {
					return errors.Wrap(err, "parse server url")
				}
				u.RawPath = "/remove_order"
				u.Path = u.RawPath
				var body = bytes.NewBuffer(nil)
				err = json.NewEncoder(body).Encode(map[string]interface{}{
					"id": c.Int("id"),
				})
				if err != nil {
					return errors.Wrap(err, "encode body")
				}
				resp, err := http.Post(u.String(), "application/json", body)
				if err != nil {
					return errors.Wrap(err, "post")
				}
				defer resp.Body.Close()
				data, _ := io.ReadAll(resp.Body)
				if resp.StatusCode != http.StatusOK {
					return errors.Errorf("http status code: %d, body is: %s", resp.StatusCode, data)
				}
				logrus.Infof("delete order #%d success", c.Int64("id"))
				return nil
			},
		},
		{
			Name:    "version",
			Usage:   "show version",
			Aliases: []string{"v"},
			Action: func(c *cli.Context) error {
				fmt.Printf("Version: %s\n", version)
				fmt.Printf("Commit: %s\n", commitMessage)
				fmt.Printf("Build time: %s\n", commitTime)
				return nil
			},
		},
		{
			Name:   "list",
			Usage:  "list supported providers and docs",
			Action: Usage,
		},
		{
			Name:  "tasks",
			Usage: "list supported tasks",
			Action: func(_ *cli.Context) error {
				fmt.Println(daemons.Help())
				return nil
			},
		},
		{
			Name:  "example",
			Usage: "create a example config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Usage:   "output file path",
					Value:   "./config.yaml.example",
					Aliases: []string{"o"},
				},
			},
			Action: func(c *cli.Context) error {
				if err := CreateExampleConfig(c.String("output")); err != nil {
					return errors.Wrap(err, "create example config")
				}
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Usage:   "config file path",
			Value:   "./config.yaml",
		},
		&cli.BoolFlag{
			Name: "placeholder",
			Usage: fmt.Sprintf("enable placeholder, you can use placeholder to replace private key, Example: Fill '{{privateKey}}' in config.yaml."+
				"Run with -p to enable placeholder, Example: SERVER_PORT=:8080 %s -c ./config.yaml -p"+
				"Waiting for 'waiting for placeholder...' log, send placeholder data according to the prompt.", os.Args[0]),
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:  "port",
			Usage: "When the placeholder parameter is set to true, you can specify and set the listening address of the HTTP server that receives the placeholder.",
			Value: ":8080",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
	}
}
