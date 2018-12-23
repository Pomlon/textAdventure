package main

import (
	"strconv"
	"time"
)

type game struct {
	logChan chan string
}

func (g *game) tick() {
	for i := 0; i < 10; i++ {
		g.logChan <- strconv.Itoa(i)
		time.Sleep(time.Second)
	}
}

func (g *game) init() {
	g.sendLog("Initialising game...")
	go g.tick()
}

func (g *game) sendLog(msg string) {
	if g.logChan != nil {
		g.logChan <- msg
	}
}

func newgame(logCh chan string) game {
	g := game{
		logChan: logCh,
	}
	g.sendLog("Creating game object")

	return g
}
