package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainPage struct {
	page        *tview.Flex
	readingBox  *tview.TextView
	infoBox     *tview.TextView
	registerBox *tview.Flex
	tab1        *tview.Table
	tab2        *tview.Table
	ramView     *tview.TextView
}

func createMainPage(tui *Tui) MainPage {
	var page = tview.NewFlex().SetDirection(tview.FlexColumn)
	var informationText = `Press (m) to open the action menu.`

	var infoCol = tview.NewFlex().SetDirection(tview.FlexRow)
	var readingBox = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	var infoBox = tview.NewTextView().SetText(informationText)
	readingBox.SetBorder(true)
	infoBox.SetBorder(true)
	infoCol.AddItem(readingBox, 0, 7, false)
	infoCol.AddItem(infoBox, 0, 1, false)

	var memCol = tview.NewFlex().SetDirection(tview.FlexRow)

	// Registers
	registerBox := tview.NewFlex().SetDirection(tview.FlexRow)
	registerBox.SetBorder(true).SetTitle("Registers")
	tab1 := tview.NewTable()
	tab2 := tview.NewTable()

	registerBox.AddItem(tab1, 0, 2, false).AddItem(tab2, 0, 1, false)

	ramView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			tui.app.Draw()
		})
	ramView.SetBorder(true)

	memCol.AddItem(registerBox, 0, 1, false).AddItem(ramView, 0, 2, false)

	page.AddItem(infoCol, 0, 1, false)
	page.AddItem(memCol, 0, 1, false)

	mainPage := MainPage{
		page:        page,
		readingBox:  readingBox,
		infoBox:     infoBox,
		registerBox: registerBox,
		tab1:        tab1,
		tab2:        tab2,
		ramView:     ramView,
	}

	tui.pages.AddPage(MAIN_PAGE, page, true, true)

	mainPage.page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if tui.state.menuOpen {
			return event
		}

		switch r := event.Rune(); r {
		case 'q':
			// if !instrField.HasFocus() {
			tui.app.Stop()
			// }
		case 'm':
			focus := tui.app.GetFocus()
			tui.state.previousFocus = &focus
			tui.pages.ShowPage(MENU_PAGE).SendToFront(MENU_PAGE)
			tui.state.menuOpen = true
			tui.app.SetFocus(tui.menuPage.page)
			return nil
			// case 's':
			// 	ramView.ScrollToHighlight()
			// case 'r':
			// 	app.SetFocus(ramView)
			// case 'i':
			// 	ramView.Highlight("123")
		}

		return event
	})
	return mainPage
}