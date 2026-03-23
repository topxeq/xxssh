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
		// setupMainView is called inside the form button handlers
		return a.app.Run()
	} else if !a.store.IsUnlocked() {
		if err := a.showMasterPasswordUnlock(); err != nil {
			return err
		}
		// setupMainView is called inside the form button handlers
		return a.app.Run()
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
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Set Master Password")

	form.AddInputField("Password", "", 40, nil, nil)
	form.AddInputField("Confirm Password", "", 40, nil, nil)

	form.AddButton("Set Password", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		confirm := form.GetFormItemByLabel("Confirm Password").(*tview.InputField).GetText()

		if password == "" {
			showErrorModal(a.pages, "Password cannot be empty")
			return
		}
		if password != confirm {
			showErrorModal(a.pages, "Passwords do not match")
			return
		}
		if len(password) < 4 {
			showErrorModal(a.pages, "Password must be at least 4 characters")
			return
		}

		if err := a.store.SetMasterPassword(password); err != nil {
			showErrorModal(a.pages, err.Error())
			return
		}

		// Success - remove all pages and continue to main view
		a.pages.RemovePage("password_setup")
		a.pages.RemovePage("error")
		a.setupMainView()
	})

	form.AddButton("Skip (Not Recommended)", func() {
		// Skip encryption - will use no master password
		a.pages.RemovePage("password_setup")
		a.setupMainView()
	})

	a.pages.AddPage("password_setup", form, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	return nil
}

func showErrorModal(pages *tview.Pages, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			pages.RemovePage("error")
		})
	pages.AddPage("error", modal, true, true)
}

// showMasterPasswordUnlock shows the password unlock dialog
func (a *App) showMasterPasswordUnlock() error {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Enter Master Password")

	form.AddInputField("Password", "", 40, nil, nil)

	form.AddButton("Unlock", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		if err := a.store.Unlock(password); err != nil {
			showErrorModal(a.pages, "Invalid password: "+err.Error())
			return
		}

		// Success - remove all pages and continue to main view
		a.pages.RemovePage("password_unlock")
		a.pages.RemovePage("error")
		a.setupMainView()
	})

	form.AddButton("Cancel", func() {
		a.app.Stop()
	})

	a.pages.AddPage("password_unlock", form, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	return nil
}
