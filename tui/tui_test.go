package tui

import (
	"strings"
	"testing"
	"yaboy/gameboy"
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
