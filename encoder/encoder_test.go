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
	t.Log(data)
	t.Log(err)
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
		t.Fatalf("Error encoding ADC input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read ADC test output")
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

func TestSubstract(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, SUB_INPUT)
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
		t.Fatalf("Error encoding SUB input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read SUB test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, SUB_OUTPUT) {
		t.Errorf("Failed SUB test\n\tExpected: %s\n\tActual: %s", SUB_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(SUB_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range SUB_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestSubstractCarry(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, SBC_INPUT)
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
		t.Fatalf("Error encoding SBC input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read SBC test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, SBC_OUTPUT) {
		t.Errorf("Failed SBC test\n\tExpected: %s\n\tActual: %s", SBC_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(SBC_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range SBC_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestAnd(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, AND_INPUT)
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
		t.Fatalf("Error encoding AND input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read AND test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, AND_OUTPUT) {
		t.Errorf("Failed AND test\n\tExpected: %s\n\tActual: %s", AND_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(AND_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range AND_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestExclusiveOr(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, XOR_INPUT)
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
		t.Fatalf("Error encoding XOR input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read XOR test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, XOR_OUTPUT) {
		t.Errorf("Failed XOR test\n\tExpected: %s\n\tActual: %s", XOR_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(XOR_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range XOR_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestOr(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, OR_INPUT)
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
		t.Fatalf("Error encoding OR input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read OR test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, OR_OUTPUT) {
		t.Errorf("Failed OR test\n\tExpected: %s\n\tActual: %s", OR_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(OR_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range OR_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestCompare(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, CP_INPUT)
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
		t.Fatalf("Error encoding CP input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read CP test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, CP_OUTPUT) {
		t.Errorf("Failed CP test\n\tExpected: %s\n\tActual: %s", CP_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(CP_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range CP_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestIncrement8(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, INC_INPUT)
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
		t.Fatalf("Error encoding INC input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read INC test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, INC_OUTPUT) {
		t.Errorf("Failed INC test\n\tExpected: %s\n\tActual: %s", INC_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(INC_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range INC_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestDecrement8(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, DEC_INPUT)
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
		t.Fatalf("Error encoding INC input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read INC test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, DEC_OUTPUT) {
		t.Errorf("Failed DEC test\n\tExpected: %s\n\tActual: %s", DEC_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(DEC_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range DEC_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestOther8bit(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, OTHER_INPUT)
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
		t.Fatalf("Error encoding other 8bit arith input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read other 8bit arith test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, OTHER_OUTPUT) {
		t.Errorf("Failed other 8bit arith test\n\tExpected: %s\n\tActual: %s", OTHER_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(OTHER_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range OTHER_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestAdd16(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, ADD16_INPUT)
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
		t.Fatalf("Error encoding ADD 16bit arith input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read ADD 16bit arith test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, ADD16_OUTPUT) {
		t.Errorf("Failed ADD 16bit arith test\n\tExpected: %s\n\tActual: %s", ADD16_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(ADD16_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range ADD16_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestIncrement16(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, INC16_INPUT)
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
		t.Fatalf("Error encoding INC 16bit arith input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read INC 16bit arith test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, INC16_OUTPUT) {
		t.Errorf("Failed INC 16bit arith test\n\tExpected: %s\n\tActual: %s", INC16_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(INC16_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range INC16_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestDecrement16(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, DEC16_INPUT)
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
		t.Fatalf("Error encoding DEC 16bit arith input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read DEC 16bit arith test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, DEC16_OUTPUT) {
		t.Errorf("Failed DEC 16bit arith test\n\tExpected: %s\n\tActual: %s", DEC16_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(DEC16_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range DEC16_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestJump(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, JUMP_INPUT)
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
		t.Fatalf("Error encoding Jump input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read Jump test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, JUMP_OUTPUT) {
		t.Errorf("Failed Jump test\n\tExpected: %s\n\tActual: %s", JUMP_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(JUMP_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range JUMP_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TesCallRestartReturn(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, CALLRSTRET_INPUT)
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
		t.Fatalf("Error encoding CALL/RST/RET/RETI input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read CALL/RST/RET/RETI test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, CALLRSTRET_OUTPUT) {
		t.Errorf("Failed CALL/RST/RET/RETI test\n\tExpected: %s\n\tActual: %s", CALLRSTRET_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(CALLRSTRET_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range CALLRSTRET_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}

func TestLoad16(t *testing.T) {
	tempDir := t.TempDir()

	// inp := "LD A,0x00"
	inputFileName, err := createInputFile(tempDir, LOAD16_INPUT)
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
		t.Fatalf("Error encoding load 16bit input: %s", err)
	}

	// decode
	data, err := enc.ReadBinaryFile(outputFileName)
	if err != nil {
		t.Fatal("Failed to read load 16bit test output")
	}

	if len(data) == 0 {
		t.Error("Read output is empty!")
	}

	// out := []uint8{0x3E, 0x00}
	// if !slices.Equal(data, ARITH8BIT_OUTPUT) {
	// 	t.Errorf("Failed LD test\n\tExpected: %s\n\tActual: %s", ARITH8BIT_OUTPUT, data)
	// }

	if !slices.Equal(data, LOAD16_OUTPUT) {
		t.Errorf("Failed load 16bit test\n\tExpected: %s\n\tActual: %s", LOAD16_OUTPUT, data)
		t.Logf("Length of expctd: %d", len(LOAD16_OUTPUT))
		t.Logf("Length of actual: %d", len(data))

		expelems := []string{}
		actelems := []string{}
		for _, bb := range LOAD16_OUTPUT {
			expelems = append(expelems, fmt.Sprintf("%#02x", bb))
		}
		for _, bb := range data {
			actelems = append(actelems, fmt.Sprintf("%#02x", bb))
		}

		t.Logf("Hex of expctd: %s", expelems)
		t.Logf("Hex of actual: %s", actelems)
	}
}
