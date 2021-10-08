package client

import (
	"github.com/bitmex-mirror/position"
)

func (c *Client) GetPositions(req *position.ReqToGetPositions) (*position.RespToGetPositions, error) {
	response := new(position.RespToGetPositions)
	err := c.request(req, response)
	return response, err
}

func (c *Client) ChangeLeverage(req *position.ReqToChangeLeverage) (*position.RespToChangeLeverage, error) {
	response := new(position.RespToChangeLeverage)
	err := c.request(req, response)
	return response, err
}

func (c *Client) ChangeRiskLimit(req *position.ReqToChangeRiskLimit) (*position.RespToChangeRiskLimit, error) {
	response := new(position.RespToChangeRiskLimit)
	err := c.request(req, response)
	return response, err
}

func (c *Client) IsolateMargin(req *position.ReqToIsolateMargin) (*position.RespToIsolateMargin, error) {
	response := new(position.RespToIsolateMargin)
	err := c.request(req, response)
	return response, err
}
