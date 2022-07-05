package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type CmdHandler func([]string)
type Client struct {
	IsConnected      bool
	message_handlers []CmdHandler
	host, port       string
	conn             net.Conn
}

func ClientNew(host, port string, handlers []CmdHandler) (c *Client, err error) {
	c = &Client{
		message_handlers: handlers,
		host:             host,
		port:             port,
	}
	c.conn, err = net.Dial("tcp", host+":"+port)

	if err == nil {
		c.IsConnected = true
		go c.clientLoop()
	}

	return c, err
}

func (c *Client) Disconnect() {
	c.SendGoodbye()
	c.IsConnected = false
	c.conn.Close()
}

func (c *Client) SendMsg(msg string) (n int, err error) {
	return fmt.Fprintf(c.conn, msg+"\n")
}

func (c *Client) SendForward() (n int, err error) {
	return c.SendMsg("w")
}

func (c *Client) SendBackward() (n int, err error) {
	return c.SendMsg("s")
}

func (c *Client) SendTurnLeft() (n int, err error) {
	return c.SendMsg("a")
}

func (c *Client) SendTurnRight() (n int, err error) {
	return c.SendMsg("d")
}

func (c *Client) SendGetItem() (n int, err error) {
	return c.SendMsg("t")
}

func (c *Client) SendShoot() (n int, err error) {
	return c.SendMsg("e")
}

func (c *Client) SendRequestObservation() (n int, err error) {
	return c.SendMsg("o")
}

func (c *Client) SendRequestGameStatus() (n int, err error) {
	return c.SendMsg("g")
}

func (c *Client) SendRequestUserStatus() (n int, err error) {
	return c.SendMsg("q")
}

func (c *Client) SendRequestPosition() (n int, err error) {
	return c.SendMsg("p")
}

func (c *Client) SendRequestScoreboard() (n int, err error) {
	return c.SendMsg("u")
}

func (c *Client) SendGoodbye() (n int, err error) {
	return c.SendMsg("quit")
}

func (c *Client) SendName(name string) (n int, err error) {
	return c.SendMsg("name;" + name)
}

func (c *Client) SendSay(msg string) (n int, err error) {
	return c.SendMsg("say;" + msg)
}

func (c *Client) SendColour(r uint8, g uint8, b uint8) (n int, err error) {
	return c.SendMsg(fmt.Sprintf("color;%d;%d;%d", r, g, b))
}

func (c *Client) processCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if len(cmd) > 0 {
		cmd := strings.Split(cmd, ";")
		for _, handler := range c.message_handlers {
			handler(cmd)
		}
	}
}

func (c *Client) clientLoop() {
	for c.IsConnected {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == nil {
			c.processCommand(msg)
		}
	}
}
