package bitmex

type Requester interface {
	path() string
	method() string
	query() (string, error)
	payload() (string, error)
}
