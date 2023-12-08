package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"sqlit-lessonTEST/lib/mistake"
)

const (
	MethodgetUpdate   = "getUpdates"
	MethodsendMessage = "sendMessage"
)

func (c *Client) Update(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.request(MethodgetUpdate, q)
	if err != nil {
		return nil, fmt.Errorf("wrong MethodgetUpdate: %w", err)
	}
	var res UpdateResponse

	if err := json.Unmarshal(data, res); err != nil {
		fmt.Errorf("wrong in Unmarshal: %w", err)
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, UrlPage string) error {
	const e = "can't send message"
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("UrlPage", UrlPage)
	_, err := c.request(MethodsendMessage, q)
	if err != nil {
		return mistake.WrapErr(e, err)
	}
	return nil

}

func (c *Client) request(method string, query url.Values) ([]byte, error) {
	const msgErr = "no request"
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("no request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, mistake.WrapErr(msgErr, err)
	}
	resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("no body read: %w", err)
	}

	return body, nil
}
