package main

import (
	"metabot/game"
	"metabot/gui"
	"metabot/server"
)

func init() {
	g := game.GetGame()
	game.HookEvents()
	gui := gui.NewGui(g)
	server.Listen(g)
	gui.Run()
}

func main() {}
