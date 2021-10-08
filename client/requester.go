package rest

type requester interface {
	Path() string
	Method() string
	Query() string
	Payload() string
	isSigned() bool
}
