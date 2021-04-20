package promptui

import (
	"errors"
	"fmt"
	"sort"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

func PromptForString(label string) (string, error) {
	validator := func(input string) error {
		if len(input) < 1 {
			return errors.New("too short")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func PromptForEnvironmentSelect(label string) (uuid.UUID, error) {
	//TODO fix it not to call config and environments here
	config, err := configuration.GetConfig()
	if err != nil {
		return uuid.Nil, err
	}
	keys := make([]string, 0)
	m := make(map[string]string)
	environments, err := environment.GetAll()
	if err != nil {
		return uuid.Nil, err
	}
	for _, e := range environments {
		keys = append(keys, e.Uuid.String())
		m[e.Uuid.String()] = e.Name
	}
	sort.Strings(keys)
	var values []string
	cursorPosition := -1
	for i, k := range keys {
		if k == config.CurrentEnvironment.String() {
			values = append(values, fmt.Sprintf("%s (%s, current)", m[k], k))
			cursorPosition = i
		} else {
			values = append(values, fmt.Sprintf("%s (%s)", m[k], k))
		}
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     values,
		Size:      len(keys),
		CursorPos: cursorPosition,
	}
	c, _, err := prompt.Run()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.MustParse(keys[c]), nil
}
