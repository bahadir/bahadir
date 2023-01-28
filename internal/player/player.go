package player

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
)

type Player struct {
	XP       int `json:"xp"`
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
}

func New() *Player {
	return &Player{
		XP: 0,
		Position: struct {
			X int `json:"x"`
			Y int `json:"y"`
		}{
			X: 0,
			Y: 0,
		},
	}
}

func NewFromJSON(filePath string) (*Player, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data Player

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Save method
func (p *Player) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(p)
}

func (p *Player) MoveCharacter() (int, int) {
	input := strings.ToLower(os.Getenv("GAME_INPUT"))

	switch input {
	case "go right":
		return p.Position.X + 1, p.Position.X
	case "go left":
		return p.Position.X - 1, p.Position.Y
	case "go up":
		return p.Position.X, p.Position.Y - 1
	case "go down":
		return p.Position.X, p.Position.Y + 1
	}

	return p.Position.X, p.Position.Y
}

func (p *Player) LevelUp() int {
	xp, err := strconv.Atoi(os.Getenv("GAME_XP"))
	if err != nil {
		return 0
	}

	p.XP += xp

	return int(math.Sqrt(float64(xp)) * 0.4)
}
