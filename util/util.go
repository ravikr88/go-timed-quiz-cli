package util

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

func ReadCSV(file io.Reader) (map[int]string, map[int]string, error) {
	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	totalQuestions := len(csvData)
	questions := make(map[int]string, totalQuestions)
	answers := make(map[int]string, totalQuestions)

	for i, data := range csvData {
		questions[i] = data[0]
		answers[i] = data[1]
	}

	return questions, answers, nil
}

func WaitForEnter() {
	bufio.NewReader(os.Stdin).ReadString('\n')
}
