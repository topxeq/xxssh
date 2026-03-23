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

	flag.Parse()

	if *showVersion {
		fmt.Println("xxssh version:", version)
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
