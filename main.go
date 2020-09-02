/*
Copyright © 2020 Mateusz Kyc
*/
package main

import (
	"github.com/epiphany-platform/cli/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cmd.Execute()
}
