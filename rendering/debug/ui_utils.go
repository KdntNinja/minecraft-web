package debug

import (
	"bytes"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// loadFont creates a font face with the specified size
func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

// createPanelBackground creates a sophisticated background with border effect
func createPanelBackground() *eimage.NineSlice {
	// Create a more sophisticated background with border
	img := ebiten.NewImage(20, 20)

	// Main background
	img.Fill(color.NRGBA{15, 15, 25, 240})

	// Border effect
	for i := 0; i < 2; i++ {
		// Top border
		for x := i; x < 20-i; x++ {
			img.Set(x, i, color.NRGBA{60, 120, 180, 200})
			img.Set(x, 19-i, color.NRGBA{60, 120, 180, 200})
		}
		// Side borders
		for y := i; y < 20-i; y++ {
			img.Set(i, y, color.NRGBA{60, 120, 180, 200})
			img.Set(19-i, y, color.NRGBA{60, 120, 180, 200})
		}
	}

	return eimage.NewNineSliceSimple(img, 3, 3)
}

// createSection creates section containers with title and separator
func createSection(title string, font text.Face, titleColor color.Color) *widget.Container {
	section := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// Section title
	titleLabel := widget.NewText(
		widget.TextOpts.Text(title, font, titleColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(280, 16),
		),
	)
	section.AddChild(titleLabel)

	// Add separator line
	separator := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{60, 120, 180, 100})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(280, 1),
		),
	)
	section.AddChild(separator)

	return section
}

// createInfoLabel creates info labels with consistent styling
func createInfoLabel(text string, font text.Face) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text(text, font, color.NRGBA{220, 220, 220, 255}),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(280, 14),
		),
	)
}

// createGraphContainer creates a container for embedded performance graphs
func createGraphContainer(title string, width, height int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{25, 25, 35, 200})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(width, height),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(4)),
		)),
	)

	// Create small font for graph title
	smallFont, _ := loadFont(10)

	// Add title for the graph
	titleLabel := widget.NewText(
		widget.TextOpts.Text(title, smallFont, color.NRGBA{180, 180, 180, 255}),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(width-8, 12),
		),
	)
	container.AddChild(titleLabel)

	return container
}
