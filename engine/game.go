package engine

import "log"

type Game struct{}

func NewGame() *Game {
	log.Println("NewGame created")
	return &Game{}
}
