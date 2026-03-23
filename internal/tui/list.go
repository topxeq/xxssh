package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/topxeq/xxssh/internal/config"
	"github.com/topxeq/xxssh/internal/ssh"
	cssh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func (a *App) createServerList(cfg *config.StoresConfig) tview.Primitive {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("xxssh - SSH Client")

	for i, srv := range cfg.Servers {
		idx := i
		list.AddItem(srv.Name, srv.Host, 0, func() {
			a.connectToServer(idx)
		})
	}

	list.AddItem("+ Add new server", "", 0, func() {
		a.showEditForm(-1) // -1 means new server
	})

	// Key handlers using SetInputCapture
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'a':
				a.showEditForm(-1)
				return nil
			case 'e':
				selected := list.GetCurrentItem()
				if selected >= 0 && selected < len(cfg.Servers) {
					a.showEditForm(selected)
				}
				return nil
			case 'd':
				selected := list.GetCurrentItem()
				if selected >= 0 && selected < len(cfg.Servers) {
					a.deleteServer(selected)
				}
				return nil
			case 'q':
				a.app.Stop()
				return nil
			}
		}
		// Let other keys (arrows, enter, etc.) pass through to the list
		return event
	})

	return list
}

func (a *App) deleteServer(idx int) {
	cfg, err := a.store.Load()
	if err != nil {
		return
	}
	if idx < 0 || idx >= len(cfg.Servers) {
		return
	}
	cfg.Servers = append(cfg.Servers[:idx], cfg.Servers[idx+1:]...)
	a.store.Save(cfg)
	a.refreshMainView()
}

func (a *App) connectToServer(idx int) {
	cfg, err := a.store.Load()
	if err != nil {
		return
	}
	if idx >= len(cfg.Servers) {
		return
	}

	srv := &cfg.Servers[idx]

	for {
		client := ssh.NewSSHClient(srv)

		if err := client.Connect(); err != nil {
			a.showError("Connection failed: " + err.Error())
			return
		}

		session, err := client.Session()
		if err != nil {
			a.showError("Session failed: " + err.Error())
			client.Close()
			return
		}

		width, height, _ := term.GetSize(int(os.Stdin.Fd()))

		session.RequestPty("xterm-256color", width, height, cssh.TerminalModes{
			cssh.ECHO: 1,
		})

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		if err := session.Shell(); err != nil {
			a.showError("Shell failed: " + err.Error())
			client.Close()
			return
		}

		err = session.Wait()
		client.Close()

		if err == nil {
			// Clean exit, return to list
			return
		}

		// Connection lost - show reconnect prompt
		if !a.promptReconnect() {
			return
		}
		// Loop and reconnect
	}
}

func (a *App) showError(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			a.pages.RemovePage("error")
		})
	a.pages.AddPage("error", modal, true, true)
}

func (a *App) promptReconnect() bool {
	done := make(chan bool, 1)

	modal := tview.NewModal().
		SetText("Connection lost. Reconnect?").
		AddButtons([]string{"Reconnect", "Back to List"}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			done <- (buttonIndex == 0)
		})
	a.pages.AddPage("reconnect", modal, true, true)
	a.pages.SwitchToPage("reconnect")

	// Wait for result
	select {
	case reconnect := <-done:
		a.pages.RemovePage("reconnect")
		return reconnect
	case <-done: // channel closed
		return false
	}
}
