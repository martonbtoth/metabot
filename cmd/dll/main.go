package main

import (
	"superbot/game"
	"superbot/gui"
)

func init() {
	gui := gui.NewGui(*game.GetGame())
	gui.Run()
}

func main() {}
