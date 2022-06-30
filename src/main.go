package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const host = "atari.icad.puc-rio.br"
const port = "8888"

const name = "Gopher"
const time_delta = 1 * time.Second

const (
	north = iota
	east  = iota
	south = iota
	west  = iota
)

const (
	gameover = iota
	ready    = iota
	game     = iota
)

type GameState struct {
	x, y   int64
	dir    int
	state  int
	score  int64
	energy uint
}

func handler(msgs chan []string, cmd []string) {
	msgs <- cmd
}

func getStateVal(state string) int {
	switch state {
	case "gameover":
		return gameover
	case "ready":
		return ready
	case "game":
		return game
	}

	return -1
}

func botLoop(msgs chan []string, c net.Conn) {
	var msgSeconds time.Duration = 0

	status := GameState{
		x:      0,
		y:      0,
		dir:    north,
		state:  ready,
		score:  0,
		energy: 0,
	}

	sendName(c, name)
	sendColour(c, 255, 0, 0)
	sendRequestGameStatus(c)

	for {
		is_msgs_empty := false

		if status.state == game {
			doDecision(c)
		}

		for !is_msgs_empty {
			select {
			case msg := <-msgs:
				switch msg[0][0] {
				case 'g':
					state := getStateVal(msg[1])
					if state != status.state {
						fmt.Println("New Game State: " + msg[1])
						status.state = state
					}
				}
				fmt.Println(strings.Join(msg, " "))
			default:
				is_msgs_empty = true
			}
		}

		if msgSeconds >= 5*time.Second {
			sendRequestGameStatus(c)
			sendRequestScoreboard(c)
			msgSeconds = 0
		}

		time.Sleep(time_delta)
		msgSeconds += time_delta
	}
}

func doDecision(c net.Conn) {
	decision := getDecision()
	switch decision {
	case turn_right:
		sendTurnRight(c)

	case turn_left:
		sendTurnLeft(c)

	case forward:
		sendForward(c)

	case backward:
		sendBackward(c)

	case attack:
		sendShoot(c)

	case take_gold:
		sendGetItem(c)

	case take_powerup:
		sendGetItem(c)
	}
}

func main() {
	messages := make(chan []string, 20)

	c, err := connect(host, port, []cmdHandler{
		func(cmd []string) {
			handler(messages, cmd)
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(1 * time.Second)
	go botLoop(messages, c)

	bufio.NewReader(os.Stdin).ReadString('\n')

	disconnect(c)
}
