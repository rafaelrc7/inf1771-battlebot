package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type cmdHandler func([]string)

func connect(host, port string, handlers []cmdHandler) (c net.Conn, err error) {
	c, err = net.Dial("tcp", host+":"+port)

	if err == nil {
		go clientLoop(c, handlers)
	}

	return c, err
}

func disconnect(c net.Conn) {
	sendGoodbye(c)
	c.Close()
}

func sendMsg(c net.Conn, msg string) (n int, err error) {
	return fmt.Fprintf(c, msg+"\n")
}

func sendForward(c net.Conn) (n int, err error) {
	return sendMsg(c, "w")
}

func sendBackward(c net.Conn) (n int, err error) {
	return sendMsg(c, "s")
}

func sendTurnLeft(c net.Conn) (n int, err error) {
	return sendMsg(c, "a")
}

func sendTurnRight(c net.Conn) (n int, err error) {
	return sendMsg(c, "d")
}

func sendGetItem(c net.Conn) (n int, err error) {
	return sendMsg(c, "t")
}

func sendShoot(c net.Conn) (n int, err error) {
	return sendMsg(c, "e")
}

func sendRequestObservation(c net.Conn) (n int, err error) {
	return sendMsg(c, "o")
}

func sendRequestGameStatus(c net.Conn) (n int, err error) {
	return sendMsg(c, "g")
}

func sendRequestUserStatus(c net.Conn) (n int, err error) {
	return sendMsg(c, "q")
}

func sendRequestPosition(c net.Conn) (n int, err error) {
	return sendMsg(c, "p")
}

func sendRequestScoreboard(c net.Conn) (n int, err error) {
	return sendMsg(c, "u")
}

func sendGoodbye(c net.Conn) (n int, err error) {
	return sendMsg(c, "quit")
}

func sendName(c net.Conn, name string) (n int, err error) {
	return sendMsg(c, "name;"+name)
}

func sendSay(c net.Conn, msg string) (n int, err error) {
	return sendMsg(c, "say;"+msg)
}

func sendColour(c net.Conn, r uint8, g uint8, b uint8) (n int, err error) {
	return sendMsg(c, fmt.Sprintf("color;%d;%d;%d", r, g, b))
}

func processCommand(cmd string, handlers []cmdHandler) {
	cmd = strings.TrimSpace(cmd)
	if len(cmd) > 0 {
		cmd := strings.Split(cmd, ";")
		for _, handler := range handlers {
			handler(cmd)
		}
	}
}

func clientLoop(c net.Conn, handlers []cmdHandler) {
	for {
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err == nil {
			processCommand(msg, handlers)
		}
	}
}
