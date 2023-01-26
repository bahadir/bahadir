package tilemap

import (
	"encoding/json"
	"os"
)

type TileMap struct {
	CompressionLevel int  `json:"compressionlevel"`
	Height           int  `json:"height"`
	Infinite         bool `json:"infinite"`
	Layers           []struct {
		Data    []int   `json:"data"`
		Height  int     `json:"height"`
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		OffsetX int     `json:"offsetx"`
		OffsetY int     `json:"offsety"`
		Opacity float64 `json:"opacity"`
		Type    string  `json:"type"`
		Visible bool    `json:"visible"`
		Width   int     `json:"width"`
		X       int     `json:"x"`
		Y       int     `json:"y"`
	} `json:"layers"`
	NextLayerId  int    `json:"nextlayerid"`
	NextObjectId int    `json:"nextobjectid"`
	Orientation  string `json:"orientation"`
	RenderOrder  string `json:"renderorder"`
	TiledVersion string `json:"tiledversion"`
	TileHeight   int    `json:"tileheight"`
	TileSets     []struct {
		Columns     int    `json:"columns"`
		FirstIndex  int    `json:"firstgid"`
		Image       string `json:"image"`
		ImageHeight int    `json:"imageheight"`
		ImageWidth  int    `json:"imagewidth"`
		Margin      int    `json:"margin"`
		Name        string `json:"name"`
		Spacing     int    `json:"spacing"`
		TileCount   int    `json:"tilecount"`
		TileHeight  int    `json:"tileheight"`
		TileWidth   int    `json:"tilewidth"`
	} `json:"tilesets"`
	TileWidth int    `json:"tilewidth"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Width     int    `json:"width"`
}

func NewFromJSON(filePath string) (*TileMap, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data TileMap

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
