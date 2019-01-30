package main

import (
	"fmt"
	"os"

	"github.com/elforg/elfplatform/core/config"
	"github.com/elforg/elfplatform/elfplatform/start"
	"github.com/spf13/cobra"
)

// Watching if the configuration has changed
const Watching = true

var mainCmd = &cobra.Command{Use: "elf"}

func main() {
	// read elf.yaml into viper
	if err := config.Initilize(mainCmd.Use, Watching); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mainCmd.AddCommand(start.Cmd())

	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
