package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Trivia struct {
	question string
	answer   int
}

type Quiz struct {
	trivias   []Trivia
	score     int
	index     int
	timeLimit int
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A CSV file in the format of 'question,answer'")
	timelimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	records, err := readCSV(*csvFilename)
	checkError(err)

	trivias, err := recordsToTrivia(records)
	checkError(err)

	quiz := Quiz{
		trivias: trivias,
		score:   0,
		index:   0,
	}

	quiz_channel := make(chan string)
	timer_channel := make(chan string)

	go playQuiz(quiz_channel, &quiz)
	go timer(timer_channel, *timelimit)

	select {
	case msg1 := <-quiz_channel:
		fmt.Println(msg1)
		fmt.Printf("You scored %d out of %d\n", quiz.score, len(quiz.trivias))
	case msg2 := <-timer_channel:
		fmt.Println(msg2)
		fmt.Printf("\nYou scored %d out of %d\n", quiz.score, len(quiz.trivias))
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)

	// if we os.Open returns an error then handle it
	if err != nil {
		return [][]string{}, err
	}

	// This will close the file for us when the function returns
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil

}

func recordsToTrivia(records [][]string) ([]Trivia, error) {
	var trivias []Trivia

	for _, record := range records {
		answer, err := strconv.Atoi(record[1])

		if err != nil {
			return []Trivia{}, err
		}

		trivia := Trivia{
			question: record[0],
			answer:   answer,
		}

		trivias = append(trivias, trivia)
	}

	return trivias, nil
}

func playQuiz(quiz_channel chan string, quiz *Quiz) {

	for quiz.index < len(quiz.trivias) {

		trivia := quiz.trivias[quiz.index]
		fmt.Printf("Problem #%d: %s = ", quiz.index+1, trivia.question)

		var answer string
		fmt.Scanln(&answer)

		answerInt, err := strconv.Atoi(answer)
		checkError(err)

		if answerInt == trivia.answer {
			quiz.score++
		}
		quiz.index++
	}
	quiz_channel <- "Quiz finished"
}

func timer(timer_channel chan string, timelimit int) {
	time.Sleep(time.Duration(timelimit) * time.Second)
	timer_channel <- "Time is up"
}
