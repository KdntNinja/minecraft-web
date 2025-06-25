package main

import (
	"log"

	"github.com/KdntNinja/webcraft/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.Println("Starting Webcraft...")
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
