package tui

import (
	"github.com/rivo/tview"
)

func createMenuPage(t *Tui) *tview.Flex {
	var menuPage = tview.NewFlex().SetDirection(tview.FlexRow)
	var menuText = `(f) Import binary instructions file
(x) Execute instructions
(l) Load instructions to specific address
(m) Close the menu`

	// Set up of menu page
	menuTextBox := tview.NewTextView()
	menuTextBox.SetBorder(true)
	menuTextBox.SetText(menuText)

	menuRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuRow.AddItem(nil, 0, 1, false).
		AddItem(menuTextBox, 0, 4, false).
		AddItem(nil, 0, 1, false)

	menuPage.AddItem(nil, 0, 1, false).
		AddItem(menuRow, 0, 4, true).
		AddItem(nil, 0, 1, false)

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
