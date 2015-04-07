package put.io

import (
    "net/http"

type Client struct {
    *http.Client
    Token   string
}

func (c *Client) 
