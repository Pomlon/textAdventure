package main

import (
	"time"

	ui "github.com/gizak/termui"
)

var guimain gameui

func bootstrapUI() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewParagraph("Tu logi bedo")
	p.Height = 30
	p.Width = 100
	p.BorderLabel = "Log"
	p.X = 0
	p.Y = 0

	txt := ""

	p.Text = txt

	guimain.AddModule(p)

	items := []string{
		"[Mob] FakinFaker",
		"[Mob] Omagader",
		"[Item] Sword",
	}

	ls := ui.NewList()
	ls.Items = items
	ls.BorderLabel = "Shit in room"
	ls.Height = 20
	ls.Width = 25
	ls.X = 0
	ls.Y = 30

	guimain.AddModule(ls)

	guimain.Rend()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * 50).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				guimain.Rend()

			}
		case msg := <-guimain.logChan:
			p.Text = msg + "\n" + p.Text
			guimain.Rend()

		case <-ticker:
			guimain.Rend()
		}
	}
}

type gameui struct {
	widgets []ui.Bufferer
	logChan chan string
}

func NewGUI(logchan chan string) {
	guimain = gameui{}
	guimain.widgets = []ui.Bufferer{}
	guimain.logChan = logchan
}

func (gui *gameui) Rend() {
	ui.Render(gui.widgets...)
}

func (gui *gameui) AddModule(mod ui.Bufferer) {
	gui.widgets = append(gui.widgets, mod)
}
