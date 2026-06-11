package main

import (
	"superbot/game"
	"superbot/gui"
	"superbot/server"
)

func init() {
	g := game.GetGame()
	game.HookEvents()
	gui := gui.NewGui(g)
	server.Listen(g)
	gui.Run()
}

func main() {}
