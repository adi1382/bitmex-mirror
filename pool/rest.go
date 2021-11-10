package pool

import (
	"context"

	"github.com/bitmex-mirror/bitmex"
	"github.com/bitmex-mirror/optional"
)

// PlaceOrder places an order and updates the state of activeOrders and positions accordingly.
func (a *account) PlaceOrder(ctx context.Context, req bitmex.ReqToPlaceOrder) error {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()
	a.positionsMu.Lock()
	defer a.ordersMu.Unlock()

	orderRes, err := a.restClient.PlaceOrder(ctx, req)
	if err != nil {
		return err
	}

	order := bitmex.Order(orderRes)

	// Update active orders
	if order.OrdStatus.ElseZero() == bitmex.OrderStatusNew ||
		order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
		a.activeOrders = append(a.activeOrders, order)
	}

	// Update open positions
	if order.OrdStatus.ElseZero() == bitmex.OrderStatusFilled ||
		order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {

		for i := range a.positions {
			if a.positions[i].Symbol == order.Symbol.ElseZero() && a.positions[i].Account == order.Account {
				switch order.Side.ElseZero() {
				// The order may be filled or partially filled, if not then orderQty will be equal to leavesQty
				// orderQty and leavesQty are always equal for a new order
				// orderQty is the total quantity of the order
				// leavesQty is the quantity left to fill
				// orderQty and leavesQty are always positive
				case bitmex.OrderSideBuy:
					a.positions[i].CurrentQty = optional.OfFloat64(a.positions[i].CurrentQty.ElseZero() +
						(order.OrderQty.ElseZero() - order.LeavesQty.ElseZero()))
				case bitmex.OrderSideSell:
					a.positions[i].CurrentQty = optional.OfFloat64(a.positions[i].CurrentQty.ElseZero() -
						(order.OrderQty.ElseZero() - order.LeavesQty.ElseZero()))
				}
				return nil
			}
		}

		// If we didn't find the position, create a new one
		var position bitmex.Position
		position.Symbol = order.Symbol.ElseZero()
		position.Account = order.Account

		switch order.Side.ElseZero() {
		case bitmex.OrderSideBuy:
			position.CurrentQty =
				optional.OfFloat64(order.OrderQty.ElseZero() - order.LeavesQty.ElseZero())
		case bitmex.OrderSideSell:
			position.CurrentQty =
				optional.OfFloat64(-(order.OrderQty.ElseZero() - order.LeavesQty.ElseZero()))
		}

		// Add the new position to the list
		a.positions = append(a.positions, position)
	}

	return nil
}

func (a *account) AmendOrder(ctx context.Context, req bitmex.ReqToAmendOrder) error {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()
	a.positionsMu.Lock()
	defer a.ordersMu.Unlock()

	orderRes, err := a.restClient.AmendOrder(ctx, req)
	if err != nil {
		return err
	}

	order := bitmex.Order(orderRes)

	// This variable will be used to track the change in cumulated quantity(orderQty-leavesQty)
	// for the amended order.
	// If the order is filled, cumulated quantity will be equal to orderQty
	// If the order is partially filled, cumulated quantity will be equal to orderQty-leavesQty
	// If the order is not filled, cumulated quantity will be equal to 0
	// Old cumulated quantity of the order is stored before the amendment in variable oldCumulatedQty
	// Note: Side cannot be changed in an amendment
	oldCumulatedQty := 0.0
	newCumulatedQty := 0.0

	isOrderExisting := false
	for i := range a.activeOrders {
		if a.activeOrders[i].OrderID == order.OrderID {
			isOrderExisting = true
			// assigning old and new cumulated quantity
			oldCumulatedQty = a.activeOrders[i].OrderQty.ElseZero() - a.activeOrders[i].LeavesQty.ElseZero()
			newCumulatedQty = order.OrderQty.ElseZero() - order.LeavesQty.ElseZero()
			// Update the order if it's active
			if order.OrdStatus.ElseZero() == bitmex.OrderStatusNew ||
				order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
				a.activeOrders[i] = order
				break
			}

			// Remove the order if it's filled or canceled
			if order.OrdStatus.ElseZero() == bitmex.OrderStatusCanceled ||
				order.OrdStatus.ElseZero() == bitmex.OrderStatusFilled {
				a.activeOrders[i] = a.activeOrders[len(a.activeOrders)-1]
				a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				break
			}
		}
	}

	// If the order wasn't found, it's a new order
	// this should never happen, but just in case
	if !isOrderExisting {
		a.activeOrders = append(a.activeOrders, order)
	}

	if order.OrdStatus.ElseZero() == bitmex.OrderStatusFilled ||
		order.OrdStatus.ElseZero() == bitmex.OrderStatusPartiallyFilled {
		// Update open positions

		if oldCumulatedQty != newCumulatedQty {
			for i := range a.positions {
				if a.positions[i].Symbol == order.Symbol.ElseZero() && a.positions[i].Account == order.Account {
					switch order.Side.ElseZero() {
					case bitmex.OrderSideBuy:
						a.positions[i].CurrentQty =
							optional.OfFloat64(a.positions[i].CurrentQty.ElseZero() + (newCumulatedQty - oldCumulatedQty))
					case bitmex.OrderSideSell:
						a.positions[i].CurrentQty =
							optional.OfFloat64(a.positions[i].CurrentQty.ElseZero() - (newCumulatedQty - oldCumulatedQty))
					}
					return nil
				}
			}
		}

		// If we didn't find the position, create a new one
		var position bitmex.Position
		position.Symbol = order.Symbol.ElseZero()
		position.Account = order.Account

		switch order.Side.ElseZero() {
		case bitmex.OrderSideBuy:
			position.CurrentQty = optional.OfFloat64(newCumulatedQty)
		case bitmex.OrderSideSell:
			position.CurrentQty = optional.OfFloat64(-newCumulatedQty)
		}

		// Add the new position to the list
		a.positions = append(a.positions, position)
	}

	return nil
}

func (a *account) CancelOrders(ctx context.Context, req bitmex.ReqToCancelOrders) error {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	orderRes, err := a.restClient.CancelOrders(ctx, req)
	if err != nil {
		return err
	}

	orders := []bitmex.Order(orderRes)

	// remove the orders from the active orders which are canceled
	for i := range orders {
		for j := range a.activeOrders {
			if a.activeOrders[j].OrderID == orders[i].OrderID {
				a.activeOrders[j] = a.activeOrders[len(a.activeOrders)-1]
				a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				break
			}
		}
	}

	// open positions will remain the same

	return nil
}

func (a *account) CancelAllOrders(ctx context.Context, req bitmex.ReqToCancelAllOrders) error {
	a.ordersMu.Lock()
	defer a.ordersMu.Unlock()

	orderRes, err := a.restClient.CancelAllOrders(ctx, req)
	if err != nil {
		return err
	}

	orders := []bitmex.Order(orderRes)

	// remove the orders from the active orders which are canceled
	for i := range orders {
		for j := range a.activeOrders {
			if a.activeOrders[j].OrderID == orders[i].OrderID {
				a.activeOrders[j] = a.activeOrders[len(a.activeOrders)-1]
				a.activeOrders = a.activeOrders[:len(a.activeOrders)-1]
				break
			}
		}
	}

	// open positions will remain the same

	return nil
}

func (a *account) ChangeLeverage(ctx context.Context, req bitmex.ReqToChangeLeverage) error {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	res, err := a.restClient.ChangeLeverage(ctx, req)
	if err != nil {
		return err
	}

	position := bitmex.Position(res)

	for i := range a.positions {
		if a.positions[i].Symbol == position.Symbol && a.positions[i].Account == position.Account {
			a.positions[i] = position
			return nil
		}
	}

	// If we didn't find the position, create a new one
	a.positions = append(a.positions, position)

	return nil
}

func (a *account) IsolateMargin(ctx context.Context, req bitmex.ReqToIsolateMargin) error {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	res, err := a.restClient.IsolateMargin(ctx, req)
	if err != nil {
		return err
	}

	position := bitmex.Position(res)

	for i := range a.positions {
		if a.positions[i].Symbol == position.Symbol && a.positions[i].Account == position.Account {
			a.positions[i] = position
			return nil
		}
	}

	// If we didn't find the position, create a new one
	a.positions = append(a.positions, position)

	return nil
}

func (a *account) ChangeRiskLimit(ctx context.Context, req bitmex.ReqToChangeRiskLimit) error {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	res, err := a.restClient.ChangeRiskLimit(ctx, req)
	if err != nil {
		return err
	}

	position := bitmex.Position(res)

	for i := range a.positions {
		if a.positions[i].Symbol == position.Symbol && a.positions[i].Account == position.Account {
			a.positions[i] = position
			return nil
		}
	}

	// If we didn't find the position, create a new one
	a.positions = append(a.positions, position)

	return nil
}

func (a *account) TransferMargin(ctx context.Context, req bitmex.ReqToTransferMargin) error {
	a.positionsMu.Lock()
	defer a.positionsMu.Unlock()

	res, err := a.restClient.TransferMargin(ctx, req)
	if err != nil {
		return err
	}

	position := bitmex.Position(res)

	for i := range a.positions {
		if a.positions[i].Symbol == position.Symbol && a.positions[i].Account == position.Account {
			a.positions[i] = position
			return nil
		}
	}

	// If we didn't find the position, create a new one
	a.positions = append(a.positions, position)

	return nil
}
