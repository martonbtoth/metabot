package gui

import (
	"fmt"
	"time"

	"superbot/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Gui interface {
	Run()
	UpdateState()
	AppendLog(s string)
}

type gui struct {
	g                  game.Game
	window             fyne.Window
	terminal           *widget.TextGrid
	log                string
	scroll             *container.Scroll
	previousPlayerGuid any
	tabs               *container.AppTabs
	treeContainer      *fyne.Container
}

func (g *gui) Run() {
	go func() {
		for range time.Tick(300 * time.Millisecond) {
			g.UpdateState()
		}
	}()

	g.window.ShowAndRun()
}

func (g *gui) AppendLog(s string) {
	g.log = g.log + s + "\n"
}

func (g *gui) UpdateState() {
	playerGuid := g.g.GetPlayerGuid()
	if playerGuid != g.previousPlayerGuid {
		g.previousPlayerGuid = playerGuid
		if playerGuid == 0 {
			g.AppendLog("Player not logged in")
		} else {
			g.AppendLog(fmt.Sprintf("Player guid is %v", playerGuid))
		}
	}
	g.terminal.SetText(g.log)

}

func (g *gui) EnumerateVisibleObjects() {
	g.g.EnumerateVisibleObjects()
	g.AppendLog(fmt.Sprintf("Objects enumerated: %v", len(g.g.GetVisibleObjects())))
	for index, wowObject := range g.g.GetVisibleObjects() {
		g.AppendLog(fmt.Sprintf("Object %v: %v", index, wowObject))
	}

	g.treeContainer.Refresh()

	g.scroll.ScrollToBottom()
}

func NewGui(game game.Game) Gui {
	gui := &gui{}
	gui.g = game
	buildGui(gui)
	return gui
}
