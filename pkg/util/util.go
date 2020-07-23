/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"fmt"
	"os"
)

func EnsureDirectory(directory string) {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		panic(fmt.Sprintf("directory creation failed: %v\n", err)) //TODO err
	}
}
