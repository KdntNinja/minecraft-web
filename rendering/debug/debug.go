package debug

import (
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// DebugUI manages an advanced EbitenUI-based debug panel
type DebugUI struct {
	UI        *ebitenui.UI
	mainPanel *widget.Container

	// Performance metrics
	fpsLabel       *widget.Text
	fpsMinMaxLabel *widget.Text
	tickTimeLabel  *widget.Text
	memoryLabel    *widget.Text
	gcLabel        *widget.Text

	// World info
	playerLabel *widget.Text
	chunkLabel  *widget.Text
	cameraLabel *widget.Text
	seedLabel   *widget.Text
	entityLabel *widget.Text

	// History data
	fpsHistory  []float64
	tickHistory []float64
	memHistory  []float64

	// UI state
	collapsed      bool
	lastUpdateTime time.Time

	// Fonts
	headerFont text.Face
	normalFont text.Face
	smallFont  text.Face
}

var debugUI *DebugUI

// InitDebugUI initializes the advanced debug UI panel
func InitDebugUI() error {
	// Load different font sizes
	headerFont, err := loadFont(14)
	if err != nil {
		return err
	}

	normalFont, err := loadFont(12)
	if err != nil {
		return err
	}

	smallFont, err := loadFont(10)
	if err != nil {
		return err
	}

	// Create root container with anchor layout
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Create main debug panel with modern styling
	mainPanel := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(createPanelBackground()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(300, 400),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(3),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(12)),
		)),
	)

	// Header section
	headerContainer := createSection("WEBCRAFT DEBUG", headerFont, color.NRGBA{100, 200, 255, 255})
	mainPanel.AddChild(headerContainer)

	// Performance section
	perfSection := createSection("PERFORMANCE", normalFont, color.NRGBA{255, 200, 100, 255})

	fpsLabel := createInfoLabel("FPS: 0.0", normalFont)
	fpsMinMaxLabel := createInfoLabel("Min/Max: 0.0/0.0", smallFont)
	tickTimeLabel := createInfoLabel("Tick: 0.0ms", smallFont)
	memoryLabel := createInfoLabel("Memory: 0MB / 0MB", smallFont)
	gcLabel := createInfoLabel("GC: 0.0%", smallFont)

	perfSection.AddChild(fpsLabel)
	perfSection.AddChild(fpsMinMaxLabel)
	perfSection.AddChild(tickTimeLabel)
	perfSection.AddChild(memoryLabel)
	perfSection.AddChild(gcLabel)
	mainPanel.AddChild(perfSection)

	// World section
	worldSection := createSection("WORLD INFO", normalFont, color.NRGBA{100, 255, 100, 255})

	playerLabel := createInfoLabel("Player: N/A", smallFont)
	chunkLabel := createInfoLabel("Chunks: N/A", smallFont)
	cameraLabel := createInfoLabel("Camera: N/A", smallFont)
	seedLabel := createInfoLabel("Seed: N/A", smallFont)
	entityLabel := createInfoLabel("Entities: 0", smallFont)

	worldSection.AddChild(playerLabel)
	worldSection.AddChild(chunkLabel)
	worldSection.AddChild(cameraLabel)
	worldSection.AddChild(seedLabel)
	worldSection.AddChild(entityLabel)
	mainPanel.AddChild(worldSection)

	// Controls section
	controlsSection := createSection("CONTROLS", normalFont, color.NRGBA{255, 150, 150, 255})

	// Add multiple control instructions
	controls := []string{
		"F3: Toggle Debug",
		"WASD: Move Player",
		"Mouse: Look/Aim",
		"Left Click: Break Block",
		"Right Click: Place Block",
		"1-9,0: Select Block",
		"ESC: Pause Game",
		"F11: Toggle Fullscreen",
	}

	for _, control := range controls {
		controlLabel := createInfoLabel(control, smallFont)
		controlsSection.AddChild(controlLabel)
	}
	mainPanel.AddChild(controlsSection)

	// Performance Graphs section
	graphsSection := createSection("PERFORMANCE GRAPHS", normalFont, color.NRGBA{150, 255, 150, 255})

	// Create containers for embedded graphs
	fpsGraphContainer := createGraphContainer("FPS Graph", 280, 80)
	tickGraphContainer := createGraphContainer("Tick Time Graph", 280, 80)
	memGraphContainer := createGraphContainer("Memory Graph", 280, 80)

	graphsSection.AddChild(fpsGraphContainer)
	graphsSection.AddChild(tickGraphContainer)
	graphsSection.AddChild(memGraphContainer)
	mainPanel.AddChild(graphsSection)

	rootContainer.AddChild(mainPanel)

	debugUI = &DebugUI{
		UI:             &ebitenui.UI{Container: rootContainer},
		mainPanel:      mainPanel,
		fpsLabel:       fpsLabel,
		fpsMinMaxLabel: fpsMinMaxLabel,
		tickTimeLabel:  tickTimeLabel,
		memoryLabel:    memoryLabel,
		gcLabel:        gcLabel,
		playerLabel:    playerLabel,
		chunkLabel:     chunkLabel,
		cameraLabel:    cameraLabel,
		seedLabel:      seedLabel,
		entityLabel:    entityLabel,
		fpsHistory:     make([]float64, 0, 120),
		tickHistory:    make([]float64, 0, 120),
		memHistory:     make([]float64, 0, 120),
		collapsed:      false,
		lastUpdateTime: time.Now(),
		headerFont:     headerFont,
		normalFont:     normalFont,
		smallFont:      smallFont,
	}
	return nil
}

// DrawDebugUI renders the advanced debug UI panel and embedded performance graphs
func DrawDebugUI(screen *ebiten.Image) {
	if debugUI != nil {
		debugUI.UI.Update()
		debugUI.UI.Draw(screen)

		// Draw embedded performance graphs inside the panel
		DrawPerformanceGraphs(screen)
	}
}
