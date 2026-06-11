package main

import (
	"superbot/game"
	"superbot/gui"
	"superbot/server"
)

func init() {
	game := game.GetGame()
	gui := gui.NewGui(game)
	server.Listen(game)
	gui.Run()
}

func main() {}
