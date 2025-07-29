//go:build js && wasm

package progress

import (
	"fmt"
	"syscall/js"
)

// updateUI sends the current progress to the JavaScript frontend
func updateUI() {
	percentage := calculateCurrentPercentage()

	// Send progress to JavaScript immediately if running in browser
	if js.Global().Get("updateLoadingProgress").Truthy() {
		js.Global().Call("updateLoadingProgress", percentage, globalTracker.CurrentStage, globalTracker.CurrentMessage)
	}

	// Log to console for debugging (after JS call to ensure UI updates first)
	fmt.Printf("[PROGRESS] %s: %s (%d%%)\n", globalTracker.CurrentStage, globalTracker.CurrentMessage, percentage)
}

// Legacy support functions for backward compatibility
func UpdateProgress(current, total int, stage, message string) {
	percentage := 0
	if total > 0 {
		percentage = (current * 100) / total
	}

	fmt.Printf("[PROGRESS] %s: %s (%d/%d - %d%%)\n", stage, message, current, total, percentage)

	if js.Global().Get("updateLoadingProgress").Truthy() {
		js.Global().Call("updateLoadingProgress", percentage, stage, message)
	}
}
