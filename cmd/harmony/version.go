package main

import (
	"fmt"
	"os"

	"github.com/PositionExchange/posichain/internal/cli"
	"github.com/spf13/cobra"
)

const (
	versionFormat = "Posichain (C) 2022. %v, version %v-%v (%v %v)"
)

// Version string variables
var (
	version string
	builtBy string
	builtAt string
	commit  string
)

var versionFlag = cli.BoolFlag{
	Name:      "version",
	Shorthand: "V",
	Usage:     "display version info",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version of the posichain binary",
	Long:  "print version of the posichain binary",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
		os.Exit(0)
	},
}

func getHarmonyVersion() string {
	return fmt.Sprintf(versionFormat, "posichain", version, commit, builtBy, builtAt)
}

func printVersion() {
	fmt.Println(getHarmonyVersion())
}
