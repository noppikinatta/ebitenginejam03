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
	langs := asset.LoadTemplates()
	english := langs["english"]

	buf := strings.Builder{}

	buf.WriteString("package name\n")
	buf.WriteString("const (\n")

	kk := make([]string, 0, len(english))
	for k := range english {
		kk = append(kk, k)
	}
	sort.Strings(kk)

	for _, k := range kk {
		s := fmt.Sprintf("%s = \"%s\"", varName(k), k)
		buf.WriteString(s)
		buf.WriteString("\n")
	}

	buf.WriteString(")\n")
	out := "../textkeys.go"
	os.WriteFile(out, []byte(buf.String()), 0644)
}

func varName(key string) string {
	ss := strings.Split(key, "-")
	buf := strings.Builder{}

	for _, s := range ss {
		if len(s) == 0 {
			continue
		}
		buf.WriteString(strings.ToUpper(s[:1]))
		buf.WriteString(s[1:])
	}

	return buf.String()
}
