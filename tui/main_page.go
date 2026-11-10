package tui

import (
	"fmt"
	"yaboy/gameboy"

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

// initRegisterTables creates the tview Tables that will show the CPU registers
func initRegisterTables(reg gameboy.CPURegisters) (*tview.Table, *tview.Table) {
	tab1 := tview.NewTable()
	tab1.SetCell(0, 0, tview.NewTableCell("A").SetAlign(tview.AlignCenter))
	tab1.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadA())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(0, 2, tview.NewTableCell("F").SetAlign(tview.AlignCenter))
	tab1.SetCell(0, 3, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadF())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(1, 0, tview.NewTableCell("B").SetAlign(tview.AlignCenter))
	tab1.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadB())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(1, 2, tview.NewTableCell("C").SetAlign(tview.AlignCenter))
	tab1.SetCell(1, 3, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadC())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(2, 0, tview.NewTableCell("D").SetAlign(tview.AlignCenter))
	tab1.SetCell(2, 1, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadD())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(2, 2, tview.NewTableCell("E").SetAlign(tview.AlignCenter))
	tab1.SetCell(2, 3, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadE())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(3, 0, tview.NewTableCell("H").SetAlign(tview.AlignCenter))
	tab1.SetCell(3, 1, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadH())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab1.SetCell(3, 2, tview.NewTableCell("L").SetAlign(tview.AlignCenter))
	tab1.SetCell(3, 3, tview.NewTableCell(fmt.Sprintf("%08b", reg.ReadL())).SetAlign(tview.AlignCenter).SetExpansion(1))

	tab2 := tview.NewTable()
	tab2.SetCell(0, 0, tview.NewTableCell("SP").SetAlign(tview.AlignCenter))
	tab2.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("%016b", reg.ReadSP())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tab2.SetCell(1, 0, tview.NewTableCell("PC").SetAlign(tview.AlignCenter))
	tab2.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%016b", reg.ReadPC())).SetAlign(tview.AlignCenter).SetExpansion(1))

	return tab1, tab2
}

func createMainPage(t *Tui) MainPage {
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
	tab1, tab2 := initRegisterTables(t.gameboy.Reg)
	registerBox.AddItem(tab1, 0, 2, false).AddItem(tab2, 0, 1, false)

	ramView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			t.app.Draw()
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

	t.pages.AddPage(MAIN_PAGE, page, true, true)

	return mainPage
}

// func (p MainPage) GetPrimitiveByName(name PrimitiveName) tview.Primitive {
// 	switch name {
// 	case InfoBox:
// 		return p.GetItem(0)
// 	}
// }

// type ReadingBox struct {
// 	app                *tview.Application
// 	text               *tview.TextView
// 	instructions       []string
// 	numLines           int
// 	positionNumber     int
// 	loadedInstructions bool
// }
//
