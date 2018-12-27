package main

import (
	"encoding/json"
)

type game struct {
	logChan   chan string
	commsChan chan string
	player    *player
	mg        MapGen
}

func (g *game) tick() {
	for {
		select {
		case incommingCommand := <-g.commsChan:
			g.HandleCommand(incommingCommand)
		}
	}
}

func (g *game) init() {
	g.sendLog("Initialising game...")
	g.mg = MapGen{}
	//g.mg.GenerateMap(100, 70, 1, 10)
	g.player = NewPlayer()
	go g.tick()
}

func (g *game) sendLog(msg string) {
	if g.logChan != nil {
		g.logChan <- msg
	}
}

func newgame(logCh, commschan chan string) game {
	g := game{
		logChan:   logCh,
		commsChan: commschan,
	}
	g.sendLog("Creating game object")

	return g
}

func (g *game) HandleCommand(c string) {
	g.logChan <- c
	var m map[string]interface{}
	err := json.Unmarshal([]byte(c), &m)
	if err != nil {
		g.logChan <- err.Error()
		g.commsChan <- ResponseJSON(false, err.Error())
	} else {
		g.logChan <- m["command"].(string)
		res := g.CommSwitch(m["command"].(string), m)
		g.commsChan <- res
	}
}

func (g *game) CommSwitch(c string, input map[string]interface{}) string {
	switch c {
	case "enter":
		return g.EnterDung()
	case "exit":
		return g.ExitDung()
	case "mov":
		return g.MovePlayer(c, input)
	default:
		return ResponseJSON(false, "You get confused.")
	}
}

func (g *game) EnterDung() string {
	if g.player.inTown == true {
		g.player.inTown = false
		if g.mg.generated != true {
			g.mg.GenerateMap(100, 30, 1, 5)
		}
		g.player.position = g.mg.Graph.nodes[0]
		paths := g.mg.Graph.GetEdges(g.player.position)
		res := jsonPaths{}
		res.Status = true
		res.Msg = "Your party enters the dungeon."
		res.AvailablePaths = paths
		return MarshUp(res)
	}
	return ResponseJSON(false, "You're already in the dungeon, ya fuckwits!")
}

func (g *game) ExitDung() string {
	if g.player.inTown == true {
		return ResponseJSON(false, "You're already in town idiot!")
	}
	if g.player.position == g.mg.Graph.nodes[0] {
		g.player.inTown = true
		return ResponseJSON(true, "Your party enters the town.")
	}
	return ResponseJSON(false, "Can't enter the town from the middle of the dungeon, ya mongloids!")
}

func (g *game) MovePlayer(c string, command map[string]interface{}) string {

	if g.player.inTown == true {
		return ResponseJSON(false, "You strut around town aimlessly.")
	}
	desired, castOK := command["path"].(float64)
	if castOK == false {
		return ResponseJSON(false, "There's no path like that. You all flop around on the dungeon floor.")
	}
	desiredPath := int(desired)
	possiblePaths := g.mg.Graph.GetEdges(g.player.position)

	for _, path := range possiblePaths {
		if path == desiredPath {
			g.player.position = g.mg.Graph.nodes[desiredPath]
			re := jsonPaths{}
			re.Status = true
			re.Msg = "You enter a room"
			re.AvailablePaths = g.mg.Graph.GetEdges(g.player.position)
			return MarshUp(re)
		}
	}

	return ResponseJSON(false, "Your whole party bouncess off a dungeon wall.")
}
