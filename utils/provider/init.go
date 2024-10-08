package provider

import (
	"omni-balance/utils/configs"
	"sync"

	"github.com/pkg/errors"
)

type InitFunc func(conf configs.Config, noInit ...bool) (Provider, error)

var (
	providers = make(map[configs.ProviderType][]InitFunc)
	m         sync.Mutex
)

func Register(providerType configs.ProviderType, provider InitFunc) {
	m.Lock()
	defer m.Unlock()
	providers[providerType] = append(providers[providerType], provider)
}

func ListProviders() map[configs.ProviderType][]InitFunc {
	var result = make(map[configs.ProviderType][]InitFunc)
	for k, v := range providers {
		result[k] = v
	}
	return result
}

func ListProvidersByConfig(conf configs.Config) map[configs.ProviderType][]InitFunc {
	m.Lock()
	defer m.Unlock()
	var (
		providerNames = make(map[string]struct{})
		result        = make(map[configs.ProviderType][]InitFunc)
	)
	for _, v := range conf.Providers {
		providerNames[v.Name] = struct{}{}
	}
	for providerType, providerInitFuncs := range providers {
		for index, fn := range providerInitFuncs {
			p, _ := fn(conf, true)

			if _, ok := providerNames[p.Name()]; !ok {
				continue
			}
			result[providerType] = append(result[providerType], providers[providerType][index])
		}
	}
	return result
}

func LiquidityProviderTypeAndConf(providerType configs.ProviderType, conf configs.Config) []InitFunc {
	return ListProvidersByConfig(conf)[providerType]
}

func GetProvider(providerType configs.ProviderType, name string) (InitFunc, error) {
	for _, fn := range providers[providerType] {
		p, err := fn(configs.Config{}, true)
		if err != nil {
			return nil, err
		}
		if p.Name() != name {
			continue
		}
		return fn, nil
	}
	return nil, errors.New("provider not found")
}

func Init(providerInitFunc InitFunc, conf configs.Config, noInit ...bool) (Provider, error) {
	bridge, err := providerInitFunc(conf, noInit...)
	if err != nil {
		return nil, errors.Wrap(err, "init bridge error")
	}
	return bridge, nil
}
