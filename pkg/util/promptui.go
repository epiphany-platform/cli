/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"errors"

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
