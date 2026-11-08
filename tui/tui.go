package tui

import (
	"fmt"
	"strings"
	"yaboy/decompiler"
	enc "yaboy/encoder"
	"yaboy/gameboy"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const MAIN_PAGE = "main-page"
const MENU_PAGE = "menu-page"
const IMPORT_BINARY_INSTRUCTIONS_PAGE = "import-bin-instructions-page"

type Tui struct {
	app      *tview.Application
	pages    *tview.Pages
	flex     *tview.Flex
	tab8bit  *tview.Table
	tab16bit *tview.Table
	ram      *tview.TextView
	gameboy  *gameboy.GameBoy
}

func (t *Tui) UpdateRAMView() {
	t.ram.Clear()
	for i, v := range t.gameboy.RAM {
		fmt.Fprintf(t.ram, "[\"%d\"][%#04x]\t%#02x\n", i, i, v)
	}
}

func loadBinaryInstructions(path string) string {
	inp_data, err := enc.ReadBinaryFile(path)
	if err != nil {
		panic(err)
	}
	instructions, err := decompiler.DecodeBytes(inp_data)
	if err != nil {
		panic(err)
	}
	return strings.Join(instructions, "\n")

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

// newTui returns a new tview application for the terminal UI
func NewTui(gb *gameboy.GameBoy) *Tui {
	// Create base TUI vars
	var app = tview.NewApplication()
	var pages = tview.NewPages()
	var mainPage = tview.NewFlex()
	var menuPage = tview.NewGrid().
		SetRows(0, 3, 0).    // 3 rows with flexible top and bottom margins
		SetColumns(0, 44, 0) // 44 columsn with flexible left and right margins
	var importBinaryInstructionsPage = tview.NewGrid().
		SetRows(0, 3, 0).
		SetColumns(0, 44, 0)

	var previousFocus tview.Primitive
	menuOpen := false

	pages.AddPage(MAIN_PAGE, mainPage, true, true)
	pages.AddPage(MENU_PAGE, menuPage, true, false)
	pages.AddPage(IMPORT_BINARY_INSTRUCTIONS_PAGE, importBinaryInstructionsPage, true, false)

	// MAIN PAGE
	// infoCol appears on the left. It has two text boxes: readingBox and infoBox.
	var infoCol = tview.NewFlex().
		SetDirection(tview.FlexRow)

	// readingBox will display different information
	var readingBox = tview.NewTextView()
	readingBox.SetBorder(true)

	// infoBox will show usage instructions
	var infoBox = tview.NewTextView().SetText("Type m to open action menu")
	infoBox.SetBorder(true)

	infoCol.AddItem(readingBox, 0, 9, false)
	infoCol.AddItem(infoBox, 0, 1, false)

	// memCol appars on the right. It will contain a registers and a RAM box.
	var memCol = tview.NewFlex().SetDirection(tview.FlexRow)

	// regs will show the values of different registers
	regs := tview.NewFlex()
	regs.SetBorder(true).SetTitle("Registers")
	tab1, tab2 := initRegisterTables(gb.Reg)
	regs.SetDirection(tview.FlexRow).AddItem(tab1, 0, 2, false).AddItem(tab2, 0, 1, false)

	// ramView will show the current RAM state.
	ramView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	ramView.SetBorder(true)

	// var ram = tview.NewBox().SetBorder(true).SetTitle("RAM")
	memCol.AddItem(regs, 0, 1, false)
	memCol.AddItem(ramView, 0, 2, false)

	mainPage.SetDirection(tview.FlexColumn)
	mainPage.AddItem(infoCol, 0, 1, false)
	mainPage.AddItem(memCol, 0, 1, false)

	// Action: Import Binary File
	importBinaryInstructionsInput := tview.NewInputField().
		SetLabel("Value: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDarkGray).
		SetFieldTextColor(tcell.ColorBlack)
	importBinaryInstructionsInput.SetBorder(true).SetTitle(" Enter a value ")

	importBinaryInstructionsPage.AddItem(importBinaryInstructionsInput, 1, 1, 1, 1, 0, 0, true)

	menuTextBox := tview.NewTextView()
	menuTextBox.SetBorder(true)
	menuTextBox.SetText("(f) Import binary instructions file")
	menuPage = menuPage.AddItem(menuTextBox, 1, 1, 1, 1, 0, 0, true)

	importBinaryInstructionsInput.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter && key != tcell.KeyEsc {
			return
		}

		value := importBinaryInstructionsInput.GetText()
		pages.HidePage(IMPORT_BINARY_INSTRUCTIONS_PAGE)
		app.SetFocus(previousFocus)

		if key == tcell.KeyEnter {
			readingBox.SetText(loadBinaryInstructions(value))
		}
	})

	tui := &Tui{
		app:      app,
		pages:    pages,
		flex:     mainPage,
		tab8bit:  tab1,
		tab16bit: tab2,
		ram:      ramView,
		gameboy:  gb,
	}

	tui.UpdateRAMView()

	menuPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !menuOpen {
			return event
		}

		if event.Key() == tcell.KeyESC {
			menuOpen = false
			pages.HidePage(MENU_PAGE)
			app.SetFocus(previousFocus)
		}
		switch r := event.Rune(); r {
		case 'f':
			importBinaryInstructionsInput.SetText("")
			pages.ShowPage(IMPORT_BINARY_INSTRUCTIONS_PAGE).SendToFront(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			pages.HidePage(MENU_PAGE)
			menuOpen = false
			app.SetFocus(importBinaryInstructionsInput)
		}
		return event
	})

	mainPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if menuOpen {
			return event
		}
		switch r := event.Rune(); r {
		case 'q':
			// if !instrField.HasFocus() {
			app.Stop()
			// }
		case 'm':
			previousFocus = app.GetFocus()
			pages.ShowPage(MENU_PAGE).SendToFront(MENU_PAGE)
			menuOpen = true
			app.SetFocus(menuTextBox)

			// app.SetFocus(menu_list)

		case 's':
			ramView.ScrollToHighlight()
		case 'r':
			app.SetFocus(ramView)
		case 'i':
			ramView.Highlight("123")
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
