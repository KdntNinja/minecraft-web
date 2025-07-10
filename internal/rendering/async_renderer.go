package rendering

import (
	"image/color"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/world/chunks"
)

// TileRenderJob represents a single tile rendering job
type TileRenderJob struct {
	blockType block.BlockType
	px, py    float64
	tile      *ebiten.Image
}

// ChunkRenderJob represents a chunk rendering job
type ChunkRenderJob struct {
	coord chunks.ChunkCoord
	chunk *block.Chunk
	jobs  []TileRenderJob
}

// AsyncRenderer handles multithreaded rendering
type AsyncRenderer struct {
	chunkJobs    chan ChunkRenderJob
	tileJobs     chan TileRenderJob
	workers      sync.WaitGroup
	numWorkers   int
	isShutdown   bool
	shutdownChan chan struct{}
}

var (
	asyncRenderer *AsyncRenderer
	rendererOnce  sync.Once
)

// GetAsyncRenderer returns the singleton async renderer
func GetAsyncRenderer() *AsyncRenderer {
	rendererOnce.Do(func() {
		asyncRenderer = &AsyncRenderer{
			chunkJobs:    make(chan ChunkRenderJob, 64),
			tileJobs:     make(chan TileRenderJob, 1024),
			numWorkers:   runtime.NumCPU(),
			shutdownChan: make(chan struct{}),
		}
		asyncRenderer.startWorkers()
	})
	return asyncRenderer
}

// startWorkers starts the rendering worker pool
func (ar *AsyncRenderer) startWorkers() {
	// Start chunk processing workers
	for i := 0; i < ar.numWorkers; i++ {
		ar.workers.Add(1)
		go ar.chunkWorker(i)
	}
}

// chunkWorker processes chunk rendering jobs
func (ar *AsyncRenderer) chunkWorker(workerID int) {
	defer ar.workers.Done()

	for {
		select {
		case <-ar.shutdownChan:
			return
		case job := <-ar.chunkJobs:
			ar.processChunkJob(job)
		}
	}
}

// processChunkJob processes a single chunk for rendering
func (ar *AsyncRenderer) processChunkJob(job ChunkRenderJob) {
	// Process the chunk and generate tile render jobs
	// This is done in parallel but the actual drawing is deferred
	for _, tileJob := range job.jobs {
		select {
		case ar.tileJobs <- tileJob:
		case <-ar.shutdownChan:
			return
		}
	}
}

// Shutdown stops all rendering workers
func (ar *AsyncRenderer) Shutdown() {
	if ar.isShutdown {
		return
	}
	ar.isShutdown = true
	close(ar.shutdownChan)
	ar.workers.Wait()
	close(ar.chunkJobs)
	close(ar.tileJobs)
}

// DrawWithCameraAsync renders the world with camera offset using async processing
func DrawWithCameraAsync(worldChunks map[chunks.ChunkCoord]*block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if !isInitialized {
		initTileImages()
	}

	// Fill sky background
	screen.Fill(color.RGBA{135, 206, 250, 255})

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	// Calculate visible chunk bounds
	startChunkX := int(cameraX / float64(chunkWidth*tileSize))
	endChunkX := int((cameraX+float64(screenWidth))/float64(chunkWidth*tileSize)) + 1
	startChunkY := int(cameraY / float64(chunkHeight*tileSize))
	endChunkY := int((cameraY+float64(screenHeight))/float64(chunkHeight*tileSize)) + 1

	// Precompute float versions for bounds
	fScreenWidth := float64(screenWidth)
	fScreenHeight := float64(screenHeight)
	fTileSize := float64(tileSize)

	// Process chunks in parallel, but collect drawing operations
	var renderJobs []TileRenderJob
	var renderMutex sync.Mutex
	var wg sync.WaitGroup

	for coord, chunk := range worldChunks {
		// Early culling: skip chunks outside the visible area
		if coord.X < startChunkX || coord.X > endChunkX || coord.Y < startChunkY || coord.Y > endChunkY {
			continue
		}

		wg.Add(1)
		go func(coord chunks.ChunkCoord, chunk *block.Chunk) {
			defer wg.Done()

			var chunkJobs []TileRenderJob
			baseTileX := coord.X * chunkWidth
			baseTileY := coord.Y * chunkHeight

			for y := 0; y < chunkHeight; y++ {
				globalTileY := baseTileY + y
				py := float64(globalTileY*tileSize) - cameraY
				// Skip entire row if off screen vertically
				if py+fTileSize < 0 || py >= fScreenHeight {
					continue
				}
				for x := 0; x < chunkWidth; x++ {
					globalTileX := baseTileX + x
					px := float64(globalTileX*tileSize) - cameraX
					// Skip tile if off screen horizontally
					if px+fTileSize < 0 || px >= fScreenWidth {
						continue
					}
					blockType := (*chunk)[y][x]
					if blockType == block.Air {
						continue // Skip air blocks
					}
					tile := tileImages[blockType]
					if tile == nil {
						continue
					}

					chunkJobs = append(chunkJobs, TileRenderJob{
						blockType: blockType,
						px:        px,
						py:        py,
						tile:      tile,
					})
				}
			}

			// Add this chunk's jobs to the global list
			renderMutex.Lock()
			renderJobs = append(renderJobs, chunkJobs...)
			renderMutex.Unlock()
		}(coord, chunk)
	}

	wg.Wait()

	// Draw all tiles (this must be done synchronously on the main thread)
	drawOpts := getDrawOptions()
	for _, job := range renderJobs {
		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Translate(job.px, job.py)
		screen.DrawImage(job.tile, drawOpts)
	}
}
