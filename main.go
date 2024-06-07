package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Define flags for command-line options
var (
	flagFilePath string
	flagRandom   bool
	flagTime     int
	wg           sync.WaitGroup
)

// Initialize flags
func init() {
	flag.StringVar(&flagFilePath, "file", "questions.csv", "path/to/csv_file")
	flag.BoolVar(&flagRandom, "random", true, "randomize order of questions")
	flag.IntVar(&flagTime, "time", 10, "test duration")
	flag.Parse()
}

func main() {
	// Read command-line flags
	csvPath, err := filepath.Abs(flagFilePath)
	if err != nil {
		log.Fatalln("Unable to parse path: " + csvPath)
	}
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Read CSV file
	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize maps to store questions, answers, and user responses
	var totalQuestions = len(csvData)
	questions := make(map[int]string, totalQuestions)
	answers := make(map[int]string, totalQuestions)
	responses := make(map[int]string, totalQuestions)

	// Parse CSV data and populate maps
	for i, data := range csvData {
		questions[i] = data[0]
		answers[i] = data[1]
	}

	respondTo := make(chan string)

	// Wait for user to press Enter to start the quiz
	fmt.Println("Press [Enter] to start test.")
	bufio.NewScanner(os.Stdin).Scan()

	// Randomize question order if specified
	if flagRandom {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	randPool := rand.Perm(totalQuestions)

	// Start quiz session
	wg.Add(1)
	timeUp := time.After(time.Second * time.Duration(flagTime))
	go func() {
	label:
		for i := 0; i < totalQuestions; i++ {
			index := randPool[i]
			go askQuestion(os.Stdout, os.Stdin, questions[index], respondTo)
			select {
			// Check if time is up
			case <-timeUp:
				fmt.Fprintln(os.Stderr, "\nTime up!")
				break label
			// Receive user response
			case ans, ok := <-respondTo:
				if ok {
					responses[index] = ans
				} else {
					break label
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()

	// Evaluate user responses
	correct := 0
	for i := 0; i < totalQuestions; i++ {
		if checkAnswer(answers[i], responses[i]) {
			correct++
		}
	}

	// Display quiz summary
	summary(correct, totalQuestions)
}

// Function to ask a question and collect user response
func askQuestion(w io.Writer, r io.Reader, question string, replyTo chan string) {
	reader := bufio.NewReader(r)
	fmt.Fprintln(w, "Question: "+question)
	fmt.Fprint(w, "Answer: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		close(replyTo)
		if err == io.EOF {
			return
		}
		log.Fatalln(err)
	}
	replyTo <- strings.TrimSpace(answer)
}

// Function to check if the user's response matches the expected answer
func checkAnswer(ans string, expected string) bool {
	if strings.EqualFold(ans, expected) {
		return true
	}
	return false
}

// Function to display quiz summary
func summary(correct, totalQuestions int) {
	fmt.Fprintf(os.Stdout, "You answered %d questions correctly (%d / %d)\n", correct,
		correct, totalQuestions)
}
