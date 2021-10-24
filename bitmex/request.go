package bitmex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func (c *restClient) request(req Requester, results interface{}) error {
	res, err := c.do(req)
	if err != nil {
		return err
	}

	if err := decode(res, results); err != nil {
		return err
	}
	return nil
}

func (c *restClient) newRequest(r Requester) (*fasthttp.Request, error) {
	u, _ := url.ParseRequestURI(c.endpoint)

	method := r.method()
	path := r.path()
	query, err := r.query()
	if err != nil {
		return nil, errors.Wrap(err, "query prepare error")
	}
	payload, err := r.payload()
	if err != nil {
		return nil, errors.Wrap(err, "payload prepare error")
	}

	u.Path = u.Path + path

	u.RawQuery = query

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(u.String())

	req.SetBody([]byte(payload))

	apiExpires := strconv.FormatInt(time.Now().Add(time.Second*60).Unix(), 10)
	req.Header.Set("Api-Expires", apiExpires)
	req.Header.Set("Api-Key", c.auth.Key)
	//TODO: ASSIGN USER AGENT
	req.Header.Set("User-Agent", "Bitmex")
	req.Header.Set("Api-Signature", c.prepareSignature(path, method, query, payload, apiExpires))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *restClient) do(r Requester) (*fasthttp.Response, error) {
	req, err := c.newRequest(r)
	if err != nil {
		return nil, err
	}

	res := fasthttp.AcquireResponse()

	//fmt.Println("Request")
	//fmt.Println(req)

	start := time.Now().UnixNano()
	err = c.httpC.DoTimeout(req, res, time.Second*5)
	fmt.Println("Request Time: ", time.Now().UnixNano()-start)
	if err != nil {
		return nil, errors.Wrap(err, "api request failed")
	}

	c.bucket1m.Take(1)

	//fmt.Println(string(res.Header.Peek("X-Ratelimit-Remaining")), time.Now())
	//fmt.Println("Available in bucket: ", c.bucket1m.Available())
	//
	//fmt.Println(res.Header.String())

	//fmt.Println(string(res.Body()))

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

func (c *restClient) prepareSignature(path, method, query, payload, apiExpires string) string {
	signatureBody := method + "/api/v1" + path

	if query != "" {
		signatureBody += "?" + query
	}

	signatureBody += apiExpires + payload

	return c.auth.Sign(signatureBody)
}
