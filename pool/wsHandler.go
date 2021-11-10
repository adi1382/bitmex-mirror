package pool

import (
	"encoding/json"
	"log"

	"github.com/bitmex-mirror/bitmex"
	"github.com/bitmex-mirror/optional"
)

//func (a *account) updateFromSocket(ctx context.Context, logger *log.Logger) {
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case msg, ok := <-a.wsUpdater:
//			if !ok {
//				logger.Println("receiver closed")
//				return
//			}
//			a.handleMessage(msg, logger)
//		}
//	}
//}

func (a *account) handleMessage(data []byte, logger *log.Logger) {
	var msg struct {
		Table  string            `json:"table"`
		Action string            `json:"action"`
		Data   []json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return
	}
	switch msg.Table {
	case bitmex.WSTableOrder:
		a.handleOrderMessage(msg.Data, msg.Action, logger)
	case bitmex.WSTablePosition:
		a.handlePositionMessage(msg.Data, msg.Action, logger)
	case bitmex.WSTableMargin:
		a.handleMarginMessage(msg.Data, msg.Action, logger)
	default:
		return
	}
}

func (a *account) handleOrderMessage(data []json.RawMessage, action string, logger *log.Logger) {
	switch action {
	case bitmex.WSDataActionPartial:
		a.handlePartialOrderMessage(data, logger)
	case bitmex.WSDataActionInsert:
		a.handleInsertOrderMessage(data, logger)
	case bitmex.WSDataActionUpdate:
		a.handleUpdateOrderMessage(data, logger)
	case bitmex.WSDataActionDelete:
		a.handleDeleteOrderMessage(data, logger)
	}
}

func (a *account) handlePartialOrderMessage(data []json.RawMessage, logger *log.Logger) {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	a.activeOrders = a.activeOrders[:0]

	for i := range data {
		var order bitmex.Order
		err := json.Unmarshal(data[i], &order)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		if order.OrdStatus.ElseZero() == bitmex.OrderStatusNew ||
			order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
			a.activeOrders = append(a.activeOrders, order)
		}
	}
}

func (a *account) handleInsertOrderMessage(data []json.RawMessage, logger *log.Logger) {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	for i := range data {
		var order bitmex.Order
		err := json.Unmarshal(data[i], &order)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		// insert order might already be in active orders because of PlaceOrder method
		// this should be used only by SubAccounts, for HostAccount no matching order should be found
		// but running this anyway for host too, because it's not a problem and just in case
		isOrderExists := false
		for j := range a.activeOrders {
			if a.activeOrders[j].OrderID == order.OrderID {
				isOrderExists = true
				// Already in active orders, update
				a.activeOrders[j] = order

				// remove order if canceled or filled, this should never happen here
				// this functionality is expected to be handled by handleUpdateOrderMessage or handleDeleteOrderMessage
				// but just in case
				if a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusCanceled ||
					a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusFilled ||
					a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusRejected {
					a.activeOrders[j] = a.activeOrders[len(a.activeOrders)-1]
					a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				}
				break
			}
		}
		if !isOrderExists {
			// Insert new order, if not already in active orders
			// This is mostly used by HostAccount, if used for SubAccount that would mean a manually placed order
			// on the SubAccount that is not placed by mirror bot
			if order.OrdStatus.ElseZero() == bitmex.OrderStatusNew ||
				order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
				a.activeOrders = append(a.activeOrders, order)
			}
		}
	}
}

func (a *account) handleUpdateOrderMessage(data []json.RawMessage, logger *log.Logger) {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	for i := range data {
		var Order bitmex.Order
		err := json.Unmarshal(data[i], &Order)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		isOrderExists := false
		for j := range a.activeOrders {
			if a.activeOrders[j].OrderID == Order.OrderID {
				isOrderExists = true
				updateOrder(&a.activeOrders[j], Order)

				// remove order if canceled or filled
				if a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusCanceled ||
					a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusFilled ||
					a.activeOrders[j].OrdStatus.ElseZero() == bitmex.OrderStatusRejected {
					a.activeOrders[j] = a.activeOrders[len(a.activeOrders)-1]
					a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				}
				break
			}
		}

		if isOrderExists {
			continue
		}

		logger.Println("Order not found:", Order.OrderID)
	}
}

func (a *account) handleDeleteOrderMessage(data []json.RawMessage, logger *log.Logger) {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	for i := range data {
		var Order bitmex.Order
		err := json.Unmarshal(data[i], &Order)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		for j := range a.activeOrders {
			if a.activeOrders[j].OrderID == Order.OrderID {
				a.activeOrders[j] = a.activeOrders[len(a.activeOrders)-1]
				a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				break
			}
		}
	}
}

func (a *account) handlePositionMessage(data []json.RawMessage, action string, logger *log.Logger) {
	switch action {
	case bitmex.WSDataActionPartial:
		a.handlePartialPositionMessage(data, logger)
	case bitmex.WSDataActionInsert:
		a.handleInsertPositionMessage(data, logger)
	case bitmex.WSDataActionUpdate:
		a.handleUpdatePositionMessage(data, logger)
	case bitmex.WSDataActionDelete:
		a.handleDeletePositionMessage(data, logger)
	}
}

func (a *account) handlePartialPositionMessage(data []json.RawMessage, logger *log.Logger) {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	a.positions = a.positions[:0]

	for i := range data {
		var position bitmex.Position
		err := json.Unmarshal(data[i], &position)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		a.positions = append(a.positions, position)
	}
}

func (a *account) handleInsertPositionMessage(data []json.RawMessage, logger *log.Logger) {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	for i := range data {
		var position bitmex.Position
		err := json.Unmarshal(data[i], &position)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		isPositionExists := false
		for j := range a.positions {
			if a.positions[j].Symbol == position.Symbol && a.positions[j].Account == position.Account {
				isPositionExists = true
				a.positions[j] = position
				break
			}
		}
		if !isPositionExists {
			a.positions = append(a.positions, position)
		}
	}
}

func (a *account) handleUpdatePositionMessage(data []json.RawMessage, logger *log.Logger) {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	for i := range data {
		var position bitmex.Position
		err := json.Unmarshal(data[i], &position)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		isPositionExists := false
		for j := range a.positions {
			if a.positions[j].Symbol == position.Symbol && a.positions[j].Account == position.Account {
				isPositionExists = true
				updatePosition(&a.positions[j], position)
				break
			}
		}
		if !isPositionExists {
			logger.Println("Position not found:", position.Symbol)
		}
	}
}

func (a *account) handleDeletePositionMessage(data []json.RawMessage, logger *log.Logger) {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	for i := range data {
		var position bitmex.Position
		err := json.Unmarshal(data[i], &position)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		for j := range a.positions {
			if a.positions[j].Symbol == position.Symbol && a.positions[j].Account == position.Account {
				a.positions[j] = a.positions[len(a.positions)-1]
				a.positions = a.positions[:len(a.positions)-1]
				break
			}
		}
	}
}

func (a *account) handleMarginMessage(data []json.RawMessage, action string, logger *log.Logger) {
	switch action {
	case bitmex.WSDataActionPartial:
		a.handlePartialMarginMessage(data, logger)
	case bitmex.WSDataActionInsert:
		a.handleInsertMarginMessage(data, logger)
	case bitmex.WSDataActionUpdate:
		a.handleUpdateMarginMessage(data, logger)
	case bitmex.WSDataActionDelete:
		a.handleDeleteMarginMessage(data, logger)
	}
}

func (a *account) handlePartialMarginMessage(data []json.RawMessage, logger *log.Logger) {
	a.marginMu.Lock()
	defer a.marginMu.Unlock()

	a.margins = a.margins[:0]

	for i := range data {
		var margin bitmex.Margin
		err := json.Unmarshal(data[i], &margin)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		a.margins = append(a.margins, margin)
	}
}

// It is literally the same as partial margin message
func (a *account) handleInsertMarginMessage(data []json.RawMessage, logger *log.Logger) {
	a.marginMu.Lock()
	defer a.marginMu.Unlock()

	for i := range data {
		var margin bitmex.Margin
		err := json.Unmarshal(data[i], &margin)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		isMarginExists := false
		for j := range a.margins {
			if a.margins[j].Currency == margin.Currency && a.margins[j].Account == margin.Account {
				isMarginExists = true
				a.margins[j] = margin
				break
			}
		}

		if !isMarginExists {
			a.margins = append(a.margins, margin)
		}
	}
}

func (a *account) handleUpdateMarginMessage(data []json.RawMessage, logger *log.Logger) {
	a.marginMu.Lock()
	defer a.marginMu.Unlock()

	for i := range data {
		var margin bitmex.Margin
		err := json.Unmarshal(data[i], &margin)
		if err != nil {
			logger.Println(err, string(data[i]))
			continue
		}

		isMarginExists := false
		for j := range a.margins {
			if a.margins[j].Currency == margin.Currency && a.margins[j].Account == margin.Account {
				isMarginExists = true
				updateMargin(&a.margins[j], margin)
				break
			}
		}

		if !isMarginExists {
			logger.Println("Margin not found:", margin.Currency)
		}
	}
}

// Hypothetical case, Bitmex would never send delete message for margin.
// The number of items in margin is equal to the number of currencies in the account.
// Currently, there are only two currencies on Bitmex : USDT and XBT.
func (a *account) handleDeleteMarginMessage(data []json.RawMessage, logger *log.Logger) {
	logger.Println("No need to delete margin data, this should have never happened")
	for i := range data {
		logger.Println("Margin Delete Data", string(data[i]))
	}
}

// This function needs dst to be held to prevent concurrent access to dst
func updateOrder(dst *bitmex.Order, src bitmex.Order) {
	if v, ok := src.ClOrdID.Get(); ok {
		dst.ClOrdID = optional.OfString(v)
	}
	if v, ok := src.Symbol.Get(); ok {
		dst.Symbol = optional.OfString(v)
	}
	if v, ok := src.Side.Get(); ok {
		dst.Side = optional.OfString(v)
	}
	if v, ok := src.OrderQty.Get(); ok {
		dst.OrderQty = optional.OfFloat64(v)
	}
	if v, ok := src.Price.Get(); ok {
		dst.Price = optional.OfFloat64(v)
	}
	if v, ok := src.DisplayQty.Get(); ok {
		dst.DisplayQty = optional.OfFloat64(v)
	}
	if v, ok := src.StopPx.Get(); ok {
		dst.StopPx = optional.OfFloat64(v)
	}
	if v, ok := src.PegOffsetValue.Get(); ok {
		dst.PegOffsetValue = optional.OfFloat64(v)
	}
	if v, ok := src.PegPriceType.Get(); ok {
		dst.PegPriceType = optional.OfString(v)
	}
	if v, ok := src.Currency.Get(); ok {
		dst.Currency = optional.OfString(v)
	}
	if v, ok := src.SettlCurrency.Get(); ok {
		dst.SettlCurrency = optional.OfString(v)
	}
	if v, ok := src.OrdType.Get(); ok {
		dst.OrdType = optional.OfString(v)
	}
	if v, ok := src.TimeInForce.Get(); ok {
		dst.TimeInForce = optional.OfString(v)
	}
	if v, ok := src.ExecInst.Get(); ok {
		dst.ExecInst = optional.OfString(v)
	}
	if v, ok := src.ExDestination.Get(); ok {
		dst.ExDestination = optional.OfString(v)
	}
	if v, ok := src.OrdStatus.Get(); ok {
		dst.OrdStatus = optional.OfString(v)
	}
	if v, ok := src.Triggered.Get(); ok {
		dst.Triggered = optional.OfString(v)
	}
	if v, ok := src.WorkingIndicator.Get(); ok {
		dst.WorkingIndicator = optional.OfBool(v)
	}
	if v, ok := src.OrdRejReason.Get(); ok {
		dst.OrdRejReason = optional.OfString(v)
	}
	if v, ok := src.LeavesQty.Get(); ok {
		dst.LeavesQty = optional.OfFloat64(v)
	}
	if v, ok := src.CumQty.Get(); ok {
		dst.CumQty = optional.OfFloat64(v)
	}
	if v, ok := src.AvgPx.Get(); ok {
		dst.AvgPx = optional.OfFloat64(v)
	}
	if v, ok := src.MultiLegReportingType.Get(); ok {
		dst.MultiLegReportingType = optional.OfString(v)
	}
	if v, ok := src.Text.Get(); ok {
		dst.Text = optional.OfString(v)
	}
	if v, ok := src.TransactTime.Get(); ok {
		dst.TransactTime = optional.OfTime(v)
	}
	if v, ok := src.Timestamp.Get(); ok {
		dst.Timestamp = optional.OfTime(v)
	}
}

// This function needs dst to be held to prevent concurrent access to dst
func updatePosition(dst *bitmex.Position, src bitmex.Position) {
	if v, ok := src.Currency.Get(); ok {
		dst.Currency = optional.OfString(v)
	}
	if v, ok := src.Underlying.Get(); ok {
		dst.Underlying = optional.OfString(v)
	}
	if v, ok := src.QuoteCurrency.Get(); ok {
		dst.QuoteCurrency = optional.OfString(v)
	}
	if v, ok := src.Commission.Get(); ok {
		dst.Commission = optional.OfFloat64(v)
	}
	if v, ok := src.InitMarginReq.Get(); ok {
		dst.InitMarginReq = optional.OfFloat64(v)
	}
	if v, ok := src.MaintMarginReq.Get(); ok {
		dst.MaintMarginReq = optional.OfFloat64(v)
	}
	if v, ok := src.RiskLimit.Get(); ok {
		dst.RiskLimit = optional.OfFloat64(v)
	}
	if v, ok := src.Leverage.Get(); ok {
		dst.Leverage = optional.OfFloat64(v)
	}
	if v, ok := src.CrossMargin.Get(); ok {
		dst.CrossMargin = optional.OfBool(v)
	}
	if v, ok := src.DeleveragePercentile.Get(); ok {
		dst.DeleveragePercentile = optional.OfFloat64(v)
	}
	if v, ok := src.RebalancedPnl.Get(); ok {
		dst.RebalancedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.PrevRealisedPnl.Get(); ok {
		dst.PrevRealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.PrevUnrealisedPnl.Get(); ok {
		dst.PrevUnrealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.PrevClosePrice.Get(); ok {
		dst.PrevClosePrice = optional.OfFloat64(v)
	}
	if v, ok := src.OpeningTimestamp.Get(); ok {
		dst.OpeningTimestamp = optional.OfTime(v)
	}
	if v, ok := src.OpeningQty.Get(); ok {
		dst.OpeningQty = optional.OfFloat64(v)
	}
	if v, ok := src.OpeningCost.Get(); ok {
		dst.OpeningCost = optional.OfFloat64(v)
	}
	if v, ok := src.OpeningComm.Get(); ok {
		dst.OpeningComm = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderBuyQty.Get(); ok {
		dst.OpenOrderBuyQty = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderBuyCost.Get(); ok {
		dst.OpenOrderBuyCost = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderBuyPremium.Get(); ok {
		dst.OpenOrderBuyPremium = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderSellQty.Get(); ok {
		dst.OpenOrderSellQty = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderSellCost.Get(); ok {
		dst.OpenOrderSellCost = optional.OfFloat64(v)
	}
	if v, ok := src.OpenOrderSellPremium.Get(); ok {
		dst.OpenOrderSellPremium = optional.OfFloat64(v)
	}
	if v, ok := src.ExecBuyQty.Get(); ok {
		dst.ExecBuyQty = optional.OfFloat64(v)
	}
	if v, ok := src.ExecBuyCost.Get(); ok {
		dst.ExecBuyCost = optional.OfFloat64(v)
	}
	if v, ok := src.ExecSellQty.Get(); ok {
		dst.ExecSellQty = optional.OfFloat64(v)
	}
	if v, ok := src.ExecSellCost.Get(); ok {
		dst.ExecSellCost = optional.OfFloat64(v)
	}
	if v, ok := src.ExecQty.Get(); ok {
		dst.ExecQty = optional.OfFloat64(v)
	}
	if v, ok := src.ExecCost.Get(); ok {
		dst.ExecCost = optional.OfFloat64(v)
	}
	if v, ok := src.ExecComm.Get(); ok {
		dst.ExecComm = optional.OfFloat64(v)
	}
	if v, ok := src.CurrentTimestamp.Get(); ok {
		dst.CurrentTimestamp = optional.OfTime(v)
	}
	if v, ok := src.CurrentQty.Get(); ok {
		dst.CurrentQty = optional.OfFloat64(v)
	}
	if v, ok := src.CurrentCost.Get(); ok {
		dst.CurrentCost = optional.OfFloat64(v)
	}
	if v, ok := src.CurrentComm.Get(); ok {
		dst.CurrentComm = optional.OfFloat64(v)
	}
	if v, ok := src.RealisedCost.Get(); ok {
		dst.RealisedCost = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedCost.Get(); ok {
		dst.UnrealisedCost = optional.OfFloat64(v)
	}
	if v, ok := src.GrossOpenCost.Get(); ok {
		dst.GrossOpenCost = optional.OfFloat64(v)
	}
	if v, ok := src.GrossOpenPremium.Get(); ok {
		dst.GrossOpenPremium = optional.OfFloat64(v)
	}
	if v, ok := src.GrossExecCost.Get(); ok {
		dst.GrossExecCost = optional.OfFloat64(v)
	}
	if v, ok := src.IsOpen.Get(); ok {
		dst.IsOpen = optional.OfBool(v)
	}
	if v, ok := src.MarkPrice.Get(); ok {
		dst.MarkPrice = optional.OfFloat64(v)
	}
	if v, ok := src.MarkValue.Get(); ok {
		dst.MarkValue = optional.OfFloat64(v)
	}
	if v, ok := src.RiskValue.Get(); ok {
		dst.RiskValue = optional.OfFloat64(v)
	}
	if v, ok := src.HomeNotional.Get(); ok {
		dst.HomeNotional = optional.OfFloat64(v)
	}
	if v, ok := src.ForeignNotional.Get(); ok {
		dst.ForeignNotional = optional.OfFloat64(v)
	}
	if v, ok := src.PosState.Get(); ok {
		dst.PosState = optional.OfString(v)
	}
	if v, ok := src.PosCost.Get(); ok {
		dst.PosCost = optional.OfFloat64(v)
	}
	if v, ok := src.PosCost2.Get(); ok {
		dst.PosCost2 = optional.OfFloat64(v)
	}
	if v, ok := src.PosCross.Get(); ok {
		dst.PosCross = optional.OfFloat64(v)
	}
	if v, ok := src.PosInit.Get(); ok {
		dst.PosInit = optional.OfFloat64(v)
	}
	if v, ok := src.PosComm.Get(); ok {
		dst.PosComm = optional.OfFloat64(v)
	}
	if v, ok := src.PosLoss.Get(); ok {
		dst.PosLoss = optional.OfFloat64(v)
	}
	if v, ok := src.PosMargin.Get(); ok {
		dst.PosMargin = optional.OfFloat64(v)
	}
	if v, ok := src.PosMaint.Get(); ok {
		dst.PosMaint = optional.OfFloat64(v)
	}
	if v, ok := src.PosAllowance.Get(); ok {
		dst.PosAllowance = optional.OfFloat64(v)
	}
	if v, ok := src.TaxableMargin.Get(); ok {
		dst.TaxableMargin = optional.OfFloat64(v)
	}
	if v, ok := src.InitMargin.Get(); ok {
		dst.InitMargin = optional.OfFloat64(v)
	}
	if v, ok := src.MaintMargin.Get(); ok {
		dst.MaintMargin = optional.OfFloat64(v)
	}
	if v, ok := src.SessionMargin.Get(); ok {
		dst.SessionMargin = optional.OfFloat64(v)
	}
	if v, ok := src.TargetExcessMargin.Get(); ok {
		dst.TargetExcessMargin = optional.OfFloat64(v)
	}
	if v, ok := src.VarMargin.Get(); ok {
		dst.VarMargin = optional.OfFloat64(v)
	}
	if v, ok := src.RealisedGrossPnl.Get(); ok {
		dst.RealisedGrossPnl = optional.OfFloat64(v)
	}
	if v, ok := src.RealisedTax.Get(); ok {
		dst.RealisedTax = optional.OfFloat64(v)
	}
	if v, ok := src.RealisedPnl.Get(); ok {
		dst.RealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedGrossPnl.Get(); ok {
		dst.UnrealisedGrossPnl = optional.OfFloat64(v)
	}
	if v, ok := src.LongBankrupt.Get(); ok {
		dst.LongBankrupt = optional.OfFloat64(v)
	}
	if v, ok := src.ShortBankrupt.Get(); ok {
		dst.ShortBankrupt = optional.OfFloat64(v)
	}
	if v, ok := src.TaxBase.Get(); ok {
		dst.TaxBase = optional.OfFloat64(v)
	}
	if v, ok := src.IndicativeTaxRate.Get(); ok {
		dst.IndicativeTaxRate = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedTax.Get(); ok {
		dst.UnrealisedTax = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedPnl.Get(); ok {
		dst.UnrealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedPnlPcnt.Get(); ok {
		dst.UnrealisedPnlPcnt = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedRoePcnt.Get(); ok {
		dst.UnrealisedRoePcnt = optional.OfFloat64(v)
	}
	if v, ok := src.AvgCostPrice.Get(); ok {
		dst.AvgCostPrice = optional.OfFloat64(v)
	}
	if v, ok := src.AvgEntryPrice.Get(); ok {
		dst.AvgEntryPrice = optional.OfFloat64(v)
	}
	if v, ok := src.BreakEvenPrice.Get(); ok {
		dst.BreakEvenPrice = optional.OfFloat64(v)
	}
	if v, ok := src.MarginCallPrice.Get(); ok {
		dst.MarginCallPrice = optional.OfFloat64(v)
	}
	if v, ok := src.LiquidationPrice.Get(); ok {
		dst.LiquidationPrice = optional.OfFloat64(v)
	}
	if v, ok := src.BankruptPrice.Get(); ok {
		dst.BankruptPrice = optional.OfFloat64(v)
	}
	if v, ok := src.Timestamp.Get(); ok {
		dst.Timestamp = optional.OfTime(v)
	}
	if v, ok := src.LastPrice.Get(); ok {
		dst.LastPrice = optional.OfFloat64(v)
	}
	if v, ok := src.LastValue.Get(); ok {
		dst.LastValue = optional.OfFloat64(v)
	}
}

// This function needs dst to be held to prevent concurrent access to dst
func updateMargin(dst *bitmex.Margin, src bitmex.Margin) {
	if v, ok := src.Action.Get(); ok {
		dst.Action = optional.OfString(v)
	}
	if v, ok := src.Amount.Get(); ok {
		dst.Amount = optional.OfFloat64(v)
	}
	if v, ok := src.AvailableMargin.Get(); ok {
		dst.AvailableMargin = optional.OfFloat64(v)
	}
	if v, ok := src.Commission.Get(); ok {
		dst.Commission = optional.OfFloat64(v)
	}
	if v, ok := src.ConfirmedDebit.Get(); ok {
		dst.ConfirmedDebit = optional.OfFloat64(v)
	}
	if v, ok := src.ExcessMargin.Get(); ok {
		dst.ExcessMargin = optional.OfFloat64(v)
	}
	if v, ok := src.ExcessMarginPcnt.Get(); ok {
		dst.ExcessMarginPcnt = optional.OfFloat64(v)
	}
	if v, ok := src.GrossComm.Get(); ok {
		dst.GrossComm = optional.OfFloat64(v)
	}
	if v, ok := src.GrossExecCost.Get(); ok {
		dst.GrossExecCost = optional.OfFloat64(v)
	}
	if v, ok := src.GrossLastValue.Get(); ok {
		dst.GrossLastValue = optional.OfFloat64(v)
	}
	if v, ok := src.GrossMarkValue.Get(); ok {
		dst.GrossMarkValue = optional.OfFloat64(v)
	}
	if v, ok := src.GrossOpenCost.Get(); ok {
		dst.GrossOpenCost = optional.OfFloat64(v)
	}
	if v, ok := src.GrossOpenPremium.Get(); ok {
		dst.GrossOpenPremium = optional.OfFloat64(v)
	}
	if v, ok := src.IndicativeTax.Get(); ok {
		dst.IndicativeTax = optional.OfFloat64(v)
	}
	if v, ok := src.InitMargin.Get(); ok {
		dst.InitMargin = optional.OfFloat64(v)
	}
	if v, ok := src.MaintMargin.Get(); ok {
		dst.MaintMargin = optional.OfFloat64(v)
	}
	if v, ok := src.MakerFeeDiscount.Get(); ok {
		dst.MakerFeeDiscount = optional.OfFloat64(v)
	}
	if v, ok := src.MarginBalance.Get(); ok {
		dst.MarginBalance = optional.OfFloat64(v)
	}
	if v, ok := src.MarginBalancePcnt.Get(); ok {
		dst.MarginBalancePcnt = optional.OfFloat64(v)
	}
	if v, ok := src.MarginLeverage.Get(); ok {
		dst.MarginLeverage = optional.OfFloat64(v)
	}
	if v, ok := src.MarginUsedPcnt.Get(); ok {
		dst.MarginUsedPcnt = optional.OfFloat64(v)
	}
	if v, ok := src.PendingCredit.Get(); ok {
		dst.PendingCredit = optional.OfFloat64(v)
	}
	if v, ok := src.PendingDebit.Get(); ok {
		dst.PendingDebit = optional.OfFloat64(v)
	}
	if v, ok := src.PrevRealisedPnl.Get(); ok {
		dst.PrevRealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.PrevState.Get(); ok {
		dst.PrevState = optional.OfString(v)
	}
	if v, ok := src.PrevUnrealisedPnl.Get(); ok {
		dst.PrevUnrealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.RealisedPnl.Get(); ok {
		dst.RealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.RiskLimit.Get(); ok {
		dst.RiskLimit = optional.OfFloat64(v)
	}
	if v, ok := src.RiskValue.Get(); ok {
		dst.RiskValue = optional.OfFloat64(v)
	}
	if v, ok := src.SessionMargin.Get(); ok {
		dst.SessionMargin = optional.OfFloat64(v)
	}
	if v, ok := src.State.Get(); ok {
		dst.State = optional.OfString(v)
	}
	if v, ok := src.SyntheticMargin.Get(); ok {
		dst.SyntheticMargin = optional.OfFloat64(v)
	}
	if v, ok := src.TakerFeeDiscount.Get(); ok {
		dst.TakerFeeDiscount = optional.OfFloat64(v)
	}
	if v, ok := src.TargetExcessMargin.Get(); ok {
		dst.TargetExcessMargin = optional.OfFloat64(v)
	}
	if v, ok := src.TaxableMargin.Get(); ok {
		dst.TaxableMargin = optional.OfFloat64(v)
	}
	if v, ok := src.Timestamp.Get(); ok {
		dst.Timestamp = optional.OfTime(v)
	}
	if v, ok := src.UnrealisedPnl.Get(); ok {
		dst.UnrealisedPnl = optional.OfFloat64(v)
	}
	if v, ok := src.UnrealisedProfit.Get(); ok {
		dst.UnrealisedProfit = optional.OfFloat64(v)
	}
	if v, ok := src.VarMargin.Get(); ok {
		dst.VarMargin = optional.OfFloat64(v)
	}
	if v, ok := src.WalletBalance.Get(); ok {
		dst.WalletBalance = optional.OfFloat64(v)
	}
	if v, ok := src.WithdrawableMargin.Get(); ok {
		dst.WithdrawableMargin = optional.OfFloat64(v)
	}
}
