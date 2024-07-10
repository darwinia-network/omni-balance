package bot

import (
	"context"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/bot"
	"omni-balance/utils/provider"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func getExistingBuyTokens() ([]*models.Order, error) {
	var existBuyTokens []*models.Order
	err := db.DB().Where("status != ? ", provider.TxStatusSuccess).Find(&existBuyTokens).Error
	if err != nil {
		return nil, err
	}
	return existBuyTokens, nil
}

func createIgnoreTokens(existBuyTokens []*models.Order) IgnoreTokens {
	var ignoreTokens []IgnoreToken
	for _, v := range existBuyTokens {
		ignoreTokens = append(ignoreTokens, IgnoreToken{
			Name:    v.TokenOutName,
			Chain:   v.TargetChainName,
			Address: v.Wallet,
		})
	}
	return ignoreTokens
}

func createOrder(ctx context.Context, tasks []bot.Task, processType bot.ProcessType) (orders []*models.Order, taskId string, err error) {
	if len(tasks) == 0 {
		return
	}
	var (
		txn = db.DB().WithContext(ctx).Begin()
	)

	taskId = uuid.NewV4().String()

	for _, v := range tasks {
		d := db.DB().Where("wallet = ? and target_chain_name = ? and token_out_name = ? and status != ? and remark = ?",
			v.Wallet, v.TokenOutChainName, v.TokenOutName, provider.TxStatusSuccess, v.Remark)
		var count int64
		if d.Model(&models.Order{}).Count(&count); count > 0 {
			logrus.Debugf("wallet: %s, target_chain_name: %s, token_out_name: %s, status != success, remark: %s order count > 1 skip it",
				v.Wallet, v.TokenOutChainName, v.TokenOutName, v.Remark)
			continue
		}
		o := &models.Order{
			Wallet:           v.Wallet,
			TokenInName:      v.TokenInName,
			TokenOutName:     v.TokenOutName,
			SourceChainName:  v.TokenInChainName,
			TargetChainName:  v.TokenOutChainName,
			CurrentChainName: v.CurrentChainName,
			Amount:           v.Amount,
			Status:           v.Status,
			ProviderType:     v.ProviderType,
			ProviderName:     v.ProviderName,
			Order:            utils.Object2Json(v.Order),
			TaskId:           taskId,
			ProcessType:      string(processType),
			Remark:           v.Remark,
		}

		if err = txn.Create(o).Error; err != nil {
			txn.Rollback()
			return
		}
		orders = append(orders, o)
	}
	if err = txn.Commit().Error; err != nil {
		txn.Rollback()
		return
	}
	return
}
