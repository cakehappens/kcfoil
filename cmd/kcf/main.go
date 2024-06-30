package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"github.com/rs/zerolog"
	"golang.org/x/term"
)

var (
	colorValues = cli.PredictSet("auto", "always", "never")
	interactive = term.IsTerminal(int(os.Stdout.Fd()))
)

func main() {
	rootCmd := &cli.Command{
		Use: "kcf",
	}

	// set default logging level early; not all commands parse --log-level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	addCommandsWithLogLevelOption(
		rootCmd,
		templateCmd(),
	)

	// Run!
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func addCommandsWithLogLevelOption(rootCmd *cli.Command, cmds ...*cli.Command) {
	for _, cmd := range cmds {
		levels := []string{
			zerolog.Disabled.String(),
			zerolog.FatalLevel.String(),
			zerolog.ErrorLevel.String(),
			zerolog.WarnLevel.String(),
			zerolog.InfoLevel.String(),
			zerolog.DebugLevel.String(),
			zerolog.TraceLevel.String(),
		}
		cmd.Flags().String("log-level", zerolog.InfoLevel.String(), "possible values: "+strings.Join(levels, ", "))

		cmdRun := cmd.Run
		cmd.Run = func(cmd *cli.Command, args []string) error {
			level, err := zerolog.ParseLevel(cmd.Flags().Lookup("log-level").Value.String())
			if err != nil {
				return err
			}
			zerolog.SetGlobalLevel(level)

			return cmdRun(cmd, args)
		}
		rootCmd.AddCommand(cmd)
	}
}
