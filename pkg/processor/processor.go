package processor

import (
	"bytes"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	logger zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Str("package", "processor").Caller().Timestamp().Logger()
}

func TemplateProcessor(c *configuration.Config, e *environment.Environment) func(s string) string {
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
					r, err := process(strconv.Itoa(ii), parts[next], c)
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
					r, err := process(strconv.Itoa(ii), parts[next], e)
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
	t, err := template.New(name).Parse(pattern)
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
