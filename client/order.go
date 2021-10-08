package client

import "github.com/bitmex-mirror/order"

func (c *Client) GetOrders(req *order.ReqToGetOrders) (*order.RespToGetOrders, error) {
	response := new(order.RespToGetOrders)
	err := c.request(req, response)
	return response, err

}

func (c *Client) PlaceOrder(req *order.ReqToPlaceOrder) (*order.RespToPlaceOrders, error) {
	response := new(order.RespToPlaceOrders)
	err := c.request(req, response)
	return response, err
}

func (c *Client) PlaceBuckOrders(req *order.ReqToPlaceBulkOrder) (*order.RespToPlaceBulkOrders, error) {
	response := new(order.RespToPlaceBulkOrders)
	err := c.request(req, response)
	return response, err
}

func (c *Client) CancelOrders(req *order.ReqToCancelOrders) (*order.RespToCancelOrders, error) {
	response := new(order.RespToCancelOrders)
	err := c.request(req, response)
	return response, err
}

func (c *Client) CancelAllOrders(req *order.ReqToCancelAllOrders) (*order.RespToCancelAllOrders, error) {
	response := new(order.RespToCancelAllOrders)
	err := c.request(req, response)
	return response, err
}

func (c *Client) AmendOrder(req *order.ReqToAmendOrder) (*order.RespToAmendOrder, error) {
	response := new(order.RespToAmendOrder)
	err := c.request(req, response)
	return response, err
}

func (c *Client) AmendBulkOrders(req *order.ReqToAmendBulkOrders) (*order.RespToAmendBulkOrders, error) {
	response := new(order.RespToAmendBulkOrders)
	err := c.request(req, response)
	return response, err
}

func (c *Client) CancelAllAfter(req *order.ReqToCancelAllAfter) (*order.RespToCancelAllAfter, error) {
	response := new(order.RespToCancelAllAfter)
	err := c.request(req, response)
	return response, err
}
