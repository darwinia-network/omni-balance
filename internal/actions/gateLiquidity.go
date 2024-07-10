package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"omni-balance/internal/handler"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"
)

func DoGateLiquidity(cli *cli.Context) error {
	server, err := url.Parse(cli.String("server"))
	if err != nil {
		return err
	}
	apiKey := cli.String("apiKey")
	args := new(handler.GateLiquidityParams)
	args.TokenName = cli.String("tokenName")
	args.FromChain = cli.String("fromChain")
	args.ToChain = cli.String("toChain")
	args.Amount = decimal.RequireFromString(cli.String("amount"))
	args.Address = cli.String("address")

	if apiKey == "" {
		return errors.New("apiKey is required")
	}
	if args.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}
	var body = bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(args); err != nil {
		return err
	}
	server.Path = "/api/gate/liquidity"
	req, err := http.NewRequest("POST", server.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-KEY", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	data, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d, body: %s", resp.StatusCode, string(data))
	}
	return nil
}
