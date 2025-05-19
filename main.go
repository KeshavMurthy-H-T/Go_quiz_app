package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	//1. get the input the name of the file
	fName := flag.String("f", "quiz.csv", "path of the csv file")
	//2. set the duration of the timer
	timer := flag.Int("t", 30, "tier for the quiz")
	flag.Parse()
	fmt.Println("Starting the Quiz....")
	//3. pull the problems from the the file (calling our problems pulller func )
	problems, err := problemPuller(*fName)
	//4. handle the error
	if err != nil {
		exit(fmt.Sprintf("Something went wrong:%s", err.Error()))
	}
	//5. create a variable to count our correct answers
	correctAns := 0
	//6. using the duration of the timer, we want to initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)
	//7. loop through the problem, print the question, we will accept the answers

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.question)

		go func() {

			fmt.Scanln(&answer)
			ansC <- answer

		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop

		case iAns := <-ansC:
			if iAns == p.answer {
				correctAns++
			}

			if i == len(problems)-1 {
				close(ansC)
			}

		}
	}
	//8. we will calculate and print out the result
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit")
	<-ansC
}

func problemPuller(fileName string) ([]problem, error) {
	// read all the problems from the quiz.csv

	//1. open the file
	if fObj, err := os.Open(fileName); err == nil {
		//2. we will create a new reader
		csvR := csv.NewReader(fObj)
		//3. it will need to read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			//4. call the parser problem function
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv"+" format from %s file: %s", fileName, err.Error())
		}

	} else {
		return nil, fmt.Errorf("error in opening %s file: %s", fileName, err.Error())
	}

}

func parseProblem(lines [][]string) []problem {
	// go over the lines and parse them, with problem struct

	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{question: lines[i][0], answer: lines[i][1]}
	}
	return r

}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
