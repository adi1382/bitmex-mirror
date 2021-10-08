package client

type Requester interface {
	Path() string
	Method() string
	Query() string
	Payload() string
}
