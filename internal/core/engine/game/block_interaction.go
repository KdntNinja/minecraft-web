package game

import (
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// handleBlockInteraction processes block interaction events from the player
func (g *Game) handleBlockInteraction(p *player.Player, interaction *player.BlockInteraction) {
	switch interaction.Type {
	case player.BreakBlock:
		g.World.BreakBlock(interaction.BlockX, interaction.BlockY)
	case player.PlaceBlock:
		g.World.PlaceBlock(interaction.BlockX, interaction.BlockY, p.SelectedBlock)
	}
}
