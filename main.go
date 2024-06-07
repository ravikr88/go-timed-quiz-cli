package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ravikr88/go-timed-quiz-cli/quiz"
	"github.com/ravikr88/go-timed-quiz-cli/util"
)

var (
	flagFilePath string
	flagRandom   bool
	flagTime     int
	wg           sync.WaitGroup
)

func init() {
	flag.StringVar(&flagFilePath, "file", "questions.csv", "path/to/csv_file")
	flag.BoolVar(&flagRandom, "random", true, "randomize order of questions")
	flag.IntVar(&flagTime, "time", 10, "test duration")
	flag.Parse()
}

func main() {
	csvPath, err := filepath.Abs(flagFilePath)
	if err != nil {
		log.Fatalln("Unable to parse path" + csvPath)
	}
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	questions, answers, err := util.ReadCSV(file)
	if err != nil {
		log.Fatalln(err)
	}

	var totalQuestions = len(questions)
	responses := make(map[int]string, totalQuestions)
	respondTo := make(chan string)

	fmt.Println("Press [Enter] to start test.")
	util.WaitForEnter()

	if flagRandom {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	randPool := rand.Perm(totalQuestions)

	wg.Add(1)
	timeUp := time.After(time.Second * time.Duration(flagTime))
	go func() {
	label:
		for i := 0; i < totalQuestions; i++ {
			index := randPool[i]
			go quiz.AskQuestion(os.Stdout, os.Stdin, questions[index], respondTo)
			select {
			case <-timeUp:
				fmt.Fprintln(os.Stderr, "\nTime up!")
				break label
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

	correct := quiz.EvaluateResponses(answers, responses)
	quiz.Summary(correct, totalQuestions)
}
