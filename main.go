package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/game"
)

func main() {
	log.Println("Starting Webcraft...")

	// Set performance options for lower-end hardware
	ebiten.SetVsyncEnabled(true) // Enable VSync to prevent screen tearing
	ebiten.SetTPS(60)            // Limit to 60 TPS for consistent performance

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
