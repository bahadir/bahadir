package game

import (
	"github.com/bahadir/bahadir/internal/sprite"
	"github.com/bahadir/bahadir/internal/tilemap"
	"image"
	"image/draw"
	"path/filepath"
	"strings"
)

type Game struct {
	pictureMap   map[int]*sprite.Picture
	playerImages map[int][]*sprite.Picture
	playerOffset image.Point

	walkableObjects map[int]bool
	obstacleObjects map[int]bool

	walkableCoords [][]bool

	tileSize int

	imgSize image.Point
	imgMap  draw.Image
}

func New() *Game {
	return &Game{
		pictureMap:   make(map[int]*sprite.Picture),
		playerImages: make(map[int][]*sprite.Picture),
		playerOffset: image.Point{X: 0, Y: 0},

		walkableObjects: make(map[int]bool),
		obstacleObjects: make(map[int]bool),

		walkableCoords: make([][]bool, 0),

		tileSize: 16,
	}
}

func (g *Game) InitUsingTimemap(filePath string) {
	fileDir := filepath.Dir(filePath)

	jsonMap, err := tilemap.NewFromJSON(filePath)
	if err != nil {
		panic(err)
	}

	g.tileSize = jsonMap.TileWidth

	g.extractTiles(jsonMap, fileDir)

	g.walkableCoords = make([][]bool, jsonMap.Height)
	for y := 0; y < jsonMap.Height; y++ {
		g.walkableCoords[y] = make([]bool, jsonMap.Width)
	}

	g.imgSize = image.Point{X: g.tileSize * jsonMap.Width, Y: g.tileSize * jsonMap.Height}

	g.imgMap = image.NewRGBA(image.Rect(0, 0, g.imgSize.X, g.imgSize.Y))

	g.loadLayers(jsonMap)
}

func (g *Game) loadLayers(jsonMap *tilemap.TileMap) {
	for _, layer := range jsonMap.Layers {
		if strings.HasPrefix(layer.Name, LayerInterractionPrefix) {
			g.loadInteractionObjects(layer)
		} else if strings.HasPrefix(layer.Name, LayerPlayerPrefix) {
			g.loadPlayerObjects(layer)
		} else if layer.Opacity > 0 {
			g.loadLevelDesign(layer)
		}
	}
}

func (g *Game) extractTiles(jsonMap *tilemap.TileMap, fileBase string) {
	for _, ts := range jsonMap.TileSets {
		s, err := sprite.NewFromFile(filepath.Join(fileBase, ts.Image), &sprite.Options{
			Width:   ts.TileWidth,
			Height:  ts.TileHeight,
			Spacing: ts.Spacing,
		})
		if err != nil {
			panic(err)
		}

		for i, p := range s.Pictures() {
			g.pictureMap[ts.FirstIndex+i] = p
		}
	}
}

func (g *Game) loadInteractionObjects(layer tilemap.TileLayer) {
	for i, tm := range layer.Data {
		if tm == 0 {
			continue
		}

		row := i / layer.Width

		switch row {
		case 0, 1:
			g.walkableObjects[tm] = true
		case 2, 3:
			g.obstacleObjects[tm] = true
		}
	}
}

func (g *Game) loadPlayerObjects(layer tilemap.TileLayer) {
	for i, tm := range layer.Data {
		if tm == 0 {
			continue
		}

		level := i / layer.Width

		g.playerOffset.X = layer.OffsetX
		g.playerOffset.Y = layer.OffsetY

		g.playerImages[level] = append(g.playerImages[level], g.pictureMap[tm])
	}
}

func (g *Game) loadLevelDesign(layer tilemap.TileLayer) {
	for i, tm := range layer.Data {
		if tm == 0 {
			continue
		}

		row := i / layer.Width
		col := i % layer.Width

		if !g.walkableCoords[row][col] {
			if walkable, ok := g.walkableObjects[tm]; ok && walkable {
				g.walkableCoords[row][col] = true
			}
		} else {
			if obstacle, ok := g.obstacleObjects[tm]; ok && obstacle {
				g.walkableCoords[row][col] = false
			}
		}

		pic := g.pictureMap[tm]

		rect := image.Rectangle{
			Min: image.Point{
				X: (col)*g.tileSize + layer.OffsetX,
				Y: (row)*g.tileSize + layer.OffsetY,
			},
			Max: image.Point{
				X: (1+col)*g.tileSize + layer.OffsetX,
				Y: (1+row)*g.tileSize + layer.OffsetY,
			},
		}

		draw.Draw(g.imgMap, rect, pic.Image, pic.Point, draw.Over)
	}
}

func (g *Game) ImageSnapshot(centerX, centerY, width, height, scale int) draw.Image {
	imgSmall := image.NewRGBA(image.Rect(0, 0, g.tileSize*width, g.tileSize*height))

	leftPos := g.tileSize*centerX - width*g.tileSize/2
	if leftPos < 0 {
		leftPos = 0
	} else if leftPos+g.tileSize*width > g.imgSize.X {
		leftPos = g.imgSize.X - g.tileSize*width
	}

	topPos := g.tileSize*centerY - height*g.tileSize/2
	if topPos < 0 {
		topPos = 0
	} else if topPos+g.tileSize*height > g.imgSize.Y {
		topPos = g.imgSize.Y - g.tileSize*height
	}

	draw.Draw(imgSmall, g.imgMap.Bounds(), g.imgMap, image.Point{leftPos, topPos}, draw.Over)

	imgBig := image.NewRGBA(image.Rect(0, 0, g.tileSize*width*scale, g.tileSize*height*scale))
	for y := 0; y < imgBig.Bounds().Dy(); y += 1 {
		for x := 0; x < imgBig.Bounds().Dx(); x += 1 {
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					imgBig.Set(x*scale+i, y*scale+j, imgSmall.At(x, y))
				}
			}
		}
	}

	return imgBig
}

func (g *Game) PlacePlayer(x, y, level int) {
	if level >= len(g.playerImages) {
		level = len(g.playerImages) - 1
	}

	for _, pl := range g.playerImages[level] {
		rect := image.Rectangle{
			Min: image.Point{
				X: g.tileSize*x + g.playerOffset.X,
				Y: g.tileSize*y + g.playerOffset.Y,
			},
			Max: image.Point{
				X: (x+1)*g.tileSize + g.playerOffset.X,
				Y: (y+1)*g.tileSize + g.playerOffset.Y,
			},
		}

		draw.Draw(g.imgMap, rect, pl.Image, pl.Point, draw.Over)
	}
}

func (g *Game) CanWalk(x, y int) bool {
	if x < 0 || y < 0 || y >= len(g.walkableCoords) || x >= len(g.walkableCoords[y]) {
		return false
	}

	return g.walkableCoords[y][x]
}
