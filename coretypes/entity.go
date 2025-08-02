package coretypes

// Entity interface that all game objects must implement
type Entity interface {
	Update()
	ClampX(min, max float64)
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
}

// Entities is a slice of all entities in the world
type Entities []Entity
