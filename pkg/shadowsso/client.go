package shadowsso

type Client struct{}

// NewClient creates a new client.
func NewClient() *Client {
	c := new(Client)

	// TODO: set remote address and the authenticator

	return c
}
