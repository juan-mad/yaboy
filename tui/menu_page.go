package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MenuPage struct {
	page *tview.Flex
}

func createMenuPage(tui *Tui) MenuPage {
	var page = tview.NewFlex().SetDirection(tview.FlexRow)
	var menuText = `(f) Import binary instructions file
(x) Execute instructions
(l) Load instructions to specific address
(m) Close the menu`

	var menuPage = MenuPage{page: page}

	// Set up of menu page
	menuTextBox := tview.NewTextView()
	menuTextBox.SetBorder(true)
	menuTextBox.SetText(menuText)

	menuRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuRow.AddItem(nil, 0, 1, false).
		AddItem(menuTextBox, 0, 4, false).
		AddItem(nil, 0, 1, false)

	page.AddItem(nil, 0, 1, false).
		AddItem(menuRow, 0, 4, true).
		AddItem(nil, 0, 1, false)

	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if (event.Key() == tcell.KeyESC || event.Rune() == 'm') && tui.state.menuOpen {
			tui.state.menuOpen = false
			tui.pages.HidePage(MENU_PAGE)
			tui.app.SetFocus(*tui.state.previousFocus)
		}
		if !tui.state.menuOpen {
			return event
		}

		switch r := event.Rune(); r {
		case 'f':
			tui.importBinaryInstructionsPage.inputBox.SetText("")
			tui.pages.ShowPage(IMPORT_BINARY_INSTRUCTIONS_PAGE).SendToFront(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			tui.pages.HidePage(MENU_PAGE)
			tui.state.menuOpen = false
			tui.app.SetFocus(tui.importBinaryInstructionsPage.inputBox)
		case 'm':
			tui.state.menuOpen = false
			tui.pages.HidePage(MENU_PAGE)
			tui.app.SetFocus(*tui.state.previousFocus)
		case 'l':
			if tui.state.importedInstructions {
				tui.inputMemoryAddressPage.inputBox.SetText("")
				tui.pages.ShowPage(INPUT_MEMORY_ADDRESS_PAGE).SendToFront(INPUT_MEMORY_ADDRESS_PAGE)
				tui.pages.HidePage(MENU_PAGE)
				tui.state.menuOpen = false
				tui.app.SetFocus(tui.inputMemoryAddressPage.inputBox)
			}
		}
		return event
	})

	// 	// Menu page user interaction
	// 	menuPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	    // Closing the menu without doing anything
	// 		if (event.Key() == tcell.KeyESC || event.Rune() == 'm') && state.menuOpen {
	// 			state.menuOpen = false
	// 			pages.HidePage(MENU_PAGE)
	// 			app.SetFocus(*state.previousFocus)
	// 			return event
	// 		}
	//
	// 		// If menu is closed or we are executing instructions, nothing else to do
	// 		if !state.menuOpen || state.executingInstructions {
	// 		    return event
	// 		}
	//
	// 		switch r := event.Rune(): r {
	// 		case 'f':
	// 		}
	// 	})

	return menuPage
}
