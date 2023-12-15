package client

type Client interface {
	Send(target string, path string) (uint64, error)
}
