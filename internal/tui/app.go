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
	// Check if master password is needed
	if a.store.IsLocked() {
		if err := a.showMasterPasswordSetup(); err != nil {
			return err
		}
	} else if !a.store.IsUnlocked() {
		if err := a.showMasterPasswordUnlock(); err != nil {
			return err
		}
	}

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

	return nil
}

func (a *App) refreshMainView() {
	// Refresh the main view by re-calling setupMainView
	a.setupMainView()
}

// showMasterPasswordSetup shows the initial password setup dialog
func (a *App) showMasterPasswordSetup() error {
	done := make(chan error, 1)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Set Master Password")

	form.AddInputField("Password", "", 40, nil, nil)
	form.AddInputField("Confirm Password", "", 40, nil, nil)

	form.AddButton("Set Password", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		confirm := form.GetFormItemByLabel("Confirm Password").(*tview.InputField).GetText()

		if password == "" {
			done <- fmt.Errorf("password cannot be empty")
			return
		}
		if password != confirm {
			done <- fmt.Errorf("passwords do not match")
			return
		}
		if len(password) < 4 {
			done <- fmt.Errorf("password must be at least 4 characters")
			return
		}

		if err := a.store.SetMasterPassword(password); err != nil {
			done <- err
			return
		}
		done <- nil
	})

	form.AddButton("Skip (Not Recommended)", func() {
		// Skip encryption - will use no master password
		done <- nil
	})

	a.pages.AddPage("password_setup", form, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	select {
	case err := <-done:
		a.pages.RemovePage("password_setup")
		return err
	}
}

// showMasterPasswordUnlock shows the password unlock dialog
func (a *App) showMasterPasswordUnlock() error {
	done := make(chan error, 1)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Enter Master Password")

	form.AddInputField("Password", "", 40, nil, nil)

	form.AddButton("Unlock", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		if err := a.store.Unlock(password); err != nil {
			done <- err
			return
		}
		done <- nil
	})

	form.AddButton("Cancel", func() {
		done <- fmt.Errorf("cancelled")
	})

	a.pages.AddPage("password_unlock", form, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	select {
	case err := <-done:
		a.pages.RemovePage("password_unlock")
		return err
	}
}
