package tui

import (
	"fmt"
	"strconv"
	"strings"
	"yaboy/decompiler"
	enc "yaboy/encoder"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ImportBinaryInstructionsPage struct {
	page         *tview.Flex
	inputBox     *tview.InputField
	opcodes      []uint8
	instructions []string
	lengths      []uint8
}

type InputMemoryAddressPage struct {
	page          *tview.Flex
	inputBox      *tview.InputField
	memoryAddress uint16
}

// loadBinaryInstructions loads a binary file containing instructions. Returns a slice of string.
func loadBinaryInstructions(path string) ([]uint8, []string, []uint8) {
	opcodes, err := enc.ReadBinaryFile(path)
	if err != nil {
		panic(err)
	}
	instructions, lengths, err := decompiler.DecodeBytes(opcodes)
	if err != nil {
		panic(err)
	}
	return opcodes, instructions, lengths

}

func createImportBinaryPage(tui *Tui) ImportBinaryInstructionsPage {
	inputBox := tview.NewInputField().
		SetLabel("Enter file name: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDarkGrey).
		SetFieldTextColor(tcell.ColorBlack)
	inputBox.SetBorder(true).SetTitle("Import Binary Instructions from File")

	inputBoxRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputBoxRow.AddItem(nil, 0, 1, false).
		AddItem(inputBox, 0, 1, true).
		AddItem(nil, 0, 1, false)

	page := tview.NewFlex().SetDirection(tview.FlexRow)
	page.AddItem(nil, 0, 1, false).
		AddItem(inputBoxRow, 0, 1, true).
		AddItem(nil, 0, 1, false)

	inputPage := ImportBinaryInstructionsPage{
		page:     page,
		inputBox: inputBox,
	}

	inputBox.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			value := inputBox.GetText()
			tui.pages.HidePage(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			tui.app.SetFocus(*tui.state.previousFocus)

			opcodes, instructions, lengths := loadBinaryInstructions(value)
			tui.importBinaryInstructionsPage.opcodes = opcodes
			tui.importBinaryInstructionsPage.instructions = instructions
			tui.importBinaryInstructionsPage.lengths = lengths

			var text strings.Builder

			var j int = 0
			for i := range len(instructions) {
				length := lengths[i] - 1

				fmt.Fprintf(&text, "[\"%d\"]%#02x\t\t%s", i, opcodes[j], tview.Escape(instructions[i]))
				for length > 0 {
					j = j + 1
					fmt.Fprintf(&text, "\n%#02x", opcodes[j])
					length = length - 1
				}
				fmt.Fprintf(&text, "[\"\"]\n")
				j = j + 1

				// fmt.Fprintf(&text, "[\"%d\"][%#04x]\t%#02x\t(%s)[\"\"]\n", i, i, opcodes[i], instructions[i])
			}

			tui.mainPage.readingBox.SetText(text.String())
			tui.state.importedInstructions = true

		} else if key == tcell.KeyESC {
			tui.pages.HidePage(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			tui.app.SetFocus(*tui.state.previousFocus)
		} else {
			return
		}
		if key != tcell.KeyEnter && key != tcell.KeyEsc {
			return
		}
	})

	return inputPage
}

func createInputMemoryAddressPage(tui *Tui) InputMemoryAddressPage {
	inputBox := tview.NewInputField().
		SetLabel("Enter address in 0xXXXX format: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDarkGrey).
		SetFieldTextColor(tcell.ColorBlack)
	inputBox.SetBorder(true).SetTitle("Load data starting at memory address")

	inputBoxRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputBoxRow.AddItem(nil, 0, 1, false).AddItem(inputBox, 0, 1, true).AddItem(nil, 0, 1, false)

	page := tview.NewFlex().SetDirection(tview.FlexRow)
	page.AddItem(nil, 0, 1, false).AddItem(inputBoxRow, 0, 1, true).AddItem(nil, 0, 1, false)

	memaddrPage := InputMemoryAddressPage{
		page:     page,
		inputBox: inputBox,
	}

	inputBox.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			value := inputBox.GetText()
			bigaddress, err := strconv.ParseUint(value, 0, 16)
			if err != nil {
				panic(err)
				// inputBox.SetPlaceholder("could not parsing")
				// return
			}
			address := uint16(bigaddress)

			if address+uint16(len(tui.importBinaryInstructionsPage.opcodes)) > uint16(len(tui.gameboy.RAM)) {
				inputBox.SetPlaceholder("no space!")
				return

			}

			inputBox.SetPlaceholder("")
			tui.pages.HidePage(INPUT_MEMORY_ADDRESS_PAGE)
			tui.app.SetFocus(*tui.state.previousFocus)

			tui.inputMemoryAddressPage.memoryAddress = address
			for i, opcode := range tui.importBinaryInstructionsPage.opcodes {
				tui.gameboy.RAM[address+uint16(i)] = opcode
			}

			tui.UpdateRAMView(true)

		}
	})

	return memaddrPage
}
