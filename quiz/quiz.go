package quiz

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func AskQuestion(w io.Writer, r io.Reader, question string, replyTo chan string) {
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
	replyTo <- answer
}

func EvaluateResponses(answers, responses map[int]string) int {
	correct := 0
	for i := 0; i < len(answers); i++ {
		if CheckAnswer(answers[i], responses[i]) {
			correct++
		}
	}
	return correct
}

func CheckAnswer(ans string, expected string) bool {
	return strings.EqualFold(strings.TrimSpace(ans), strings.TrimSpace(expected))
}

func Summary(correct, totalQuestions int) {
	fmt.Fprintf(os.Stdout, "You answered %d questions correctly (%d / %d)\n", correct, correct, totalQuestions)
}
