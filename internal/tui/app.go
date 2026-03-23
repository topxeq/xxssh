package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/topxeq/xxssh/internal/config"
	"github.com/topxeq/xxssh/internal/store"
)

type App struct {
	app   *tview.Application
	pages *tview.Pages
	store *store.Store
}

func NewApp(s *store.Store) *App {
	app := tview.NewApplication()
	pages := tview.NewPages()

	return &App{
		app:   app,
		pages: pages,
		store: s,
	}
}

func (a *App) Run() error {
	if err := a.setupMainView(); err != nil {
		return err
	}
	return a.app.Run()
}

func (a *App) setupMainView() error {
	cfg, err := a.store.Load()
	if err != nil {
		cfg = &config.StoresConfig{}
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] setupMainView: %d servers\n", len(cfg.Servers))

	list := a.createServerList(cfg)

	// Wrap list in a Flex with black background for visibility
	flex := tview.NewFlex().
		AddItem(list, 0, 1, true)
	flex.SetBackgroundColor(tcell.ColorBlack)

	a.pages.AddPage("main", flex, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(list)
	a.app.Draw()

	return nil
}

func (a *App) refreshMainView() {
	// Refresh the main view by re-calling setupMainView
	a.setupMainView()
}
