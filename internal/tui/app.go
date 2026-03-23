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
		return a.app.Run()
	} else if !a.store.IsUnlocked() {
		if err := a.showMasterPasswordUnlock(); err != nil {
			return err
		}
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

	form.AddPasswordField("Password", "", 30, '*', nil)
	form.AddPasswordField("Confirm", "", 30, '*', nil)

	form.AddButton("Set Password", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		confirm := form.GetFormItemByLabel("Confirm").(*tview.InputField).GetText()

		if password == "" {
			showMessage(a.app, a.pages, "Error", "Password cannot be empty")
			return
		}
		if password != confirm {
			showMessage(a.app, a.pages, "Error", "Passwords do not match")
			return
		}
		if len(password) < 4 {
			showMessage(a.app, a.pages, "Error", "Password must be at least 4 characters")
			return
		}

		if err := a.store.SetMasterPassword(password); err != nil {
			showMessage(a.app, a.pages, "Error", err.Error())
			return
		}

		// Success - switch to main view
		a.pages.SwitchToPage("main")
	})

	form.AddButton("Skip", func() {
		// Skip encryption - continue to main view
		a.pages.SwitchToPage("main")
	})

	centered := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)

	a.pages.AddPage("setup", centered, true, true)
	a.pages.SwitchToPage("setup")
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	return nil
}

// showMasterPasswordUnlock shows the password unlock dialog
func (a *App) showMasterPasswordUnlock() error {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Enter Master Password")

	form.AddPasswordField("Password", "", 30, '*', nil)

	form.AddButton("Unlock", func() {
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		if err := a.store.Unlock(password); err != nil {
			showMessage(a.app, a.pages, "Error", "Invalid password")
			return
		}

		// Success - switch to main view
		a.pages.SwitchToPage("main")
	})

	form.AddButton("Cancel", func() {
		a.app.Stop()
	})

	centered := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)

	a.pages.AddPage("unlock", centered, true, true)
	a.pages.SwitchToPage("unlock")
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)

	return nil
}

// showMessage shows an error/info message using Form as modal
func showMessage(app *tview.Application, pages *tview.Pages, title, message string) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(title)

	label := tview.NewTextView()
	label.SetText(message)
	label.SetTextAlign(tview.AlignCenter)

	form.AddButton("OK", func() {
		pages.RemovePage("message")
	})

	centered := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(label, 0, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)

	pages.AddPage("message", centered, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(form)
}
