package sprite

import (
	"image"
	"os"

	_ "image/png"
)

type Sprite struct {
	img     image.Image
	options *Options
}

type Options struct {
	Width  int
	Height int
	Margin int
}

type Picture struct {
	Point image.Point
	Image image.Image
}

func NewFromFile(filePath string, options *Options) (*Sprite, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if options == nil {
		options = &Options{
			Width:  16,
			Height: 16,
			Margin: 0,
		}
	}

	s := &Sprite{
		img:     img,
		options: options,
	}

	return s, nil
}

func (s *Sprite) Pictures() []*Picture {
	var pictures []*Picture

	for y := 0; y < s.img.Bounds().Dy(); y += s.options.Height + s.options.Margin {
		for x := 0; x < s.img.Bounds().Dx(); x += s.options.Width + s.options.Margin {
			pictures = append(pictures, &Picture{
				Image: s.img,
				Point: image.Point{
					X: x,
					Y: y,
				},
			})
		}
	}

	return pictures
}
