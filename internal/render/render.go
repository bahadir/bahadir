package render

import (
	_ "embed"
	"html/template"
	"os"
)

//go:embed README.md.tmpl
var readmeTmpl string

func Readme(imagePath, readmePath string) error {
	// render readme.tmpl to filepath
	tmpl, err := template.New("readme").Parse(readmeTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(readmePath)
	if err != nil {
		return err
	}

	data := struct {
		ImagePath string
	}{
		ImagePath: imagePath,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
