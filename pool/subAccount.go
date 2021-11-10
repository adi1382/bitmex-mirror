package pool

type subAccount struct {
	hostReceiver <-chan []byte
	*account
}
