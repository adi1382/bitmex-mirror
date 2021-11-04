package bitmex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func (c *restClient) request(req Requester, results interface{}) error {
	res, err := c.do(req)

	defer fasthttp.ReleaseResponse(res)

	if err != nil {
		return err
	}

	if err := decode(res, results); err != nil {
		return err
	}
	return nil
}

func (c *restClient) newRequest(r Requester) (*fasthttp.Request, error) {
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

	uri := fasthttp.AcquireURI() // Acquiring URL from Pool

	uri.SetHost(c.endpoint)
	uri.SetQueryString(query)
	uri.SetPath(path)
	uri.SetScheme("https")

	fmt.Println("Calculated: ", uri.String())

	req := fasthttp.AcquireRequest() // Acquiring Request from the Pool, released in do()

	req.Header.SetMethod(method)
	req.SetRequestURI(uri.String())

	fasthttp.ReleaseURI(uri) // Releasing URL back to Pool

	req.SetBodyString(payload)

	apiExpires := strconv.FormatInt(time.Now().Add(time.Second*60).Unix(), 10)
	req.Header.Set("Api-Expires", apiExpires)
	req.Header.Set("Api-Key", c.auth.Key)
	//TODO: ASSIGN USER AGENT
	//req.Header.Set("User-Agent", "Bitmex")
	req.Header.Set("Api-Signature", c.prepareSignature(path, method, query, payload, apiExpires))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *restClient) do(r Requester) (*fasthttp.Response, error) {
	req, err := c.newRequest(r)
	if err != nil {
		return nil, err
	}

	res := fasthttp.AcquireResponse() // Acquiring Response from Pool

	//fmt.Println("Request")
	fmt.Println(req)

	start := time.Now().UnixNano()
	err = c.httpC.DoTimeout(req, res, time.Second*5)

	fasthttp.ReleaseRequest(req) // Releasing Request back to Pool

	fmt.Println("Request Time: ", time.Now().UnixNano()-start)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}

	fmt.Println("Rate limit available: ", string(res.Header.Peek("X-Ratelimit-Remaining")), time.Now())
	//fmt.Println("Available in bucket: ", c.bucketM.Available())
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