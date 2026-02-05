package html

import (
	"embed"
	"strings"
)

//go:embed *.html
var templateHTML embed.FS

func GetTemplateHTML(templateName string) ([]byte, error) {
	content, err := templateHTML.ReadFile(strings.Join([]string{templateName, "html"}, "."))
	if err != nil {
		return nil, err
	}
	return content, nil
}
