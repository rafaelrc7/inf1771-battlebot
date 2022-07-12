package main

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rafaelrc7/inf1771-battlebot/ai"
	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

const host = "atari.icad.puc-rio.br"
const port = "8888"

const time_delta = 100 * time.Millisecond

const width = 59
const height = 34

const (
	gameover = iota
	ready
	game
	dead
)

const (
	Xi = iota + 1
	Yi
	Di
	SCOREi
	ENERGYi
	GAMESTATUSi
	TIMEi
	HITi
	DAMAGEi
	OBSi
	ENEMYi
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
	rand.Seed(time.Now().UnixMilli())
	messages := make(chan message, 100)

	f, _ := os.Create(time.Now().Format("06-02-01_15-04-05.log"))
	defer f.Close()

	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)

	c, err := ClientNew(host, port, []CmdHandler{
		func(cmd []string) {
			handler(messages, cmd)
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	defer c.Disconnect()

	time.Sleep(1 * time.Second)
	go botLoop(messages, c)

	bufio.NewReader(os.Stdin).ReadString('\n')
}

func botLoop(msgs chan message, c *Client) {
	var msgSeconds time.Duration = 0

	status := GameState{
		state: ready,
		score: 0,
	}

	c.SendName(getRandomName())
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
							status.state = msg.info
						}
					case OBSi:
						status.ai.Observations |= msg.infou
					case ENEMYi:
						status.ai.EnemyDetected = true
					case HITi:
						status.ai.Observations |= gamemap.HIT
					case DAMAGEi:
						status.ai.Observations |= gamemap.DAMAGE
					}
				default:
					is_msgs_empty = true
				}
			}

			hasChanged := status.ai.Gamemap.Tick()
			hasChanged = hasChanged || status.ai.Gamemap.VisitCell(status.ai.Coord, status.ai.Observations)

			if status.lastAction == ai.BACKWARD || status.lastAction == ai.FORWARD {
				if status.ai.Coord == status.lastCoord {
					hasChanged = hasChanged || status.ai.Gamemap.MarkWall(status.ai.Coord, status.lastAction == ai.FORWARD)
				}
			}

			status.lastCoord = status.ai.Coord

			log.Infof("Tick: %d", status.tick)

			if status.ai.Energy > 0 {
				status.ai.Think(hasChanged)
				decision := status.ai.GetDecision(hasChanged)
				doDecision(c, decision)
				status.lastAction = decision
				log.Infof("State:  %s", stateStr(status.ai.State))
				log.Infof("Action: %s", actionStr(decision))

			} else {
				if msgSeconds >= 10*time.Second {
					c.SendRequestGameStatus()
					c.SendRequestScoreboard()
					msgSeconds = 0
				}
			}

			status.ai.Observations = 0
			c.SendRequestUserStatus()
			c.SendRequestObservation()

			end := time.Now()
			if diff := end.Sub(start); diff < time_delta {
				log.WithFields(log.Fields{
					"Took":       diff,
					"Will sleep": (time_delta - diff),
				}).Print("Sleep")
				time.Sleep(time_delta - diff)
			}
		} else {
			status.ingame = false
			c.SendRequestGameStatus()
			if msgSeconds >= 10*time.Second {
				c.SendRequestScoreboard()
				msgSeconds = 0
			}
			for is_msgs_empty := false; !is_msgs_empty; {
				select {
				case msg := <-msgs:
					switch msg.t {
					case GAMESTATUSi:
						if msg.info != status.state {
							status.state = msg.info
							log.Infof("New game status %d\n", status.state)
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
	switch strings.ToLower(cmd[0]) {
	case "notification":
		log.WithFields(log.Fields{
			"command": cmd,
		}).Printf("%s", strings.Join(cmd[1:], " "))

	case "hello":
		log.WithFields(log.Fields{
			"command": cmd,
			"player":  cmd[1],
		}).Print("joined the game")

	case "goodbye":
		log.WithFields(log.Fields{
			"command": cmd,
			"player":  cmd[1],
		}).Print("left the game")

	case "changename":
		log.WithFields(log.Fields{
			"command": cmd,
		}).Printf("'%s' is now known as '%s'", cmd[1], cmd[2])

	case "player":
		if len(cmd) == 7 {
			log.WithFields(log.Fields{
				"name":  cmd[2],
				"node":  cmd[1],
				"x":     cmd[3],
				"y":     cmd[4],
				"dir":   cmd[5],
				"state": cmd[6],
			})
		}

	case "u":
		printScoreboard(cmd[1:])

	case "g":
		log.WithFields(log.Fields{
			"time":  cmd[2],
			"state": cmd[1],
		}).Print("Game Status")

		if st := getStateVal(cmd[1]); st != -1 {
			msgs <- message{t: GAMESTATUSi, info: st}
		} else {
			log.Errorf("getStateVal(): invalid state %s", cmd[1])
		}

		if time, err := strconv.Atoi(cmd[2]); err == nil {
			msgs <- message{t: TIMEi, info: time}
		} else {
			log.Errorf("Atoi(): %s", err)
		}

	case "s":
		log.WithFields(log.Fields{
			"x":      cmd[1],
			"y":      cmd[2],
			"d":      cmd[3],
			"state":  cmd[4],
			"score":  cmd[5],
			"energy": cmd[6],
		}).Print("player stats")

		if x, err := strconv.Atoi(cmd[1]); err == nil {
			msgs <- message{t: Xi, info: x}
		} else {
			log.Errorf("Atoi(): %s", err)
		}

		if y, err := strconv.Atoi(cmd[2]); err == nil {
			msgs <- message{t: Yi, info: y}
		} else {
			log.Errorf("Atoi(): %s", err)
		}

		if d := getDirVal(cmd[3]); d != -1 {
			msgs <- message{t: Di, info: d}
		} else {
			log.Errorf("getDirVal(): invalid dir %s", cmd[3])
		}

		if st := getStateVal(cmd[4]); st != -1 {
			msgs <- message{t: GAMESTATUSi, info: st}
		} else {
			log.Errorf("getStateVal(): invalid state %s", cmd[4])
		}

		if s, err := strconv.Atoi(cmd[5]); err == nil {
			msgs <- message{t: SCOREi, info: s}
		} else {
			log.Errorf("Atoi(): %s", err)
		}

		if e, err := strconv.Atoi(cmd[6]); err == nil {
			msgs <- message{t: ENERGYi, info: e}
		} else {
			log.Errorf("Atoi(): %s", err)
		}

	case "h":
		log.Print("hit received")
		msgs <- message{t: HITi, info: 1}

	case "d":
		log.Print("damage inflicted")
		msgs <- message{t: DAMAGEi, info: 1}

	case "o":
		obs := strings.Split(cmd[1], ",")

		log.WithFields(log.Fields{
			"senses": obs,
		}).Printf("Observations")

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
						log.Errorf("Atoi(): %s", err)
					}
				}
			}
		}

		if msg.infou != 0 {
			msgs <- msg
		}

	default:
		log.WithFields(log.Fields{
			"command": cmd,
		}).Print("Unknown command")
	}
}

func doDecision(c *Client, decision int) {
	switch decision {
	case ai.TURN_RIGHT:
		c.SendTurnRight()

	case ai.TURN_LEFT:
		c.SendTurnLeft()

	case ai.FORWARD:
		c.SendForward()

	case ai.BACKWARD:
		c.SendBackward()

	case ai.ATTACK:
		c.SendShoot()

	case ai.TAKE:
		c.SendGetItem()

	case ai.NOTHING:
	}
}

func printScoreboard(scoreboard []string) {
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

		log.WithFields(log.Fields{
			"name":   info[0],
			"status": info[1],
			"hp":     hp,
			"score":  score,
		}).Printf("score %d", i)

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

func stateStr(state int) string {
	switch state {
	case ai.STOP:
		return "STOP"
	case ai.EXPLORING:
		return "EXPLORING"
	case ai.FETCHING_GOLD:
		return "FETCHING_GOLD"
	case ai.FETCHING_PU:
		return "FETCHING_PU"
	case ai.ATTACKING:
		return "ATTACKING"
	case ai.FLEEING:
		return "FLEEING"
	default:
		return "UNK"
	}
}

func actionStr(action int) string {
	switch action {
	case ai.NOTHING:
		return "NOTHING"
	case ai.TURN_RIGHT:
		return "TURN_RIGHT"
	case ai.TURN_LEFT:
		return "TURN_LEFT"
	case ai.FORWARD:
		return "FORWARD"
	case ai.BACKWARD:
		return "BACKWARD"
	case ai.ATTACK:
		return "ATTACK"
	case ai.TAKE:
		return "TAKE"
	default:
		return "UNK"
	}
}

func getRandomName() string {
	names := []string{"Crusader", "Covenanter", "Cavalier", "Centaur",
		"Cromwell", "Challenger", "Comet", "Centurion", "Conqueror",
		"Chieftain", "Churchill", "TOG", "Valentine", "Matilda", "Contentious"}
	return names[rand.Intn(len(names))]
}
