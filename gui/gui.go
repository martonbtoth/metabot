package gui

import (
	"fmt"
	"time"

	"superbot/game"
	"superbot/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Gui interface {
	Run()
	UpdateState()
}

type gui struct {
	g                  game.Game
	window             fyne.Window
	terminal           *widget.TextGrid
	l                  *logger.Logger
	scroll             *container.Scroll
	previousPlayerGuid any
	tabs               *container.AppTabs
	treeContainer      *fyne.Container
}

func (g *gui) Run() {
	logger.GetLogger().Listener = func(logBuffer string) {
		g.terminal.SetText(logBuffer)
		g.scroll.ScrollToBottom()
	}

	go func() {
		for range time.Tick(300 * time.Millisecond) {
			g.UpdateState()
		}
	}()

	g.window.ShowAndRun()
}

func (g *gui) UpdateState() {
	playerGuid := g.g.GetPlayerGuid()
	if playerGuid != g.previousPlayerGuid {
		g.previousPlayerGuid = playerGuid
		if playerGuid == 0 {
			g.l.Log("Player not logged in")
		} else {
			g.l.Log(fmt.Sprintf("Player guid is %v", playerGuid))
		}
	}
}

func (g *gui) EnumerateVisibleObjects() {
	g.g.EnumerateVisibleObjects()
	g.l.Log(fmt.Sprintf("%v objects enumerated", len(g.g.GetVisibleObjects())))

	g.treeContainer.Refresh()

	g.scroll.ScrollToBottom()
}

func NewGui(game game.Game) Gui {
	gui := &gui{}
	gui.g = game
	gui.l = logger.GetLogger()
	buildGui(gui)
	return gui
}
