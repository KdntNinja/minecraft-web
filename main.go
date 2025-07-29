package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	game "github.com/KdntNinja/webcraft/engine"
	"github.com/KdntNinja/webcraft/gameplay/world"
	"github.com/KdntNinja/webcraft/generation"
	"github.com/KdntNinja/webcraft/settings"
	"github.com/KdntNinja/webcraft/worldgen"
)

func main() {
	log.Println("Starting Webcraft...")

	// Graphics settings for performance
	ebiten.SetVsyncEnabled(true) // Prevent screen tearing
	ebiten.SetTPS(60)            // 60 ticks per second

	// Set window size hint for better performance
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Webcraft")

	chunkManager := generation.NewChunkManager(settings.ChunkViewDistance)
	spawn := worldgen.FindSafeSpawnPoint()
	g := game.NewGame(world.NewWorld(time.Now().UnixNano(), chunkManager, spawn))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
