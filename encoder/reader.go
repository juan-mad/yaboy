package encoder

import "os"
import "io"
import "fmt"
import "errors"
import "encoding/binary"

func ReadBinaryFile(inputPath string) ([]uint8, error) {
	data := []uint8{}
	file, err := os.OpenFile(inputPath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return data, fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close() // Ensure file is closed when we're done

	for {
		var instruction uint8

		err = binary.Read(file, binary.LittleEndian, &instruction)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return data, fmt.Errorf("error reading binary data: %w", err)
		}

		data = append(data, instruction)
	}

	return data, nil
}
