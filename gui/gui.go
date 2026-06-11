package gui

import (
	"fmt"
	"strings"
	"time"

	"metabot/game"
	"metabot/logger"

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
	logger.GetLogger().AddListener(
		func(logBuffer string, s string) {
			split := strings.Split(logBuffer, "\n")
			splitLength := len(split) - 40
			if splitLength > 0 {
				split = split[splitLength:]
			}
			joined := strings.Join(split, "\n")
			g.terminal.SetText(joined)
			g.scroll.ScrollToBottom()
		},
	)

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
