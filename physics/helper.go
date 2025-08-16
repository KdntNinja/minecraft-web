package physics

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// IsSolid checks if a block at grid coordinates is solid, using grid offset
func IsSolid(blocks [][]int, x, y int, offsetX, offsetY int) bool {
	x -= offsetX
	y -= offsetY
	if y < 0 || x < 0 || y >= len(blocks) || x >= len(blocks[0]) {
		return false
	}
	return blocks[y][x] > 0
}
