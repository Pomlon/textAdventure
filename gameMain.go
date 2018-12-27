package main

import (
	"encoding/json"
	"errors"

	"github.com/Pomlon/textAdventure/utils"
	"github.com/jinzhu/copier"
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
	g.player.GenerateHeroes()
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
		g.commsChan <- ResponseJSON(errcodes.JSONParseErr, err.Error())
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
	case "partyStatus":
		return g.PartyStatus()
	case "attack":
		return g.Attack(c, input)
	case "roomStatus":
		return g.RoomStatus()
	case "endTurn":
		return g.NextTurn()
	default:
		return ResponseJSON(errcodes.UnkownCommand, "You get confused.")
	}
}

func (g *game) EnterDung() string {
	if g.player.inTown == true {
		g.player.inTown = false
		if g.mg.generated != true {
			g.mg.GenerateMap(100, 30, 1, 5)
		}
		g.player.position = g.mg.Graph.nodes[0]
		g.player.position.visited = true
		paths := g.mg.Graph.GetEdges(g.player.position)
		res := jsonPaths{}
		res.Status = errcodes.OK
		res.Msg = "Your party enters the dungeon."
		res.AvailablePaths = paths
		return MarshUp(res)
	}
	return ResponseJSON(errcodes.AlreadyInside, "You're already in the dungeon, ya fuckwits!")
}

func (g *game) ExitDung() string {
	if g.player.inTown == true {
		return ResponseJSON(errcodes.AlreadyInside, "You're already in town idiot!")
	}
	if g.player.position == g.mg.Graph.nodes[0] {
		g.player.inTown = true
		return ResponseJSON(errcodes.OK, "Your party enters the town.")
	}
	return ResponseJSON(errcodes.UnfulfilledReqs, "Can't enter the town from the middle of the dungeon, ya mongloids!")
}

func (g *game) MovePlayer(c string, command map[string]interface{}) string {

	if g.player.inTown == true {
		return ResponseJSON(errcodes.CannotInTown, "You strut around town aimlessly.")
	}
	desired, castOK := command["path"].(float64)
	if castOK == false {
		return ResponseJSON(errcodes.JSONParseErr, "There's no path like that. You all flop around on the dungeon floor.")
	}
	desiredPath := int(desired)
	if g.player.position.monsterHandler.DeathStatus() == false && g.mg.Graph.nodes[desiredPath].visited == false {
		return ResponseJSON(errcodes.MonstersStillLive, "With enemies around, only already visited areas can be travelled to.")
	}

	possiblePaths := g.mg.Graph.GetEdges(g.player.position)

	for _, path := range possiblePaths {
		if path == desiredPath {
			if g.player.position.monsterHandler.DeathStatus() == true {
				g.player.position.visited = true
			}
			g.player.position = g.mg.Graph.nodes[desiredPath]
			g.player.position.monsterHandler.GenerateMobs()
			re := jsonPaths{}
			re.Status = errcodes.OK
			re.Msg = "You enter a room"
			re.AvailablePaths = g.mg.Graph.GetEdges(g.player.position)
			re.Monsters = g.player.position.monsterHandler.Monsters
			return MarshUp(re)
		}
	}

	return ResponseJSON(errcodes.DoesNotExist, "Your whole party bouncess off a dungeon wall.")
}

func (g *game) PartyStatus() string {
	res := jsonPartyStatus{}
	res.Status = errcodes.OK
	res.Msg = "You check the wellbeing of your party"
	res.Heroes = g.player.heroes
	return MarshUp(res)
}

func (g *game) RoomStatus() string {
	res := jsonPaths{}
	res.AvailablePaths = g.mg.Graph.GetEdges(g.player.position)
	res.Monsters = g.player.position.monsterHandler.Monsters
	res.Status = errcodes.OK
	res.Msg = "You look around and check what's up in the dungeon."
	return MarshUp(res)
}

func (g *game) Attack(c string, input map[string]interface{}) string {
	var heropos int
	msg, err := JSONtoInt(input["heroPos"], &heropos)
	if err != nil {
		return msg
	}
	if heropos >= len(g.player.heroes) || heropos < 0 {
		return ResponseJSON(errcodes.DoesNotExist, "That hero does not exist!")
	}
	if g.player.position.monsterHandler.DeathStatus() == true {
		return ResponseJSON(errcodes.UnfulfilledReqs, "Everything is dead already you sick fuck!")
	}
	mon, err := g.player.position.monsterHandler.GetFirstLivingThing()
	if err != nil {
		return ResponseJSON(errcodes.UnfulfilledReqs, "Everything is dead already you sick fuck!")
	}
	if g.player.heroes[heropos].Action == false {
		return ResponseJSON(errcodes.OutOfResource, "Ya don't have actions left mate!")
	}
	dmg := mon.getHit(g.player.heroes[heropos].Atk)
	g.player.heroes[heropos].Action = false
	res := jsonAttack{}
	res.Status = errcodes.OK
	res.Msg = "You hit tha monstah!"
	res.Damage = dmg
	res.Monster = mon

	return MarshUp(res)
}

func (g *game) NextTurn() string {
	ret := make(map[string]interface{})
	if g.player.position.monsterHandler.DeathStatus() == false {
		re, err := g.MonsterAttack()
		if err != nil {
			return ResponseJSON(errcodes.OutOfResource, "All heroes ded")
		}
		ret["monsterAttacks"] = re
	}
	g.RefreshActions()
	return MarshUp(ret)
}

func (g *game) RefreshActions() {
	for i := 0; i < len(g.player.heroes); i++ {
		g.player.heroes[i].Action = true
	}

	for i := 0; i < len(g.player.position.monsterHandler.Monsters); i++ {
		g.player.position.monsterHandler.Monsters[i].Action = true
	}
}

func (g *game) MonsterAttack() ([]jsonDamaged, error) {
	var dmg []jsonDamaged
	for i := 0; i < len(g.player.position.monsterHandler.Monsters); i++ {
		if g.player.position.monsterHandler.Monsters[i].Dead == false && g.player.position.monsterHandler.Monsters[i].Action == true {
			her, err := g.player.GetFirstLivingHero()
			if err != nil {
				return []jsonDamaged{}, errors.New("All heroes dead")
			}
			dmgrcvd := her.getHit(g.player.position.monsterHandler.Monsters[i].Atk)
			g.player.position.monsterHandler.Monsters[i].Action = false
			re := jsonDamaged{}
			re.Status = errcodes.OK
			re.Damage = dmgrcvd
			re.Hero = Hero{}
			copier.Copy(&re.Hero, &her)
			re.Msg = "Ya got hit by a nasty monstah!"
			dmg = append(dmg, re)
		}
	}
	return dmg, nil
}

func JSONtoInt(val interface{}, to *int) (string, error) {
	desired, castOK := val.(float64)
	if castOK == false {
		return ResponseJSON(errcodes.JSONParseErr, "You've tried to apply a wrong set of commands :(."),
			errors.New("Cast fail")
	}

	num := int(desired)
	*to = num
	return "", nil
}
