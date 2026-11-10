package tui

import (
	"fmt"
	"strings"
	"yaboy/gameboy"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const MAIN_PAGE = "main-page"
const MENU_PAGE = "menu-page"
const IMPORT_BINARY_INSTRUCTIONS_PAGE = "import-bin-instructions-page"

type State struct {
	menuOpen             bool
	importBinaryPageOpen bool
	previousFocus        *tview.Primitive
}

type Tui struct {
	app                          *tview.Application
	pages                        *tview.Pages
	gameboy                      *gameboy.GameBoy
	state                        State
	mainPage                     MainPage
	importBinaryInstructionsPage ImportBinaryInstructionsPage
}

func (t *Tui) UpdateRAMView() {
	// Clear current text
	t.mainPage.ramView.Clear()

	// Update RAM text to current state
	var text strings.Builder
	for i, v := range t.gameboy.RAM {
		fmt.Fprintf(&text, "[\"%d\"][%#04x]\t%#02x\n", i, i, v)
	}
	t.mainPage.ramView.SetText(text.String())
}

// newTui returns a new tview application for the terminal UI
func NewTui(gb *gameboy.GameBoy) *Tui {
	// Create base TUI vars
	var app = tview.NewApplication()
	var pages = tview.NewPages()

	tui := &Tui{
		app:     app,
		pages:   pages,
		gameboy: gb,
		state:   State{},
	}

	mainPage := createMainPage(tui)
	tui.mainPage = mainPage

	menuPage := createMenuPage(tui)

	importBinaryInstructionsPage := createImportBinaryPage(tui)
	tui.importBinaryInstructionsPage = importBinaryInstructionsPage

	pages.AddPage(MAIN_PAGE, mainPage.page, true, true)
	pages.AddPage(MENU_PAGE, menuPage, true, false)
	pages.AddPage(IMPORT_BINARY_INSTRUCTIONS_PAGE, importBinaryInstructionsPage.page, true, false)

	tui.UpdateRAMView()

	menuPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if (event.Key() == tcell.KeyESC || event.Rune() == 'm') && tui.state.menuOpen {
			tui.state.menuOpen = false
			pages.HidePage(MENU_PAGE)
			app.SetFocus(*tui.state.previousFocus)
		}
		if !tui.state.menuOpen {
			return event
		}

		switch r := event.Rune(); r {
		case 'f':
			tui.importBinaryInstructionsPage.inputBox.SetText("")
			pages.ShowPage(IMPORT_BINARY_INSTRUCTIONS_PAGE).SendToFront(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			pages.HidePage(MENU_PAGE)
			tui.state.menuOpen = false
			app.SetFocus(tui.importBinaryInstructionsPage.inputBox)

		// case 'x':
		// 	if readingBox.loadedInstructions {
		// 		state.executingInstructions = true
		// 		state.menuOpen = false
		// 		pages.HidePage(MENU_PAGE)
		// 		app.SetFocus(*state.previousFocus)
		// 	}
		case 'm':
			tui.state.menuOpen = false
			pages.HidePage(MENU_PAGE)
			app.SetFocus(*tui.state.previousFocus)
		}
		return event
	})

	mainPage.page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if tui.state.menuOpen {
			return event
		}

		// if state.executingInstructions {
		// 	switch r := event.Rune(); r {
		// 	case 'n':
		// 		readingBox.positionNumber = readingBox.positionNumber + 1
		// 		if readingBox.positionNumber == len(readingBox.instructions)+1 {
		// 			state.executingInstructions = false
		// 		}
		// 		readingBox.DisplayInstructions(state)
		// 		return event
		// 	}
		// }

		switch r := event.Rune(); r {
		case 'q':
			// if !instrField.HasFocus() {
			app.Stop()
			// }
		case 'm':
			focus := app.GetFocus()
			tui.state.previousFocus = &focus
			pages.ShowPage(MENU_PAGE).SendToFront(MENU_PAGE)
			tui.state.menuOpen = true
			app.SetFocus(menuPage)
			return nil
			// app.SetFocus(menuTextBox)

			// app.SetFocus(menu_list)

			// case 's':
			// 	ramView.ScrollToHighlight()
			// case 'r':
			// 	app.SetFocus(ramView)
			// case 'i':
			// 	ramView.Highlight("123")
		}

		return event
	})

	return tui
}

func (t *Tui) Run() {
	if err := t.app.SetRoot(t.pages, true).
		// SetFocus(t.instrField).
		EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
