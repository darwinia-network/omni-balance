package token_price

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"omni-balance/internal/daemons"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"omni-balance/utils/token_price"
	"sync"
	"time"
)

func init() {
	daemons.RegisterIntervalTask(daemons.Task{
		Name:            "getTokenPriceInUsdt",
		Description:     "Responsible for obtaining the token price, denominated in USDT.",
		TaskFunc:        Run,
		DefaultInterval: time.Minute * 3,
		// RunOnStart:      true,
	})
}

func Run(ctx context.Context, conf configs.Config) error {
	var (
		w          sync.WaitGroup
		m          sync.Mutex
		tokenPrice []models.TokenPrice
	)
	providers := token_price.ListTokenPriceProviders()
	for _, provider := range providers {
		w.Add(1)
		go func(provider token_price.TokenPrice) {
			defer utils.Recover()
			defer w.Done()
			result, err := provider.GetTokenPriceInUSDT(ctx, conf.SourceTokens...)
			if err != nil {
				logrus.Warnf("%s get token price error: %s", provider.Name(), err)
				return
			}
			m.Lock()
			defer m.Unlock()
			for _, r := range result {
				if r.Price.LessThan(decimal.Zero) {
					logrus.Warnf("%s get %s token price is less than 0", provider.Name(), r.TokeName)
					continue
				}
				tokenPrice = append(tokenPrice, models.TokenPrice{
					TokenName: r.TokeName,
					Price:     r.Price,
					Source:    provider.Name(),
				})
			}

		}(provider)
	}
	w.Wait()
	return models.SaveTokenPrice(ctx, db.DB().Begin(), tokenPrice)
}
