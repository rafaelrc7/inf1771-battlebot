package main

import (
	"bufio"
	"fmt"
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

const name = "xxxxxxx"
const time_delta = 100 * time.Millisecond

const width = 59
const height = 34

const (
	gameover = iota
	ready    = iota
	game     = iota
	dead     = iota
)

const (
	Xi          = iota
	Yi          = iota
	Di          = iota
	SCOREi      = iota
	ENERGYi     = iota
	GAMESTATUSi = iota
	TIMEi       = iota
	HITi        = iota
	DAMAGEi     = iota
	OBSi        = iota
	ENEMYi      = iota
)

type GameState struct {
	state      int
	score      int
	tick       int
	lastAction int
	ai         ai.AI
	ingame     bool
	lastCoord  gamemap.Coord
}

type message struct {
	t, info int
	infou   uint
}

func main() {
	messages := make(chan message, 100)

	c, err := ClientNew(host, port, []CmdHandler{
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

	c.Disconnect()
}

func botLoop(msgs chan message, c *Client) {
	var msgSeconds time.Duration = 0

	status := GameState{
		state: ready,
		score: 0,
	}

	//sendName(c, name)
	c.SendColour(255, 0, 0)

	for c.IsConnected {

		if status.state == game {
			if !status.ingame {
				time.Sleep(time_delta)
				status.tick = 0
				status.ingame = true
				c.SendRequestUserStatus()
				c.SendRequestObservation()
				status.ai.Gamemap = gamemap.NewMap(height, width)
				time.Sleep(time_delta)
				status.ai.State = ai.EXPLORING
			}

			start := time.Now()

			status.tick++
			status.ai.EnemyDetected = false

			for is_msgs_empty := false; !is_msgs_empty; {
				select {
				case msg := <-msgs:
					switch msg.t {
					case Xi:
						status.ai.Coord.X = msg.info
					case Yi:
						status.ai.Coord.Y = msg.info
					case Di:
						status.ai.Coord.D = msg.info
					case SCOREi:
						status.score = msg.info
					case ENERGYi:
						status.ai.Energy = msg.info
					case GAMESTATUSi:
						if status.state != msg.info {
							fmt.Printf("[SERVER]: New game state %d\n", msg.info)
							status.state = msg.info
						}
					case OBSi:
						status.ai.Observations = msg.infou
					case ENEMYi:
						status.ai.EnemyDetected = true
					}
				default:
					is_msgs_empty = true
				}
			}

			hasChanged := status.ai.Gamemap.Tick()
			hasChanged = hasChanged || status.ai.Gamemap.VisitCell(status.ai.Coord, status.ai.Observations)
			status.ai.Gamemap.Print(status.ai.Coord)

			if status.lastAction == ai.BACKWARD || status.lastAction == ai.FORWARD {
				if status.ai.Coord == status.lastCoord {
					hasChanged = hasChanged || status.ai.Gamemap.MarkWall(status.ai.Coord, status.lastAction == ai.FORWARD)
				}
			}

			status.lastCoord = status.ai.Coord

			if status.ai.Energy > 0 {
				status.ai.Think(hasChanged)
				decision := status.ai.GetDecision(hasChanged)
				doDecision(c, decision)
				status.lastAction = decision

			} else {
				if msgSeconds >= 5*time.Second {
					c.SendRequestGameStatus()
					msgSeconds = 0
				}
			}

			status.ai.Observations = 0
			c.SendRequestUserStatus()
			c.SendRequestObservation()

			end := time.Now()
			if diff := end.Sub(start); diff < time_delta {
				fmt.Println(diff)
				fmt.Printf("STATE: %d\n", status.ai.State)
				fmt.Printf("Computation took: %v\nWill wait: %v\n", diff, (time_delta - diff))
				time.Sleep(time_delta - diff)
			}
		} else {
			status.ingame = false
			c.SendRequestGameStatus()
			c.SendRequestScoreboard()
			for is_msgs_empty := false; !is_msgs_empty; {
				select {
				case msg := <-msgs:
					switch msg.t {
					case GAMESTATUSi:
						if msg.info != status.state {
							status.state = msg.info
							fmt.Printf("[SERVER]: New game status %d\n", status.state)
						}
					}
				default:
					is_msgs_empty = true
				}
			}
			if status.state != game {
				time.Sleep(1 * time.Second)
			}
		}

		msgSeconds += time_delta
	}
}

func handler(msgs chan message, cmd []string) {
	fmt.Println(cmd)
	switch strings.ToLower(cmd[0]) {
	case "notification":
		fmt.Printf("[SERVER] %s\n", strings.Join(cmd[1:], " "))

	case "hello":
		fmt.Printf("[SERVER] %s has joined the game!\n", cmd[1])

	case "goodbye":
		fmt.Printf("[SERVER] %s has left the game!\n", cmd[1])

	case "changename":
		fmt.Printf("[SERVER] %s is now known as %s\n", cmd[1], cmd[2])

	case "player":
		fmt.Println(cmd)
		if len(cmd) == 7 {
			node, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
				return
			}

			name := cmd[2]

			x, err := strconv.Atoi(cmd[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
				return
			}

			y, err := strconv.Atoi(cmd[4])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
				return
			}

			dir := getDirVal(cmd[5])
			if dir == -1 {
				fmt.Fprintf(os.Stderr, "[ERROR] getDirVal(): invalid state %s\n", cmd[5])
				return
			}

			state := getStateVal(cmd[6])
			if dir == -1 {
				fmt.Fprintf(os.Stderr, "[ERROR] getStateVal(): invalid state %s\n", cmd[6])
				return
			}

			fmt.Printf("[PLAYER] %s (%d) %d %d,%d,%d\n", name, node, state, x, y, dir)
		}

	case "u":
		printScoreboard(cmd[1:])

	case "g":
		if st := getStateVal(cmd[1]); st != -1 {
			msgs <- message{t: GAMESTATUSi, info: st}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] getDirVal(): invalid state %s\n", cmd[4])
		}

		if time, err := strconv.Atoi(cmd[2]); err == nil {
			msgs <- message{t: TIMEi, info: time}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
		}

	case "s":
		if x, err := strconv.Atoi(cmd[1]); err == nil {
			msgs <- message{t: Xi, info: x}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
		}

		if y, err := strconv.Atoi(cmd[2]); err == nil {
			msgs <- message{t: Yi, info: y}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
		}

		if d := getDirVal(cmd[3]); d != -1 {
			msgs <- message{t: Di, info: d}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] getDirVal(): invalid dir %s\n", cmd[3])
		}

		if st := getStateVal(cmd[4]); st != -1 {
			msgs <- message{t: GAMESTATUSi, info: st}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] getDirVal(): invalid state %s\n", cmd[4])
		}

		if s, err := strconv.Atoi(cmd[5]); err == nil {
			msgs <- message{t: SCOREi, info: s}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
		}

		if e, err := strconv.Atoi(cmd[6]); err == nil {
			msgs <- message{t: ENERGYi, info: e}
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] Atoi(): %s\n", err)
		}

	case "h":
		msgs <- message{t: HITi, info: 1}

	case "d":
		msgs <- message{t: DAMAGEi, info: 1}

	case "o":
		obs := strings.Split(cmd[1], ",")
		msg := message{t: OBSi, infou: 0}
		for _, o := range obs {
			switch o {
			case "breeze":
				msg.infou |= gamemap.BREEZE

			case "flash":
				msg.infou |= gamemap.FLASH

			case "steps":
				msg.infou |= gamemap.STEPS

			case "redLight":
				msg.infou |= gamemap.REDLIGHT

			case "blueLight":
				msg.infou |= gamemap.BLUELIGHT

			case "blocked":
				msg.infou |= gamemap.BLOCKED

			default:
				enemy := strings.Split(o, "#")
				if len(enemy) > 1 {
					if dist, err := strconv.Atoi(enemy[1]); err == nil {
						msgs <- message{t: ENEMYi, info: dist}
					} else {
						fmt.Fprintf(os.Stderr, "Atoi(): %s\n", err)
					}
				}
			}
		}

		if msg.infou != 0 {
			msgs <- msg
		}

	default:
		fmt.Println(cmd)
	}
}

func doDecision(c *Client, decision int) {
	switch decision {
	case ai.TURN_RIGHT:
		c.SendTurnRight()
		fmt.Println("RIGHT")

	case ai.TURN_LEFT:
		c.SendTurnLeft()
		fmt.Println("LEFT")

	case ai.FORWARD:
		c.SendForward()
		fmt.Println("FORWARD")

	case ai.BACKWARD:
		c.SendBackward()
		fmt.Println("BACKWARD")

	case ai.ATTACK:
		c.SendShoot()
		fmt.Println("SHOOT")

	case ai.TAKE:
		c.SendGetItem()
		fmt.Println("TAKE")

	case ai.NOTHING:
		fmt.Println("NOTHING")
	}
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

func getStateVal(state string) int {
	switch strings.ToLower(state) {
	case "gameover":
		return gameover
	case "ready":
		return ready
	case "game":
		return game
	case "dead":
		return dead
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
