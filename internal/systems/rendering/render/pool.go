package render

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Object pooling for performance - increased pool size
	drawOptionsPool []*ebiten.DrawImageOptions
	poolIndex       int
)

func initObjectPool() {
	// Initialize object pool for draw options - increased size for better performance
	drawOptionsPool = make([]*ebiten.DrawImageOptions, 100) // Increased from 50
	for i := range drawOptionsPool {
		drawOptionsPool[i] = &ebiten.DrawImageOptions{}
	}
	poolIndex = 0
}

func getDrawOptions() *ebiten.DrawImageOptions {
	// Simple object pooling
	if poolIndex >= len(drawOptionsPool) {
		poolIndex = 0
	}
	options := drawOptionsPool[poolIndex]
	options.GeoM.Reset() // Reset transform
	poolIndex++
	return options
}
