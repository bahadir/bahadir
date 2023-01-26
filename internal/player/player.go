package player

import (
	"encoding/json"
	"os"
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
