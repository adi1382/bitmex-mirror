package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

func (c *Client) request(req Requester, results interface{}) error {
	res, err := c.do(req)
	if err != nil {
		return err
	}

	if err := decode(res, results); err != nil {
		return err
	}
	return nil
}

func (c *Client) newRequest(r Requester) *fasthttp.Request {
	// avoid Pointer's butting
	u, _ := url.ParseRequestURI(c.auth.Endpoint)
	u.Path = u.Path + r.Path()
	u.RawQuery = r.Query()

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(r.Method())
	req.SetRequestURI(u.String())
	body := r.Payload()
	req.SetBody([]byte(body))

	apiExpires := strconv.FormatInt(time.Now().Add(time.Second*60).Unix(), 10)
	req.Header.Set("Api-Expires", apiExpires)
	req.Header.Set("Api-Key", c.auth.Key)
	//TODO: ASSIGN USER AGENT
	//req.Header.Set("User-Agent", "bitmexgo")
	req.Header.Set("Api-Signature", c.PrepareSignature(r, apiExpires))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (c *Client) do(r Requester) (*fasthttp.Response, error) {
	req := c.newRequest(r)
	res := fasthttp.AcquireResponse()

	err := c.httpC.DoTimeout(req, res, time.Second*5)
	if err != nil {
		return nil, errors.Wrap(err, "api request failed")
	}

	fmt.Println(string(res.Body()))

	if res.StatusCode() >= 400 {
		apiError := new(APIError)
		err := json.Unmarshal(res.Body(), apiError)
		fmt.Println(apiError)
		apiError.StatusCode = res.StatusCode()
		fmt.Println(apiError)
		if err != nil {
			fmt.Println("decode failed")
			return res, errors.Wrap(err, "api error decode failed")
		}
		return res, apiError
	}

	return res, nil
}

func decode(res *fasthttp.Response, out interface{}) error {

	if err := json.Unmarshal(res.Body(), out); err != nil {
		return nil
	} else {
		return errors.Wrap(err, "response decode error")
	}
}

func (c *Client) PrepareSignature(r Requester, apiExpires string) string {
	signatureBody := r.Method() + "/api/v1" + r.Path()

	if r.Query() != "" {
		signatureBody += "?" + r.Query()
	}

	signatureBody += apiExpires + r.Payload()

	return c.auth.Sign(signatureBody)
}
