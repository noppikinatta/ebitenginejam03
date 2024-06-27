package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/noppikinatta/ebitenginejam03/asset"
)

//go:generate go run .
//go:generate go fmt ../...

func main() {
	generateTextKeys()
	generateImageKeys()
}

func generateTextKeys() {
	langs := asset.LoadTemplates()
	english := langs["english"]

	generateKeys("text", english)
}

func generateImageKeys() {
	imgs := asset.Images()
	generateKeys("img", imgs)
}

func generateKeys[T any](prefix string, data map[string]T) {
	buf := strings.Builder{}

	buf.WriteString("package name\n")
	buf.WriteString("const (\n")

	kk := make([]string, 0, len(data))
	for k := range data {
		kk = append(kk, k)
	}
	sort.Strings(kk)

	prefixTitle := title(prefix)
	for _, k := range kk {
		s := fmt.Sprintf("%sKey%s = \"%s\"", prefixTitle, varName(k), k)
		buf.WriteString(s)
		buf.WriteString("\n")
	}

	buf.WriteString(")\n")
	out := fmt.Sprintf("../%skeys.go", prefix)
	os.WriteFile(out, []byte(buf.String()), 0644)
}

func varName(key string) string {
	ss := strings.Split(key, "-")
	buf := strings.Builder{}

	for _, s := range ss {
		if len(s) == 0 {
			continue
		}
		buf.WriteString(title(s))
	}

	return buf.String()
}

func title(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
