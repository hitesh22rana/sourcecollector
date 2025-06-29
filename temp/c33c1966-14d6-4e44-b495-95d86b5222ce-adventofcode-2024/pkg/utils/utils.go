package utils

import (
	"bufio"
	"os"
)

func ReadInput(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	sc := bufio.NewScanner(file)
	var data []string
	for sc.Scan() {
		data = append(data, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
