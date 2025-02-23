package encoder_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"
	enc "yaboy/encoder"
)

func createInputFile(tempDir string, content string) (_ string, err error) {
	f, err := os.CreateTemp(tempDir, "tempFile")
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("Error writing test string for input: %w", err)
	}

	writer.Flush()
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

func TestLoad(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, LOAD_INPUT)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	outputFileName, err := createFile(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// testing
	enc.EncodeTextFile(inputFileName, outputFileName)

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read LOAD test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, LOAD_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", LOAD_OUTPUT, data)
	// }

	if !slices.Equal(data, LOAD_OUTPUT) {
		t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", LOAD_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(LOAD_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range LOAD_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestAdd(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, ADD_INPUT)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	outputFileName, err := createFile(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// testing
	err = enc.EncodeTextFile(inputFileName, outputFileName)
	if err != nil {
		t.Fatalf("Error encoding ADD input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read ADD arithmetic test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, ADD_OUTPUT) {
		t.Errorf("Failed ADD test\n\tExpected: %s\n\tActual: %s", ADD_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(ADD_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range ADD_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestAddCarry(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, ADC_INPUT)
	// inputFileName, err := createInputFile(tempDir, inp)
	if err != nil {
		t.Fatal(err)
	}
	outputFileName, err := createFile(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// testing
	err = enc.EncodeTextFile(inputFileName, outputFileName)
	if err != nil {
		t.Fatalf("Error encoding ADD input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read ADD test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, ADC_OUTPUT) {
		t.Errorf("Failed ADC test\n\tExpected: %s\n\tActual: %s", ADC_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(ADC_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range ADC_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}
