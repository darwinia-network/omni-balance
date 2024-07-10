package daemons

import (
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils/configs"
	"omni-balance/utils/provider"
	"omni-balance/utils/wallets"

	"github.com/sirupsen/logrus"
)

func CreateSwapParams(order models.Order, orderProcess models.OrderProcess, log *logrus.Entry, wallet wallets.Wallets) provider.SwapParams {
	sourceChain := order.CurrentChainName
	if order.SourceChainName != "" {
		sourceChain = order.SourceChainName
	}
	return provider.SwapParams{
		OrderId:     order.ID,
		SourceChain: sourceChain,
		Sender:      wallet,
		Receiver:    order.Wallet,
		TargetChain: order.TargetChainName,
		SourceToken: order.TokenInName,
		TargetToken: order.TokenOutName,
		Amount:      order.Amount,
		LastHistory: createLastHistory(orderProcess),
		RecordFn:    createRecordFunction(order, log),
		Order:       order.Order,
		Remark:      order.Remark,
	}
}

func createLastHistory(orderProcess models.OrderProcess) provider.SwapHistory {
	return provider.SwapHistory{
		Actions:      orderProcess.Action,
		Status:       orderProcess.Status,
		CurrentChain: orderProcess.CurrentChainName,
		Amount:       orderProcess.Amount,
		Tx:           orderProcess.Tx,
	}
}

func createRecordFunction(order models.Order, log *logrus.Entry) func(s provider.SwapHistory, errs ...error) {
	return func(s provider.SwapHistory, errs ...error) {
		op := createOrderProcess(order, s)
		if len(errs) != 0 && errs[0] != nil {
			op.Error = errs[0].Error()
		}
		if err := db.DB().Create(op).Error; err != nil {
			log.Errorf("save %s bridge provider error: %s", order.TokenOutName, err.Error())
			return
		}
		log.Infof("order #%d action %s status is %s",
			order.ID, s.Actions, s.Status)
	}
}

func createOrderProcess(order models.Order, s provider.SwapHistory) *models.OrderProcess {
	var (
		providerType = order.ProviderType
		providerName = order.ProviderName
	)
	if s.ProviderType != "" {
		providerType = configs.ProviderType(s.ProviderType)
	}
	if s.ProviderName != "" {
		providerName = s.ProviderName
	}
	return &models.OrderProcess{
		OrderId:          order.ID,
		ProviderType:     providerType,
		ProviderName:     providerName,
		Status:           s.Status,
		Action:           s.Actions,
		Amount:           s.Amount,
		CurrentChainName: s.CurrentChain,
		Tx:               s.Tx,
	}
}
