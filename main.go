package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/rafaelrc7/inf1771-battlebot/ai"
	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

const host = "atari.icad.puc-rio.br"
const port = "8888"

const name = "Centurion"
const time_delta = 1 * time.Second

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
	energy int
}

func printScoreboard(scoreboard []string) {
	fmt.Println("----- SCOREBOARD -----")
	for i, s := range scoreboard {
		info := strings.Split(s, "#")

		score, err := strconv.Atoi(info[2])
		if err != nil {
			score = 0
		}

		hp, err := strconv.Atoi(info[3])
		if err != nil {
			hp = 0
		}

		fmt.Printf("%3d. %s (%s)\nHP = %3d - SCORE: %-6d\n\n", i, info[0], info[1], hp, score)
	}
	fmt.Println("----------------------")
}

func handler(msgs chan []string, cmd []string) {
	switch strings.ToLower(cmd[0]) {
	case "notification":
		fmt.Printf("[SERVER] %s\n", strings.Join(cmd[1:], " "))

	case "hello":
		fmt.Printf("[SERVER] %s has joined the game!\n", cmd[1])

	case "goodbye":
		fmt.Printf("[SERVER] %s has left the game!\n", cmd[1])

	case "changename":
		fmt.Printf("[SERVER] %s is now known as %s\n", cmd[1], cmd[2])

	case "u":
		printScoreboard(cmd[1:])

	default:
		msgs <- cmd
	}
}

func getStateVal(state string) int {
	switch strings.ToLower(state) {
	case "gameover":
		return gameover
	case "ready":
		return ready
	case "game":
		return game
	}

	return -1
}

func getDirVal(state string) int {
	switch strings.ToLower(state) {
	case "north":
		return gamemap.NORTH
	case "east":
		return gamemap.EAST
	case "south":
		return gamemap.SOUTH
	case "west":
		return gamemap.WEST
	}

	return -1
}

func botLoop(msgs chan []string, c net.Conn) {
	var msgSeconds time.Duration = 0

	status := GameState{
		x:      0,
		y:      0,
		dir:    gamemap.NORTH,
		state:  ready,
		score:  0,
		energy: 0,
	}

	sendName(c, name)
	sendColour(c, 255, 0, 0)

	sendRequestGameStatus(c)

	for {
		is_msgs_empty := false

		for !is_msgs_empty {
			select {
			case msg := <-msgs:
				switch strings.ToLower(msg[0]) {
				case "g":
					state := getStateVal(msg[1])
					if state != status.state {
						fmt.Println("New Game State: " + msg[1])
						status.state = state
						if state == game {
							sendRequestUserStatus(c)
							sendRequestObservation(c)
						}
					}
				case "s":
					status.x, _ = strconv.ParseInt(msg[1], 10, 64)
					status.y, _ = strconv.ParseInt(msg[2], 10, 64)
					status.dir = getDirVal(msg[3])
					status.score, _ = strconv.ParseInt(msg[5], 10, 64)
					status.energy, _ = strconv.Atoi(msg[6])
				}
				fmt.Println(strings.Join(msg, " "))
			default:
				is_msgs_empty = true
			}
		}

		if status.state == game && status.energy > 0 {
			doDecision(c)
		} else {
			if msgSeconds >= 5*time.Second {
				sendRequestGameStatus(c)
				sendRequestScoreboard(c)
				msgSeconds = 0
			}
		}

		time.Sleep(time_delta)
		msgSeconds += time_delta
	}
}

func doDecision(c net.Conn) {
	decision := ai.GetDecision()
	switch decision {
	case ai.TURN_RIGHT:
		sendTurnRight(c)

	case ai.TURN_LEFT:
		sendTurnLeft(c)

	case ai.FORWARD:
		sendForward(c)

	case ai.BACKWARD:
		sendBackward(c)

	case ai.ATTACK:
		sendShoot(c)

	case ai.TAKE_GOLD:
		sendGetItem(c)

	case ai.TAKE_POWERUP:
		sendGetItem(c)
	}

	sendRequestUserStatus(c)
	sendRequestObservation(c)
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
