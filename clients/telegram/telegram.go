package telegram

import (
	"net/http"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newbasePath(token),
		client:   http.Client{},
	}

}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int           `json:"update_id"`
	Message *InboxMessage `json:"message"`
}

type InboxMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}
type Chat struct {
	Id int `json:"id"`
}

func newbasePath(token string) string {
	return "bot" + token
}
