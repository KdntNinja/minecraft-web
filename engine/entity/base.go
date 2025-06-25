package entity

type EntityBase struct {
	X, Y          float64
	VX, VY        float64
	Width, Height int
}

func (e *EntityBase) GetPosition() (float64, float64) {
	return e.X, e.Y
}

func (e *EntityBase) SetPosition(x, y float64) {
	e.X = x
	e.Y = y
}

func (e *EntityBase) Update() {}

// Entities slice for world

type Entities []Entity
