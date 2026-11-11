package tui

import (
	"fmt"
	"strings"
	"yaboy/decompiler"
	"yaboy/gameboy"

	"github.com/rivo/tview"
)

const MAIN_PAGE = "main-page"
const MENU_PAGE = "menu-page"
const IMPORT_BINARY_INSTRUCTIONS_PAGE = "import-bin-instructions-page"
const INPUT_MEMORY_ADDRESS_PAGE = "memory-address-page"
const SET_STACK_POINTER_MEMORY_ADDRESS_PAGE = "set-stack-pointer-at-memory-address-page"

type State struct {
	menuOpen             bool
	importBinaryPageOpen bool
	importedInstructions bool
	previousFocus        *tview.Primitive
}

type Tui struct {
	app                          *tview.Application
	pages                        *tview.Pages
	gameboy                      *gameboy.GameBoy
	state                        State
	mainPage                     MainPage
	menuPage                     MenuPage
	importBinaryInstructionsPage ImportBinaryInstructionsPage
	inputMemoryAddressPage       InputMemoryAddressPage
	setStackPointerPage          SetStackPointerAtMemoryAddressPage
}

func (tui *Tui) UpdateRAMView(includeInstructions bool) {
	// Clear current text
	tui.mainPage.ramView.Clear()

	// Update RAM text to current state
	var text strings.Builder
	if !includeInstructions {
		for i, v := range tui.gameboy.RAM {
			fmt.Fprintf(&text, "[\"%d\"][%#04x]\t%#02x\n", i, i, v)
		}
	} else {
		opcodes := tui.gameboy.RAM[:]
		instructions, lengths, err := decompiler.DecodeBytes(opcodes)
		if err != nil {
			panic(err)
		}
		address := 0
		for i, instruction := range instructions {
			fmt.Fprintf(&text, "[\"%d\"][%#04x]\t%#02x\t\t%s", address, address, opcodes[address], tview.Escape(instruction))
			for offset := 1; offset < int(lengths[i]); offset++ {
				operandAddress := address + offset
				fmt.Fprintf(&text, "\n[%#04x]\t%#02x", operandAddress, opcodes[operandAddress])
			}
			fmt.Fprintf(&text, "[\"\"]\n")
			address += int(lengths[i])
		}
	}
	tui.mainPage.ramView.SetText(text.String())
}
func (tui *Tui) UpdateRegisterTables() {
	tui.mainPage.tab1.SetCell(0, 0, tview.NewTableCell("A").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadA())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(0, 2, tview.NewTableCell("F").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(0, 3, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadF())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(1, 0, tview.NewTableCell("B").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadB())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(1, 2, tview.NewTableCell("C").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(1, 3, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadC())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(2, 0, tview.NewTableCell("D").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(2, 1, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadD())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(2, 2, tview.NewTableCell("E").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(2, 3, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadE())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(3, 0, tview.NewTableCell("H").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(3, 1, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadH())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab1.SetCell(3, 2, tview.NewTableCell("L").SetAlign(tview.AlignCenter))
	tui.mainPage.tab1.SetCell(3, 3, tview.NewTableCell(fmt.Sprintf("%08b", tui.gameboy.Reg.ReadL())).SetAlign(tview.AlignCenter).SetExpansion(1))

	tui.mainPage.tab2.SetCell(0, 0, tview.NewTableCell("SP").SetAlign(tview.AlignCenter))
	tui.mainPage.tab2.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("%016b", tui.gameboy.Reg.ReadSP())).SetAlign(tview.AlignCenter).SetExpansion(1))
	tui.mainPage.tab2.SetCell(1, 0, tview.NewTableCell("PC").SetAlign(tview.AlignCenter))
	tui.mainPage.tab2.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%016b", tui.gameboy.Reg.ReadPC())).SetAlign(tview.AlignCenter).SetExpansion(1))
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
	tui.menuPage = menuPage

	importBinaryInstructionsPage := createImportBinaryPage(tui)
	tui.importBinaryInstructionsPage = importBinaryInstructionsPage

	inputMemoryAddressPage := createInputMemoryAddressPage(tui)
	tui.inputMemoryAddressPage = inputMemoryAddressPage

	setStackPointerPage := createSetStackPointerAtMemoryAddress(tui)
	tui.setStackPointerPage = setStackPointerPage

	pages.AddPage(MAIN_PAGE, mainPage.page, true, true)
	pages.AddPage(MENU_PAGE, menuPage.page, true, false)
	pages.AddPage(IMPORT_BINARY_INSTRUCTIONS_PAGE, importBinaryInstructionsPage.page, true, false)
	pages.AddPage(INPUT_MEMORY_ADDRESS_PAGE, inputMemoryAddressPage.page, true, false)
	pages.AddPage(SET_STACK_POINTER_MEMORY_ADDRESS_PAGE, setStackPointerPage.page, true, false)

	tui.UpdateRAMView(false)
	tui.UpdateRegisterTables()

	return tui
}

func (t *Tui) Run() {
	if err := t.app.SetRoot(t.pages, true).
		// SetFocus(t.instrField).
		EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
