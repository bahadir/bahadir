package render

import "math/rand"

var tips = []string{
	"Star [this repository](https://github.com/bahadir/bahadir) to give XP to the hero. Hero needs XP to level up.",
	"Feel free to open a PR to add new elements to the map. Take a look at [Map Design](data/README.md) document.",
}

func Tip() string {
	return tips[rand.Intn(len(tips))]
}
