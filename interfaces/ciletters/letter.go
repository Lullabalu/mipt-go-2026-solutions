//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed letter.tmpl
var letTem string

func MakeLetter(n *Notification) (string, error) {
	funcs := template.FuncMap{
		"gethash": func(s string) string {
			if len(s) <= 8 {
				return s
			}
			return s[:8]
		},
		"splitl": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"lastl": func(lines []string, n int) []string {
			if len(lines) <= n {
				return lines
			}
			return lines[len(lines)-n:]
		},
	}

	tem, err := template.New("letter").Funcs(funcs).Parse(letTem)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tem.Execute(&buf, n)
	if err != nil {
		return "something wrong", err
	}
	return buf.String(), nil

}
