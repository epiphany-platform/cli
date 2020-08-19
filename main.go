/*
Copyright © 2020 Mateusz Kyc
*/
package main

import (
	"github.com/mkyc/epiphany-wrapper-poc/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cmd.Execute()
}
