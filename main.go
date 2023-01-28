package main

import (
	"fmt"
	"github.com/bahadir/bahadir/internal/game"
	"github.com/bahadir/bahadir/internal/player"
	"github.com/bahadir/bahadir/internal/render"
	"image/png"
	"os"
	"time"
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
	playerLevel := p.LevelUp()

	x, y := p.MoveCharacter()
	if g.CanWalk(x, y) {
		p.Position.X = x
		p.Position.Y = y
	}

	err = p.Save(fileNamePlayer)
	if err != nil {
		panic(err)
	}

	// Draw player on map
	g.PlacePlayer(p.Position.X, p.Position.Y, playerLevel)

	// Get the scaled image of the map
	imgBig := g.ImageSnapshot(p.Position.X, p.Position.Y, 15, 8, 3)

	// Create random file name
	fileName := fmt.Sprintf("data/map-%d.png", time.Now().UTC().UnixMilli())

	// Save genarated image
	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	err = png.Encode(out, imgBig)
	if err != nil {
		panic(err)
	}

	// Render readme
	err = render.Readme(fileName, "README.md")
	if err != nil {
		panic(err)
	}
}
