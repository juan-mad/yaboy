package decompiler_test

import "os"
import "fmt"
import "errors"
import "encoding/binary"
import "testing"
import "slices"
import "yaboy/decompiler"
import enc "yaboy/encoder"

func createInputFile(tempDir string, content []uint8) (_ string, err error) {
	f, err := os.CreateTemp(tempDir, "tempFile")
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	// Write the actual data slice
	err = binary.Write(f, binary.LittleEndian, content)
	if err != nil {
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	return f.Name(), nil
}

func createFile(tempDir string) (_ string, err error) {
	f, err := os.CreateTemp(tempDir, "tempFile")
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	return f.Name(), nil
}

func Test0x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, INPUT_0x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_0x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual: %s", OUTPUT_0x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_0x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test1x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,1x00"
	inputFileName, err := createInputFile(tempDir, INPUT_1x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_1x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_1x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_1x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test2x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,2x00"
	inputFileName, err := createInputFile(tempDir, INPUT_2x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_2x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_2x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_2x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test3x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,3x00"
	inputFileName, err := createInputFile(tempDir, INPUT_3x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_3x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_3x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_3x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test4x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,4x00"
	inputFileName, err := createInputFile(tempDir, INPUT_4x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_4x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_4x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_4x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test5x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,5x00"
	inputFileName, err := createInputFile(tempDir, INPUT_5x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_5x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_5x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_5x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test6x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,6x00"
	inputFileName, err := createInputFile(tempDir, INPUT_6x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_6x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_6x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_6x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test7x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,7x00"
	inputFileName, err := createInputFile(tempDir, INPUT_7x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_7x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_7x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_7x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test8x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,8x00"
	inputFileName, err := createInputFile(tempDir, INPUT_8x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_8x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_8x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_8x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func Test9x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,9x00"
	inputFileName, err := createInputFile(tempDir, INPUT_9x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_9x) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_9x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_9x))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestAx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Ax00"
	inputFileName, err := createInputFile(tempDir, INPUT_Ax)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Ax) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_Ax, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Ax))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestBx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Bx00"
	inputFileName, err := createInputFile(tempDir, INPUT_Bx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Bx) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_Bx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Bx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Cx00"
	inputFileName, err := createInputFile(tempDir, INPUT_Cx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Cx) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_Cx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Cx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestDx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Dx00"
	inputFileName, err := createInputFile(tempDir, INPUT_Dx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Dx) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_Dx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Dx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestEx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Ex00"
	inputFileName, err := createInputFile(tempDir, INPUT_Ex)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Ex) {
		t.Errorf("Failed test\n\tExpected: %s\n\tActual:   %s", OUTPUT_Ex, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Ex))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestFx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, INPUT_Fx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, OUTPUT_Fx) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", OUTPUT_Fx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB0x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_0x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_0x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_0x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB1x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_1x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_1x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_1x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB2x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_2x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_2x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_2x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB3x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_3x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_3x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_3x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB4x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_4x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_4x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_4x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB5x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_5x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_5x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_5x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB6x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_6x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_6x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_6x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB7x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_7x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_7x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_7x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB8x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_8x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_8x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_8x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCB9x(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_9x)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_9x) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_9x, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBAx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Ax)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Ax) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Ax, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBBx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Bx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Bx) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Bx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBCx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Cx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Cx) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Cx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBDx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Dx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Dx) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Dx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBEx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Ex)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Ex) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Ex, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}

func TestCBFx(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,Fx00"
	inputFileName, err := createInputFile(tempDir, CB_INPUT_Fx)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	// decode
	inp, err := enc.ReadBinaryFile(inputFileName)
	if err != nil {
		t.Fatal("Failed to read input expected test file")
	}

	// testing
	actual, err := decompiler.DecodeBytes(inp)
	if err != nil {
		t.Errorf("Error decoding bytes: %s", err)
	}

	if !slices.Equal(actual, CB_OUTPUT_Fx) {
		t.Errorf("Failed test\n\tFxpected: %s\n\tActual:   %s", CB_OUTPUT_Fx, actual)
		t.Logf("Length of expctd: %d", len(OUTPUT_Fx))
		t.Logf("Length of actual: %d", len(actual))
		t.Log("Input: ", inp)
	}
}
