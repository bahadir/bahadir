package main

import (
	"flag"
	"github.com/bahadir/bahadir/internal/player"
	"github.com/bahadir/bahadir/internal/sprite"
	"github.com/bahadir/bahadir/internal/tilemap"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fileName := "data/map.json"
	fileBase := filepath.Dir(fileName)

	jsonMap, err := tilemap.NewFromJSON(fileName)
	if err != nil {
		panic(err)
	}

	t := make(map[int]*sprite.Picture)

	for _, ts := range jsonMap.TileSets {
		s, err := sprite.NewFromFile(filepath.Join(fileBase, ts.Image), &sprite.Options{
			Width:  16,
			Height: 16,
			Margin: 1,
		})
		if err != nil {
			panic(err)
		}

		for i, p := range s.Pictures() {
			t[ts.FirstIndex+i] = p
		}
	}

	wk := make(map[int]bool)
	playerImages := make(map[int][]*sprite.Picture, 0)
	playerOffset := image.Point{X: 0, Y: 0}

	img := image.NewRGBA(image.Rect(0, 0, jsonMap.TileWidth*jsonMap.Width, jsonMap.TileHeight*jsonMap.Height))

	for _, layer := range jsonMap.Layers {
		if layer.Name == "Walkable" {
			for _, tm := range layer.Data {
				if tm == 0 {
					continue
				}

				wk[tm] = true
			}
		} else if strings.HasPrefix(layer.Name, "Player") {
			for i, tm := range layer.Data {
				if tm == 0 {
					continue
				}

				level := i / layer.Width

				playerOffset.X = layer.OffsetX
				playerOffset.Y = layer.OffsetY

				playerImages[level] = append(playerImages[level], t[tm])
			}
			continue
		} else if layer.Opacity == 0 {
			continue
		}

		for i, tm := range layer.Data {
			if tm == 0 {
				continue
			}

			pic := t[tm]
			rect := image.Rectangle{
				Min: image.Point{
					X: (i%layer.Width)*16 + layer.OffsetX,
					Y: (i/layer.Width)*16 + layer.OffsetY,
				},
				Max: image.Point{
					X: 16 + (i%layer.Width)*16 + layer.OffsetX,
					Y: 16 + (i/layer.Width)*16 + layer.OffsetY,
				},
			}
			draw.Draw(img, rect, pic.Image, pic.Point, draw.Over)
		}
	}

	fileNamePlayer := "data/player.json"
	p, err := player.NewFromJSON(fileNamePlayer)
	if err != nil {
		p = player.New()
	}

	// run commands

	var xp int
	var direction int
	flag.IntVar(&xp, "xp", 0, "set xp")
	flag.IntVar(&direction, "dir", 0, "move player to direction")
	flag.Parse()

	if xp > 0 {
		p.XP = xp
	}

	switch direction {
	case 1:
		p.Position.X++
	case 2:
		p.Position.Y++
	case 3:
		p.Position.X--
	case 4:
		p.Position.Y--
	}

	playerLevel := p.XP

	err = p.Save(fileNamePlayer)
	if err != nil {
		panic(err)
	}

	// draw player
	for _, pl := range playerImages[playerLevel] {
		rect := image.Rectangle{
			Min: image.Point{
				X: 16*(p.Position.X) + playerOffset.X,
				Y: 16*(p.Position.Y) + playerOffset.Y,
			},
			Max: image.Point{
				X: (p.Position.X+1)*16 + playerOffset.X,
				Y: (p.Position.Y+1)*16 + playerOffset.Y,
			},
		}

		draw.Draw(img, rect, pl.Image, pl.Point, draw.Over)
	}

	// generate small map
	imgSmall := image.NewRGBA(image.Rect(0, 0, 16*10, 16*6))

	playerX := 16 * p.Position.X
	leftPos := playerX - 16*5
	if leftPos < 0 {
		leftPos = 0
	}

	playerY := 16 * p.Position.Y
	topPos := playerY - 16*3
	if topPos < 0 {
		topPos = 0
	}

	pos := image.Point{X: leftPos, Y: topPos}
	draw.Draw(imgSmall, img.Bounds(), img, pos, draw.Over)

	imgBig := image.NewRGBA(image.Rect(0, 0, 16*10*4, 16*6*4))
	for y := 0; y < imgBig.Bounds().Dy(); y += 1 {
		for x := 0; x < imgBig.Bounds().Dx(); x += 1 {
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					imgBig.Set(x*4+i, y*4+j, imgSmall.At(x, y))
				}
			}
		}
	}

	// save img to file
	out, err := os.Create("data/map.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(out, imgBig)
	if err != nil {
		panic(err)
	}
}
