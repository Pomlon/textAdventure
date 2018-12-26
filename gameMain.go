package main

type game struct {
	logChan   chan string
	commsChan chan string
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
	mg := MapGen{}
	mg.GenerateMap(100, 70, 1, 10)
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
}
