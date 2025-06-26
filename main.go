package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/game"
)

func main() {
	log.Println("Starting Webcraft...")
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
