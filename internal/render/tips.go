package render

import (
	"time"
)

var tips = []string{
	"Star [this repository](https://github.com/bahadir/bahadir) to give XP to the hero. Hero needs XP to level up.",
	"Feel free to open a PR to add new elements to the map. Take a look at [Map Design](data/README.md) document.",
	"Map is larger than the viewport. Hero can move around the map.",
}

func Tip() string {
	days := int(time.Now().Unix() / (24 * 60 * 60))

	return tips[days%len(tips)]
}
