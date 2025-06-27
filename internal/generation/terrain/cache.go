package terrain

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

var (
	// Chunk cache for improved performance
	chunkCache   = make(map[string]block.Chunk)
	maxCacheSize = 100
)

// getCachedChunk retrieves a chunk from cache if it exists
func getCachedChunk(chunkX, chunkY int) (block.Chunk, bool) {
	cacheKey := fmt.Sprintf("%d,%d", chunkX, chunkY)
	if cached, exists := chunkCache[cacheKey]; exists {
		return cached, true
	}
	return block.Chunk{}, false
}

// cacheChunk stores a chunk in the cache
func cacheChunk(chunkX, chunkY int, chunk block.Chunk) {
	// Limit cache size to prevent memory bloat
	if len(chunkCache) >= maxCacheSize {
		clearOldCacheEntries()
	}

	cacheKey := fmt.Sprintf("%d,%d", chunkX, chunkY)
	chunkCache[cacheKey] = chunk
}

// clearOldCacheEntries removes half of the cache entries when full
func clearOldCacheEntries() {
	// Simple approach - clear half the cache
	for k := range chunkCache {
		delete(chunkCache, k)
		if len(chunkCache) <= maxCacheSize/2 {
			break
		}
	}
}
