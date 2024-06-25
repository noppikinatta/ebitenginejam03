package asset

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed font/Mplus2-Regular.ttf
	mplus2regularttf []byte

	Mplus2RegularFaceSource *text.GoTextFaceSource

	fontFaces map[float64]text.Face

	//go:embed lang/*.csv
	langDir embed.FS

	langTamplates map[string]map[string]string
)

func init() {
	var err error

	Mplus2RegularFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(mplus2regularttf))
	if err != nil {
		log.Fatal("cannot read embedded font:", err)
	}

	fontFaces = make(map[float64]text.Face)

	langTamplates = make(map[string]map[string]string)
	const langDirPath = "lang"
	entries, err := langDir.ReadDir(langDirPath)
	if err != nil {
		log.Fatal("canot open language directory:", err)
	}

	for _, e := range entries {
		err := initTemplates(e, langDirPath)
		if err != nil {
			log.Fatal("cannot load templates for ", e.Name(), ":", err)
		}
	}
}

func FontFace(size float64) text.Face {
	f, ok := fontFaces[size]
	if !ok {
		f = &text.GoTextFace{
			Source: Mplus2RegularFaceSource,
			Size:   size,
		}
		fontFaces[size] = f
	}

	return f
}

func initTemplates(entry fs.DirEntry, dirPath string) error {
	fpath := path.Join(dirPath, entry.Name())
	data, err := loadCSV(fpath)
	if err != nil {
		return err
	}

	langTamplates[strings.TrimSuffix(entry.Name(), ".csv")] = data
	return nil
}

func loadCSV(filepath string) (map[string]string, error) {
	f, err := langDir.Open(filepath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	r := csv.NewReader(f)
	r.Comment = '#'
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	allContents, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, line := range allContents {
		if len(line) != 2 {
			return nil, fmt.Errorf("line %d length should 2 but %d, data: %v", i, len(line), line)
		}
		m[line[0]] = strings.ReplaceAll(line[1], `\n`, "\n")
	}

	return m, nil
}

func LoadTemplates() map[string]map[string]string {
	return langTamplates
}
