package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"yaboy/gameboy"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestNewTuiPopulatesRAMView(t *testing.T) {
	gb := gameboy.NewGameBoy()
	gb.RAM[0] = 0x42
	gb.RAM[len(gb.RAM)-1] = 0xab

	ui := NewTui(&gb)
	if ui.mainPage.ramView == nil {
		t.Fatal("NewTui did not initialize the RAM view")
	}

	text := ui.mainPage.ramView.GetText(true)
	if !strings.HasPrefix(text, "[0x0000]\t0x42\n") {
		t.Error("RAM view does not show the first RAM address and its value")
	}
	if !strings.HasSuffix(text, "[0xfffe]\t0xab\n") {
		t.Error("RAM view does not show the last RAM address and its value")
	}
	if got := strings.Count(text, "\n"); got != len(gb.RAM) {
		t.Errorf("RAM view has %d rows, want %d", got, len(gb.RAM))
	}
}

func TestUpdateRAMViewWithInstructions(t *testing.T) {
	gb := gameboy.NewGameBoy()
	copy(gb.RAM[0x100:], []byte{0x01, 0x34, 0x12, 0xcb, 0x11, 0x77, 0x00})
	ui := NewTui(&gb)
	ui.UpdateRAMView(true)

	want := "[0x0100]\t0x01\t\tLD BC,0x1234\n" +
		"[0x0101]\t0x34\n" +
		"[0x0102]\t0x12\n" +
		"[0x0103]\t0xcb\t\tRL C\n" +
		"[0x0104]\t0x11\n" +
		"[0x0105]\t0x77\t\tLD [HL],A\n" +
		"[0x0106]\t0x00\t\tNOP\n"
	text := ui.mainPage.ramView.GetText(true)
	if !strings.Contains(text, want) {
		t.Fatal("RAM view does not align instructions and operand bytes with their addresses")
	}
	if got := strings.Count(text, "\n"); got != len(gb.RAM) {
		t.Errorf("RAM view has %d rows, want %d", got, len(gb.RAM))
	}
}

func TestImportAndLoadInstructionsIntoRAM(t *testing.T) {
	opcodes := []byte{0x01, 0x34, 0x12, 0xcb, 0x11, 0x77}
	path := filepath.Join(t.TempDir(), "instructions.bin")
	if err := os.WriteFile(path, opcodes, 0o600); err != nil {
		t.Fatal(err)
	}

	gb := gameboy.NewGameBoy()
	ui := NewTui(&gb)
	var focus tview.Primitive = ui.mainPage.page
	ui.state.previousFocus = &focus
	setFocus := func(p tview.Primitive) { ui.app.SetFocus(p) }
	ui.importBinaryInstructionsPage.inputBox.SetText(path)
	ui.importBinaryInstructionsPage.inputBox.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), setFocus)
	if !bytes.Equal(ui.importBinaryInstructionsPage.opcodes, opcodes) {
		t.Fatal("imported bytes were not saved to the TUI page")
	}

	ui.inputMemoryAddressPage.inputBox.SetText("0x0100")
	ui.inputMemoryAddressPage.inputBox.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), setFocus)
	if !bytes.Equal(gb.RAM[0x100:0x100+len(opcodes)], opcodes) {
		t.Fatal("imported bytes were not loaded into RAM at the selected address")
	}
	if !strings.Contains(ui.mainPage.ramView.GetText(true), "[0x0100]\t0x01\t\tLD BC,0x1234\n") {
		t.Fatal("loading instructions did not refresh the RAM view")
	}
}
