package bungee

import (
	"omni-balance/utils/configs"
	"omni-balance/utils/provider"
)

func init() {
	provider.Register(configs.Bridge, New)
}
