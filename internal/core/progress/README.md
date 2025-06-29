# Progress Tracking System

The WebCraft progress tracking system provides a dynamic, extensible way to report loading progress during game initialization. It's designed specifically for WASM deployment where users need visual feedback during potentially long loading processes.

## Overview

The progress system consists of three main components:

1. **ProgressTracker**: Core system that manages weighted progress steps
2. **Progress Steps**: Individual loading phases with configurable weights and substeps
3. **UI Integration**: Real-time progress updates to the web interface

## Architecture

### Core Components

- `progress_common.go`: Platform-independent progress tracking logic
- `progress.go`: WASM-specific UI integration using JavaScript calls
- `README.md`: This documentation file

### Key Structures

```go
type ProgressStep struct {
    Name        string  // Display name for the step
    Weight      float64 // Relative importance (higher = more progress bar movement)
    SubSteps    int     // Number of substeps within this step
    CurrentSub  int     // Current substep progress
    Description string  // Current status message
}

type ProgressTracker struct {
    Steps          []ProgressStep
    CurrentStep    int
    TotalWeight    float64
    CurrentStage   string
    CurrentMessage string
}
```

## Usage Guide

### Basic Setup

1. **Initialize Progress System**

   ```go
   steps := []progress.ProgressStep{
       {Name: "Loading Assets", Weight: 2.0, SubSteps: 10, Description: "Loading game assets..."},
       {Name: "World Generation", Weight: 5.0, SubSteps: 1, Description: "Generating terrain..."},
       {Name: "Finalization", Weight: 1.0, SubSteps: 3, Description: "Finalizing setup..."},
   }
   progress.InitializeProgress(steps)
   ```

2. **Update Progress Within Steps**

   ```go
   // Update substeps dynamically (e.g., when you know how many chunks to generate)
   progress.SetCurrentStepSubSteps(25, "Generating 25 terrain chunks...")
   
   // Update progress within current step
   for i := 0; i < 25; i++ {
       // Do actual work here...
       progress.UpdateCurrentStepProgress(i+1, fmt.Sprintf("Generated chunk %d/%d", i+1, 25))
   }
   ```

3. **Complete Steps**

   ```go
   // Move to next step
   progress.CompleteCurrentStep()
   ```

### Weight System

The weight system allows you to allocate progress bar movement proportionally:

- **Weight 1.0**: Standard step
- **Weight 2.0**: Takes twice as much progress bar space
- **Weight 0.5**: Takes half as much progress bar space

Example:

```go
steps := []progress.ProgressStep{
    {Name: "Quick Setup", Weight: 1.0, SubSteps: 1, Description: "Fast initialization"},
    {Name: "Heavy Processing", Weight: 8.0, SubSteps: 100, Description: "Resource-intensive task"},
    {Name: "Cleanup", Weight: 1.0, SubSteps: 1, Description: "Final cleanup"},
}
// Total weight: 10.0
// Quick Setup: 10% of progress bar
// Heavy Processing: 80% of progress bar  
// Cleanup: 10% of progress bar
```

## Advanced Features

### Dynamic Substep Allocation

You can change the number of substeps during execution:

```go
// Initially unknown number of substeps
{Name: "Dynamic Task", Weight: 3.0, SubSteps: 1, Description: "Calculating work..."}

// Later, when you know the actual work required:
chunkCount := calculateChunksNeeded()
progress.SetCurrentStepSubSteps(chunkCount, fmt.Sprintf("Processing %d chunks...", chunkCount))
```

### Progress Monitoring

Get current progress state:

```go
percentage, stage, message := progress.GetProgress()
fmt.Printf("Progress: %d%% - %s: %s\n", percentage, stage, message)
```

### Reset System

Reset for new loading sequences:

```go
progress.Reset()
```

## Integration with WASM UI

### JavaScript Integration

The progress system automatically communicates with the web interface through:

1. **Console Logging**: `[PROGRESS]` tagged messages for debugging
2. **JavaScript Calls**: Direct `updateLoadingProgress()` function calls
3. **Real-time Updates**: Immediate UI updates on every progress change

### UI Components

The WASM `index.html` includes:

- **Progress Bar**: Visual percentage indicator
- **Stage Display**: Current step name
- **Message Display**: Current detailed message
- **Log Area**: Scrollable progress history
- **Debug Controls**: Manual skip option for testing

### Message Format

Console messages follow this pattern:

```text
[PROGRESS] Stage Name: Detailed message (XX%)
```

Example:

```text
[PROGRESS] Generating Terrain: Generated chunk 15/25 (67%)
```

## Extending the System

1. **Define the step** in your initialization:

   ```go
   {Name: "Custom Process", Weight: 2.5, SubSteps: 10, Description: "Running custom logic..."}
   ```

2. **Implement progress updates** in your code:

   ```go
   for i := 0; i < customWorkItems; i++ {
       // Perform work
       doCustomWork(i)
       
       // Report progress
       progress.UpdateCurrentStepProgress(i+1, 
           fmt.Sprintf("Processed custom item %d/%d", i+1, customWorkItems))
   }
   progress.CompleteCurrentStep()
   ```

### Custom Progress Patterns

#### File Loading Pattern

```go
files := []string{"texture1.png", "sound1.wav", "model1.obj"}
progress.SetCurrentStepSubSteps(len(files), "Loading game assets...")

for i, file := range files {
    loadFile(file)
    progress.UpdateCurrentStepProgress(i+1, fmt.Sprintf("Loaded %s", file))
}
progress.CompleteCurrentStep()
```

#### Batch Processing Pattern

```go
batchSize := 10
totalItems := 100
batches := (totalItems + batchSize - 1) / batchSize

progress.SetCurrentStepSubSteps(batches, "Processing items in batches...")

for batch := 0; batch < batches; batch++ {
    start := batch * batchSize
    end := min(start+batchSize, totalItems)
    
    processBatch(start, end)
    progress.UpdateCurrentStepProgress(batch+1, 
        fmt.Sprintf("Processed batch %d/%d (%d-%d)", batch+1, batches, start+1, end))
}
progress.CompleteCurrentStep()
```

## Best Practices

### 1. Meaningful Weights

Assign weights based on actual expected time:

```go
// Good: Reflects actual time requirements
{Name: "Config Loading", Weight: 0.5, SubSteps: 1},     // Very fast
{Name: "Terrain Generation", Weight: 8.0, SubSteps: 1}, // Slow process
{Name: "UI Setup", Weight: 1.0, SubSteps: 1},          // Moderate

// Avoid: All equal weights when tasks take very different times
{Name: "Config Loading", Weight: 1.0, SubSteps: 1},     // Misleading
{Name: "Terrain Generation", Weight: 1.0, SubSteps: 1}, // Misleading
```

### 2. Informative Messages

Provide specific, actionable information:

```go
// Good: Specific and informative
progress.UpdateCurrentStepProgress(15, "Generated chunk 15/25 at coordinates (240, 160)")

// Avoid: Vague or generic
progress.UpdateCurrentStepProgress(15, "Working...")
```

### 3. Appropriate Substep Granularity

Balance update frequency with performance:

```go
// Good: Reasonable update frequency
for chunkIndex := 0; chunkIndex < 25; chunkIndex++ {
    generateChunk(chunkIndex)
    progress.UpdateCurrentStepProgress(chunkIndex+1, ...)
}

// Avoid: Too frequent updates (performance impact)
for blockX := 0; blockX < 1000; blockX++ {
    for blockY := 0; blockY < 1000; blockY++ {
        progress.UpdateCurrentStepProgress(...) // Called 1M times!
    }
}
```

### 4. Error Handling

Include error states in progress reporting:

```go
if err := riskyOperation(); err != nil {
## Debugging

### Console Output

Monitor progress in browser console:
```

## Debugging

### Console Output

Monitor progress in browser console:

```log
[PROGRESS] Initializing: Starting game initialization... (5%)
[PROGRESS] Initializing: Generated new world seed (10%)
[PROGRESS] World Setup: Setting up world structure... (15%)
```

### Manual Testing

Use the "Skip Loading Screen" button in development to bypass long loading times.

### Progress Validation

Verify progress calculations:

```go
If migrating from a simpler progress system:

### Old Code
```

## Migration from Old System

If migrating from a simpler progress system:

### Old Code

```go
UpdateProgress(5, 10, "Loading", "Loading textures...")  // 50%
UpdateProgress(10, 10, "Loading", "Loading complete!")  // 100%
```

### New Code

```go
// Setup
steps := []progress.ProgressStep{
    {Name: "Loading", Weight: 1.0, SubSteps: 10, Description: "Loading game assets..."},
}
progress.InitializeProgress(steps)

// Usage
progress.UpdateCurrentStepProgress(5, "Loading textures...")  // 50%
progress.UpdateCurrentStepProgress(10, "Loading complete!")  // 100%
progress.CompleteCurrentStep()
```

## Future Enhancements

Potential improvements to consider:

- Parallel step execution tracking
- Progress persistence across sessions
- Estimated time remaining calculations
- Progress analytics and optimization
- Custom UI theme support

Potential improvements to consider:

- Parallel step execution tracking
- Progress persistence across sessions
- Estimated time remaining calculations
- Progress analytics and optimization
- Custom UI theme support

---

*This documentation covers the complete WebCraft progress tracking system. For implementation examples, see `game.go` and `world.go` for real-world usage patterns.*
