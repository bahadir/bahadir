package main

import (
	"github.com/bahadir/bahadir/internal/game"
	"github.com/bahadir/bahadir/internal/player"
	"image/png"
	"os"
)

func main() {
	g := game.New()

	// Load the map and draw base elements
	g.InitUsingTimemap("data/map.json")

	// Load player state
	fileNamePlayer := "data/player.json"
	p, err := player.NewFromJSON(fileNamePlayer)
	if err != nil {
		p = player.New()
	}

	// Process input
	p.ApplyInput()
	playerLevel := p.LevelUp()

	if !g.CanWalk(p.Position.X, p.Position.Y) {
		return
	}

	err = p.Save(fileNamePlayer)
	if err != nil {
		panic(err)
	}

	// Draw player on map
	g.PlacePlayer(p.Position.X, p.Position.Y, playerLevel)

	// Get the scaled image of the map
	imgBig := g.ImageSnapshot(p.Position.X, p.Position.Y, 15, 8, 3)

	// Save genarated image
	out, err := os.Create("data/map.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(out, imgBig)
	if err != nil {
		panic(err)
	}
}
