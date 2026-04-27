package batch_mail

import (
	"regexp"
	"strings"
)

var (
	fallbackRegex = regexp.MustCompile(`\{\{\s*(\w+)\s*\|\s*([^}]+?)\s*\}\}`)
	simpleVarRegex = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
)

type ColdTemplatePreprocessor struct{}

var defaultColdPreprocessor = &ColdTemplatePreprocessor{}

func GetColdTemplatePreprocessor() *ColdTemplatePreprocessor {
	return defaultColdPreprocessor
}

// Preprocess converts cold-mail syntax to Go template syntax
// {{ first_name | Friend }} → {{ with getCustom . "first_name" }}{{ . }}{{ else }}Friend{{ end }}
// {{ first_name }} → {{ with getCustom . "first_name" }}{{ . }}{{ end }}
func (p *ColdTemplatePreprocessor) Preprocess(content string) string {
	if strings.Contains(content, "{{ .") {
		return content
	}

	content = fallbackRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := fallbackRegex.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		varName := groups[1]
		fallback := strings.TrimSpace(groups[2])
		return `{{ with getCustom . "` + varName + `" }}{{ . }}{{ else }}` + fallback + `{{ end }}`
	})

	content = simpleVarRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := simpleVarRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		varName := groups[1]
		return `{{ with getCustom . "` + varName + `" }}{{ . }}{{ end }}`
	})

	return content
}
