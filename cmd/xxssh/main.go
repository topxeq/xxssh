package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/topxeq/xxssh/internal/store"
	"github.com/topxeq/xxssh/internal/tui"
)

var version = "dev"

func main() {
	noColor := flag.Bool("no-color", false, "Disable color output")
	showVersion := flag.Bool("version", false, "Show version")
	resetConfig := flag.Bool("reset", false, "Reset all configuration (delete all servers and master password)")

	flag.Parse()

	if *showVersion {
		fmt.Println("xxssh version:", version)
		os.Exit(0)
	}

	if *resetConfig {
		st, err := store.NewStore()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to initialize store:", err)
			os.Exit(1)
		}
		if err := st.Reset(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to reset configuration:", err)
			os.Exit(1)
		}
		fmt.Println("Configuration reset successfully.")
		os.Exit(0)
	}

	if *noColor {
		tui.SetForceColor(false)
	}

	st, err := store.NewStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize store:", err)
		os.Exit(1)
	}

	app := tui.NewApp(st)
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running app:", err)
		os.Exit(1)
	}
}
