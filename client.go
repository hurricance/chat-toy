package main

import (
	"net"
)

type Client struct {
	Name string
	Addr string
	Channel chan string
	coon net.Conn
}

func NewClient(coon net.Conn) *Client {
	clientAddr := coon.RemoteAddr().String()
	client := &Client {
		clientAddr,
		clientAddr,
		make(chan string),
		coon,
	}

	go client.ListenMessage()
	return client
}

func (client *Client) ListenMessage() {
	for {
		msg, ok := <- client.Channel
		if !ok {
			return
		}
		client.coon.Write([]byte(msg + "\n"))
	}
}
