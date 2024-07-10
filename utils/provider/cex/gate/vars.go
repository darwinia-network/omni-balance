package gate

var (
	ProvideLiquidityActionSourceChainTxsendingAction  = "provide_liquidity_source_chain_txsending"
	ProvideLiquidityActionSourceChainTxComplete       = "provide_liquidity_source_chain_tx_complete"
	ProvideLiquidityActionGateRechargeTxWaiting       = "provide_liquidity_gate_recharge_tx_waiting"
	ProvideLiquidityActionGateRechargeTxComplete      = "provide_liquidity_gate_recharge_tx_complete"
	ProvideLiquidityActionGateWithdrawSending         = "provide_liquidity_gate_withdraw_sending"
	ProvideLiquidityActionGateWithdrawWaiting         = "provide_liquidity_gate_withdraw_waiting"
	ProvideLiquidityActionGateWithdrawComplete        = "provide_liquidity_gate_withdraw_complete"
	ProvideLiquidityActionTargetChainWithdrawComplete = "provide_liquidity_target_chain_withdraw_complete"
)

func ProvideLiquidityAction2Number(action string) int {
	switch action {
	case ProvideLiquidityActionSourceChainTxsendingAction:
		return 1
	case ProvideLiquidityActionSourceChainTxComplete:
		return 2
	case ProvideLiquidityActionGateRechargeTxWaiting:
		return 3
	case ProvideLiquidityActionGateRechargeTxComplete:
		return 4
	case ProvideLiquidityActionGateWithdrawSending:
		return 5
	case ProvideLiquidityActionGateWithdrawWaiting:
		return 6
	case ProvideLiquidityActionGateWithdrawComplete:
		return 7
	case ProvideLiquidityActionTargetChainWithdrawComplete:
		return 8
	default:
		return 0
	}
}
