package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/game"
)

func main() {
	log.Println("Starting Webcraft...")

	// Graphics settings for performance
	ebiten.SetVsyncEnabled(true) // Prevent screen tearing
	ebiten.SetTPS(60)            // 60 ticks per second

	// Set window size hint for better performance
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Webcraft")

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
