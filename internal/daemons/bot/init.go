package bot

import (
	"omni-balance/internal/daemons"
	_ "omni-balance/utils/bot/balance_on_chain"
	_ "omni-balance/utils/bot/gate_liquidity"
	_ "omni-balance/utils/bot/helix_liquidity"
	"time"
)

func init() {
	daemons.RegisterIntervalTask(daemons.Task{
		Name:            "bot",
		Description:     "Check the balance based on the specific bot and create an order",
		TaskFunc:        Run,
		DefaultInterval: time.Minute * 3,
	})
}
