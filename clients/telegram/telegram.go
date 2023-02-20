package telegram

import "net/http"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: Token(token),
		client:   http.Client{},
	}
}

func Token(token string) string {
	return "bot" + token
}

// Updates gets messages
func (c *Client) Updates() {

}

func (c *Client) SendMessage() {

}
