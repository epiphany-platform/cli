package processor

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/configuration"
	environments "github.com/epiphany-platform/cli/pkg/environment"
)

func init() {
	logger.Initialize()
}

func TemplateProcessor(config *configuration.Config, environment *environments.Environment) func(s string) string {
	return func(s string) string {
		if strings.Contains(s, "#") {
			parts := strings.Split(s, "#")
			logger.Debug().Msgf("and value has parts: %#v", parts)
			if len(parts) < 2 {
				return s
			}
			result := ""
			for i := 0; i < len(parts); i++ {
				switch p := parts[i]; p {
				case "Config":
					ii := i
					next := i + 1
					i++
					r, err := process(strconv.Itoa(ii), parts[next], config)
					if err != nil {
						logger.Error().Err(err)
						break
					}
					logger.Debug().Msgf("result value: %#v", r)
					result = result + r
				case "Environment":
					ii := i
					next := i + 1
					i++
					r, err := process(strconv.Itoa(ii), parts[next], environment)
					if err != nil {
						logger.Error().Err(err)
						break
					}
					logger.Debug().Msgf("result value: %#v", r)
					result = result + r
				default:
					result = result + p
				}
			}
			return result
		}
		return s
	}
}

func process(name, pattern string, data interface{}) (string, error) {
	t, err := template.New(name).Option("missingkey=error").Parse(pattern)
	if err != nil {
		logger.Error().Err(err)
		return "", err
	}
	logger.Debug().Msgf("parsed template: %#v", t)
	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		logger.Error().Err(err)
		return "", err
	}
	r := b.String()
	logger.Debug().Msgf("result value: %#v", r)
	return r, nil
}
