package progress

import "time"

// ProgressStep represents a single step in the loading process
type ProgressStep struct {
	Name        string  // Name of the step
	Weight      float64 // Relative weight (e.g., 1.0 for equal weight, 2.0 for double weight)
	SubSteps    int     // Number of substeps (e.g., number of chunks to generate)
	CurrentSub  int     // Current substep progress
	Description string  // Current description of what's happening
}

// ProgressTracker manages dynamic progress allocation
type ProgressTracker struct {
	Steps          []ProgressStep
	CurrentStep    int
	TotalWeight    float64
	CurrentStage   string
	CurrentMessage string
}

var globalTracker ProgressTracker

// InitializeProgress sets up the progress tracker with predefined steps
func InitializeProgress(steps []ProgressStep) {
	globalTracker = ProgressTracker{
		Steps:       steps,
		CurrentStep: 0,
		TotalWeight: 0,
	}

	// Calculate total weight
	for _, step := range steps {
		globalTracker.TotalWeight += step.Weight
	}

	// Start with the first step
	if len(steps) > 0 {
		globalTracker.CurrentStage = steps[0].Name
		globalTracker.CurrentMessage = steps[0].Description
	}

	updateUI()
}

// SetCurrentStepSubSteps updates the number of substeps for the current step
func SetCurrentStepSubSteps(subSteps int, description string) {
	if globalTracker.CurrentStep < len(globalTracker.Steps) {
		globalTracker.Steps[globalTracker.CurrentStep].SubSteps = subSteps
		globalTracker.Steps[globalTracker.CurrentStep].CurrentSub = 0
		globalTracker.Steps[globalTracker.CurrentStep].Description = description
		globalTracker.CurrentMessage = description
		updateUI()
	}
}

// UpdateCurrentStepProgress updates the progress within the current step
func UpdateCurrentStepProgress(currentSub int, description string) {
	if globalTracker.CurrentStep < len(globalTracker.Steps) {
		globalTracker.Steps[globalTracker.CurrentStep].CurrentSub = currentSub
		globalTracker.Steps[globalTracker.CurrentStep].Description = description
		globalTracker.CurrentMessage = description
		time.Sleep(100 * time.Millisecond) // Short pause to simulate loading
		updateUI()
	}
}

// CompleteCurrentStep marks the current step as complete and moves to the next
func CompleteCurrentStep() {
	if globalTracker.CurrentStep < len(globalTracker.Steps) {
		// Mark current step as complete
		step := &globalTracker.Steps[globalTracker.CurrentStep]
		step.CurrentSub = step.SubSteps

		// Move to next step
		globalTracker.CurrentStep++

		if globalTracker.CurrentStep < len(globalTracker.Steps) {
			nextStep := globalTracker.Steps[globalTracker.CurrentStep]
			globalTracker.CurrentStage = nextStep.Name
			globalTracker.CurrentMessage = nextStep.Description
		}

		updateUI()
	}
}

// calculateCurrentPercentage calculates the overall progress percentage
func calculateCurrentPercentage() int {
	if globalTracker.TotalWeight == 0 {
		return 0
	}

	completedWeight := 0.0

	// Add weight from completed steps
	for i := 0; i < globalTracker.CurrentStep && i < len(globalTracker.Steps); i++ {
		completedWeight += globalTracker.Steps[i].Weight
	}

	// Add partial weight from current step
	if globalTracker.CurrentStep < len(globalTracker.Steps) {
		currentStep := globalTracker.Steps[globalTracker.CurrentStep]
		if currentStep.SubSteps > 0 {
			stepProgress := float64(currentStep.CurrentSub) / float64(currentStep.SubSteps)
			completedWeight += stepProgress * currentStep.Weight
		}
	}

	percentage := int((completedWeight / globalTracker.TotalWeight) * 100)
	if percentage > 100 {
		percentage = 100
	}

	return percentage
}

// GetProgress returns the current progress state
func GetProgress() (int, string, string) {
	percentage := calculateCurrentPercentage()
	return percentage, globalTracker.CurrentStage, globalTracker.CurrentMessage
}

// Reset resets the progress tracker
func Reset() {
	globalTracker = ProgressTracker{}
}
