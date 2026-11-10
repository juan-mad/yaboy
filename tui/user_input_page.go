package tui

import (
	"fmt"
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

func createImportBinaryPage(t *Tui) ImportBinaryInstructionsPage {
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
			t.pages.HidePage(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			t.app.SetFocus(*t.state.previousFocus)

			opcodes, instructions, lengths := loadBinaryInstructions(value)
			inputPage.opcodes = opcodes
			inputPage.instructions = instructions
			inputPage.lengths = lengths

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

			t.mainPage.readingBox.SetText(text.String())

		} else if key == tcell.KeyESC {
			t.pages.HidePage(IMPORT_BINARY_INSTRUCTIONS_PAGE)
			t.app.SetFocus(*t.state.previousFocus)
		} else {
			return
		}
		if key != tcell.KeyEnter && key != tcell.KeyEsc {
			return
		}
	})

	return inputPage
}
