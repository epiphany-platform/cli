/*
 * Copyright © 2020 Mateusz Kyc
 */

package promptui

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"
	"sort"

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
	config, err := configuration.GetConfig()
	if err != nil {
		return uuid.Nil, err
	}
	keys := make([]string, 0)
	m := make(map[string]string)
	environments, err := environment.GetAllEnvironments()
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
