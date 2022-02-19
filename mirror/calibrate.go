package mirror

import (
	"context"
	"github.com/bitmex-mirror/bitmex"
	"github.com/bitmex-mirror/optional"
	"github.com/pkg/errors"
	"log"
	"math"
	"sync"
)

func (m *mirror) calibrate(ctx context.Context, logger *log.Logger) {
	hostOrders := m.host.OrdersCopy()
	subOrders := m.sub.OrdersCopy()
	hostPositions := m.host.PositionsCopy()
	subPositions := m.sub.PositionsCopy()

	ordersToPlace, ordersToAmend, orderIDsToCancel := m.calibrateOrders(hostOrders, subOrders, logger)
	leverageChanges := m.calibrateLeverage(hostPositions, subPositions)
	posIncMarketOrders, posDecMarketOrders := m.calibratePositions(hostPositions, subPositions)

	err := m.sub.CancelOrders(ctx, bitmex.ReqToCancelOrders{OrderIDs: orderIDsToCancel})
	m.handleError(err, logger)

	var wgPosDecMarketOrders sync.WaitGroup
	for i := range posDecMarketOrders {
		wgPosDecMarketOrders.Add(1)
		go func(i int) {
			defer wgPosDecMarketOrders.Done()
			err := m.sub.PlaceOrder(ctx, posDecMarketOrders[i])
			m.handleError(err, logger)
		}(i)
	}
	wgPosDecMarketOrders.Wait()

	var wgLeverageChanges sync.WaitGroup
	for i := range leverageChanges {
		wgLeverageChanges.Add(1)
		go func(i int) {
			defer wgLeverageChanges.Done()
			err := m.sub.ChangeLeverage(ctx, leverageChanges[i])
			m.handleError(err, logger)
		}(i)
	}
	wgLeverageChanges.Wait()

	var wgPosIncMarketOrders sync.WaitGroup
	for i := range posIncMarketOrders {
		wgPosIncMarketOrders.Add(1)
		go func(i int) {
			defer wgPosIncMarketOrders.Done()
			err := m.sub.PlaceOrder(ctx, posIncMarketOrders[i])
			m.handleError(err, logger)
		}(i)
	}
	wgPosIncMarketOrders.Wait()

	var wgOrdersToPlace sync.WaitGroup
	for i := range ordersToPlace {
		wgOrdersToPlace.Add(1)
		go func(i int) {
			defer wgOrdersToPlace.Done()
			err := m.sub.PlaceOrder(ctx, ordersToPlace[i])
			m.handleError(err, logger)
		}(i)
	}

	var wgOrdersToAmend sync.WaitGroup
	for i := range ordersToAmend {
		wgOrdersToAmend.Add(1)
		go func(i int) {
			defer wgOrdersToAmend.Done()
			err := m.sub.AmendOrder(ctx, ordersToAmend[i])
			m.handleError(err, logger)
		}(i)
	}

	wgOrdersToPlace.Wait()
	wgOrdersToAmend.Wait()
}

func (m *mirror) handleError(err error, logger *log.Logger) {
	if err == nil {
		return
	}
	bitmexError, ok := errors.Cause(err).(bitmex.APIError)
	if !ok {
		logger.Println("ILLEGAL: A non-api error was passed to handle error")
		return
	}
	logger.Printf("ERROR on mirror %s: %s", m.subId, bitmexError.Error())

	if bitmexError.StatusCode == 401 || bitmexError.StatusCode == 403 {
		logger.Printf("ERROR: Invalid APIs on mirror: %s\n", m.subId)
		m.close("handleError: Invalid APIs")
		return
	}

}

func (m *mirror) calibratePositions(hostPositions,
	subPositions []bitmex.Position) ([]bitmex.ReqToPlaceOrder, []bitmex.ReqToPlaceOrder) {

	// Market orders that will increase the position size
	var posIncMarketOrders []bitmex.ReqToPlaceOrder
	// Market orders that will decrease the position size
	var posDecMarketOrders []bitmex.ReqToPlaceOrder

	for i := range hostPositions {
		positionFound := false
		for j := range subPositions {
			if hostPositions[i].Symbol == subPositions[j].Symbol {
				positionFound = true
				m.calibrateExistingPosition(hostPositions[i], subPositions[j], &posIncMarketOrders, &posDecMarketOrders)
				break
			}
		}
		if !positionFound {
			m.calibrateNewPosition(hostPositions[i], &posIncMarketOrders)
		}
	}

	for i := range subPositions {
		positionFound := false
		for j := range hostPositions {
			if subPositions[i].Symbol == hostPositions[j].Symbol {
				positionFound = true
				break
			}
		}
		if !positionFound && subPositions[i].CurrentQty.ElseZero() != 0 {
			// close sub position if host position does not exist for symbol
			var side string

			if subPositions[i].CurrentQty.ElseZero() > 0 {
				side = "Sell"
			} else {
				side = "Buy"
			}

			posDecMarketOrders = append(posDecMarketOrders, bitmex.ReqToPlaceOrder{
				Symbol:   subPositions[i].Symbol,
				Side:     side,
				ClOrdID:  clOrdId(randomGuid()),
				OrdTyp:   bitmex.OrderTypeMarket,
				ExecInst: bitmex.ExecInstClose,
			})
		}
	}

	return posIncMarketOrders, posDecMarketOrders
}

func (m *mirror) calibrateExistingPosition(hostPosition, subPosition bitmex.Position,
	posIncMarketOrders, posDecMarketOrders *[]bitmex.ReqToPlaceOrder) {
	balanceRatio := m.BalanceRatio(hostPosition.Currency.ElseZero())
	adjustedPositionSize := m.roundOrderSize(hostPosition.Symbol,
		hostPosition.CurrentQty.ElseZero()*balanceRatio)

	if adjustedPositionSize != subPosition.CurrentQty.ElseZero() {
		qty := adjustedPositionSize - subPosition.CurrentQty.ElseZero()
		var side string
		if qty > 0 {
			side = bitmex.OrderSideBuy
		} else {
			side = bitmex.OrderSideSell
		}

		if adjustedPositionSize == 0 {
			*posDecMarketOrders = append(*posDecMarketOrders, bitmex.ReqToPlaceOrder{
				Symbol:   hostPosition.Symbol,
				Side:     side,
				ClOrdID:  clOrdId(randomGuid()),
				OrdTyp:   bitmex.OrderTypeMarket,
				ExecInst: bitmex.ExecInstClose,
			})
		} else {
			if qty*subPosition.CurrentQty.ElseZero() > 0 ||
				(qty*subPosition.CurrentQty.ElseZero() < 0 && qty > subPosition.CurrentQty.ElseZero()) {
				*posIncMarketOrders = append(*posIncMarketOrders, bitmex.ReqToPlaceOrder{
					Symbol:   hostPosition.Symbol,
					Side:     side,
					OrderQty: optional.OfFloat64(math.Abs(qty)),
					ClOrdID:  clOrdId(randomGuid()),
					OrdTyp:   bitmex.OrderTypeMarket,
				})
			} else {
				*posDecMarketOrders = append(*posDecMarketOrders, bitmex.ReqToPlaceOrder{
					Symbol:   hostPosition.Symbol,
					Side:     side,
					OrderQty: optional.OfFloat64(math.Abs(qty)),
					ClOrdID:  clOrdId(randomGuid()),
					OrdTyp:   bitmex.OrderTypeMarket,
				})
			}
		}
	}
}

func (m *mirror) calibrateNewPosition(hostPosition bitmex.Position, marketOrders *[]bitmex.ReqToPlaceOrder) {
	balanceRatio := m.BalanceRatio(hostPosition.Currency.ElseZero())
	adjustedPositionSize := m.roundOrderSize(hostPosition.Symbol,
		hostPosition.CurrentQty.ElseZero()*balanceRatio)

	var side string
	if adjustedPositionSize > 0 {
		side = bitmex.OrderSideBuy
	} else {
		side = bitmex.OrderSideSell
	}
	qty := math.Abs(adjustedPositionSize)

	*marketOrders = append(*marketOrders, bitmex.ReqToPlaceOrder{
		Symbol:   hostPosition.Symbol,
		Side:     side,
		OrderQty: optional.OfFloat64(qty),
		ClOrdID:  clOrdId(randomGuid()),
		OrdTyp:   bitmex.OrderTypeMarket,
	})
}

// calibrateOrders compares the orders on the host and the sub, and creates place, amends, cancels order requests
// and append these requests to the respective slices.
// It does not call any rest API function, just create requests and modifies slices.
func (m *mirror) calibrateOrders(hostOrders,
	subOrders []bitmex.Order, logger *log.Logger) ([]bitmex.ReqToPlaceOrder, []bitmex.ReqToAmendOrder, []string) {

	var ordersToPlace []bitmex.ReqToPlaceOrder
	var ordersToAmend []bitmex.ReqToAmendOrder
	var orderIDsToCancel []string

	for h := range hostOrders {
		orderFound := false
		// searching through subOrders
		for s := range subOrders {
			// if the order is found check for any amendments
			if matchIDs(hostOrders[h].OrderID, subOrders[s].ClOrdID.ElseZero()) {

				if !orderFound {
					// sanity check if the order is desirable or forged by some other application
					// sanity check note: An Extreme Case Scenario:
					// If the host order is already triggered before being copied by sub, then sub will copy a simple
					// limit order, in that situation OrdType, ExecInst, PegPriceType, etc. could be different.
					// Therefore, a sanity check on symbol and side should provide sufficient level of protection.
					if hostOrders[h].Symbol != subOrders[s].Symbol || hostOrders[h].Side != subOrders[s].Side {
						logger.Printf("WARNING: Non-Amendable values of the order do not match")
						orderIDsToCancel = append(orderIDsToCancel, subOrders[s].OrderID)
						continue
					}

					// if the order is found and the values match, check for any amendments
					orderFound = true
					m.calibrateAmendOrders(hostOrders[h], subOrders[s], &ordersToAmend)
				} else {
					// there might be multiple sub orders for single host order
					// in such case, we need to cancel duplicate sub orders
					orderIDsToCancel = append(orderIDsToCancel, subOrders[s].OrderID)
				}
			}
		}

		// If order is not found, place a new order
		if !orderFound {
			m.calibratePlaceOrders(hostOrders[h], &ordersToPlace)
		}
	}

	// delete sub orders for which host order does not exist
	for i := range subOrders {
		orderFound := false
		for j := range hostOrders {
			if matchIDs(hostOrders[j].OrderID, subOrders[i].ClOrdID.ElseZero()) {
				orderFound = true
				break
			}
		}
		if !orderFound {
			orderIDsToCancel = append(orderIDsToCancel, subOrders[i].OrderID)
		}
	}
	return ordersToPlace, ordersToAmend, orderIDsToCancel
}

// calibratePlaceOrders appends a new order to the ordersToPlace slice if needed
func (m *mirror) calibratePlaceOrders(hostOrder bitmex.Order, ordersToPlace *[]bitmex.ReqToPlaceOrder) {
	balanceRatio := m.BalanceRatio(hostOrder.Currency.ElseZero())
	adjustedOrderSize := m.roundOrderSize(hostOrder.Symbol.ElseZero(), hostOrder.LeavesQty.ElseZero()*balanceRatio)

	if adjustedOrderSize == 0 {
		return
	}

	var placeOrder bitmex.ReqToPlaceOrder

	if (hostOrder.OrdType.ElseZero() == bitmex.OrderTypeStopLimit ||
		hostOrder.OrdType.ElseZero() == bitmex.OrderTypeLimitIfTouched) &&
		hostOrder.Triggered.ElseZero() != "" {

		placeOrder.ClOrdID = clOrdId(hostOrder.OrderID)
		placeOrder.Symbol = hostOrder.Symbol.ElseZero()
		placeOrder.Side = hostOrder.Side.ElseZero()
		placeOrder.TimeInForce = hostOrder.TimeInForce.ElseZero()
		placeOrder.ExecInst = hostOrder.ExecInst.ElseZero()
		placeOrder.Text = hostOrder.Text.ElseZero()

		if v, ok := hostOrder.LeavesQty.Get(); ok {
			placeOrder.OrderQty = optional.OfFloat64(v)
		}

		if v, ok := hostOrder.Price.Get(); ok {
			placeOrder.Price = optional.OfFloat64(v)
		}

		if v, ok := hostOrder.DisplayQty.Get(); ok {
			placeOrder.DisplayQty = optional.OfFloat64(v)
		}

		*ordersToPlace = append(*ordersToPlace, placeOrder)
		return

	}

	placeOrder.ClOrdID = clOrdId(hostOrder.OrderID)
	placeOrder.Symbol = hostOrder.Symbol.ElseZero()
	placeOrder.Side = hostOrder.Side.ElseZero()
	placeOrder.OrdTyp = hostOrder.OrdType.ElseZero()
	placeOrder.TimeInForce = hostOrder.TimeInForce.ElseZero()
	placeOrder.ExecInst = hostOrder.ExecInst.ElseZero()
	placeOrder.Text = hostOrder.Text.ElseZero()
	placeOrder.PegPriceType = hostOrder.PegPriceType.ElseZero()

	if v, ok := hostOrder.LeavesQty.Get(); ok {
		placeOrder.OrderQty = optional.OfFloat64(v)
	}

	if v, ok := hostOrder.Price.Get(); ok {
		placeOrder.Price = optional.OfFloat64(v)
	}

	if v, ok := hostOrder.DisplayQty.Get(); ok {
		placeOrder.DisplayQty = optional.OfFloat64(v)
	}

	if v, ok := hostOrder.StopPx.Get(); ok {
		placeOrder.StopPx = optional.OfFloat64(v)
	}

	if v, ok := hostOrder.PegOffsetValue.Get(); ok {
		placeOrder.PegOffsetValue = optional.OfFloat64(v)
	}

	if v, ok := hostOrder.DisplayQty.Get(); ok {
		placeOrder.DisplayQty = optional.OfFloat64(v)
	}

	*ordersToPlace = append(*ordersToPlace, placeOrder)

}

// calibrateAmendOrders appends amend requests to the ordersToAmend slice if there are any amendment required
func (m *mirror) calibrateAmendOrders(hostOrder, subOrder bitmex.Order,
	ordersToAmend *[]bitmex.ReqToAmendOrder) {

	balanceRatio := m.BalanceRatio(hostOrder.Currency.ElseZero())
	adjustedOrderSize := m.roundOrderSize(hostOrder.Symbol.ElseZero(), hostOrder.LeavesQty.ElseZero()*balanceRatio)

	// If the found order is already triggered on host, it does not guarantee that the order is triggered on sub.
	// Leaves Qty must match as per the proportional formula, as positions calibrated separately.
	if (hostOrder.OrdType.ElseZero() == bitmex.OrderTypeStopLimit ||
		hostOrder.OrdType.ElseZero() == bitmex.OrderTypeLimitIfTouched) &&
		hostOrder.Triggered.ElseZero() != "" {

		// After trigger order will be converted to Limit order so theoretically only price or quantity can be changed
		if hostOrder.Price.ElseZero() != subOrder.Price.ElseZero() ||
			adjustedOrderSize != subOrder.LeavesQty.ElseZero() {

			var amendOrder bitmex.ReqToAmendOrder

			amendOrder.OrderID = subOrder.OrderID

			// if host price does not match sub price, amend sub price
			if hostOrder.Price.ElseZero() != subOrder.Price.ElseZero() {
				amendOrder.Price = hostOrder.Price
			}

			// Amend LeavesQty if partially filled else amend OrderQty
			if hostOrder.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
				amendOrder.LeavesQty = optional.OfFloat64(adjustedOrderSize)
			} else {
				amendOrder.OrderQty = optional.OfFloat64(adjustedOrderSize)
			}
			*ordersToAmend = append(*ordersToAmend, amendOrder)
		}
		return
	}

	if hostOrder.Price.ElseZero() != subOrder.Price.ElseZero() ||
		adjustedOrderSize != subOrder.LeavesQty.ElseZero() ||
		hostOrder.StopPx.ElseZero() != subOrder.StopPx.ElseZero() ||
		hostOrder.PegOffsetValue.ElseZero() != subOrder.PegOffsetValue.ElseZero() {

		var amendOrder bitmex.ReqToAmendOrder

		amendOrder.OrderID = subOrder.OrderID
		if subOrder.Price.Present() {
			amendOrder.Price = hostOrder.Price
		}
		if subOrder.StopPx.Present() {
			amendOrder.StopPx = hostOrder.StopPx
		}
		if subOrder.PegOffsetValue.Present() {
			amendOrder.PegOffsetValue = hostOrder.PegOffsetValue
		}
		if subOrder.LeavesQty.Present() {
			if subOrder.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
				amendOrder.LeavesQty = optional.OfFloat64(adjustedOrderSize)
			} else {
				amendOrder.OrderQty = optional.OfFloat64(adjustedOrderSize)
			}
		}
		*ordersToAmend = append(*ordersToAmend, amendOrder)
	}
	return
}

// calibrateLeverage calibrates all leverage of the host with sub. It calls rest api to change the leverage of the sub
func (m *mirror) calibrateLeverage(hostPositions, subPositions []bitmex.Position) []bitmex.ReqToChangeLeverage {

	var leverageChanges []bitmex.ReqToChangeLeverage

	for i := range hostPositions {
		positionFound := false
		for j := range subPositions {
			if hostPositions[i].Symbol == subPositions[j].Symbol {
				positionFound = true
				if hostPositions[i].CrossMargin.ElseZero() && !subPositions[j].CrossMargin.ElseZero() {
					leverageChanges = append(leverageChanges, bitmex.ReqToChangeLeverage{
						Symbol:   hostPositions[i].Symbol,
						Leverage: optional.OfFloat64(0),
					})
				} else if !hostPositions[i].CrossMargin.ElseZero() && subPositions[j].CrossMargin.ElseZero() {
					leverageChanges = append(leverageChanges, bitmex.ReqToChangeLeverage{
						Symbol:   hostPositions[i].Symbol,
						Leverage: optional.OfFloat64(hostPositions[i].Leverage.ElseZero()),
					})
				} else if hostPositions[i].Leverage.ElseZero() != subPositions[j].Leverage.ElseZero() {
					leverageChanges = append(leverageChanges, bitmex.ReqToChangeLeverage{
						Symbol:   hostPositions[i].Symbol,
						Leverage: optional.OfFloat64(hostPositions[i].Leverage.ElseZero()),
					})
				}
			}
		}

		if !positionFound {
			if hostPositions[i].CrossMargin.ElseZero() {
				leverageChanges = append(leverageChanges, bitmex.ReqToChangeLeverage{
					Symbol:   hostPositions[i].Symbol,
					Leverage: optional.OfFloat64(0),
				})
			} else {
				leverageChanges = append(leverageChanges, bitmex.ReqToChangeLeverage{
					Symbol:   hostPositions[i].Symbol,
					Leverage: optional.OfFloat64(hostPositions[i].Leverage.ElseZero()),
				})
			}
		}
	}

	return leverageChanges

}
