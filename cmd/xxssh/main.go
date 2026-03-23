package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/topxeq/xxssh/internal/store"
	"github.com/topxeq/xxssh/internal/termui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	resetConfig := flag.Bool("reset", false, "Reset all configuration")

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

	st, err := store.NewStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize store:", err)
		os.Exit(1)
	}

	ui := termui.NewUI(st)
	if err := ui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running app:", err)
		os.Exit(1)
	}
}
