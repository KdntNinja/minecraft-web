package main

import (
	"log"

	"github.com/KdntNinja/webcraft/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.Println("Starting Webcraft...")
	ebiten.SetWindowSize(engine.ScreenWidth, engine.ScreenHeight)
	game := engine.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
