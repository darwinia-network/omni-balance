package main

import (
	"context"
	"flag"
	"fmt"
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
	_ "omni-balance/utils/provider/bridge/bungee"
	_ "omni-balance/utils/provider/bridge/darwinia"
	_ "omni-balance/utils/provider/bridge/helix_liquidity_claim"
	_ "omni-balance/utils/provider/bridge/li"
	_ "omni-balance/utils/provider/bridge/okx"
	_ "omni-balance/utils/provider/bridge/routernitro"
	_ "omni-balance/utils/provider/cex/binance"
	_ "omni-balance/utils/provider/cex/gate"
	_ "omni-balance/utils/provider/dex/uniswap"
	_ "omni-balance/utils/wallets/multisig"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "omni-balance/utils/logging"

	"github.com/pkg/errors"

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
				log.Fatalf("init provider error: %v", err)
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
	if err := initConfig(ctx, cli.Bool("placeholder"), cli.String("conf"), cli.String("port")); err != nil {
		return errors.Wrap(err, "init config")
	}
	log.Infof("init config success")
	log.SetLevel(log.LevelInfo)
	if config.Debug {
		log.SetLevel(log.LevelDebug)
	}
	if version != "" && commitTime != "" {
		log.Infof("version: %s, commit: %s, commitTime: %s", version, commitMessage, commitTime)
	}

	if err := notice.Init(notice.Type(config.Notice.Type), config.Notice.Config, config.Notice.Interval); err != nil {
		log.Warnf("init notice error: %v", err)
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
	app.Commands = Commands
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
		log.Error(err)
	}
}
