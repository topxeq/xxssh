package tui

import (
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

	list := a.createServerList(cfg)

	a.pages.AddPage("main", list, true, true)
	a.app.SetRoot(a.pages, true)

	return nil
}

func (a *App) refreshMainView() {
	// Refresh the main view by re-calling setupMainView
	a.setupMainView()
}
